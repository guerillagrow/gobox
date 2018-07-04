package logger

import (
	"fmt"
	"log/syslog"
	"os"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

type Package struct {
	sync.RWMutex
	parent      *Logger
	packageName string
	level       LogLevel
	syslog      *syslog.Writer
}

func (v *Package) Close() error {
	v.Lock()
	defer v.Unlock()
	if v.syslog != nil {
		err := v.syslog.Close()
		v.syslog = nil
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Package) SetLogLevel(level LogLevel) {
	v.Lock()
	defer v.Unlock()
	v.level = level
}

func (v *Package) GetLogLevel() LogLevel {
	v.RLock()
	defer v.RUnlock()
	return v.level
}

func (v *Package) getSyslog(level LogLevel, levelFormat LevelFormat,
	appName string) (*syslog.Writer, error) {
	v.Lock()
	defer v.Unlock()
	if v.syslog == nil {
		tag := fmtStr(false, level, levelFormat, appName,
			v.packageName, -1, "", "%[2]s-%[3]s")
		sl, err := syslog.New(syslog.LOG_DEBUG, tag)
		if err != nil {
			err = spew.Errorf("Failed to connect to syslog: %v\n", err)
			return nil, err
		}
		v.syslog = sl
	}
	return v.syslog, nil
}

func (v *Package) writeToSyslog(level LogLevel,
	levelFormat LevelFormat, appName string, msg string) error {
	sl, err := v.getSyslog(level, levelFormat, appName)
	if err != nil {
		return err
	}
	switch level {
	case DebugLevel:
		return sl.Debug(msg)
	case InfoLevel:
		return sl.Info(msg)
	case WarnLevel:
		return sl.Warning(msg)
	case ErrorLevel:
		return sl.Err(msg)
	case PanicLevel:
		return sl.Crit(msg)
	case FatalLevel:
		return sl.Emerg(msg)
	default:
		return sl.Debug(msg)
	}
}

func (v *Package) print(level LogLevel, msg string) {
	lvl := v.GetLogLevel()
	if lvl >= level {
		levelFormat := v.parent.GetLevelFormat()
		packagePrintLen := v.parent.GetPackagePrintLength()
		appName := v.parent.GetApplicationName()
		if appName == "" {
			appName = os.Args[0]
		}
		out1 := fmtStr(true, level, levelFormat, appName,
			v.packageName, packagePrintLen, msg, "%[1]s [%[3]s] %[4]s  %[5]s")
		// File output
		if lf := v.parent.GetLogFileInfo(); lf != nil {
			rotateMaxSize := v.parent.GetRotateMaxSize()
			rotateMaxCount := v.parent.GetRotateMaxCount()
			out2 := fmtStr(false, level, levelFormat, appName,
				v.packageName, packagePrintLen, msg, "%[1]s [%[3]s] %[4]s  %[5]s")
			if err := lf.writeToFile(out2, rotateMaxSize, rotateMaxCount); err != nil {
				err = spew.Errorf("Failed to report syslog message %q: %v\n", out2, err)
				v.parent.log.Fatal(err)
			}
		}
		// Syslog output
		if v.parent.GetSyslogEnabled() {
			if err := v.writeToSyslog(level, levelFormat, appName, msg); err != nil {
				err = spew.Errorf("Failed to report syslog message %q: %v\n", msg, err)
				v.parent.log.Fatal(err)
			}
		}
		// Console output
		v.parent.log.Print(out1 + fmt.Sprintln())
		// Check critical events
		if level == PanicLevel {
			panic(out1)
		} else if level == FatalLevel {
			os.Exit(1)
		}
	}
}

func (v *Package) Printf(level LogLevel, format string, args ...interface{}) {
	lvl := v.GetLogLevel()
	if lvl >= level {
		msg := spew.Sprintf(format, args...)
		v.print(level, msg)
	}
}

func (v *Package) Print(level LogLevel, args ...interface{}) {
	lvl := v.GetLogLevel()
	if lvl >= level {
		msg := fmt.Sprint(args...)
		v.print(level, msg)
	}
}

func (v *Package) Debugf(format string, args ...interface{}) {
	v.Printf(DebugLevel, format, args...)
}

func (v *Package) Debug(args ...interface{}) {
	v.Print(DebugLevel, args...)
}

func (v *Package) Infof(format string, args ...interface{}) {
	v.Printf(InfoLevel, format, args...)
}

func (v *Package) Info(args ...interface{}) {
	v.Print(InfoLevel, args...)
}

func (v *Package) Warningf(format string, args ...interface{}) {
	v.Printf(WarnLevel, format, args...)
}

func (v *Package) Warnf(format string, args ...interface{}) {
	v.Printf(WarnLevel, format, args...)
}

func (v *Package) Warning(args ...interface{}) {
	v.Print(WarnLevel, args...)
}

func (v *Package) Warn(args ...interface{}) {
	v.Print(WarnLevel, args...)
}

func (v *Package) Errorf(format string, args ...interface{}) {
	v.Printf(ErrorLevel, format, args...)
}

func (v *Package) Error(args ...interface{}) {
	v.Print(ErrorLevel, args...)
}

func (v *Package) Panicf(format string, args ...interface{}) {
	v.Printf(PanicLevel, format, args...)
}

func (v *Package) Panic(args ...interface{}) {
	v.Print(PanicLevel, args...)
}

func (v *Package) Fatalf(format string, args ...interface{}) {
	v.Printf(FatalLevel, format, args...)
}

func (v *Package) Fatal(args ...interface{}) {
	v.Print(FatalLevel, args...)
}
