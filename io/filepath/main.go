package main

import (
	"fmt"
	"github.com/falconray0704/u4go/sysCfg"
	slog "github.com/falconray0704/u4go/sysLogger"
	"log"
	"path/filepath"
)

func initSys() (clean func()){
	var (
		errOnce error
	)

	sysLoggerCfg := slog.NewSysLogCfg()
	errOnce = sysCfg.LoadFileCfgs("./sysDatas/cfgs/appCfgs.yaml", "sysLogger", &sysLoggerCfg)
	if errOnce != nil {
		log.Fatalf("Loading system configs fail:%s\n", errOnce.Error())
	}
	_, _, err := slog.Init(sysLoggerCfg)
	if err != nil {
		log.Fatalf("Init system logger fail: %s.\n", err.Error())
	}
	slog.Debug("System init finished.")

	return func() {
		slog.Sync()
		slog.Close()
	}

}

var (
	sysClean func()
)

func init() {
	sysClean = initSys()
}

func main() {
	defer sysClean()

	fmt.Println("------- filepath.Dir()")
	fmt.Println(filepath.Dir("/home/ray/hello.go"))
	fmt.Println(filepath.Dir("/home/ray/hello/"))
	fmt.Println(filepath.Dir("/home/ray/hello"))
	fmt.Println(filepath.Dir(""))

	fmt.Println("------- filepath.Base()")
	fmt.Println(filepath.Base("/home/ray/hello.go"))
	fmt.Println(filepath.Base("/home/ray/hello/"))
	fmt.Println(filepath.Base("/home/ray/hello"))
	fmt.Println(filepath.Base(""))

	fmt.Println("------- filepath.Rel() -------")
	fmt.Println(filepath.Rel("/home/ray/example", "/home/ray/example/src/logic/topic.go"))
	fmt.Println(filepath.Rel("/home/ray/example", "/data/example"))

}

