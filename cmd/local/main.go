package local

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"sync"
	"time"

	logger "github.com/isayme/go-logger"
	"github.com/isayme/go-wstunnel/wstunnel"
	"github.com/isayme/go-wstunnel/wstunnel/conf"
	"golang.org/x/net/websocket"
)

func Run() {
	config := conf.Get()

	if len(config.Services) == 0 {
		logger.Warnw("no service configured")
		os.Exit(0)
	}

	logger.Debugf("ws addr: %s", config.Local.WebsocketAddr)

	for _, serviceConfig := range config.Services {
		go func(serviceConfig conf.ServiceConfig) {
			client, err := NewClient(config.Local.WebsocketAddr, serviceConfig)
			if err != nil {
				logger.Panicw("newClient fail", "err", err)
			}
			err = client.ListenAndServe()
			if err != nil {
				logger.Panicw("newClient fail", "err", err)
			}
		}(serviceConfig)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

type Client struct {
	name          string
	timeout       time.Duration
	listenAddress string
	wsAddress     string
	originAddress string
}

func NewClient(wsAddr string, config conf.ServiceConfig) (*Client, error) {
	client := &Client{
		name:          config.Name,
		timeout:       config.Timeout.Duration(),
		listenAddress: config.LocalAddress,
		wsAddress:     wsAddr,
	}
	if client.timeout <= 0 {
		client.timeout = time.Second * 60
	}

	URL, err := url.Parse(client.wsAddress)
	if err != nil {
		return nil, fmt.Errorf("url.Parse '%s' fail: %w", client.wsAddress, err)
	}

	switch URL.Scheme {
	case "ws":
		URL.Scheme = "http"
	case "wss":
		URL.Scheme = "https"
	default:
		logger.Panicw("URL.Scheme invalid", "schema", URL.Scheme, "ws", client.wsAddress)
	}
	client.originAddress = URL.String()

	return client, nil
}

func (c *Client) ListenAndServe() error {
	logger.Debugw("start listen", "address", c.listenAddress)
	l, err := net.Listen("tcp", c.listenAddress)
	if err != nil {
		logger.Panicw("listen fail", "err", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Warnw("accept fail", "err", err)
			continue
		}

		go c.handleConnection(conn)
	}
}

func (c *Client) handleConnection(conn net.Conn) {
	logger.Debugw("new connection", "address", conn.RemoteAddr().String())
	conn = wstunnel.NewTimeoutConn(conn, c.timeout)
	defer conn.Close()

	url := fmt.Sprintf("%s?serviceName=%s", c.wsAddress, c.name)
	ws, err := websocket.Dial(url, "", c.originAddress)
	if err != nil {
		logger.Warnw("dial websocket server fail", "url", c.wsAddress, "err", err)
		return
	}
	defer ws.Close()

	wstunnel.Proxy(conn, ws)

	logger.Debugw("connection close", "address", conn.RemoteAddr().String())
}
