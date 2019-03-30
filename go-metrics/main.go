package main

import (
	"github.com/falconray0704/u4go/sysCfg"
	slog "github.com/falconray0704/u4go/sysLogger"
	"github.com/rcrowley/go-metrics"
	"log"
	"os"
	"time"
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

var (
	opsRate = metrics.NewRegisteredTimer("ops", nil)
)

func init() {
	//sysClean = initSys()
}

func main() {
	//defer sysClean()

	go func() {
		/*
		startPoint, err := time.Parse("2006-01-02T15:04:05 -0700", *startMetric)
		if err != nil {
			panic(err)
		}
		time.Sleep(startPoint.Sub(time.Now()))
		*/
		//time.Sleep(time.Second * 10)

		//metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
		metrics.Log(metrics.DefaultRegistry, 2*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
	}()

	for {
		pre := time.Now().UnixNano()
		time.Sleep(time.Second)
		opsRate.Update(time.Duration(time.Now().UnixNano() - pre))
	}
}

