package overview

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"
)

type KPI struct {
	ID    string  `json:"id"`
	Label string  `json:"label"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
	Delta float64 `json:"delta"`
	// Better is the direction a healthy value moves: "up" or "down".
	Better string `json:"better"`
}

type Point struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

type Status string

const (
	StatusDone    Status = "done"
	StatusRunning Status = "running"
	StatusFailed  Status = "failed"
)

type Activity struct {
	ID     string    `json:"id"`
	Title  string    `json:"title"`
	Status Status    `json:"status"`
	At     time.Time `json:"at"`
}

type Overview struct {
	KPIs []KPI `json:"kpis"`
	// SeriesMetric names what Series measures, as a translation key suffix.
	SeriesMetric string     `json:"seriesMetric"`
	Series       []Point    `json:"series"`
	Recent       []Activity `json:"recent"`
	// Demo marks synthetic figures so the UI can label them.
	Demo bool `json:"demo"`
}

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/overview", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(demo(time.Now()))
	})
}

func demo(now time.Time) Overview {
	series := make([]Point, 30)
	v := 1200.0
	for i := range series {
		v += rand.Float64()*160 - 60
		series[i] = Point{Date: now.AddDate(0, 0, i-29).Format("2006-01-02"), Value: v}
	}

	statuses := []Status{StatusDone, StatusDone, StatusRunning, StatusFailed}
	recent := make([]Activity, 6)
	for i := range recent {
		recent[i] = Activity{
			ID:     fmt.Sprintf("job-%d", 100+i),
			Title:  fmt.Sprintf("Sync task #%d", 100+i),
			Status: statuses[rand.IntN(len(statuses))],
			At:     now.Add(-time.Duration(i*37) * time.Minute),
		}
	}

	return Overview{
		KPIs: []KPI{
			{ID: "users", Label: "Active users", Value: 3421, Delta: 4.2, Better: "up"},
			{ID: "revenue", Label: "Revenue", Value: 18250, Unit: "USD", Delta: -1.3, Better: "up"},
			{ID: "latency", Label: "p95 latency", Value: 182, Unit: "ms", Delta: -8.5, Better: "down"},
			{ID: "errors", Label: "Error rate", Value: 0.42, Unit: "%", Delta: 0.1, Better: "down"},
		},
		SeriesMetric: "dailyUsers",
		Series:       series,
		Recent:       recent,
		Demo:         true,
	}
}
