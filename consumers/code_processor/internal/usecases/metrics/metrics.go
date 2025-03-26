package metrics

import (
	"code_processor/internal/usecases"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)


var (
	opsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "consumer_processed_ops_total",
		Help: "The total number of processed events",
	}, []string{"translator"})

    opsTimes = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "consumer_ops_processing_time",
        Help: "Processing time of events",
    })
)

func CollectMetrics(resp *usecases.ProcessingServiceResponse) error {
    cntr, err := opsProcessed.GetMetricWith(prometheus.Labels{
        "translator": resp.Translator,
    })

    if err != nil {
        return err
    }

    cntr.Inc()
    opsTimes.Set(resp.ProcessingTime.Seconds())

    return nil
}
