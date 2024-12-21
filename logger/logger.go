package logger

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Transport interface {
	Log(level string, message string) error
}

type Logger struct {
	level       LogLevel
	serviceName string
	environment string
	transports  []Transport
}

func NewLogger(serviceName string, env string) *Logger {
	opts := TransportOption{
		BatchSize: 20,
		Interval:  1 * time.Second,
	}
	transport := NewTransport(opts)
	return &Logger{
		level:       INFO,
		serviceName: serviceName,
		environment: env,
		transports:  []Transport{transport},
	}
}

func (l *Logger) getInternalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "Unknown"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return "Unknown"
}

func (l *Logger) formatMessage(level LogLevel, message string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := l.logLevelToString(level)
	// TODO implement env
	connection := ""
	internalIP := l.getInternalIP()
	return fmt.Sprintf("%s | %s | %s - %s %s %s", timestamp, message, l.serviceName, levelStr, connection, internalIP)
}

func (l *Logger) logLevelToString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func (l *Logger) log(level LogLevel, message string, serviceLevel string) {
	if level < l.level {
		return
	}

	formatMessage := l.formatMessage(level, message)
	log.Println(formatMessage)

	for _, transport := range l.transports {
		transport.Log(serviceLevel, formatMessage)
	}

	if level == FATAL {
		os.Exit(1)
	}
}

func (l *Logger) DEBUG(message string) {
	l.log(DEBUG, message, "")
}

func (l *Logger) Info(message string) {
	l.log(INFO, message, "notify")
}

func (l *Logger) Warn(message string) {
	l.log(WARN, message, "notify")
}

func (l *Logger) Error(message string) {
	l.log(ERROR, message, "error")
}

func (l *Logger) Fatal(message string) {
	l.log(FATAL, message, "error")
}
