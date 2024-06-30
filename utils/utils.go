package utils

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var gzPool = sync.Pool{
	New: func() interface{} {
		w := gzip.NewWriter(io.Discard)
		return w
	},
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(status)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)

		gz.Reset(w)
		defer gz.Close()

		next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
func GzipF(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "text/html")
		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)

		gz.Reset(w)
		defer gz.Close()

		next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

var out = createOutput()

type fileWriter struct {
}

var ansi = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

func (w fileWriter) Write(data []byte) (n int, err error) {
	file, err := os.OpenFile("tmp/stdout.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return len(data), err
	}
	defer file.Close()
	_, err = file.Write(ansi.ReplaceAll(data, nil))
	return len(data), err
}
func createOutput() zerolog.ConsoleWriter {
	file := fileWriter{}

	writer := io.MultiWriter(os.Stdout, file)

	o := zerolog.ConsoleWriter{Out: writer, TimeFormat: time.TimeOnly}

	o.FormatLevel = func(i interface{}) string {
		return ""
	}
	o.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m", 34, i)
	}
	o.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("\x1b[%dm%s:\x1b[0m", 2, i)
	}
	o.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprint(i)
	}
	return o
}

var log = zerolog.New(out).With().Timestamp().Logger()

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := time.Now()
		next.ServeHTTP(w, r)
		log.Info().Str("Path", r.URL.Path).Str("Time", time.Since(s).String()).Str("IP", getIp(r)).Msg(r.Method)
	})

}
func getIp(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip != "" && ip != r.RemoteAddr {
		return r.RemoteAddr + "_" + ip
	}
	if ip != "" {
		return ip
	}
	return r.RemoteAddr
}
