package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	prommodel "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/assert"
)

func Test_dups(t *testing.T) {
	reg := prometheus.NewRegistry()
	reg.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "prom_test",
		Subsystem:   "",
		Name:        "a",
		Help:        "help1",
		ConstLabels: nil,
	}, func() float64 {
		return 1
	}))
	assert.Panics(t, func() {
		reg.MustRegister(prometheus.NewCounterFunc(prometheus.CounterOpts{
			Namespace:   "prom_test",
			Subsystem:   "",
			Name:        "a",
			Help:        "help2",
			ConstLabels: nil,
		}, func() float64 {
			return float64(time.Now().Unix())
		}))
	})
}

func Test_export(t *testing.T) {
	reg := prometheus.NewRegistry()
	reg.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "prom_test",
		Subsystem:   "",
		Name:        "fixed_gauge",
		Help:        "help1",
		ConstLabels: nil,
	}, func() float64 {
		return 1
	}))
	reg.MustRegister(prometheus.NewCounterFunc(prometheus.CounterOpts{
		Namespace:   "prom_test",
		Subsystem:   "sub",
		Name:        "time_counter",
		Help:        "help2",
		ConstLabels: nil,
	}, func() float64 {
		return float64(time.Now().Unix())
	}))
	srv := httptest.NewServer(promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	defer srv.Close()

	resp, err := http.Get(fmt.Sprintf("%s/metrics", srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	dec := expfmt.NewDecoder(resp.Body, expfmt.ResponseFormat(resp.Header))

	ms := make(map[string]*prommodel.MetricFamily, 0)
	for {
		var m prommodel.MetricFamily
		err = dec.Decode(&m)
		if err == io.EOF {
			break // all metrics decoded
		} else if err != nil {
			t.Fatal(err)
		}

		ms[m.GetName()] = &m
	}

	{
		m, ok := ms["prom_test_fixed_gauge"]
		assert.True(t, ok)
		assert.Equal(t, "help1", m.GetHelp())
		assert.Equal(t, prommodel.MetricType_GAUGE, m.GetType())
	}

	{
		m, ok := ms["prom_test_sub_time_counter"]
		assert.True(t, ok)
		assert.Equal(t, "help2", m.GetHelp())
		assert.Equal(t, prommodel.MetricType_COUNTER, m.GetType())
	}
}
