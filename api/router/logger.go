package router

import (
	"fmt"
	"net/http"
	"time"
	
	"github.com/MadAppGang/httplog"
	"github.com/gorilla/mux"
	
	log "github.com/sirupsen/logrus"
)

// logger
func logger(h http.HandlerFunc) http.Handler {
	return httplog.LoggerWithConfig(
		httplog.LoggerConfig{
			RouterName:  "FillBodyFormatter",
			Formatter:   httplog.DefaultLogFormatter,
			CaptureBody: true,
		},
		http.HandlerFunc(h),
	)
}

// CustomLogFormatter
func CustomLogFormatter(param httplog.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("%s |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n",
		param.RouterName,
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
	)
}

// RequestLoggerMiddleware
func RequestLoggerMiddleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			defer func() {
				log.Infof(
					"Handled request [%s] %s %s %s",
					req.Method,
					req.Host,
					req.URL.Path,
					req.URL.RawQuery,
				)
			}()
			next.ServeHTTP(w, req)
		})
	}
}
