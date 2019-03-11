package main

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"github.com/falconray0704/u4go/sysCfg"
	slog "github.com/falconray0704/u4go/sysLogger"
	"io"
	"log"
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

	fmt.Printf("------- sha1 string -------\n\n")
	TestString := "Hello world!"
	Sha1Inst := md5.New()
	Sha1Inst.Write([]byte(TestString))
	Result := Sha1Inst.Sum([]byte(""))
	fmt.Printf("%x\n", Result)

	fmt.Printf("------- sha1 file -------\n\n")
	infile, inerr := os.Open("./main.go")
	if inerr == nil {
		sha1h := sha1.New()
		io.Copy(sha1h, infile)
		fmt.Printf("sha1 of file ./main.go is: %x", sha1h.Sum([]byte("")))
	} else {
		fmt.Println(inerr)
		os.Exit(1)
	}
}

