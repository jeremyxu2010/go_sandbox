package gokitlogadapter


import (
	simplelog "github.com/jeremyxu2010/log"
	stdlog "log"
	"github.com/go-kit/kit/log"
	"os"
	"time"
	"bytes"
	"sync"
	"github.com/go-logfmt/logfmt"
)

var (
	simpleLogger *simplelog.Logger
)

func main() {
	initSimpleLog()
	logger := getGoKitLogger()
	logger = log.With(logger, "ts", log.TimestampFormat(time.Now, "2006-01-02 15:04:05.000"), "caller", log.DefaultCaller)
	logger.Log("msg", "xxxx");
}

func initSimpleLog(){
	simplelog.AddBracket()
	simpleLogger = simplelog.New(os.Stdout, "[test] ", stdlog.LstdFlags | stdlog.Lshortfile, simplelog.LevelNotice)
}

type GoKitLogAdapter struct {
	simpleLogger *simplelog.Logger
}

type logfmtEncoder struct {
	*logfmt.Encoder
	buf bytes.Buffer
}

func (l *logfmtEncoder) Reset() {
	l.Encoder.Reset()
	l.buf.Reset()
}

var logfmtEncoderPool = sync.Pool{
	New: func() interface{} {
		var enc logfmtEncoder
		enc.Encoder = logfmt.NewEncoder(&enc.buf)
		return &enc
	},
}

func (a *GoKitLogAdapter) Log(keyvals ...interface{}) error {
	enc := logfmtEncoderPool.Get().(*logfmtEncoder)
	enc.Reset()
	defer logfmtEncoderPool.Put(enc)

	if err := enc.EncodeKeyvals(keyvals...); err != nil {
		return err
	}

	// Add newline to the end of the buffer
	if err := enc.EndRecord(); err != nil {
		return err
	}
	simpleLogger.Error(enc.buf.String())
	return nil
}

func getGoKitLogger()log.Logger{
	return &GoKitLogAdapter{simpleLogger: simpleLogger}
}
