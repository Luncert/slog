package log

/*
Author: Luncert
*/

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func infoF(format string, a ...interface{}) {
	fmt.Printf("[INFO] %s\n", fmt.Sprintf(format, a...))
}

func errorF(format string, a ...interface{}) {
	fmt.Printf("[ERROR] %s\n", fmt.Sprintf(format, a...))
}

func fatalF(format string, a ...interface{}) {
	fmt.Printf("[FATAL] %s\n", fmt.Sprintf(format, a...))
	os.Exit(1)
}

// log level
type logLevel int

const (
	debugLevel logLevel = iota
	infoLevel
	warnLevel
	errorLevel
	fatalLevel
)

var logLevelNameArray = []string{"DEBUG", "INFO",
	"WARN", "ERROR", "FATAL",
}

func (l logLevel) string() string {
	return logLevelNameArray[l]
}

type logAppender interface {
	Write(data []byte) error
	Close() error
}

type logger struct {
	level     logLevel
	formatter *logFormatter
	appender  logAppender
}

func (l *logger) log(level logLevel, v ...interface{}) {
	if level >= l.level {
		if err := l.appender.Write(l.formatter.format(level, v...)); err != nil {
			errorF("Failed to write log: %v", err)
		}
		if level == logLevel.debugLevel {
			os.Exit(1)
		}
	}
}

func whenNotInitialized(level logLevel, v ...interface{}) {
	fatalF("Please invoke InitLogger first")
}

var log *logger
var logFunc = whenNotInitialized

func InitLogger(configFile string) {
	// read config
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		fatalF("Failed to read %s: %v", configFile, err)
	}

	config := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fatalF("Failed ot unmarshal configuration: %v", err)
		return
	}

	// create logger
	log = &logger{}

	if level, ok := config["level"]; !ok || level == nil {
		fatalF("Config missing: level")
	} else {
		switch strings.ToLower(level.(string)) {
		case "debug":
			log.level = debugLevel
		case "info":
			log.level = infoLevel
		case "warn":
			log.level = warnLevel
		case "error":
			log.level = errorLevel
		case "fatal":
			log.level = fatalLevel
		default:
			fatalF("Unknown log level: %s", level)
		}
	}

	if format, ok := config["format"]; !ok || format == nil {
		fatalF("Config missing: format")
	} else {
		log.formatter = newLogFormatter(format.(string))
	}

	appenderType, ok := config["appender"]
	if !ok || appenderType == nil {
		infoF("No appender specified, using default one: stdout appender")
	}
	tmp, ok := config["appenderConfig"]
	var appenderConfig map[interface{}]interface{}
	if ok {
		appenderConfig = tmp.(map[interface{}]interface{})
	}
	switch strings.ToLower(appenderType.(string)) {
	case "tcp":
		if appenderConfig == nil {
			fatalF("File Appender config missing")
		}
		serverAddr := getConfig("serverAddr", appenderConfig).(string)
		log.appender = newTcpAppender(serverAddr)
	case "udp":
		if appenderConfig == nil {
			fatalF("File Appender config missing")
		}
		serverAddr := getConfig("serverAddr", appenderConfig).(string)
		log.appender = newUdpAppender(serverAddr)
	case "file":
		if appenderConfig == nil {
			fatalF("File Appender config missing")
		}
		logPath := getConfig("logPath", appenderConfig).(string)
		logFileNamePrefix := getOptionalConfig("logFileNamePrefix", appenderConfig, "").(string)
		maxSingleFileSize := getOptionalConfig("maxSingleFileSize", appenderConfig, "").(string)
		log.appender = newFileAppender(logPath, logFileNamePrefix, maxSingleFileSize)
	case "stdout":
		fallthrough
	default:
		log.appender = newStdoutAppender()
	}

	logFunc = log.log
}

func DestroyLogger() {
	if err := log.appender.Close(); err != nil {
		fatalF("failed to destroy logger: %v ", err)
	}
	log = nil
}

func getConfig(key string, config map[interface{}]interface{}) interface{} {
	value, ok := config[key]
	if !ok || value == nil {
		fatalF("Missing appender config `%s`", key)
	}
	return value
}

func getOptionalConfig(key string, config map[interface{}]interface{}, defaultValue interface{}) interface{} {
	value, ok := config[key]
	if !ok || value == nil {
		return defaultValue
	}
	return value
}

func Debug(v ...interface{}) {
	logFunc(debugLevel, v...)
}

func Info(v ...interface{}) {
	logFunc(infoLevel, v...)
}

func Warn(v ...interface{}) {
	logFunc(warnLevel, v...)
}

func Error(v ...interface{}) {
	logFunc(errorLevel, v...)
}

func Fatal(v ...interface{}) {
	logFunc(fatalLevel, v...)
}
