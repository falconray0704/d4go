package httpHandler_test

import (
	"github.com/falconray0704/d4go/unitTest/httpHandler"
	"github.com/falconray0704/u4go/sysCfg"
	slog "github.com/falconray0704/u4go/sysLogger"
	"log"
	"net/http"
	"net/http/httptest"
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
	httpHandler.Routes()
}

func TestSendJSON(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/sendjson", nil)
	if err != nil {
		t.Fatal("Create request fail.")
	}

	rw := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rw, req)

	log.Println("code:", rw.Code)
	log.Println("body:", rw.Body.String())
}


