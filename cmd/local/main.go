package main

import (
	"flag"
	"net"
	"net/url"
	"os"
	"time"

	logger "github.com/isayme/go-logger"
	"github.com/isayme/go-wstunnel/wstunnel"
	"golang.org/x/net/websocket"
)

var listenAddress = flag.String("listen", ":8388", "listen address")
var remoteAddress = flag.String("ws", "", "remote websocket tunnel address")
var clientTimeout = flag.String("timeout", "30s", "timeout for client connection")
var loggerLevel = flag.String("level", "info", "log level")
var showHelp = flag.Bool("h", false, "show help")
var showVersion = flag.Bool("v", false, "show version")

func main() {
	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *showVersion {
		wstunnel.PrintVersion()
		os.Exit(0)
	}

	logger.SetLevel(*loggerLevel)

	timeout, err := time.ParseDuration(*clientTimeout)
	if err != nil {
		logger.Panicw("time.ParseDuration fail", "err", err)
	}

	URL, err := url.Parse(*remoteAddress)
	if err != nil {
		logger.Panicw("url.Parse", "err", err)
	}
	switch URL.Scheme {
	case "ws":
		URL.Scheme = "http"
	case "wss":
		URL.Scheme = "https"
	default:
		logger.Panicw("URL.Scheme invalid", "sceham", URL.Scheme)
	}
	origin := URL.String()

	logger.Debugw("start listen", "address", listenAddress)
	l, err := net.Listen("tcp", *listenAddress)
	if err != nil {
		logger.Panicw("listen fail", "err", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Warnw("accept fail", "err", err)
			continue
		}

		handleConnection(conn, *remoteAddress, origin, timeout)
	}
}

func handleConnection(conn net.Conn, wsurl, origin string, timeout time.Duration) {
	logger.Debugw("new connection", "address", conn.RemoteAddr().String())
	conn = wstunnel.NewTimeoutConn(conn, timeout)
	defer conn.Close()

	ws, err := websocket.Dial(wsurl, "", origin)
	if err != nil {
		logger.Warnw("dial websocket server fail", "url", wsurl)
		return
	}
	defer ws.Close()

	wstunnel.Proxy(conn, ws)

	logger.Debugw("connection close", "address", conn.RemoteAddr().String())
}
