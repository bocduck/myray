this is AI work , use with caution

Build:
```
mkdir test
cd test
curl -o main.go https://raw.githubusercontent.com/bocduck/myray/refs/heads/main/main.txt
go mod init test
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w"
```

Usage:
```
./test -net unix -addr @myray2
```

add following proxy to caddy web server
```
handle_path /secret_path/* {
 reverse_proxy unix/@myray2
}
```

以下代码是一个最简go vless httpupgrade transport 服务端实现，其以反代形式工作在安全的tls web服务器背后，且在反代层以path做了安全认证，与现有客户端兼容。

审计有无安全问题、性能问题。
