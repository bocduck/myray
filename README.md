**This is AI work,use with caution**

# Go VLESS HTTPUpgrade Server

>以下代码是一个最简go vless httpupgrade transport 服务端实现，其以反代形式工作在安全的tls web服务器背后，且在反代层以path做了安全认证，与现有客户端兼容。审计有无安全问题、性能问题。

## Usage - Quick Start
Quick start with free IPv4 & IPv6 cert if you don't have a domain name
```
curl https://raw.githubusercontent.com/bocduck/myray/refs/heads/main/install.sh | bash
```
Client connect via vless link output in console

## Usage - Manual
### Setup Caddy Web TLS Reverse Proxy
Caddyfile
```
YOUR_SERVER_DOMAIN_OR_IP {
 handle /secret_path {
  reverse_proxy unix/@myray
 }
}
```
You may want replace secret_path with real secret, generate one using the code below
```
openssl rand 16 | basenc --base64url | tr -d =
```
### Run Server
```
curl -o myray https://raw.githubusercontent.com/bocduck/myray/refs/heads/main/linux_amd64

chmod +x myray

nohup ./myray -s unix/@myray -path /secret_path >/dev/null 2>&1 &
```
### Client Connect
```
vless://00000000-0000-0000-0000-000000000000@example.com:443?security=tls&type=httpupgrade&path=%2F$secret_path%2F&sni=example.com&udp=0
```
## Debug
### Check if raw server works
```
curl --abstract-unix-socket myray ws://test -v
```
### Check if final server works
```
curl wss://example.com/secret_path -v
```
## Build
```
mkdir test
cd test
curl -o main.go https://raw.githubusercontent.com/bocduck/myray/refs/heads/main/main.go
go mod init test
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w"
```
