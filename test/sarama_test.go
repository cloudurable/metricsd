package test

import (
    //c "github.com/cloudurable/metricsd/common"
    sa "gopkg.in/Shopify/sarama.v1"
    "testing"
    "fmt"
    rcm "github.com/rcrowley/go-metrics"
)

func TestSarama(test *testing.T) {
    broker := sa.NewBroker("localhost:9092")
    err := broker.Open(nil)
    if err != nil {
        panic(err)
    }

    request := sa.MetadataRequest{}
    response, err := broker.GetMetadata(&request)
    if err != nil {
        _ = broker.Close()
        panic(err)
    }

    fmt.Println("There are", len(response.Topics), "topics active in the cluster.\n")

    for _,t := range response.Topics {
        fmt.Println("Name: " + t.Name)
    }

    if err = broker.Close(); err != nil {
        panic(err)
    }
}

func meterStr(m rcm.Meter) string {
    return fmt.Sprint("%d %2.2f %2.2f %2.2f %2.2f",
        m.Count(),m.Rate1(),m.Rate5(), m.Rate15(), m.RateMean(),
    )
}

func histStr(h rcm.Histogram) string {
    return fmt.Sprint("%d %d %2.2f %d",
        h.Count(), h.Max(), h.Mean(), h.Min(),
    )
}