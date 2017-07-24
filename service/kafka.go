package service

import (
    k "github.com/optiopay/kafka"
    proto "github.com/optiopay/kafka/proto"
    c "github.com/cloudurable/metricsd/common"
    l "github.com/cloudurable/simplelog/logging"
    "hash/fnv"
)

type Kafka struct {
    logger l.Logger
    topic string
    partitions int32
    broker *k.Broker
    producer k.Producer
}

func NewKafka(logger l.Logger, config *c.Config) *Kafka {
    logger = c.EnsureLogger(logger, config.Debug, "kafka-service")
    addresses := config.KafkaConfig.ClusterAddresses
    broker, err := k.Dial(addresses, k.NewBrokerConf("metricsd"))
    if err != nil {
        panic(err)
    }
    defer broker.Close()

    conf := k.NewProducerConf()
    conf.RequiredAcks = proto.RequiredAcksLocal

    producer := broker.Producer(conf)

    return &Kafka{
        logger: logger,
        topic: config.KafkaConfig.Topic,
        partitions: int32(config.KafkaConfig.Partitions),
        broker: broker,
        producer: producer,
    }
}

func (this *Kafka) Produce(messages []*proto.Message) {
    for _, msg := range messages {
        if _, err := this.producer.Produce(this.topic, hash(msg.Key, this.partitions), msg); err != nil {
            panic(err)
        }
    }
}

func hash(keyBytes []byte, partitions int32) int32 {
    hasher := fnv.New32()
    if _, err := hasher.Write(keyBytes); err != nil {
        panic(err)
    }

    hash := int32(hasher.Sum32())
    if hash < 0 {
        hash = -hash
    }
    mod := hash % partitions

    return mod
}
