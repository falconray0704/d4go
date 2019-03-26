package main

import (
	"flag"
	"fmt"
	"github.com/falconray0704/u4go/sysCfg"
	slog "github.com/falconray0704/u4go/sysLogger"
	"log"
	"net"
	"time"
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
	//sysClean = initSys()
}

var (
	ip          = flag.String("ip", "127.0.0.1", "server IP")
	connections = flag.Int("conn", 1, "number of tcp connections")
)

func main() {
	flag.Parse()
	addr := *ip + ":8972"
	log.Printf("Connect to %s", addr)
	var conns []net.Conn
	for i := 0; i < *connections; i++ {
		c, err := net.DialTimeout("tcp", addr, 10*time.Second)
		if err != nil {
			fmt.Println("failed to connect", i, err)
			i--
			continue
		}
		conns = append(conns, c)
		time.Sleep(time.Millisecond)
	}
	defer func() {
		for _, c := range conns {
			c.Close()
		}
	}()
	log.Printf("Accomplished %d connections", len(conns))
	tts := time.Second
	if *connections > 100 {
		tts = time.Millisecond * 0
	}
	for {
		for i := 0; i < len(conns); i++ {
			time.Sleep(tts)
			conn := conns[i]
			conn.Write([]byte("hello world\r\n"))
		}
	}
}
