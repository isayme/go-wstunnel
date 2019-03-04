package wstunnel

import (
	"io"
	"net"

	logger "github.com/isayme/go-logger"
	"github.com/pkg/errors"
)

func Proxy(client, server net.Conn) {
	defer client.Close()
	defer server.Close()

	// any of remote/client closed, the other one should close with quiet
	closed := false

	go func() {
		_, err := Copy(server, client)
		if err != nil && !closed {
			if errors.Cause(err) != io.EOF {
				logger.Errorf("[%s] Copy from client to server fail, err: %#v", server.RemoteAddr(), err)
			}
		}
		closed = true
		logger.Debug("client read end")
	}()

	_, err := Copy(client, server)
	if err != nil && !closed {
		if errors.Cause(err) != io.EOF {
			logger.Errorf("[%s] Copy from server to client fail, err: %#v", server.RemoteAddr(), err)
		}
	}
	closed = true
	logger.Debug("remote read end")
}
