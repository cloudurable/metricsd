package repeater

import (
	lg "github.com/cloudurable/simplelog/logging"
	c "github.com/cloudurable/metricsd/common"
    s "github.com/cloudurable/metricsd/service"
    proto "github.com/optiopay/kafka/proto"
)

type KafkaMetricsRepeater struct {
	logger lg.Logger
    service *s.Kafka
    instanceId string
}

func (this KafkaMetricsRepeater) Repeat(metrics []c.Metric) error {
    messages := []*proto.Message{}
    for _,m := range metrics {
        messages = append(messages,
            &proto.Message{
                Key: []byte(this.instanceId + c.DOT + m.QualifiedName()),
                Value: c.ToJsonBytes(m)})
    }

    this.service.Produce(messages)

    return nil
}

func (this KafkaMetricsRepeater) Verify() bool { return true }

func (this KafkaMetricsRepeater) Name() string { return this.logger.Name() } // already there so just use it

func NewKafkaMetricsRepeater(config *c.Config) *KafkaMetricsRepeater {
    logger := c.GetLogger(config.Debug, c.REPEATER + c.REPEATER_KAFKA)
    awsContext := c.GetAwsContext(logger, config)

    return &KafkaMetricsRepeater{
        logger: logger,
        service: s.NewKafka(logger, config),
        instanceId: awsContext.EC2InstanceId,
    }
}
