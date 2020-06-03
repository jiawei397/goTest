package main

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	socks5 := flag.String("s", ":1081", "socks5 server addr")
	http := flag.String("h", ":1082", "http server addr")
	flag.DurationVar(&readTimeout, "t", time.Minute, "read timeout")
	flag.Parse()

	go func() {
		err := handleProxyBase(proxySocks5, *socks5)
		if err != nil {
			log.Fatal(err)
		}
	}()
	err := handleProxyBase(proxyHttp, *http)
	if err != nil {
		log.Fatal(err)
	}
}

const (
	socksVer5       = 5
	socksCmdConnect = 1

	cliIdVer     = 0
	cliIdNMethod = 1
	cliIdCmd     = 1
	cliIdType    = 3
	cliIdIP0     = 4
	cliIdDmLen   = 4
	cliIdDm0     = 5

	cliTypeIPv4 = 1
	cliTypeDm   = 3
	cliTypeIPv6 = 4

	cliLenIPv4   = 3 + 1 + net.IPv4len + 2
	cliLenIPv6   = 3 + 1 + net.IPv6len + 2
	cliLenDmBase = 3 + 1 + 1 + 2

	proxySocks5 = "socks5"
	proxyHttp   = "http"
)

var (
	readTimeout time.Duration
	bytePool    = &sync.Pool{New: func() interface{} {
		return make([]byte, 4108)
	}}

	errMode          = errors.New("must socks5 or http")
	errVer           = errors.New("socks5 version not supported")
	errAuthExtraData = errors.New("socks5 authentication get extra data")
	errCmd           = errors.New("socks5 command not supported")
	errAddrType      = errors.New("socks5 addr type not supported")
	errReqExtraData  = errors.New("socks5 request get extra data")

	established = []byte("HTTP/1.1 200 Connection Established\r\n\r\n")
	httpConnect = []byte("CONNECT")

	socksVer5Established = []byte{socksVer5, 0, 0, 1, 0, 0, 0, 0, 8, 0x43}
)

func handleProxyBase(mode, addr string) error {
	ser, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := ser.Accept()
		if err != nil {
			return err
		}
		go func() {
			if err := handleProxy(mode, conn); err != nil {
				conn.Close()
				if err != io.EOF {
					log.Println("handle socks client:", err)
				}
			}
		}()
	}
}

func handleProxy(mode string, conn net.Conn) error {
	var (
		n          = 0
		err        error
		notConnect = true
		rawAddr    string
		buf0       = bytePool.Get().([]byte)
	)
	defer bytePool.Put(buf0)

	switch mode {
	case proxySocks5:
		notConnect = false // socks5后面不需要发送
		err = handleCliShake(conn, buf0[:258])
		if err != nil {
			return err
		}
		rawAddr, err = getCliRequest(conn, buf0[:269])
		if err != nil {
			return err
		}
		if _, err = conn.Write(socksVer5Established); err != nil {
			return err
		}
	case proxyHttp:
		SetReadTimeout(conn)
		n, err = conn.Read(buf0)
		if err != nil {
			return err
		}

		var httpUrl *url.URL // 从第一行中得到有用信息
		for s, i := 0, 0; buf0[i] != '\n'; i++ {
			if buf0[i] == 0x20 {
				if s == 0 {
					if bytes.Compare(httpConnect, buf0[:i]) == 0 {
						conn.Write(established) // CONNECT模式,返回代理成功
						notConnect = false
					}
					s = i + 1
				} else {
					if httpUrl, err = url.Parse(string(buf0[s:i])); err != nil {
						return err
					}
					break
				}
			}
		}
		if httpUrl.Opaque == "" {
			rawAddr = httpUrl.Host
			if strings.LastIndex(httpUrl.Host, ":") < 0 {
				rawAddr += ":80" // 没有端口则默认80
			}
		} else {
			rawAddr = httpUrl.Scheme + ":" + httpUrl.Opaque
		}
	default:
		return errMode
	}
	log.Println("handle proxy:", mode, ",addr:", rawAddr)

	// 如果想FQ,可以自定义[remote net.Conn]接口的对象,并编写服务器
	// 自定义Read和Write方法,在里面实现加密传输,即可达到FQ需求
	// 参考https://github.com/shadowsocks/shadowsocks-go
	remote, err := net.DialTimeout("tcp", rawAddr, time.Minute)
	if err != nil {
		return err
	}
	if notConnect { // 只有http代理,非CONNECT模式
		remote.Write(buf0[:n])
	}
	go CopyBufferThenClose(conn, remote, buf0)
	buf := bytePool.Get().([]byte)
	CopyBufferThenClose(remote, conn, buf)
	bytePool.Put(buf)
	return nil
}

/*---------------------------socks5 handle------------------------------------*/
func handleCliShake(conn net.Conn, buf []byte) error {
	SetReadTimeout(conn)
	n, err := io.ReadFull(conn, buf[:cliIdNMethod+1])
	if err != nil {
		return err
	}
	if buf[cliIdVer] != socksVer5 {
		return errVer
	}
	msgLen := int(buf[cliIdNMethod]) + 2
	if n == msgLen {
	} else if n < msgLen {
		if _, err = io.ReadFull(conn, buf[n:msgLen]); err != nil {
			return err
		}
	} else {
		return errAuthExtraData
	}
	_, err = conn.Write(socksVer5Established[:2])
	return err
}

func getCliRequest(conn net.Conn, buf []byte) (string, error) {
	SetReadTimeout(conn)
	n, err := io.ReadFull(conn, buf[:cliIdDmLen+1])
	if err != nil {
		return "", err
	}
	if buf[cliIdVer] != socksVer5 {
		return "", errVer
	}
	if buf[cliIdCmd] != socksCmdConnect {
		return "", errCmd
	}

	reqLen := -1
	switch buf[cliIdType] {
	case cliTypeDm:
		reqLen = int(buf[cliIdDmLen]) + cliLenDmBase
	case cliTypeIPv4:
		reqLen = cliLenIPv4
	case cliTypeIPv6:
		reqLen = cliLenIPv6
	default:
		return "", errAddrType
	}

	if n == reqLen {
	} else if n < reqLen {
		if _, err = io.ReadFull(conn, buf[n:reqLen]); err != nil {
			return "", err
		}
	} else {
		return "", errReqExtraData
	}
	var host string
	switch buf[cliIdType] {
	case cliTypeDm:
		host = string(buf[cliIdDm0 : cliIdDm0+int(buf[cliIdDmLen])])
	case cliTypeIPv4:
		host = net.IP(buf[cliIdIP0 : cliIdIP0+net.IPv4len]).String()
	case cliTypeIPv6:
		host = net.IP(buf[cliIdIP0 : cliIdIP0+net.IPv6len]).String()
	}
	port := int64(uint16(buf[reqLen-1]) | uint16(buf[reqLen-2])<<8)
	host = net.JoinHostPort(host, strconv.FormatInt(port, 10))
	return host, nil
}

/*---------------------------------tools--------------------------------------*/
func SetReadTimeout(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(readTimeout))
}

func CopyBufferThenClose(dst, src net.Conn, buf []byte) {
	var (
		n      int
		er, ew error
	)
	defer dst.Close()
	for {
		SetReadTimeout(src)
		n, er = src.Read(buf)
		if n > 0 {
			if _, ew = dst.Write(buf[:n]); ew != nil {
				break
			}
		}
		if er != nil {
			break
		}
	}
}
