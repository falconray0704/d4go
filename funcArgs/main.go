package main

import (
	"fmt"
	"github.com/falconray0704/u4go/sysCfg"
	slog "github.com/falconray0704/u4go/sysLogger"
	"log"
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

	ageArr := []int{7, 9, 3, 5, 1}
	f1(ageArr...)

}

func f1(arr ...int) {
	f2(arr ...)
	fmt.Println("")
	f3(arr)
}

func f2(arr ...int) {
	for _, char := range arr {
		fmt.Printf("%d ", char)
	}
}

func f3(arr []int) {
	for _, char := range arr {
		fmt.Printf("%d ", char)
	}
}

