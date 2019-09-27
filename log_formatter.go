package log

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// log formatter
type logFormatter struct {
	providers []logPartProvider
}

type logPartProvider func(logLevel logLevel, v interface{}) (string, bool)

func newLogFormatter(format string) *logFormatter {
	l := &logFormatter{
		providers: make([]logPartProvider, 0),
	}
	var preRune rune
	builder := strings.Builder{}
	for i, r := range format {
		if preRune == '%' {
			switch r {
			case 'T':
				l.providers = append(l.providers, timestampProvider)
			case 'y':
				l.providers = append(l.providers, timeYearProvider)
			case 'M':
				l.providers = append(l.providers, timeMonthProvider)
			case 'd':
				l.providers = append(l.providers, timeDayProvider)
			case 'h':
				l.providers = append(l.providers, timeHourProvider)
			case 'm':
				l.providers = append(l.providers, timeMinuteProvider)
			case 's':
				l.providers = append(l.providers, timeSecondProvider)
			case 'L':
				l.providers = append(l.providers, logLevelProvider)
			case 'S':
				l.providers = append(l.providers, placeholderProvider)
			default:
				fatalF("Invalid control character at pos %d of `%s`", i, format)
			}
		} else if r == '%' {
			tmp := &plainTextProvider{content: builder.String()}
			l.providers = append(l.providers, tmp.consume)
			builder.Reset()
		} else {
			builder.WriteRune(r)
		}
		preRune = r
	}
	return l
}

/*
ctrl literal:
%T timestamp
%y years
%M months
%d days
%h hours
%m minutes
%s seconds
%L log level
%S placeholder
*/
func (l *logFormatter) format(logLevel logLevel, v ...interface{}) []byte {
	builder := strings.Builder{}
	i := 0
	tmp := v[i]
	for _, provider := range l.providers {
		part, ok := provider(logLevel, tmp)
		builder.WriteString(part)
		if ok {
			i++
			if i < len(v) {
				tmp = v[i]
			} else {
				tmp = nil
			}
		}
	}
	builder.WriteRune('\n')
	return []byte(builder.String())
}

func timestampProvider(_ logLevel, _ interface{}) (string, bool) {
	n := time.Now()
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", n.Year(), n.Month(), n.Day(),
		n.Hour(), n.Minute(), n.Second()), false
}

func timeYearProvider(_ logLevel, _ interface{}) (string, bool) {
	n := time.Now()
	return strconv.FormatInt(int64(n.Year()), 10), false
}

func timeMonthProvider(_ logLevel, _ interface{}) (string, bool) {
	n := time.Now()
	return strconv.FormatInt(int64(n.Month()), 10), false
}

func timeDayProvider(_ logLevel, _ interface{}) (string, bool) {
	n := time.Now()
	return strconv.FormatInt(int64(n.Day()), 10), false
}

func timeHourProvider(_ logLevel, _ interface{}) (string, bool) {
	n := time.Now()
	return strconv.FormatInt(int64(n.Hour()), 10), false
}

func timeMinuteProvider(_ logLevel, _ interface{}) (string, bool) {
	n := time.Now()
	return strconv.FormatInt(int64(n.Minute()), 10), false
}

func timeSecondProvider(_ logLevel, _ interface{}) (string, bool) {
	n := time.Now()
	return strconv.FormatInt(int64(n.Second()), 10), false
}

func logLevelProvider(logLevel logLevel, _ interface{}) (string, bool) {
	return logLevel.string(), false
}

func placeholderProvider(_ logLevel, v interface{}) (string, bool) {
	if v != nil {
		return fmt.Sprintf("%v", v), true
	} else {
		return "", false
	}
}

type plainTextProvider struct {
	content string
}

func (p *plainTextProvider) consume(logLevel logLevel, v interface{}) (string, bool) {
	return p.content, false
}
