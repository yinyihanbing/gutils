package rpc

import (
	"sync"
	"fmt"
	"time"
	"errors"
	"runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"golang.org/x/net/context"

	"github.com/yinyihanbing/gutils/logs"
)

// RPC客户端
type RpcCli struct {
	wg              sync.WaitGroup
	conn            sync.Map                              // 连接列表 k=地址, v=*connInfo
	f               func(cc *grpc.ClientConn) interface{} // 执行函数
	queueLimitCount int                                   // 异步执行函数队列上限数, 超过上限必需等待
}

// 连接信息1
type connInfo struct {
	addr       string                                               // 远程服务地址
	cc         *grpc.ClientConn                                     // 客户端连接
	sc         interface{}                                          // 远程服务客户端
	chanAsyncF chan func(ctx context.Context, sc interface{}) error // 异步执行函数
}

// 启动队列执行任务
func (c *connInfo) startQueueTask() {
	for {
		select {
		case f := <-c.chanAsyncF:
			if f == nil && len(c.chanAsyncF) == 0 {
				logs.Info("rpc queue stopped: %s", c.addr)
				return
			}
			c.exec(f)
		}
	}
}

// 执行
func (c *connInfo) exec(f func(ctx context.Context, sc interface{}) error) {
	defer panicError()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	if err := f(ctx, c.sc); err != nil {
		cancel()
		state := c.cc.GetState()
		switch state {
		case connectivity.Connecting, connectivity.TransientFailure, connectivity.Shutdown:
			tryExecTimeInterval := []int{3, 5, 10, 30, 60, 60, 60, 60}
			for i := 0; i < len(tryExecTimeInterval); i++ {
				waitTime := tryExecTimeInterval[i]
				logs.Error("cannot connect rpc server:%v, execute again after trying to reconnect...(%ds)(%d/%d)", c.addr, waitTime, i+1, len(tryExecTimeInterval))
				time.Sleep(time.Second * time.Duration(waitTime))

				tryCtx, tryCtxCancel := context.WithTimeout(context.Background(), 10*time.Minute)
				if err := f(tryCtx, c.sc); err != nil {
					tryCtxCancel()
					switch state {
					case connectivity.Connecting, connectivity.TransientFailure, connectivity.Shutdown:
					default:
						logs.Error(err)
					}
				} else {
					tryCtxCancel()
					break
				}
			}
		default:
			logs.Error(err)
		}
	} else {
		cancel()
	}
}

// 新建RPC客户端
func NewRpcCli(f func(cc *grpc.ClientConn) interface{}) *RpcCli {
	cli := RpcCli{}
	cli.f = f
	cli.queueLimitCount = 100000
	return &cli
}

// 崩溃错误处理
func panicError() {
	if r := recover(); r != nil {
		buf := make([]byte, 4096)
		l := runtime.Stack(buf, false)
		err := fmt.Errorf("%v: %s", r, buf[:l])
		logs.Error(err)
	}
}

// 获取连接信息
func (this *RpcCli) getConn(addr string) (*connInfo, error) {
	// 查找缓存连接
	v, ok := this.conn.Load(addr)
	if ok {
		return v.(*connInfo), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 建立连接
	cc, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		err = errors.New(fmt.Sprintf("can not connect to rpc server, addr:%v, err:%v", addr, err))
		logs.Error(err)
		return nil, err
	}

	// 新连接信息
	conn := new(connInfo)
	conn.addr = addr
	conn.cc = cc
	conn.sc = this.f(cc)
	conn.chanAsyncF = make(chan func(ctx context.Context, sc interface{}) error, this.queueLimitCount)
	go conn.startQueueTask()

	this.conn.Store(conn.addr, conn)
	logs.Debug("new rpc cli conn, service addr=%s", addr)

	return conn, nil
}

// 同步执行( 连接服务器 - 执行函数 - 关闭连接 )
func (this *RpcCli) SyncCall(addr string, f func(ctx context.Context, sc interface{}) error) error {
	defer panicError()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 建立连接
	cc, err := grpc.DialContext(ctx, addr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		err = errors.New(fmt.Sprintf("can not connect to rpc server, addr:%v, err:%v", addr, err))
		logs.Error(err)
		return err
	}

	// 执行函数
	if err = f(ctx, this.f(cc)); err != nil {
		cc.Close()
		logs.Error(err)
		return err
	}

	// 关闭连接
	cc.Close()

	return nil
}

// 异步执行
func (this *RpcCli) AsyncCall(addr string, f func(ctx context.Context, sc interface{}) error) (err error) {
	defer panicError()

	logs.Debug("开始异步执行RPC, 地址=%v", addr)

	// 获取连接信息
	conn, err := this.getConn(addr)
	if err != nil || conn == nil {
		err := errors.New(fmt.Sprintf("unaccessible connection, addr=%s, error=%v", addr, err))
		logs.Error(err)
		return err
	}

	// 将要执行的函数放入队列
	conn.chanAsyncF <- f

	// 如果队列数量超过一半则报错
	if len(conn.chanAsyncF) > this.queueLimitCount/2 {
		logs.Error("rpc queue [%v] wait exec func count = %v", conn.addr, len(conn.chanAsyncF))
	}

	return nil
}