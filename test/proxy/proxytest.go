package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/jeremyxu2010/goproxy"
)

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}

var alwaysHTTPMitm goproxy.FuncHttpsHandler = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	return goproxy.HTTPMitmConnect, host
}

var alwaysAccept goproxy.FuncHttpsHandler = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	return goproxy.OkConnect, host
}

func main() {
	proxy := goproxy.NewProxyHttpServer()
	//proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile(`^.*baidu\.com:80$`))).
	//	HandleConnect(goproxy.AlwaysReject)

	proxy.OnRequest().HandleConnect(alwaysAccept)
	//proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*:80$"))).
	//	HandleConnect(alwaysAccept)

	// enable curl -p for all hosts on port 80
	//proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*:80$"))).
	//	HijackConnect(func(req *http.Request, client net.Conn, ctx *goproxy.ProxyCtx) {
	//	defer func() {
	//		if e := recover(); e != nil {
	//			ctx.Logf("error connecting to remote: %v", e)
	//			client.Write([]byte("HTTP/1.1 500 Cannot reach destination\r\n\r\n"))
	//		}
	//		client.Close()
	//	}()
	//	clientBuf := bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client))
	//	remote, err := net.Dial("tcp", req.URL.Host)
	//	orPanic(err)
	//	remoteBuf := bufio.NewReadWriter(bufio.NewReader(remote), bufio.NewWriter(remote))
	//	for {
	//		req, err := http.ReadRequest(clientBuf.Reader)
	//		orPanic(err)
	//		orPanic(req.Write(remoteBuf))
	//		orPanic(remoteBuf.Flush())
	//		resp, err := http.ReadResponse(remoteBuf.Reader, req)
	//		orPanic(err)
	//		orPanic(resp.Write(clientBuf.Writer))
	//		orPanic(clientBuf.Flush())
	//	}
	//})
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8888", "proxy listen address")
	flag.Parse()
	proxy.Verbose = *verbose
	log.Fatal(http.ListenAndServe(*addr, proxy))
}