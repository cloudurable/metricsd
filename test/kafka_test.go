package test

import (
    //c "github.com/cloudurable/metricsd/common"
    g "github.com/cloudurable/metricsd/gatherer"
    s "github.com/cloudurable/metricsd/service"
    proto "github.com/optiopay/kafka/proto"
    "testing"
    "github.com/cloudurable/metricsd/common"
)

func TestKafka(test *testing.T) {
    logger := GetTestLogger(test, "kafka")

    config := GetTestConfig(logger)
    //config := c.Config{
    //    Debug: true,
    //    DiskConfig: c.DiskGathererConfig{
    //        Command: "df",
    //        FileSystems: []string{"/dev/*", "C:*"},
    //    },
    //
    //}

    gatherer := g.NewDiskMetricsGatherer(nil, config)

    metrics, _:= gatherer.Gather()
    messages := []*proto.Message{}
    for _,m := range metrics {
        messages = append(messages, &proto.Message{Key: []byte(m.QualifiedName()), Value: common.ToJsonBytes(m)} )
    }

    service := s.NewKafka(logger, config)
    service.Produce(messages)

    metrics, _= gatherer.Gather()
    messages = []*proto.Message{}
    for _,m := range metrics {
        messages = append(messages, &proto.Message{Key: []byte(m.QualifiedName()), Value: common.ToJsonBytes(m)} )
    }
    service.Produce(messages)
}
