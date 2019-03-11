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

	var (
		name string
		age	int
		n int
	)

	fmt.Println("--------- Sscan -------")
	n, _ = fmt.Sscan("张三 28", &name, &age)
	//slog.Info("Sscan result", zap.Int("n", n), zap.String("name", name), zap.Int("age", age))
	fmt.Println(n, name, age)
	n, _ = fmt.Sscan("张三\n28", &name, &age)
	//slog.Info("Sscan result", zap.Int("n", n), zap.String("name", name), zap.Int("age", age))
	fmt.Println(n, name, age)

	fmt.Println("--------- Sscanf -------")
	n, _ = fmt.Sscanf("张三 28", "%s%d", &name, &age)
	//slog.Info("Sscan result", zap.Int("n", n), zap.String("name", name), zap.Int("age", age))
	fmt.Println(n, name, age)
	n, _ = fmt.Sscan("张三\n28", "%s%d", &name, &age)
	//slog.Info("Sscan result", zap.Int("n", n), zap.String("name", name), zap.Int("age", age))
	fmt.Println(n, name, age)

	fmt.Println("--------- Sscanln -------")
	n, _ = fmt.Sscanln("张三 28", "%s%d", &name, &age)
	//slog.Info("Sscan result", zap.Int("n", n), zap.String("name", name), zap.Int("age", age))
	fmt.Println(n, name, age)
	n, _ = fmt.Sscan("张三\n28", "%s%d", &name, &age)
	//slog.Info("Sscan result", zap.Int("n", n), zap.String("name", name), zap.Int("age", age))
	fmt.Println(n, name, age)
}

