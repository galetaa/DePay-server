package observability

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"shared/logging"
	"shared/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		duration := time.Since(start)
		recorder.record(key{
			Service: service,
			Method:  c.Request.Method,
			Path:    path,
			Status:  c.Writer.Status(),
		}, duration)
		if logging.Logger != nil {
			logging.Logger.Info("http request",
				zap.String("service", service),
				zap.String("request_id", middleware.RequestID(c)),
				zap.String("method", c.Request.Method),
				zap.String("route", path),
				zap.Int("status", c.Writer.Status()),
				zap.Duration("duration", duration),
			)
		}
	}
}

func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(200, "text/plain; version=0.0.4; charset=utf-8", []byte(recorder.render()))
	}
}

func DatabaseHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var b strings.Builder
		b.WriteString(recorder.render())
		if db != nil {
			b.WriteString(renderDatabaseMetrics(c.Request.Context(), db))
		}
		c.Data(200, "text/plain; version=0.0.4; charset=utf-8", []byte(b.String()))
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

func renderDatabaseMetrics(ctx context.Context, db *sql.DB) string {
	var b strings.Builder
	writeGroupedCount(ctx, &b, db, "depay_transactions_total", "status", `SELECT status::text, count(*) FROM payment_transactions GROUP BY status`)
	writeGroupedCount(ctx, &b, db, "depay_invoices_total", "status", `SELECT status::text, count(*) FROM payment_invoices GROUP BY status`)
	writeGroupedCount(ctx, &b, db, "depay_webhook_deliveries_total", "status", `SELECT status, count(*) FROM merchant_webhook_deliveries GROUP BY status`)
	writeRiskAlerts(ctx, &b, db)
	writeKYCApplications(ctx, &b, db)
	writeRPCLatency(ctx, &b, db)
	return b.String()
}

func writeGroupedCount(ctx context.Context, b *strings.Builder, db *sql.DB, metric string, labelName string, query string) {
	fmt.Fprintf(b, "# TYPE %s gauge\n", metric)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		fmt.Fprintf(b, "# %s unavailable: %s\n", metric, escape(err.Error()))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var label string
		var count int64
		if err := rows.Scan(&label, &count); err != nil {
			continue
		}
		fmt.Fprintf(b, `%s{%s="%s"} %d`+"\n", metric, labelName, escape(label), count)
	}
}

func writeRiskAlerts(ctx context.Context, b *strings.Builder, db *sql.DB) {
	b.WriteString("# TYPE depay_risk_alerts_total gauge\n")
	rows, err := db.QueryContext(ctx, `SELECT risk_level::text, status::text, count(*) FROM risk_alerts GROUP BY risk_level, status`)
	if err != nil {
		fmt.Fprintf(b, "# depay_risk_alerts_total unavailable: %s\n", escape(err.Error()))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var level string
		var status string
		var count int64
		if err := rows.Scan(&level, &status, &count); err != nil {
			continue
		}
		fmt.Fprintf(b, `depay_risk_alerts_total{level="%s",status="%s"} %d`+"\n", escape(level), escape(status), count)
	}
}

func writeKYCApplications(ctx context.Context, b *strings.Builder, db *sql.DB) {
	writeGroupedCount(ctx, b, db, "depay_kyc_applications_total", "status", `SELECT status::text, count(*) FROM kyc_applications GROUP BY status`)
}

func writeRPCLatency(ctx context.Context, b *strings.Builder, db *sql.DB) {
	b.WriteString("# TYPE depay_rpc_node_latency_ms gauge\n")
	rows, err := db.QueryContext(ctx, `SELECT node_name, chain_name, COALESCE(avg_latency_ms, 0) FROM vw_rpc_node_status`)
	if err != nil {
		fmt.Fprintf(b, "# depay_rpc_node_latency_ms unavailable: %s\n", escape(err.Error()))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var node string
		var chain string
		var latency float64
		if err := rows.Scan(&node, &chain, &latency); err != nil {
			continue
		}
		fmt.Fprintf(b, `depay_rpc_node_latency_ms{node="%s",chain="%s"} %.2f`+"\n", escape(node), escape(chain), latency)
	}
}
