package log

import "net"

type tcpAppender struct {
	serverAddr string
	conn       net.Conn
}

func newTcpAppender(serverAddr string) *tcpAppender {
	t := &tcpAppender{
		serverAddr: serverAddr,
	}
	t.connect()
	return t
}

func (t *tcpAppender) Write(data []byte) (err error) {
	_, err = t.conn.Write(data)
	return
}

func (t *tcpAppender) Close() (err error) {
	t.disconnect()
	return
}

func (t *tcpAppender) connect() {
	if t.conn != nil {
		fatalF("The last tcp connection hasn't been closed")
	}
	var err error
	if t.conn, err = net.Dial("tcp", t.serverAddr); err != nil {
		fatalF("Failed to connect to tcp server: %v", err)
	}
}

func (t *tcpAppender) disconnect() {
	if t.conn == nil {
		fatalF("There is no tcp connection could be closed.")
	}
	if err := t.conn.Close(); err != nil {
		fatalF("Failed to close udp connection: %v", err)
	} else {
		t.conn = nil
	}
}
