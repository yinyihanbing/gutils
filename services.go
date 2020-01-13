package goutils

import (
	"github.com/kardianos/service"
	"fmt"
	"os"
)

// 服务
func RunGoService(serv service.Interface, name string) {
	serviceName := fmt.Sprintf("%v_%v", name, EncodeMd5(CurrentPath()))
	svcConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: fmt.Sprintf("This is an %v service.", name),
	}

	s, err := service.New(serv, svcConfig)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	logger, err := s.Logger(nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "install":
		err = s.Install()
	case "uninstall":
		err = s.Uninstall()
	case "start":
		err = s.Start()
	case "stop":
		err = s.Stop()
	case "restart":
		err = s.Restart()
	default:
		err = s.Run()
	}
	if err != nil {
		logger.Errorf("Failed to %s: %s\n", cmd, err)
	} else {
		if cmd != "" {
			logger.Infof("%v success", cmd)
		}
	}
}
