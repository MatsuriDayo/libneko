package neko_log

import (
	"io"
	"log"
	"os"

	"github.com/matsuridayo/libneko/neko_common"
)

var LogWriter *logWriter

func SetupLog(maxSize int, path string) {
	if LogWriter != nil {
		return
	}
	// Truncate mod from libcore, simplify because only 1 proccess.
	oldBytes, err := os.ReadFile(path)
	if err == nil && len(oldBytes) > maxSize {
		if os.Truncate(path, 0) == nil {
			oldBytes = oldBytes[len(oldBytes)-maxSize:]
		}
	}
	// open
	f_neko_log, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err == nil {
		f_neko_log.Write(oldBytes)
	} else {
		log.Println("error open log", err)
	}
	//
	LogWriter = &logWriter{}
	if neko_common.RunMode == neko_common.RunMode_NekoBoxForAndroid {
		LogWriter.writers = []io.Writer{neko_common.NB4A_GuiLogWriter, f_neko_log}
	} else {
		LogWriter.writers = []io.Writer{os.Stdout, f_neko_log}
	}
	// setup std log
	log.SetFlags(log.LstdFlags | log.LUTC)
	log.SetOutput(LogWriter)
}

type logWriter struct {
	writers []io.Writer
}

func (w *logWriter) Write(p []byte) (n int, err error) {
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
			f.Truncate(0)
		}
	}
}

func (w *logWriter) Close() error {
	for _, w := range w.writers {
		if w == nil {
			continue
		}
		if f, ok := w.(*os.File); ok {
			f.Close()
		}
	}
	return nil
}
