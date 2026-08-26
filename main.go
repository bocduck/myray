package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"flag"
	"net"
	"net/http"
	"strings"
)


func main() {
	var listenNetwork string
	var listenAddress string
	flag.StringVar(&listenAddress,"s","127.0.0.1:8080","Listen address")
	flag.Parse()

	if strings.HasPrefix(listenAddress, "unix/"){
		listenAddress = strings.TrimPrefix(listenAddress,"unix/")
		listenNetwork = "unix"
	}else{
		listenNetwork = "tcp"
	}

	ln, err := net.Listen(listenNetwork, listenAddress)
	if err != nil {
		panic(err)
	}
	fmt.Println("listen", listenNetwork, listenAddress)

	for {
		c, err := ln.Accept()
		if err != nil {
			continue
		}
		go handle(c)
	}
}

func handle(c net.Conn) {
	defer c.Close()

	br := bufio.NewReader(c)

	// 读取 HTTP Upgrade 请求。
	req, err := http.ReadRequest(br)
	if err != nil {
		return
	}

	if req.URL.Path != "/" ||
		!strings.EqualFold(req.Header.Get("Upgrade"), "websocket") ||
		!strings.Contains(strings.ToLower(req.Header.Get("Connection")), "upgrade") {
		return
	}

	// HTTP 101。
	_, err = io.WriteString(c,
		"HTTP/1.1 101 Switching Protocols\r\n"+
			"Connection: Upgrade\r\n"+
			"Upgrade: websocket\r\n"+
			"\r\n")
	if err != nil {
		return
	}

	// Upgrade 后，br 中可能已经缓存了客户端发送的 VLESS 数据。
	conn := &bufferedConn{
		Conn: c,
		r:    br,
	}

	if err := serveVLESS(conn); err != nil {
		fmt.Println("vless:", err)
	}
}

type bufferedConn struct {
	net.Conn
	r *bufio.Reader
}

func (c *bufferedConn) Read(p []byte) (int, error) {
	return c.r.Read(p)
}

// VLESS UUID 为 16 字节。
// 最简实现这里不做 UUID 校验，只读取并丢弃。
func serveVLESS(c net.Conn) error {
	// Version
	var version [1]byte
	if _, err := io.ReadFull(c, version[:]); err != nil {
		return err
	}

	// UUID
	var uuid [16]byte
	if _, err := io.ReadFull(c, uuid[:]); err != nil {
		return err
	}

	// Addon length
	var addonLen [1]byte
	if _, err := io.ReadFull(c, addonLen[:]); err != nil {
		return err
	}

	// Addon
	if _, err := io.CopyN(io.Discard, c, int64(addonLen[0])); err != nil {
		return err
	}

	// Command
	var cmd [1]byte
	if _, err := io.ReadFull(c, cmd[:]); err != nil {
		return err
	}

	// 1 = TCP, 2 = UDP, 3 = Mux
	if cmd[0] != 1 {
		return fmt.Errorf("unsupported command: %d", cmd[0])
	}

	// Port
	var portBuf [2]byte
	if _, err := io.ReadFull(c, portBuf[:]); err != nil {
		return err
	}
	port := binary.BigEndian.Uint16(portBuf[:])

	// Address type
	var atyp [1]byte
	if _, err := io.ReadFull(c, atyp[:]); err != nil {
		return err
	}

	var host string

	switch atyp[0] {
	case 1: // IPv4
		var ip [4]byte
		if _, err := io.ReadFull(c, ip[:]); err != nil {
			return err
		}
		host = net.IP(ip[:]).String()

	case 2: // Domain
		var n [1]byte
		if _, err := io.ReadFull(c, n[:]); err != nil {
			return err
		}

		name := make([]byte, n[0])
		if _, err := io.ReadFull(c, name); err != nil {
			return err
		}
		host = string(name)

	case 3: // IPv6
		var ip [16]byte
		if _, err := io.ReadFull(c, ip[:]); err != nil {
			return err
		}
		host = net.IP(ip[:]).String()

	default:
		return fmt.Errorf("unsupported address type: %d", atyp[0])
	}

	target := net.JoinHostPort(host, fmt.Sprint(port))
	fmt.Println("connect", target)

	dst, err := net.Dial("tcp", target)
	if err != nil {
		return err
	}
	defer dst.Close()

	// VLESS response:
	// Version + AddonLen + Addon
	if _, err := c.Write([]byte{version[0], 0}); err != nil {
		return err
	}

	// 双向转发。
	errCh := make(chan error, 2)

	go func() {
		_, err := io.Copy(dst, c)
		errCh <- err
	}()

	go func() {
		_, err := io.Copy(c, dst)
		errCh <- err
	}()

	return <-errCh
}
