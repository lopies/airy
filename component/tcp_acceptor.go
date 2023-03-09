// Copyright (c) nano Author and TFG Co and Airy. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package component

import (
	"crypto/tls"
	"fmt"
	"github.com/airy/config"
	"github.com/airy/constants"
	"github.com/airy/interfaces"
	"github.com/airy/logger"
	"github.com/airy/mgr"
	"net"
	"strconv"
)

// TCPAcceptor struct
type TCPAcceptor struct {
	addr     string
	connChan chan interfaces.PlayerConn
	listener net.Listener
	running  bool
	certFile string
	keyFile  string
	BaseComponent
}

type tcpPlayerConn struct {
	net.Conn
}

func NewTCPAcceptor() *TCPAcceptor {
	t := new(TCPAcceptor)
	t.SetType(constants.TCPComponent)
	return t
}

// Init inits acceptor
func (a *TCPAcceptor) Init(conf *config.AiryConfig) {
	keyFile := ""
	certFile := ""
	if conf.TLS.Cert != "" && conf.TLS.Key != "" {
		certFile = conf.TLS.Cert
		keyFile = conf.TLS.Key
	}
	port := strconv.Itoa(mgr.GetPort(conf.Port.TCPPort))
	a.addr = ":" + port
	a.connChan = make(chan interfaces.PlayerConn)
	a.running = false
	a.certFile = certFile
	a.keyFile = keyFile

	if conf.Server.Metadata != nil {
		conf.Server.Metadata[constants.AcceptorPort] = port
	} else {
		conf.Server.Metadata = map[string]string{
			constants.AcceptorPort: port,
		}
	}
}

// AfterInit runs after initialization
func (a *TCPAcceptor) AfterInit() {
	go a.ListenAndServe()
	logger.Infof("tcp server is running on :%s", a.addr)
}

// Shutdown stops acceptor
func (a *TCPAcceptor) Shutdown() error {
	a.Stop()
	return nil
}

// ReadFrom return the next packet sent by the client
func (t *tcpPlayerConn) ReadFrom(codec interfaces.PacketCodec) (b []byte, err error) {
	return codec.Decode(t.Conn)
}

// Write send packet to client
func (t *tcpPlayerConn) WriteTo(b []byte) error {
	if _, err := t.Conn.Write(b); err != nil {
		logger.Errorf("Failed to write in conn: %s", err.Error())
		return constants.ErrWriteConn
	}
	return nil
}

// GetAddr returns the addr the acceptor will listen on
func (a *TCPAcceptor) GetAddr() string {
	if a.listener != nil {
		return a.listener.Addr().String()
	}
	return ""
}

// GetConnChan gets a connection channel
func (a *TCPAcceptor) GetConnChan() chan interfaces.PlayerConn {
	return a.connChan
}

// Stop stops the acceptor
func (a *TCPAcceptor) Stop() {
	a.running = false
	_ = a.listener.Close()
}

func (a *TCPAcceptor) hasTLSCertificates() bool {
	return a.certFile != "" && a.keyFile != ""
}

// ListenAndServe using tcp acceptor
func (a *TCPAcceptor) ListenAndServe() {
	if a.hasTLSCertificates() {
		a.ListenAndServeTLS(a.certFile, a.keyFile)
		return
	}

	listener, err := net.Listen("tcp", a.addr)
	if err != nil {
		logger.Fatalf("Failed to listen: %s", err.Error())
		panic(fmt.Sprintf("Failed to listen: %s", err.Error()))
	}
	a.listener = listener
	a.running = true
	a.serve()
}

// ListenAndServeTLS listens using tls
func (a *TCPAcceptor) ListenAndServeTLS(cert, key string) {
	crt, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		logger.Fatalf("Failed to listen: %s", err.Error())
	}

	tlsCfg := &tls.Config{Certificates: []tls.Certificate{crt}}

	listener, err := tls.Listen("tcp", a.addr, tlsCfg)
	if err != nil {
		logger.Fatalf("Failed to listen: %s", err.Error())
	}
	a.listener = listener
	a.running = true
	a.serve()
}

func (a *TCPAcceptor) serve() {
	defer a.Stop()
	for a.running {
		conn, err := a.listener.Accept()
		if err != nil {
			logger.Errorf("tcp listener acceptor err : %s", err.Error())
			continue
		}
		a.connChan <- &tcpPlayerConn{
			Conn: conn,
		}
	}
}
