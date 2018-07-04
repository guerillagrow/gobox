package shell

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"
)

func isLinuxMacOSFreeBSD() bool {
	return runtime.GOOS == "linux" || runtime.GOOS == "darwin" ||
		runtime.GOOS == "freebsd"
}

func CloseKillChannelOnKillSignal(kill chan struct{}) {
	// Set up channel on which to send signal notifications
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal
	c := make(chan os.Signal, 1)
	sigs := []os.Signal{os.Interrupt, os.Kill}
	if isLinuxMacOSFreeBSD() {
		sigs = append(sigs, syscall.SIGTERM)
	}
	signal.Notify(c, sigs...)
	// run gorutine and block until a signal is received
	go func() {
		<-c
		// send signal to threads about pending to close
		log.Println("Signal received, close kill channel")
		close(kill)
	}()
}

func CloseContextOnKillSignal(cancel context.CancelFunc) {
	// Set up channel on which to send signal notifications
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal
	c := make(chan os.Signal, 1)
	sigs := []os.Signal{os.Interrupt, os.Kill}
	if isLinuxMacOSFreeBSD() {
		sigs = append(sigs, syscall.SIGTERM)
	}
	signal.Notify(c, sigs...)
	// run gorutine and block until a signal is received
	go func() {
		<-c
		// send signal to threads about pending to close
		log.Println("Signal received, cancel context")
		if cancel != nil {
			cancel()
		}
	}()
}

type ExitCodeOrError struct {
	ExitCode int
	Error    error
}

type App struct {
	cmd             *exec.Cmd
	waitCh          chan ExitCodeOrError
	exitCodeOrError atomic.Value
	// stdOut          atomic.Value
}

func NewApp(name string, args ...string) *App {
	cmd := exec.Command(name, args...)

	// cmd.Env = append(os.Environ())
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	app := &App{cmd: cmd}
	return app
}

func (app *App) Run(stdOut *bytes.Buffer, stdErr *bytes.Buffer) ExitCodeOrError {
	_, err := app.Start(stdOut, stdErr)
	if err != nil {
		return ExitCodeOrError{0, err}
	}
	st := app.Wait()
	return st
}

func (app *App) sendExitCodeOrError(exitCode int, err error) {
	state := &ExitCodeOrError{ExitCode: exitCode, Error: err}
	app.exitCodeOrError.Store(state)
	app.waitCh <- *state
}

func readFromIo(read io.ReadCloser, buf *bytes.Buffer) error {
	b := make([]byte, 4096)
	for {
		n, err := read.Read(b)
		if n > 0 {
			buf.Write(b[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *App) asyncWait(stdOut, stdErr *bytes.Buffer,
	readOut, readErr io.ReadCloser) {
	defer close(app.waitCh)
	if readOut != nil {
		err2 := readFromIo(readOut, stdOut)
		if err2 != nil {
			app.sendExitCodeOrError(0, err2)
			return
		}
	}
	if readErr != nil {
		err2 := readFromIo(readErr, stdErr)
		if err2 != nil {
			app.sendExitCodeOrError(0, err2)
			return
		}
	}
	err := app.cmd.Wait()
	var exitCode int
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if stat, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = stat.ExitStatus()
				// reset error, since exitCode already not equal to zero
				err = nil
			}
		}
	}
	app.sendExitCodeOrError(exitCode, err)
}

func (app *App) Start(stdOut *bytes.Buffer,
	stdErr *bytes.Buffer) (chan ExitCodeOrError, error) {
	var readOut io.ReadCloser
	var readErr io.ReadCloser
	var err error
	if stdOut != nil {
		readOut, err = app.cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
	}
	if stdErr != nil {
		readErr, err = app.cmd.StderrPipe()
		if err != nil {
			return nil, err
		}
	}
	err = app.cmd.Start()
	if err != nil {
		return nil, err
	}
	app.waitCh = make(chan ExitCodeOrError)
	go app.asyncWait(stdOut, stdErr, readOut, readErr)
	return app.waitCh, nil
}

func (app *App) CheckIsInstalled() error {
	whApp := NewApp("which", app.cmd.Path)
	st := whApp.Run(nil, nil)
	if st.Error != nil {
		return st.Error
	}
	if st.ExitCode != 0 {
		return fmt.Errorf("App \"%s\" does not exist", app.cmd.Path)
	}
	return nil
}

func (app *App) ExitCodeOrError() *ExitCodeOrError {
	ref := app.exitCodeOrError.Load()
	return ref.(*ExitCodeOrError)
}

func (app *App) Wait() ExitCodeOrError {
	st, ok := <-app.waitCh
	if ok {
		return st
	} else {
		return ExitCodeOrError{ExitCode: 0, Error: fmt.Errorf("Exited already")}
	}
}

func (app *App) Kill() error {
	log.Println(fmt.Sprintf("Start killing app: %v", app.cmd))
	if isLinuxMacOSFreeBSD() {
		// Kill not only main but all children processes,
		// so extract for this purpose group id.
		pgid, err := syscall.Getpgid(app.cmd.Process.Pid)
		if err != nil {
			return err
		}
		// Specifiing gid with negative sign led to killing children processes as well.
		err = syscall.Kill(-pgid, syscall.SIGKILL)
		if err != nil {
			return err
		}
	} else {
		// Kill only mother process
		err := app.cmd.Process.Kill()
		if err != nil {
			return err
		}
	}
	state := app.Wait()
	log.Println(fmt.Sprintf("Done killing app: %v", app.cmd))
	return state.Error
}

func CheckRunAsRoot() bool {
	uid := os.Geteuid()
	if uid == 0 {
		return true
	}
	return false
}
