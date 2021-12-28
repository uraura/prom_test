package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main()  {
	reg := prometheus.NewRegistry()
	reg.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "prom_test",
		Subsystem:   "",
		Name:        "fixed1",
		Help:        "",
	}, func() float64 {
		return 1
	}))
	reg.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "prom_test",
		Subsystem:   "",
		Name:        "fixed2",
		Help:        "",
	}, func() float64 {
		return 2
	}))
	i := 0.0
	reg.MustRegister(prometheus.NewCounterFunc(prometheus.CounterOpts{
		Namespace:   "prom_test",
		Subsystem:   "",
		Name:        "count1",
		Help:        "",
	}, func() float64 {
		i++
		return i
	}))

	err := http.ListenAndServe(":8080", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	if err != nil  {
		log.Fatalln(err)
	}
}