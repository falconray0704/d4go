package main

import (
	"encoding/base64"
	"github.com/falconray0704/u4go/sysCfg"
	slog "github.com/falconray0704/u4go/sysLogger"
	"go.uber.org/zap"
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

	msg := "Hello world"

	// encode msg
	encMsg := base64.StdEncoding.EncodeToString([]byte(msg))
	slog.Info("Encode msg", zap.String("encMsg", encMsg))

	// decode msg
	data, err := base64.StdEncoding.DecodeString(encMsg)

	if err != nil {
		slog.Fatal("Decode msg fail", zap.Error(err))
	} else {
		slog.Info("Decode msg", zap.ByteString("decMsg", data))
	}

}

