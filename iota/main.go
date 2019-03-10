package main

import (
	"fmt"
	"go.uber.org/zap"
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

const (
	a		= iota					// a == 0
	b							// b == 1, implicit using iota, equal to b = iota
	c							// c == 2, equal to c = iota
	d, e, f = iota, iota, iota	// d = 3, e = 3, f = 3, same value in the same line, must use 3 iota
	g		= iota				// g == 4
	h		= "h"				// h == "h", assign value to h individually, iota still increase to 5
	i							// i == "h", using the previous value assigned, iota still increase to 6
	j		= iota				// j == 7
)

const z = iota 					// reset iota outside previous const, z == 0

func main() {
	defer sysClean()

	fmt.Println("iota demo:", a, b, c, d, e, f, g, h, i, j, z)

	slog.Info("iota demo:",
		zap.Int("a", a),
		zap.Int("b", b),
		zap.Int("c", c),
		zap.Int("d", d),
		zap.Int("e", e),
		zap.Int("f", f),
		zap.Int("g", g),
		zap.String("h", h),
		zap.String("i", i),
		zap.Int("j", j))

}

