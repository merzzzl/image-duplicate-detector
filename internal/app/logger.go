package app

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type Logger struct {
	out *stageLoggerOut
	zerolog.Logger
}

var logger = initLogger()

func (l *Logger) SetStage(s string) {
	l.out.mux.Lock()
	defer l.out.mux.Unlock()

	l.out.stage = s
	l.out.buff.Write([]byte("\r"))
}

func (l *Logger) SetStageStatus(s string) {
	l.out.mux.Lock()
	defer l.out.mux.Unlock()

	l.out.status = s
	l.out.buff.Write([]byte("\r"))
}

func GetLogger() Logger {
	return logger
}

func EnableDebug() {
	logger.Logger = logger.Level(zerolog.DebugLevel)
}

func initLogger() Logger {
	zerolog.LevelColors[zerolog.DebugLevel] = 35

	o := stageLoggerOut{
		stage:  "Initialize",
		status: "running...",
		buff:   new(bytes.Buffer),
	}

	go o.lazyWriter()

	l := zlog.Output(zerolog.ConsoleWriter{
		Out: &o,
		FormatMessage: func(i any) string {
			return fmt.Sprintf("\033[37m%s\033[0m", i)
		},
		FormatFieldName: func(i any) string {
			return fmt.Sprintf("\033[34m%s=\033[0m", i)
		},
		FormatFieldValue: func(i any) string {
			return fmt.Sprintf("\033[34m%s\033[0m", i)
		},
		FormatErrFieldName: func(any) string {
			return ""
		},
		FormatErrFieldValue: func(i any) string {
			return fmt.Sprintf("\033[31m(%s)\033[0m", i)
		},
	}).Level(zerolog.InfoLevel)

	return Logger{
		out:    &o,
		Logger: l,
	}
}

type stageLoggerOut struct {
	stage  string
	status string
	mux    sync.Mutex
	buff   *bytes.Buffer
}

func (o *stageLoggerOut) Write(p []byte) (n int, err error) {
	return o.buff.Write(p)
}

func (o *stageLoggerOut) lazyWriter() {
	ticker := time.NewTicker(time.Microsecond * 100)

	for range ticker.C {
		if err := o.writeAll(); err != nil {
			panic(err)
		}
	}
}

func (o *stageLoggerOut) writeAll() error {
	o.mux.Lock()
	defer o.mux.Unlock()

	bytes, err := io.ReadAll(o.buff)
	if err != nil {
		return err
	}

	if len(bytes) == 0 {
		return nil
	}

	status := fmt.Sprintf("\033[K\033[1m%s \033[36m%s: \033[37m%s\033[0m\n\033[A\r", time.Now().Format(time.Kitchen), o.stage, o.status)

	bytes = append(bytes, []byte(status)...)

	_, err = os.Stdout.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
