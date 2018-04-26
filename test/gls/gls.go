package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jtolds/gls"
	"sync"
	"errors"
	"fmt"
)

var (
	srv = &service1{}
)

func main() {
	//初始化路由信息
	r := gin.New()
	r.GET("/test", handleTest)
	//r.Run(":13333")
	http.ListenAndServe(":13333", WrapSetRequestKeyHandler(r))
}

func handleTest(context *gin.Context) {
	s := srv.bizMethod()
	context.Writer.Write([]byte(s))
	context.Writer.Write([]byte(fmt.Sprintf("controller get request_key: '%s'\n", GetRequestKey())))
}

type service1 struct{
}

func (s *service1) bizMethod() string{
	return fmt.Sprintf("service get request_key: '%s'\n", GetRequestKey())
}

func WrapSetRequestKeyHandler(handler http.Handler) http.Handler{
	ctxHelper := getCtxHelper()
	wrapper := &setRequestKeyHandler{
		ctxHelper: ctxHelper,
		next: handler,
	}
	return http.Handler(wrapper)
}

func GetRequestKey() string{
	ctxHelper := getCtxHelper()
	val, ok := ctxHelper.mgr.GetValue(ctxHelper.requestKey)
	if ok {
		return val.(string)
	} else {
		panic(errors.New("do not get request key"))
	}
}

func getCtxHelper() *contextHelper {
	once.Do(func() {
		ctxHelper = &contextHelper{
			mgr: gls.NewContextManager(),
			requestKey: gls.GenSym(),
		}
	})
	return ctxHelper
}

var (
	once sync.Once
	ctxHelper *contextHelper
)

type contextHelper struct {
	mgr *gls.ContextManager
	requestKey gls.ContextKey
}

type setRequestKeyHandler struct {
	ctxHelper *contextHelper
	next http.Handler
}

func (w *setRequestKeyHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request){

	key := "request_key from http.Request" // get key from http.Request

	w.ctxHelper.mgr.SetValues(gls.Values{w.ctxHelper.requestKey: key}, func(){
		w.next.ServeHTTP(writer, req)
	})
}





