package log

import "testing"

func TestStdoutLogger(t *testing.T) {
	InitLogger("logger-stdout.yml")
	Info("This", "is", "test")
	Debug("This should not be logged")
	DestroyLogger()
}

func TestTcpLogger(t *testing.T) {
	InitLogger("logger-tcp.yml")
	Info("This", "is", "test")
	DestroyLogger()
}

func TestUdpLogger(t *testing.T) {
	InitLogger("logger-udp.yml")
	Info("This", "is", "test")
	DestroyLogger()
}

func TestFileLogger(t *testing.T) {
	InitLogger("logger-file.yml")
	Info("This", "is", "test")
	DestroyLogger()
}
