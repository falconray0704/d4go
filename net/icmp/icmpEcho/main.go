package main

import (
	"bytes"
	"fmt"
	"github.com/falconray0704/u4go/sysCfg"
	slog "github.com/falconray0704/u4go/sysLogger"
	"go.uber.org/zap"
	"io"
	"log"
	"net"
	"os"
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
		fmt.Println("Usage: ", os.Args[0], "hostname")
	}

	service := os.Args[1]

	conn, err := net.Dial("ip4:icmp", service)
	checkError(err)
	slog.Info("net.Dial() is done.")
	defer conn.Close()

	var msg[512]byte
	msg[0] = 8	// echo
	msg[1] = 0	// code 0
	msg[2] = 0	// checksum
	msg[3] = 0	// checksum
	msg[4] = 0	// identifier[0]
	msg[5] = 13	// identifier[1]
	msg[6] = 0	// sequence[0]
	msg[7] = 37 // sequence[1]
	len := 8
	check := checkSum(msg[0:len])
	msg[2] = byte(check >> 8)
	msg[3] = byte(check & 255)

	_, err = conn.Write(msg[0:len])
	checkError(err)
	slog.Info("conn.Write() is done")

	var resp [512]byte
	_, err = conn.Read(resp[0:])
	checkError(err)
	/*
	resp, err = readFully(conn)
	checkError(err)
	slog.Info("readFully() is done.")
	*/

	slog.Info("Got response")
	if resp[5] == 13 {
		slog.Info("Identifier matched!")
	}
	if resp[7] == 37 {
		slog.Info("Sequence matched!")
	}

}

func checkSum(msg []byte) uint16 {
	sum := 0

	for n := 1; n < len(msg) - 1; n += 2 {
		sum += int(msg[n])*256 + int(msg[n+1])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var answer uint16 = uint16(^sum)
	return answer
}

func checkError(err error) {
	if err != nil {
		slog.Fatal("", zap.String("Err", err.Error()))
	}
}

func readFully(conn net.Conn) ([]byte, error) {
	defer conn.Close()

	result := bytes.NewBuffer(nil)
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return result.Bytes(), nil
}


