package tm_http_redirect

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

const DefaultInfoPrefix string = "INFO"
const DefaultWarnPrefix string = "WARN"
const DefaultErrorPrefix string = "ERROR"
const DefaultFatalPrefix string = "FATAL"
const DefaultPanicPrefix string = "PANIC"

const DefaultNamePrefix string = ""

const DefaultCollecting bool = true
const DefaultCollectionLength uint16 = 25

type FormattingLogMethod func(string, ...any)

var DefaultLogMethod FormattingLogMethod = log.Printf
var DefaultFatalLogMethod FormattingLogMethod = log.Fatalf
var DefaultPanicLogMethod FormattingLogMethod = log.Panicf

type Logger struct {
	NamePrefix       string
	InfoPrefix       string
	WarnPrefix       string
	ErrorPrefix      string
	FatalPrefix      string
	PanicPrefix      string
	LogMethod        FormattingLogMethod
	FatalLogMethod   FormattingLogMethod
	PanicLogMethod   FormattingLogMethod
	Collecting       bool
	CollectedLogs    []string
	CollectionLength uint16
	collectionMutex  sync.Mutex
}

func DefaultLogger() *Logger {
	return &Logger{
		NamePrefix:       DefaultNamePrefix,
		InfoPrefix:       DefaultInfoPrefix,
		WarnPrefix:       DefaultWarnPrefix,
		ErrorPrefix:      DefaultErrorPrefix,
		FatalPrefix:      DefaultFatalPrefix,
		PanicPrefix:      DefaultPanicPrefix,
		LogMethod:        DefaultLogMethod,
		FatalLogMethod:   DefaultFatalLogMethod,
		PanicLogMethod:   DefaultPanicLogMethod,
		Collecting:       true,
		CollectedLogs:    make([]string, 0),
		CollectionLength: DefaultCollectionLength,
		collectionMutex:  sync.Mutex{},
	}
}

func NamedLogger(name string) *Logger {
	logger := DefaultLogger()
	logger.NamePrefix = name
	return logger
}

func (l *Logger) collect(format string, args ...any) {
	if !l.Collecting {
		return
	}
	l.collectionMutex.Lock()
	defer l.collectionMutex.Unlock()
	for len(l.CollectedLogs) >= int(l.CollectionLength) {
		l.CollectedLogs = l.CollectedLogs[1:]
	}
	l.CollectedLogs = append(l.CollectedLogs, fmt.Sprintf(format, args...))
}

func (l *Logger) CollectedError(baseError error, flush bool) error {
	if !l.Collecting {
		return errors.Join(ErrLogCollectionInactive, baseError)
	}
	l.collectionMutex.Lock()
	defer l.collectionMutex.Unlock()
	errorObj := []error{}
	for _, m := range l.CollectedLogs {
		errorObj = append(errorObj, errors.New(m))
	}
	if flush {
		l.CollectedLogs = []string{}
	}
	errorObj = append(errorObj, baseError)
	return errors.Join(errorObj...)
}

func (l *Logger) log(method FormattingLogMethod, prefix string, format string, args ...any) {
	logMethod := method
	if l.Collecting {
		logMethod = l.collect
	}
	if l.NamePrefix != "" {
		logMethod("%s [%s] %s", l.NamePrefix, prefix, fmt.Sprintf(format, args...))
	} else {
		logMethod("[%s] %s", prefix, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Info(format string, args ...any) {
	l.log(l.LogMethod, l.InfoPrefix, format, args...)
}

func (l *Logger) Warn(format string, args ...any) {
	l.log(l.LogMethod, l.WarnPrefix, format, args...)
}

func (l *Logger) Error(format string, args ...any) {
	l.log(l.LogMethod, l.ErrorPrefix, format, args...)
}

func (l *Logger) Fatal(format string, args ...any) {
	l.log(l.FatalLogMethod, l.FatalPrefix, format, args...)
}

func (l *Logger) Panic(format string, args ...any) {
	l.log(l.PanicLogMethod, l.PanicPrefix, format, args...)
}
