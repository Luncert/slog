package log

import (
	"fmt"
	"testing"
)

func TestLogFormatter(t *testing.T) {
	formatter := newLogFormatter("timestamp = %T, custom timestamp = %y/%M/%d-%h:%m:%s" +
		" log level = %L, placeholder = %S %S %S")
	result := formatter.format(debugLevel, "this", "is", "test")
	fmt.Println(string(result))
}
