package observability

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type key struct {
	Service string
	Method  string
	Path    string
	Status  int
}

type requestMetrics struct {
	Count       int64
	DurationSum float64
}

var recorder = &metricsRecorder{
	startedAt: time.Now().UTC(),
	requests:  make(map[key]*requestMetrics),
}

type metricsRecorder struct {
	mu        sync.RWMutex
	startedAt time.Time
	requests  map[key]*requestMetrics
}

func Middleware(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		recorder.record(key{
			Service: service,
			Method:  c.Request.Method,
			Path:    path,
			Status:  c.Writer.Status(),
		}, time.Since(start))
	}
}

func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(200, "text/plain; version=0.0.4; charset=utf-8", []byte(recorder.render()))
	}
}

func (r *metricsRecorder) record(k key, duration time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.requests[k]
	if !ok {
		item = &requestMetrics{}
		r.requests[k] = item
	}
	item.Count++
	item.DurationSum += duration.Seconds()
}

func (r *metricsRecorder) render() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var b strings.Builder
	b.WriteString("# HELP depay_process_uptime_seconds Service uptime in seconds.\n")
	b.WriteString("# TYPE depay_process_uptime_seconds gauge\n")
	fmt.Fprintf(&b, "depay_process_uptime_seconds %.0f\n", time.Since(r.startedAt).Seconds())
	b.WriteString("# HELP depay_http_requests_total Total HTTP requests by service, method, path, and status.\n")
	b.WriteString("# TYPE depay_http_requests_total counter\n")
	b.WriteString("# HELP depay_http_request_duration_seconds_sum Total HTTP request duration in seconds.\n")
	b.WriteString("# TYPE depay_http_request_duration_seconds_sum counter\n")

	keys := make([]key, 0, len(r.requests))
	for k := range r.requests {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Service != keys[j].Service {
			return keys[i].Service < keys[j].Service
		}
		if keys[i].Path != keys[j].Path {
			return keys[i].Path < keys[j].Path
		}
		if keys[i].Method != keys[j].Method {
			return keys[i].Method < keys[j].Method
		}
		return keys[i].Status < keys[j].Status
	})

	for _, k := range keys {
		item := r.requests[k]
		labels := fmt.Sprintf(`service="%s",method="%s",path="%s",status="%d"`, escape(k.Service), escape(k.Method), escape(k.Path), k.Status)
		fmt.Fprintf(&b, "depay_http_requests_total{%s} %d\n", labels, item.Count)
		fmt.Fprintf(&b, "depay_http_request_duration_seconds_sum{%s} %.6f\n", labels, item.DurationSum)
	}
	return b.String()
}

func escape(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `"`, `\"`)
	return value
}
