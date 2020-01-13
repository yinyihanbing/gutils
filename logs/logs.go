package logs

import (
	"encoding/json"
	"os"
	"path"
)

type LogConfig struct {
	FileName string   `json:"filename"`
	MaxLines int      `json:"maxLines"`
	MaxSize  int      `json:"maxsize"`
	Daily    bool     `json:"daily"`
	MaxDays  int      `json:"maxDays"`
	Rotate   bool     `json:"rotate"`
	Perm     string   `json:"perm"`
	Separate []string `json:"separate"`
	path     string
	debug    bool
	level    int
}

type Option func(*LogConfig)

func Init(opts ...Option) {
	c := LogConfig{MaxDays: 30, Daily: true, Perm: "0600"}

	for _, option := range opts {
		option(&c)
	}

	if len(c.path) > 0 {
		if _, err := os.Stat(c.path); err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(c.path, os.ModePerm)
			}
			if err != nil {
				panic(err)
			}
		}
		c.FileName = path.Join(c.path, c.FileName)
	}

	if c.debug {
		Reset()
		if err := SetLogger(AdapterConsole, ""); err != nil {
			panic(err)
		}
		SetLogFuncCall(true)
	} else {
		SetLevel(c.level)
	}

	data, err := json.Marshal(&c)
	if err != nil {
		panic(err)
	}

	if err := SetLogger(AdapterMultiFile, string(data)); err != nil {
		panic(err)
	}
}

func WithPath(path string) Option {
	return func(opt *LogConfig) {
		opt.path = path
	}
}

func WithFileName(fileName string) Option {
	return func(opt *LogConfig) {
		opt.FileName = fileName
	}
}

func WithMaxLines(maxLines int) Option {
	return func(opt *LogConfig) {
		opt.MaxLines = maxLines
	}
}

func WithMaxSize(maxSize int) Option {
	return func(opt *LogConfig) {
		opt.MaxSize = maxSize
	}
}

func WithDaily(daily bool) Option {
	return func(opt *LogConfig) {
		opt.Daily = daily
	}
}

func WithMaxDays(maxDays int) Option {
	return func(opt *LogConfig) {
		opt.MaxDays = maxDays
	}
}

func WithRotate(rotate bool) Option {
	return func(opt *LogConfig) {
		opt.Rotate = rotate
	}
}

func WithPerm(perm string) Option {
	return func(opt *LogConfig) {
		opt.Perm = perm
	}
}

func WithSeparate(separate []string) Option {
	return func(opt *LogConfig) {
		opt.Separate = separate
	}
}

func WithLevel(level int) Option {
	return func(opt *LogConfig) {
		opt.level = level
	}
}

func WithDebug(debug bool) Option {
	return func(opt *LogConfig) {
		opt.debug = debug
	}
}
