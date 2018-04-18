package main

import (
	stdlog "log"
	"github.com/jeremyxu2010/log"
	"github.com/pkg/errors"
)

func main() {
	w, err := log.OpenDaily("test.log")
	if err != nil {
		panic(err)
	}
	log.AddBracket()
	logger := log.New(w, "[test] ", stdlog.LstdFlags | stdlog.Lshortfile, log.LevelNotice)
	logger.Info("xxx")
	logger.Error(errors.New("xxx"))
	logger.ErrWarning(errors.New("some error"))
	logger.ErrWarning(nil)
	logger.Errorf("%s", "zzz")
	logger.Fatal("qqq")
}