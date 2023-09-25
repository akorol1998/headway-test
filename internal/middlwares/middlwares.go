package middlwares

import (
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// ResponseWriteWrapper is wrapper around http.ResponseWriter interface
// the reason we have to create it is because for now there is no way of
// writing to header after Write() or WriteHead() have been performed.
// there are trailers though, but i dont think X-Response-Time should be there
type responseWriteWrapper struct {
	http.ResponseWriter
	isHeaderWritten bool
	start           time.Time
}

func (w *responseWriteWrapper) WriteHeader(statusCode int) {
	w.Header().Set("X-Response-Time", time.Since(w.start).String())
	w.ResponseWriter.WriteHeader(statusCode)
	w.isHeaderWritten = true
}

func (w *responseWriteWrapper) Write(b []byte) (int, error) {
	if !w.isHeaderWritten {
		w.WriteHeader(200)
	}

	return w.ResponseWriter.Write(b)
}

func LogMiddlware(log *zap.SugaredLogger) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Infow("Request",
				"method", r.Method,
				"path", r.URL.Path)
			h.ServeHTTP(w, r)
		}
	}
}

func HeaderMiddlware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := time.Now()
		w.Header().Set("X-Server-Name", r.Host)
		h.ServeHTTP(&responseWriteWrapper{w, false, time.Now()}, r)
		// always converting to microseconds
		w.Header().Set("X-Response-Time", strconv.Itoa(int(time.Since(d).Microseconds())))
	}
}
