package util

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	loggerInstance *zap.Logger
	loggerOnce     sync.Once
)

func GetLogger() *zap.Logger {
	loggerOnce.Do(func() {
		core := zapcore.NewCore(newZapLogger(), zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
		loggerInstance = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	})
	return loggerInstance
}

func Infof(format string, args ...interface{}) {
	GetLogger().Info(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...interface{}) {
	GetLogger().Error(fmt.Sprintf(format, args...))
}

func Debugf(format string, args ...interface{}) {
	GetLogger().Debug(fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...interface{}) {
	GetLogger().Warn(fmt.Sprintf(format, args...))
}

func newZapLogger() zapcore.Encoder {
	return &CustomZapLogger{}
}

type CustomZapLogger struct{}

func (e *CustomZapLogger) Clone() zapcore.Encoder { return &CustomZapLogger{} }

func (e *CustomZapLogger) EncodeEntry(entry zapcore.Entry, _ []zapcore.Field) (*buffer.Buffer, error) {
	buf := buffer.NewPool().Get()
	timestamp := entry.Time.Format("2006-01-02 15:04:05.000000")
	level := fmt.Sprintf("%-5s", entry.Level.CapitalString())

	caller := ""
	if entry.Caller.Defined {
		_, file := filepath.Split(entry.Caller.File)
		caller = fmt.Sprintf("[%s:%d]", file, entry.Caller.Line)
	}

	_, err := fmt.Fprintf(buf, "%s %s %s %s\n", timestamp, level, caller, entry.Message)

	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (e *CustomZapLogger) AddArray(string, zapcore.ArrayMarshaler) error   { return nil }
func (e *CustomZapLogger) AddObject(string, zapcore.ObjectMarshaler) error { return nil }
func (e *CustomZapLogger) AddBinary(string, []byte)                        {}
func (e *CustomZapLogger) AddByteString(string, []byte)                    {}
func (e *CustomZapLogger) AddBool(string, bool)                            {}
func (e *CustomZapLogger) AddComplex128(string, complex128)                {}
func (e *CustomZapLogger) AddComplex64(string, complex64)                  {}
func (e *CustomZapLogger) AddDuration(string, time.Duration)               {}
func (e *CustomZapLogger) AddFloat64(string, float64)                      {}
func (e *CustomZapLogger) AddFloat32(string, float32)                      {}
func (e *CustomZapLogger) AddInt(string, int)                              {}
func (e *CustomZapLogger) AddInt64(string, int64)                          {}
func (e *CustomZapLogger) AddInt32(string, int32)                          {}
func (e *CustomZapLogger) AddInt16(string, int16)                          {}
func (e *CustomZapLogger) AddInt8(string, int8)                            {}
func (e *CustomZapLogger) AddString(string, string)                        {}
func (e *CustomZapLogger) AddTime(string, time.Time)                       {}
func (e *CustomZapLogger) AddUint(string, uint)                            {}
func (e *CustomZapLogger) AddUint64(string, uint64)                        {}
func (e *CustomZapLogger) AddUint32(string, uint32)                        {}
func (e *CustomZapLogger) AddUint16(string, uint16)                        {}
func (e *CustomZapLogger) AddUint8(string, uint8)                          {}
func (e *CustomZapLogger) AddUintptr(string, uintptr)                      {}
func (e *CustomZapLogger) AddReflected(string, interface{}) error          { return nil }
func (e *CustomZapLogger) OpenNamespace(string)                            {}
