package log

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	logPath           = "./test-output"
	logFileNamePrefix = "test-log-"
	maxSingleFileSize = "0.5kb"
	testLog           = "test log"
)

var expectLogContent string

func TestBase(t *testing.T) {
	testPart(t)

	testPart(t)

	// clean
	if data, err := ioutil.ReadFile(filepath.Join(logPath, metadataFileName)); err != nil {
		t.Error(err)
	} else {
		filePath := filepath.Join(logPath, logFileNamePrefix+string(data)+".log")
		if err := os.Remove(filePath); err != nil {
			t.Error(err)
		}
		if err := os.Remove(filepath.Join(logPath, metadataFileName)); err != nil {
			t.Error(err)
		}
		if err := os.Remove(logPath); err != nil {
			t.Error(err)
		}
	}

}

func testPart(t *testing.T) {
	f := newFileAppender(logPath, logFileNamePrefix, maxSingleFileSize)

	if f.logPath != logPath {
		t.Errorf("logPath should be `%s`, got `%s`", logPath, f.logPath)
	}
	if f.logFileNamePrefix != logFileNamePrefix {
		t.Errorf("logFileNamePrefix should be `%s`, got `%s`", logPath, f.logPath)
	}
	if f.maxSingleFileSize != 0.5*1024 {
		t.Errorf("maxSingleFileSize should be 2*1024*1024, got %d", f.maxSingleFileSize)
	}
	if f.current == nil {
		t.Error("failed to open log file")
	}

	if err := f.Write([]byte(testLog)); err != nil {
		t.Error(err)
	}

	if err := f.Close(); err != nil {
		t.Error(err)
	}
	if f.current != nil {
		t.Error("log file should be closed")
	}

	expectLogContent += testLog
	if data, err := ioutil.ReadFile(f.getLogFilePath()); err != nil {
		t.Error(err)
	} else if string(data) != expectLogContent {
		t.Errorf("log output should be `%s`, got `%s`", expectLogContent, string(data))
	}

	// check if file_appender creates the log path correctly
	if _, err := os.Stat(logPath); err != nil {
		t.Error(err)
	}
}
