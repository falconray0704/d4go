package main

import (
	"fmt"
	"github.com/falconray0704/u4go/sysCfg"
	slog "github.com/falconray0704/u4go/sysLogger"
	"go.uber.org/zap"
	"log"
	"net"
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

func init() {
	sysClean = initSys()
}

func main() {
	defer sysClean()

	address := net.TCPAddr{
		IP: net.ParseIP("0.0.0.0"),
		Port: 8000,
	}
	listener, err := net.ListenTCP("tcp4", &address)
	if err != nil {
		slog.Fatal(err.Error())
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			slog.Fatal(err.Error())
		}
		slog.Info("New connection ", zap.String("Remote address:", conn.RemoteAddr().String()))
		go echo(conn)
	}
}

func echo(conn *net.TCPConn) {
	tick := time.Tick(5 * time.Second)
	for now := range tick {
		n, err := conn.Write([]byte(now.String()))
		if err != nil {
			slog.Error("Response write error", zap.String("Err", err.Error()))
			conn.Close()
			return
		}
		fmt.Printf("Send %d bytes to %s\n", n, conn.RemoteAddr())
	}
}

func checkError(err error) {
	if err != nil {
		slog.Fatal("", zap.String("Err", err.Error()))
	}
}
