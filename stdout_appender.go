package log

import "os"

type stdoutAppender struct {
}

func newStdoutAppender() *stdoutAppender {
	return &stdoutAppender{}
}

func (s *stdoutAppender) Write(data []byte) (err error) {
	_, err = os.Stdout.Write(data)
	return
}

func (s *stdoutAppender) Close() (err error) {
	return
}
