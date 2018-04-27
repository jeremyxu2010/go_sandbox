package main


import (
	simplelog "github.com/jeremyxu2010/log"
	stdlog "log"
	golog "github.com/go-log/log"
	"os"
	"github.com/micro/go-log"
)

var (
	simpleLogger *simplelog.Logger
)

func main() {
	initSimpleLog()
	logger := NewGoLogLoggerAdapter()

	log.SetLogger(logger)
	log.Log("xxxx")
}

func initSimpleLog(){
	simplelog.AddBracket()
	simpleLogger = simplelog.New(os.Stdout, "[test] ", stdlog.LstdFlags | stdlog.Lshortfile, simplelog.LevelNotice)
}

type GoLogAdapter struct {
	logger *simplelog.Logger
}

func (a *GoLogAdapter) Log(v ...interface{}) {
	a.logger.Error(v...)
}

func (a *GoLogAdapter) Logf(format string, v ...interface{}) {
	a.logger.Errorf(format, v...)
}

func NewGoLogLoggerAdapter() golog.Logger{
	return &GoLogAdapter{logger: simpleLogger}
}
