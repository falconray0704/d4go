package main

import (
	"fmt"
	"github.com/falconray0704/u4go/sysCfg"
	slog "github.com/falconray0704/u4go/sysLogger"
	"go.uber.org/zap"
	"log"
	"net/rpc/jsonrpc"
	"os"

	. "github.com/falconray0704/d4go/RPC"
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

	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "server:port")
		os.Exit(1)
	}
	service := os.Args[1]

	client, err := jsonrpc.Dial("tcp", service)
	if err != nil {
		slog.Fatal("Dial error", zap.String("err", err.Error()))
	}

	// Synchronous call
	args := Args{17, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		slog.Fatal("arith error", zap.String("err", err.Error()))
	}
	slog.Info("Arith: A * B = ", zap.Int("", reply))

	var quot Quotient
	err = client.Call("Arith.Divide", args, &quot)
	if err != nil {
		slog.Fatal("arith error", zap.String("err", err.Error()))
	}
	slog.Info("Arith: A / B = ", zap.Int("", quot.Quo))
	slog.Info("Arith: A % B = ", zap.Int("", quot.Rem))

}

