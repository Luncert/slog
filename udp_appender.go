package log

import "net"

type udpAppender struct {
	serverAddr string
	conn       net.Conn
}

func newUdpAppender(serverAddr string) *udpAppender {
	u := &udpAppender{
		serverAddr: serverAddr,
	}
	u.connect()
	return u
}

func (u *udpAppender) Write(data []byte) (err error) {
	_, err = u.conn.Write(data)
	return
}

func (u *udpAppender) Close() (err error) {
	u.disconnect()
	return
}

func (u *udpAppender) connect() {
	if u.conn != nil {
		fatalF("The last udp connection hasn't been closed")
	}
	var err error
	if u.conn, err = net.Dial("udp", u.serverAddr); err != nil {
		fatalF("Failed to connect to udp server: %v", err)
	}
}

func (u *udpAppender) disconnect() {
	if u.conn == nil {
		fatalF("There is no udp connection could be closed.")
	}
	if err := u.conn.Close(); err != nil {
		fatalF("Failed to close udp connection: %v", err)
	} else {
		u.conn = nil
	}
}
