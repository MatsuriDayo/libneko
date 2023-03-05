package neko_log

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/matsuridayo/libneko/neko_common"
)

var LogWriter *logWriter
var LogWriterDisable = false
var NB4AGuiLogWriter io.Writer

func SetupLog(maxSize int, path string) (err error) {
	if LogWriter != nil {
		return
	}
	// Truncate mod from libcore, simplify because only 1 proccess.
	oldBytes, err := os.ReadFile(path)
	needTruncate := len(oldBytes) > maxSize
	if err == nil && needTruncate {
		if os.Truncate(path, 0) == nil {
			oldBytes = oldBytes[len(oldBytes)-maxSize:]
		}
	}
	// open
	var f *os.File
	f, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err == nil {
		if needTruncate {
			_, _ = f.Write(oldBytes)
		}
	} else {
		err = fmt.Errorf("error open log: %v", err)
	}
	//
	LogWriter = &logWriter{}
	if neko_common.RunMode == neko_common.RunMode_NekoBoxForAndroid {
		LogWriter.writers = []io.Writer{NB4AGuiLogWriter, f}
	} else {
		LogWriter.writers = []io.Writer{os.Stdout, f}
	}
	// setup std log
	log.SetFlags(log.LstdFlags | log.LUTC)
	log.SetOutput(LogWriter)
	return
}

type logWriter struct {
	writers []io.Writer
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	if LogWriterDisable {
		return len(p), nil
	}

	for _, w := range w.writers {
		if w == nil {
			continue
		}
		n, err = w.Write(p)
		if err != nil {
			return
		}
	}
	return
}

func (w *logWriter) Truncate() {
	for _, w := range w.writers {
		if w == nil {
			continue
		}
		if f, ok := w.(*os.File); ok {
			_ = f.Truncate(0)
		}
	}
}

func (w *logWriter) Close() error {
	for _, w := range w.writers {
		if w == nil {
			continue
		}
		if f, ok := w.(*os.File); ok {
			_ = f.Close()
		}
	}
	return nil
}
