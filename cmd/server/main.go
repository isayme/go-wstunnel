package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	logger "github.com/isayme/go-logger"
	"github.com/isayme/go-wstunnel/wstunnel"
	"github.com/isayme/go-wstunnel/wstunnel/conf"
	"golang.org/x/net/websocket"
)

var showVersion = flag.Bool("v", false, "show version")

func main() {
	flag.Parse()
	if *showVersion {
		wstunnel.PrintVersion()
		os.Exit(0)
	}

	config := conf.Get()

	targets := map[string]string{}
	for _, service := range config.Services {
		targets[service.Name] = service.RemoteAddress
	}

	http.Handle("/", websocket.Server{
		Handshake: handshakeWebsocket,
		Handler:   handleWebsocket(targets),
	})

	listenAddress := config.Server.Addr
	logger.Debugw("start listen", "address", listenAddress)
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
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

// get serviceName from query
func getServiceName(ws *websocket.Conn) string {
	URI := ws.Request().URL
	return URI.Query().Get("serviceName")
}

func handleWebsocket(targets map[string]string) func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		defer ws.Close()

		serviceName := getServiceName(ws)
		address := targets[serviceName]
		logger.Debugw("new connection", "address", ws.RemoteAddr().String(), "serviceName", serviceName, "address", address)
		if address == "" {
			logger.Debugf("service(%s) not found", serviceName)
			return
		}

		conn, err := net.Dial("tcp", address)
		if err != nil {
			logger.Warnw("dial service fail", "err", err)
			return
		}
		defer conn.Close()

		wstunnel.Proxy(ws, conn)

		logger.Debugw("connection close", "address", ws.RemoteAddr().String())
	}
}
