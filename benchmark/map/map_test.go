package map_test

import (
	"github.com/falconray0704/u4go/sysCfg"
	slog "github.com/falconray0704/u4go/sysLogger"
	"log"
	"strconv"
	"testing"
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

func BenchmarkMapInt2String(b *testing.B) {
	mm := make(map[int]int, b.N)
	km := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		val := i //strconv.Itoa(i)
		mm[i] = val
		km[i] = val
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := km[i]
		if i != mm[k] {
			break
		}
	}
}

func BenchmarkMapString2String(b *testing.B) {
	mm := make(map[string]int, b.N)
	km := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		val := strconv.Itoa(i)
		km[i] = val
		mm[val] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i != mm[km[i]] {
			break
		}
	}
}





