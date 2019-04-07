package main

import (
	"flag"
	"github.com/falconray0704/u4go/sysCfg"
	slog "github.com/falconray0704/u4go/sysLogger"
	"github.com/rcrowley/go-metrics"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"syscall"
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
	// sysClean = initSys()
}


var (
	c = flag.Int("c", 10, "concurrency")
	ec = flag.Int("ec", 4, "epoller concurrency")
)
var (
	opsRate = metrics.NewRegisteredMeter("ops", nil)
)

var epollers []*epoll
var workerPool *pool

func main() {
	var (
		idx int
	)
	flag.Parse()

	setLimit()
	go metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	ln, err := net.Listen("tcp", ":8972")
	if err != nil {
		panic(err)
	}

	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	workerPool = newPool(*c, 1000000)
	workerPool.start()

	epollers = make([]*epoll, *ec)

	for idx = 0; idx < *ec; idx++ {
		epollers[idx], err = MkEpoll()
		if err != nil {
			panic(err)
		}
		go start(epollers[idx])
	}

	//isSwap := false
	swapCnt := (10000 * (*c)) / (*ec)
	idx = 0

	for {
		conn, e := ln.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				log.Printf("accept temp err: %v", ne)
				continue
			}

			log.Printf("accept err: %v", e)
			return
		}

		if swapCnt == 0 {
			swapCnt = (10000 * (*c)) / (*ec)
			idx++
		} else {
				swapCnt--
		}

		if err := epollers[idx].Add(conn); err != nil {
			log.Printf("failed to add connection %v", err)
			conn.Close()
		}

	}

	workerPool.Close()
}

func start(epoller *epoll) {
	for {
		connections, err := epoller.Wait()
		if err != nil {
			log.Printf("failed to epoll wait %v", err)
			continue
		}
		for _, conn := range connections {
			if conn == nil {
				break
			}

			workerPool.addTask(conn, epoller)
		}
	}
}

func setLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	log.Printf("set cur limit: %d", rLimit.Cur)
}

