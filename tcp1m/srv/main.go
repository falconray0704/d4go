package main

import (
	"github.com/falconray0704/u4go/sysCfg"
	slog "github.com/falconray0704/u4go/sysLogger"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func initSys() (clean func()) {
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

	ln, err := net.Listen("tcp", ":8972")
	if err != nil {
		panic(err)
	}
	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			slog.Fatal("pprof failed: ", zap.Error(err))
		}
	}()
	var connections []net.Conn
	defer func() {
		for _, conn := range connections {
			conn.Close()
		}
	}()

	for {
		conn, e := ln.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				slog.Error("accept temp err: ", zap.Error(ne))
				continue
			}
			slog.Error("accept err: ", zap.Error(e))
			return
		}

		go handleConn(conn)
		connections = append(connections, conn)
		if len(connections)%200 == 0 {
			slog.Info("total connections: ", zap.Int("conCnt", len(connections)))
		}
	}
}

func handleConn(conn net.Conn) {
	io.Copy(ioutil.Discard, conn)
}
