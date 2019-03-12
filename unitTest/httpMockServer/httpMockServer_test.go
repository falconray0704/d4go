package httpMockServer_test

import (
	"github.com/falconray0704/d4go/unitTest/httpMockServer"
	"github.com/falconray0704/u4go/sysCfg"
	slog "github.com/falconray0704/u4go/sysLogger"
	"go.uber.org/zap"
	"io/ioutil"
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
	httpMockServer.Routes()
}

func mockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(httpMockServer.SendJSON))
}

func TestSendJSON(t *testing.T) {

	server := mockServer()
	defer server.Close()

	resq, err := http.Get(server.URL)

	if err != nil {
		t.Fatal("Create Get fail!")
	}
	defer resq.Body.Close()

	log.Println("code:", resq.StatusCode)
	json, err := ioutil.ReadAll(resq.Body)
	if err != nil {
		t.Fatal(err)
	}
	//log.Printf("body:%s\n", json)
	slog.Info("Get response:", zap.ByteString("json", json))

}


