package kafka

import (
	cluster "github.com/bsm/sarama-cluster"
	"github.com/hero1s/golib/log"
)

// Job is push job.
type ConsumerClient struct {
	c        *KafkaConf
	consumer *cluster.Consumer
}

// New new a push job.
func New(c *KafkaConf) *ConsumerClient {
	j := &ConsumerClient{
		c:        c,
		consumer: newKafkaSub(c),
	}
	return j
}

func newKafkaSub(c *KafkaConf) *cluster.Consumer {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	consumer, err := cluster.NewConsumer(c.Brokers, c.Group, []string{c.Topic}, config)
	if err != nil {
		panic(err)
	}
	return consumer
}

// Close close resounces.
func (j *ConsumerClient) Close() error {
	if j.consumer != nil {
		return j.consumer.Close()
	}
	return nil
}

// Consume messages, watch signals
func (j *ConsumerClient) Consume(process func(key, msg []byte)) {
	for {
		select {
		case err := <-j.consumer.Errors():
			log.Errorf("consumer error(%v)", err)
		case n := <-j.consumer.Notifications():
			log.Infof("consumer rebalanced(%v)", n)
		case msg, ok := <-j.consumer.Messages():
			if !ok {
				return
			}
			j.consumer.MarkOffset(msg, "")
			// process push message
			process(msg.Key, msg.Value)
			log.Infof("consume: %s/%d/%d\t%s\t%s", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
		}
	}
}
