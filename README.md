# (Deprecated) Websocket Tunnel

[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/isayme/wstunnel?sort=semver&style=flat-square)](https://hub.docker.com/r/isayme/wstunnel)
![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/isayme/wstunnel?sort=semver&style=flat-square)
![Docker Pulls](https://img.shields.io/docker/pulls/isayme/wstunnel?style=flat-square)

Proxy a TCP protocol with websocket.
HTTP/HTTPS protocol not supported.

# Why or When

- proxy server that cannot reach directly, like socks5 but no client proxy config;
- access server with wss to avoid Man-in-the-middle attack;
- use cloudflare to protect/hide real server;

# Example Configure

- [server](./example/server.cfg.yaml)
- [local](./example/local.cfg.yaml)
