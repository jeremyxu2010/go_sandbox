package main

import (
	stdlog "log"
	"github.com/cxr29/log"
	"github.com/pkg/errors"
)

func main() {
	w, err := log.OpenDaily("test.log")
	if err != nil {
		panic(err)
	}
	logger := log.New(w, "[test] ", stdlog.LstdFlags | stdlog.Llongfile | stdlog.LUTC,  log.LevelNotice)
	logger.Info("xxx")
	logger.Error("yyy")
	logger.ErrWarning(errors.New("some error"))
	logger.ErrWarning(nil)
	logger.Errorf("%s", "zzz")
	logger.Fatal("qqq")
}