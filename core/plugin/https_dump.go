package plugin

import (
	"crypto/tls"
	"os"
	"net"
	"fmt"
	"reflect"
	"crypto/x509"
	"unsafe"
	"encoding/base64"
	"crypto/rsa"
	"crypto/ecdsa"
	"crypto/rand"
	"golang.org/x/crypto/pkcs12"
)

func HttpsDump(lc net.Conn, host string) {
	conf, sc := ClientToServer()
	conn := tls.Server(lc, conf)
	go send(conn, sc, false)
	go send(sc, conn, true)
}

func ClientToServer() (*tls.Config, net.Conn) {
	conf := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
		KeyLogWriter:       os.Stdout,
	}
	co, err := net.Dial("tcp", "www.baidu.com:443")
	if err != nil {
		panic(err)
	}
	conn := tls.Client(co, conf)
	fmt.Println("request => ")
	_, err = conn.Write([]byte("GET / HTTP/1.1\r\nHost: www.baidu.com\r\nUser-Agent: curl/7.54.0\r\nAccept: */*\r\n\r\n"))
	if err != nil {
		conn.Close()
		panic(err)
	}
	buf := make([]byte, 8190)
	var n int
	for {
		n, err = conn.Read(buf)
		if err != nil {
			conn.Close()
			break
		}
		fmt.Println(string(buf[:n]))
		if n < cap(buf) {
			break
		}
	}

	rt := reflect.TypeOf(conn).Elem()
	filed, ok := rt.FieldByName("peerCertificates")
	var cert *x509.Certificate
	if ok {
		ptr := (uintptr)(unsafe.Pointer(conn))
		cert = (*(*[]*x509.Certificate)(unsafe.Pointer(ptr + filed.Offset)))[0]
	}
	derBytes, privateKey := makeCert(cert)
	//conn.Close()
	if err != nil {
		panic(err)
	}
	conf = &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
		KeyLogWriter:       os.Stdout,
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{derBytes},
				PrivateKey:  privateKey,
			},
		},
	}
	return conf, conn
}
func makeCert(cert *x509.Certificate) ([]byte, interface{}) {
	passphrase := "9C454C51"
	p12Base64 := "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"
	buf, err := base64.StdEncoding.DecodeString(p12Base64)
	if err != nil {
		panic(err)
	}
	privateKey, certificate, err := pkcs12.Decode(buf, passphrase)
	if err != nil {
		panic(err)
	}
	var publicKey interface{}
	switch k := privateKey.(type) {
	case *rsa.PrivateKey:
		publicKey = &k.PublicKey
	case *ecdsa.PrivateKey:
		publicKey = &k.PublicKey
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, cert, certificate, publicKey, privateKey)
	if err != nil {
		panic(err)
	}
	return derBytes, privateKey
}
