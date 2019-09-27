package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Luncert/RedisExample/redishadow/util"
)

const (
	defaultLogFileNamePrefix = "log-"
	defaultMaxSingleFileSize = 1024 * 1024 // 1MB
	metadataFileName         = ".log-metadata"
)

type fileAppender struct {
	logPath           string
	logFileNamePrefix string
	maxSingleFileSize int64
	lastLogFileTag    string
	logFileSequence   int
	current           *os.File
	currentFileSize   int64
}

func newFileAppender(logPath, logFileNamePrefix, maxSingleFileSize string) *fileAppender {
	// check whether logPath refers to a directory
	if fileInfo, err := os.Stat(logPath); err != nil {
		if err = os.MkdirAll(logPath, os.ModePerm); err != nil {
			fatalF("Failed to initial log path: %s", logPath)
		}
	} else if !fileInfo.IsDir() {
		fatalF("Target log path is not a directory: %s", logPath)
	}

	f := &fileAppender{
		logPath:           logPath,
		logFileNamePrefix: defaultLogFileNamePrefix,
		maxSingleFileSize: defaultMaxSingleFileSize,
		lastLogFileTag:    "",
		logFileSequence:   0,
		current:           nil,
		currentFileSize:   0,
	}

	if len(logFileNamePrefix) != 0 {
		f.logFileNamePrefix = logFileNamePrefix
	}

	// read config: maxSingleFileSize
	if len(maxSingleFileSize) != 0 {
		readingNumber := true
		numPart := strings.Builder{}
		unitPart := strings.Builder{}
		for _, i := range maxSingleFileSize {
			if readingNumber {
				if util.IsDigit(i) || i == '.' {
					numPart.WriteByte(byte(i))
				} else if util.IsAlpha(i) {
					readingNumber = false
					unitPart.WriteByte(byte(i))
				} else {
					fatalF("Failed to parse config `maxSingleFileSize` = %s", maxSingleFileSize)
				}
			} else if util.IsAlpha(i) {
				unitPart.WriteByte(byte(i))
			} else {
				fatalF("Failed to parse config `maxSingleFileSize` = %s", maxSingleFileSize)
			}
		}
		tmp, err := strconv.ParseFloat(numPart.String(), 64)
		if err != nil {
			fatalF("Failed to parse config `maxSingleFileSize` = %s: %v", maxSingleFileSize, err)
		}
		switch strings.ToUpper(unitPart.String()) {
		case "MB":
			tmp *= 1024 * 1024
		case "KB":
			tmp *= 1024
		case "":
			fallthrough
		case "B":
			// no action
		default:
			fatalF("Not supported unit in `maxSingleFileSize` = %s", maxSingleFileSize)
		}
		f.maxSingleFileSize = int64(tmp)
	}

	// read metadata and open log file
	if data, err := ioutil.ReadFile(f.getMetadataFilePath()); err == nil {
		f.unmarshalLogFileFullTag(string(data))
		logFilePath := f.getLogFilePath()
		f.current, err = os.OpenFile(logFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			fatalF("Failed to open log file `%s`: %v", logFilePath, err)
		} else {
			if fileInfo, err := f.current.Stat(); err != nil {
				fatalF("Failed to fetch file stats for `%s`", logFilePath)
			} else {
				f.currentFileSize = fileInfo.Size()
			}
		}
	} else {
		f.openNewLogFile()
	}

	return f
}

func (f *fileAppender) Write(data []byte) (err error) {
	// open new log file if current log file is too large
	if int64(len(data))+f.currentFileSize > f.maxSingleFileSize {
		if err = f.current.Close(); err != nil {
			return
		} else {
			f.openNewLogFile()
		}
	}
	_, err = f.current.Write(data)
	return
}

func (f *fileAppender) Close() (err error) {
	err = f.current.Close()
	f.current = nil
	return
}

func (f *fileAppender) getLogFilePath() string {
	return filepath.Join(f.logPath,
		fmt.Sprint(f.logFileNamePrefix, f.lastLogFileTag, "#", f.logFileSequence, ".log"))
}

func (f *fileAppender) unmarshalLogFileFullTag(logFileFullTag string) {
	var err error
	i := strings.LastIndex(logFileFullTag, "#")
	f.lastLogFileTag = logFileFullTag[:i]
	f.logFileSequence, err = strconv.Atoi(logFileFullTag[i+1:])
	if err != nil {
		fatalF("Failed to read metadata, could not unmarshal `lastLogFileFullTag`")
	}
}

func (f *fileAppender) getLogFileFullTag() string {
	return fmt.Sprintf("%s#%d", f.lastLogFileTag, f.logFileSequence)
}

func (f *fileAppender) getMetadataFilePath() string {
	return filepath.Join(f.logPath, metadataFileName)
}

func (f *fileAppender) openNewLogFile() {
	var err error

	now := time.Now()
	logFileTag := fmt.Sprintf("%d-%d-%dT%d-%d-%d", now.Year(),
		now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	if f.lastLogFileTag == logFileTag {
		f.logFileSequence++
	} else {
		f.logFileSequence = 0
	}
	f.lastLogFileTag = logFileTag

	// persist logFileFullTag to metadata file
	logFileFullTag := []byte(f.getLogFileFullTag())
	if err = ioutil.WriteFile(f.getMetadataFilePath(), logFileFullTag, 0666); err != nil {
		fatalF("Failed to persist metadata: %v", err)
	}

	// open new log file
	newFilePath := f.getLogFilePath()
	f.current, err = os.OpenFile(newFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fatalF("Failed to open new log file: %s", newFilePath)
	}
	f.currentFileSize = 0
}
