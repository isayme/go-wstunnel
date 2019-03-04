package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	logger "github.com/isayme/go-logger"
	"github.com/isayme/go-wstunnel/wstunnel"
	"golang.org/x/net/websocket"
)

var listenAddress = flag.String("listen", ":8388", "listen address")
var proxyAddress = flag.String("proxy", "", "proxy(remote) address")
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

	http.Handle("/", websocket.Server{
		Handshake: handshakeWebsocket,
		Handler:   handleWebsocket(*proxyAddress),
	})
	// http.Handle("/", websocket.Handler(handleWebsocket(*proxyAddress)))

	logger.Debugw("start listen", "address", listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		logger.Panicw("ListenAndServe fail", "err", err)
	}
}

func handshakeWebsocket(config *websocket.Config, req *http.Request) error {
	var err error
	config.Origin, err = websocket.Origin(config, req)
	if err == nil && config.Origin == nil {
		return fmt.Errorf("null origin")
	}
	return err
}

func handleWebsocket(address string) func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		logger.Debugw("new connection", "address", ws.RemoteAddr().String())
		defer ws.Close()

		conn, err := net.Dial("tcp", address)
		if err != nil {
			logger.Warnw("dial ssserver fail", "err", err)
			return
		}
		defer conn.Close()

		wstunnel.Proxy(ws, conn)

		logger.Debugw("connection close", "address", ws.RemoteAddr().String())
	}
}
