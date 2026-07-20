package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	mu     sync.Mutex
	logger = log.New(os.Stdout, "", 0)
	now    = time.Now
)

func SetOutput(w io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	logger.SetOutput(w)
}

func SetNow(fn func() time.Time) {
	mu.Lock()
	defer mu.Unlock()
	now = fn
}

func Reset() {
	mu.Lock()
	defer mu.Unlock()
	logger.SetOutput(os.Stdout)
	now = time.Now
}

func Infof(format string, args ...any) {
	logf("INFO", format, args...)
}

func Errorf(format string, args ...any) {
	logf("ERROR", format, args...)
}

func Fatalf(format string, args ...any) {
	logf("ERROR", format, args...)
	os.Exit(1)
}

func logf(level string, format string, args ...any) {
	mu.Lock()
	defer mu.Unlock()
	timestamp := now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logger.Printf("time=\"%s\" level=%s msg=\"%s\"", timestamp, strings.ToLower(level), escape(message))
}

func escape(message string) string {
	message = strings.ReplaceAll(message, `\`, `\\`)
	return strings.ReplaceAll(message, `"`, `\"`)
}
