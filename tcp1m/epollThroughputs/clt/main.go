package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/falconray0704/u4go/sysCfg"
	slog "github.com/falconray0704/u4go/sysLogger"
	"github.com/rcrowley/go-metrics"
	"log"
	"net"
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
	//sysClean = initSys()
}


var (
	ip          = flag.String("ip", "127.0.0.1", "server IP")
	connections = flag.Int("conn", 1, "number of tcp connections")
	startMetric = flag.String("sm", time.Now().Format("2006-01-02T15:04:05 -0700"), "start time point of all clients")
)

var (
	opsRate = metrics.NewRegisteredTimer("ops", nil)
)
var epoller *epoll

func main() {
	flag.Parse()

	setLimit()

	go func() {
		startPoint, err := time.Parse("2006-01-02T15:04:05 -0700", *startMetric)
		if err != nil {
			panic(err)
		}
		time.Sleep(startPoint.Sub(time.Now()))

		metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
	}()

	var err error
	epoller, err = MkEpoll()
	if err != nil {
		panic(err)
	}

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
		if err := epoller.Add(c); err != nil {
			log.Printf("failed to add connection %v", err)
			c.Close()
		}
		conns = append(conns, c)
	}

	log.Printf("Accomplished %d connections", len(conns))

	tts := time.Second
	if *connections > 100 {
		tts = time.Millisecond * 5
	}

	go start()

	time.Sleep(time.Minute * 3)
	log.Println("--- streaming start ---")
	for i := 0; i < len(conns); i++ {
		time.Sleep(tts)
		conn := conns[i]
		err = binary.Write(conn, binary.BigEndian, time.Now().UnixNano())
		if err != nil {
			log.Printf("failed to write timestamp %v", err)
			if err := epoller.Remove(conn); err != nil {
				if err := epoller.Remove(conn); err != nil {
					log.Printf("failed to remove %v", err)
				}
			}
		}
	}

	select {}
}

func start() {
	var nano int64
	initNano := time.Now().UnixNano()
	log.Printf("------------ start(), initNano:%d \n", initNano)
	for {
		epConns, err := epoller.Wait()
		if err != nil {
			log.Printf("failed to epoll wait %v", err)
			continue
		}
		/*
		if cnt != 0 {
			log.Printf("------------ start() looping wait, initNano:%d len:%d cnt:%d\n", initNano, len(epConns), cnt)
		}
		*/
		for _, conn := range epConns {
			//log.Printf("------------ start() looping connections, initNano:%d \n", initNano)
			if conn == nil {
				break
			}

			if err := binary.Read(conn, binary.BigEndian, &nano); err != nil {
				log.Printf("failed to read %v", err)
				if err := epoller.Remove(conn); err != nil {
					log.Printf("failed to remove %v", err)
				}

				conn.Close()
				continue
			} else {
				/*
				if idx % 5000 == 0 {
					log.Println("latency: ", time.Now().UnixNano() - nano)
				}
				*/
				opsRate.Update(time.Duration(time.Now().UnixNano() - nano))
				//opsRate.Update(time.Millisecond * 500)
			}

			//log.Printf("latency2: %d", time.Now().UnixNano() - initNano)

			err = binary.Write(conn, binary.BigEndian, time.Now().UnixNano())
			if err != nil {
				log.Printf("failed to write %v", err)
				if err := epoller.Remove(conn); err != nil {
					log.Printf("failed to remove %v", err)
				}
				conn.Close()
			}
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
}
