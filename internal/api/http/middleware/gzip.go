package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	gz           *gzip.Writer
	needCompress bool
}

func newGzipResponseWriter(w http.ResponseWriter) *gzipResponseWriter {
	return &gzipResponseWriter{
		ResponseWriter: w,
		gz:             nil,
		needCompress:   false,
	}
}

func (gw *gzipResponseWriter) WriteHeader(code int) {
	contentType := gw.Header().Get("Content-Type")

	gw.needCompress = strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "text/html")

	if gw.needCompress {
		gw.Header().Set("Content-Encoding", "gzip")
		gw.gz = gzip.NewWriter(gw.ResponseWriter)
	}

	gw.ResponseWriter.WriteHeader(code)
}

func (gw *gzipResponseWriter) Write(b []byte) (int, error) {
	if gw.needCompress && gw.gz != nil {
		return gw.gz.Write(b)
	}
	return gw.ResponseWriter.Write(b)
}

func (gw *gzipResponseWriter) Close() error {
	if gw.gz != nil {
		return gw.gz.Close()
	}
	return nil
}

// GzipCompress создаёт middleware для сжатия ответов.
func GzipCompress() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			gz := newGzipResponseWriter(w)
			defer gz.Close()

			next.ServeHTTP(gz, r)
		})
	}
}

type gzipReadCloser struct {
	io.ReadCloser
	gzReader *gzip.Reader
}

func (g *gzipReadCloser) Read(p []byte) (n int, err error) {
	return g.gzReader.Read(p)
}

func (g *gzipReadCloser) Close() error {
	if err := g.gzReader.Close(); err != nil {
		return err
	}
	return g.ReadCloser.Close()
}

// GzipDecompress создаёт middleware для декомпрессии запросов.
func GzipDecompress() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Content-Encoding") != "gzip" {
				next.ServeHTTP(w, r)
				return
			}

			gzReader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "ошибка декомпрессии", http.StatusBadRequest)
				return
			}

			r.Body = &gzipReadCloser{
				ReadCloser: r.Body,
				gzReader:   gzReader,
			}

			next.ServeHTTP(w, r)
		})
	}
}
