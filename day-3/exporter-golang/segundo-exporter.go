package main 

import (
    "log"
    "net/http"
    "github.com/pbnjay/memory"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func memorialivre() float64{
    memoria_livre :=  memory.FreeMemory()
    return float64(memoria_livre) 

}

func totalmemoria() float64 {
    memoria_total := memory.TotalMemory()
    return float64(memoria_total)
}

var (
    memorialivreBytesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "memoria_livre_bytes",
        Help: "Quantidade de memoria livre em bytes"

    })

    memorialivremegasGauge = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "memoria_livre_megas",
        Help: "Quantidade de memoria livre em megas"

    })

    totalMemoryBytesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "total_memoria_bytes",
        Help: "total de memoria  em bytes"

    })

    totalMemoryGigasGauge = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "total_memoria_gigas",
        Help: "total de memoria  em gigas"

    })
)

func init() {
    prometheus.MustRegister(memorialivreBytesGauge)
    prometheus.MustRegister(memorialivremegasGauge)
    prometheus.MustRegister(totalMemoryBytesGauge)
    prometheus.MustRegister(totalMemoryGigasGauge)
}

func main(){
    memorialivreBytesGauge.Set(memorialivre())
    memorialivremegasGauge.Set(memorialivre() / 1024 / 1024 )
    totalMemoryBytesGauge.Set(totalmemoria())
    totalMemoryGigasGauge.Set(totalmemoria() / 1024 / 1024 / 1024)
    http.Handle("/metrics", promhttp.Handler())

    log.Fatal(http.ListenAndServe(":7788", nil))
}