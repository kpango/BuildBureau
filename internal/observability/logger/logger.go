// Copyright (C) 2024 BuildBureau team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package logger provides a structured logging abstraction
// inspired by vdaas/vald logging patterns
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// Level represents the log level
type Level int

const (
	// LevelDebug is for debug messages
	LevelDebug Level = iota
	// LevelInfo is for informational messages
	LevelInfo
	// LevelWarn is for warning messages
	LevelWarn
	// LevelError is for error messages
	LevelError
	// LevelFatal is for fatal messages
	LevelFatal
)

// String returns the string representation of the level
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Format represents the log format
type Format int

const (
	// FormatText is plain text format
	FormatText Format = iota
	// FormatJSON is JSON format
	FormatJSON
)

// Logger is the interface for logging
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	WithFields(fields ...Field) Logger
	Close() error
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// String creates a string field
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// Int creates an int field
func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// Error creates an error field
func Error(err error) Field {
	return Field{Key: "error", Value: err}
}

// Any creates a field with any value
func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// standardLogger implements Logger using the standard library
type standardLogger struct {
	mu     sync.Mutex
	level  Level
	format Format
	out    io.Writer
	logger *log.Logger
	fields []Field
}

// NewStandardLogger creates a new standard logger
func NewStandardLogger(opts ...Option) Logger {
	l := &standardLogger{
		level:  LevelInfo,
		format: FormatText,
		out:    os.Stdout,
		fields: make([]Field, 0),
	}

	for _, opt := range opts {
		opt(l)
	}

	l.logger = log.New(l.out, "", 0)
	return l
}

// Option is a function that configures a logger
type Option func(*standardLogger)

// WithLevel sets the log level
func WithLevel(level Level) Option {
	return func(l *standardLogger) {
		l.level = level
	}
}

// WithFormat sets the log format
func WithFormat(format Format) Option {
	return func(l *standardLogger) {
		l.format = format
	}
}

// WithOutput sets the output writer
func WithOutput(out io.Writer) Option {
	return func(l *standardLogger) {
		l.out = out
	}
}

func (l *standardLogger) log(level Level, msg string, fields ...Field) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	allFields := append(l.fields, fields...)

	timestamp := time.Now().Format(time.RFC3339)
	levelStr := level.String()

	switch l.format {
	case FormatJSON:
		l.logJSON(timestamp, levelStr, msg, allFields)
	default:
		l.logText(timestamp, levelStr, msg, allFields)
	}
}

func (l *standardLogger) logText(timestamp, level, msg string, fields []Field) {
	output := fmt.Sprintf("[%s] %s: %s", timestamp, level, msg)

	if len(fields) > 0 {
		output += " |"
		for _, f := range fields {
			output += fmt.Sprintf(" %s=%v", f.Key, f.Value)
		}
	}

	l.logger.Println(output)
}

func (l *standardLogger) logJSON(timestamp, level, msg string, fields []Field) {
	output := fmt.Sprintf(`{"timestamp":"%s","level":"%s","message":"%s"`, timestamp, level, msg)

	if len(fields) > 0 {
		for _, f := range fields {
			output += fmt.Sprintf(`,"%s":"%v"`, f.Key, f.Value)
		}
	}

	output += "}"
	l.logger.Println(output)
}

// Debug logs a debug message
func (l *standardLogger) Debug(msg string, fields ...Field) {
	l.log(LevelDebug, msg, fields...)
}

// Info logs an info message
func (l *standardLogger) Info(msg string, fields ...Field) {
	l.log(LevelInfo, msg, fields...)
}

// Warn logs a warning message
func (l *standardLogger) Warn(msg string, fields ...Field) {
	l.log(LevelWarn, msg, fields...)
}

// Error logs an error message
func (l *standardLogger) Error(msg string, fields ...Field) {
	l.log(LevelError, msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *standardLogger) Fatal(msg string, fields ...Field) {
	l.log(LevelFatal, msg, fields...)
	os.Exit(1)
}

// WithFields returns a new logger with the given fields
func (l *standardLogger) WithFields(fields ...Field) Logger {
	newLogger := &standardLogger{
		level:  l.level,
		format: l.format,
		out:    l.out,
		logger: l.logger,
		fields: append(l.fields, fields...),
	}
	return newLogger
}

// Close closes the logger
func (l *standardLogger) Close() error {
	return nil
}

// Global logger instance
var (
	globalLogger Logger
	once         sync.Once
)

// Init initializes the global logger
func Init(opts ...Option) {
	once.Do(func() {
		globalLogger = NewStandardLogger(opts...)
	})
}

// GetLogger returns the global logger
func GetLogger() Logger {
	if globalLogger == nil {
		Init()
	}
	return globalLogger
}

// Debug logs a debug message using the global logger
func Debug(msg string, fields ...Field) {
	GetLogger().Debug(msg, fields...)
}

// Info logs an info message using the global logger
func Info(msg string, fields ...Field) {
	GetLogger().Info(msg, fields...)
}

// Warn logs a warning message using the global logger
func Warn(msg string, fields ...Field) {
	GetLogger().Warn(msg, fields...)
}

// ErrorLog logs an error message using the global logger
func ErrorLog(msg string, fields ...Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal logs a fatal message using the global logger and exits
func Fatal(msg string, fields ...Field) {
	GetLogger().Fatal(msg, fields...)
}

// WithFields returns a logger with the given fields
func WithFields(fields ...Field) Logger {
	return GetLogger().WithFields(fields...)
}
