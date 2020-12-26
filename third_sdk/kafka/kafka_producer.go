package kafka

import (
	"encoding/json"
	"git.moumentei.com/plat_go/golib/log"
	"gopkg.in/Shopify/sarama.v1"
)

// Kafka is kafka config.
type KafkaConf struct {
	Topic   string
	Group   string
	Brokers []string
}

type KafkaPubClient struct {
	c        *KafkaConf
	kafkaPub sarama.SyncProducer
}

func NewKafkaPubClient(c *KafkaConf) *KafkaPubClient {
	k := &KafkaPubClient{
		c:        c,
		kafkaPub: newKafkaPub(c),
	}
	return k
}

func newKafkaPub(c *KafkaConf) sarama.SyncProducer {
	kc := sarama.NewConfig()
	kc.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	kc.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	kc.Producer.Return.Successes = true
	pub, err := sarama.NewSyncProducer(c.Brokers, kc)
	if err != nil {
		panic(err)
	}
	return pub
}

func (cli *KafkaPubClient) PushMsg(key string, pushMsg interface{}) error {
	b, err := json.Marshal(pushMsg)
	if err != nil {
		return err
	}
	m := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(key),
		Topic: cli.c.Topic,
		Value: sarama.ByteEncoder(b),
	}
	if _, _, err = cli.kafkaPub.SendMessage(m); err != nil {
		log.Errorf("PushMsg.send(push pushMsg:%v) error(%v)", pushMsg, err)
		return err
	}
	return nil
}

func (cli *KafkaPubClient) PushMsgString(key string, msg []byte) error {
	m := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(key),
		Topic: cli.c.Topic,
		Value: sarama.ByteEncoder(msg),
	}
	if _, _, err := cli.kafkaPub.SendMessage(m); err != nil {
		log.Errorf("PushMsg.send(push pushMsg:%v) error(%v)", string(msg), err)
		return err
	}
	return nil
}
