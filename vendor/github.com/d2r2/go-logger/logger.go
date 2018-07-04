package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/d2r2/go-shell"
)

type LogLevel int

const (
	FatalLevel LogLevel = iota
	PanicLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

func (v LogLevel) String() string {
	switch v {
	case FatalLevel:
		return "Fatal"
	case PanicLevel:
		return "Panic"
	case ErrorLevel:
		return "Error"
	case WarnLevel:
		return "Warning"
	case InfoLevel:
		return "Information"
	case DebugLevel:
		return "Debug"
	default:
		return "<undefined>"
	}
}

func (v *LogLevel) LongStr() string {
	return v.String()
}

func (v LogLevel) ShortStr() string {
	switch v {
	case FatalLevel:
		return "Fatal"
	case PanicLevel:
		return "Panic"
	case ErrorLevel:
		return "Error"
	case WarnLevel:
		return "Warn"
	case InfoLevel:
		return "Info"
	case DebugLevel:
		return "Debug"
	default:
		return "undef"
	}
}

type LevelFormat int

const (
	LevelShort LevelFormat = iota
	LevelLong
)

const (
	ShortLevelLen = 5
	LongLevelLen  = 11
)

type Logger struct {
	sync.RWMutex
	log                *log.Logger
	packages           []*Package
	packagePrintLength int
	levelFormat        LevelFormat
	logFile            *File
	rotateMaxSize      int64
	rotateMaxCount     int
	appName            string
	enableSyslog       bool
}

func NewLogger() *Logger {
	log := log.New(os.Stdout, "", 0)
	l := &Logger{
		log:                log,
		levelFormat:        LevelShort,
		packagePrintLength: 8,
		rotateMaxSize:      1024 * 1024 * 512,
		rotateMaxCount:     3,
	}
	return l
}

func (v *Logger) Close() error {
	v.Lock()
	defer v.Unlock()

	for _, pack := range v.packages {
		pack.Close()
	}
	v.packages = nil

	if v.logFile != nil {
		v.logFile.Close()
	}
	return nil
}

func (v *Logger) SetRotateParams(rotateMaxSize int64, rotateMaxCount int) {
	v.Lock()
	defer v.Unlock()
	v.rotateMaxSize = rotateMaxSize
	v.rotateMaxCount = rotateMaxCount
}

func (v *Logger) GetRotateMaxSize() int64 {
	v.Lock()
	defer v.Unlock()
	return v.rotateMaxSize
}

func (v *Logger) GetRotateMaxCount() int {
	v.Lock()
	defer v.Unlock()
	return v.rotateMaxCount
}

func (v *Logger) SetLevelFormat(levelFormat LevelFormat) {
	v.Lock()
	defer v.Unlock()
	v.levelFormat = levelFormat
}

func (v *Logger) GetLevelFormat() LevelFormat {
	v.RLock()
	defer v.RUnlock()
	return v.levelFormat
}

func (v *Logger) SetApplicationName(appName string) {
	v.Lock()
	defer v.Unlock()
	v.appName = appName
}

func (v *Logger) GetApplicationName() string {
	v.RLock()
	defer v.RUnlock()
	return v.appName
}

func (v *Logger) EnableSyslog(enable bool) {
	v.Lock()
	defer v.Unlock()
	v.enableSyslog = enable
}

func (v *Logger) GetSyslogEnabled() bool {
	v.RLock()
	defer v.RUnlock()
	return v.enableSyslog
}

func (v *Logger) SetPackagePrintLength(packagePrintLength int) {
	v.Lock()
	defer v.Unlock()
	v.packagePrintLength = packagePrintLength
}

func (v *Logger) GetPackagePrintLength() int {
	v.RLock()
	defer v.RUnlock()
	return v.packagePrintLength
}

func (v *Logger) SetLogFileName(logFilePath string) error {
	if path.Ext(logFilePath) == "" {
		logFilePath += ".log"
	}
	fp, err := filepath.Abs(logFilePath)
	if err != nil {
		return err
	}
	v.Lock()
	defer v.Unlock()
	lf := &File{Path: fp}
	v.logFile = lf
	return nil
}

func (v *Logger) GetLogFileInfo() *File {
	v.RLock()
	defer v.RUnlock()
	return v.logFile
}

func (v *Logger) NewPackageLogger(packageName string, level LogLevel) *Package {
	v.Lock()
	defer v.Unlock()
	p := &Package{parent: v, packageName: packageName, level: level}
	v.packages = append(v.packages, p)
	return p
}

func (v *Logger) ChangePackageLogLevel(packageName string, level LogLevel) error {
	var p *Package
	for _, item := range v.packages {
		if item.packageName == packageName {
			p = item
			break
		}
	}
	if p != nil {
		p.SetLogLevel(level)
	} else {
		err := fmt.Errorf("Package log %q is not found", packageName)
		return err
	}
	return nil
}

var (
	globalLock sync.RWMutex
	lgr        *Logger
)

func SetLevelFormat(levelFormat LevelFormat) {
	globalLock.RLock()
	defer globalLock.RUnlock()
	lgr.SetLevelFormat(levelFormat)
}

func SetPackagePrintLength(packagePrintLength int) {
	globalLock.RLock()
	defer globalLock.RUnlock()
	lgr.SetPackagePrintLength(packagePrintLength)
}

func SetRotateParams(rotateMaxSize int64, rotateMaxCount int) {
	globalLock.RLock()
	defer globalLock.RUnlock()
	lgr.SetRotateParams(rotateMaxSize, rotateMaxCount)
}

func NewPackageLogger(module string, level LogLevel) *Package {
	globalLock.RLock()
	defer globalLock.RUnlock()
	return lgr.NewPackageLogger(module, level)
}

func ChangePackageLogLevel(packageName string, level LogLevel) error {
	globalLock.RLock()
	defer globalLock.RUnlock()
	return lgr.ChangePackageLogLevel(packageName, level)
}

func SetLogFileName(logFilePath string) error {
	globalLock.RLock()
	defer globalLock.RUnlock()
	return lgr.SetLogFileName(logFilePath)
}

func SetApplicationName(appName string) {
	globalLock.RLock()
	defer globalLock.RUnlock()
	lgr.SetApplicationName(appName)
}

func EnableSyslog(enable bool) {
	globalLock.RLock()
	defer globalLock.RUnlock()
	lgr.EnableSyslog(enable)
}

func FinalizeLogger() error {
	var err error
	if lgr != nil {
		err = lgr.Close()
	}
	globalLock.Lock()
	defer globalLock.Unlock()
	lgr = nil
	return err
}

func init() {
	lgr = NewLogger()
	ctx, cancel := context.WithCancel(context.Background())
	shell.CloseContextOnKillSignal(cancel)

	go func(logger *Logger) {
		<-ctx.Done()
		lg := logger.NewPackageLogger("logger", InfoLevel)
		lg.Info("Finalizing logger, due to termination pending request")
		logger.Close()
		globalLock.Lock()
		defer globalLock.Unlock()
		lgr = nil
	}(lgr)
}
