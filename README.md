**This is AI work,use with caution**

# Go VLESS HTTPUpgrade Server

>以下代码是一个最简go vless httpupgrade transport 服务端实现，其以反代形式工作在安全的tls web服务器背后，且在反代层以path做了安全认证，与现有客户端兼容。审计有无安全问题、性能问题。

## Build&Usage
### Build
```
mkdir test
cd test
curl -o main.go https://raw.githubusercontent.com/bocduck/myray/refs/heads/main/main.txt
go mod init test
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w"
```
### Usage
#### Run Server
```
./test -net unix -addr @myray2
```
#### Setup Proxy
add proxy to your website block in caddy web server
```
handle_path /secret_path/* {
 reverse_proxy unix/@myray2
}
```
#### Secret Path
You may want replace secret_path with real secret, generate one using the code below
```
openssl rand 16 | basenc --base64url | tr -d =
```
### Debug
#### Check if raw server works
```
curl --abstract-unix-socket myray2 ws://test -v
```
#### Check if final server works
```
curl wss://example.com/secret_path -v
```
