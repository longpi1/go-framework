package extractor

import (
	"strings"

	"github.com/pkg/errors"

	"go-pulgin/internal/datasource/kafka"
	"go-pulgin/internal/infra/extractor"
	"go-pulgin/internal/infra/plugin"
)

const (
	kafkaTopicsTitle    = "topics"
	kafkaGroupIdTitle   = "groupId"
	kafkaGroupIdDefault = "go-pulgin"
	kafkaContextSep     = ","
)

type KafkaExtractor struct {
	groupId       string
	topics        []string
	kafkaConsumer *kafka.Consumer
}

func (k *KafkaExtractor) Install() error {
	if len(k.topics) == 0 {
		return errors.New("topics not specified")
	}

	kafkaConsumer, err := kafka.NewKafkaConsumer(k.groupId, k.topics)
	if err != nil {
		return errors.Errorf("init kafka consumer failed, topics: %v, error: %v", k.topics, err)
	}

	k.kafkaConsumer = kafkaConsumer
	return nil
}

func (k *KafkaExtractor) Uninstall() {
	k.kafkaConsumer.Close()
}

func (k *KafkaExtractor) SetContext(ctx plugin.Context) {
	topicStr, ok := ctx.GetString(kafkaTopicsTitle)
	if !ok {
		return
	}
	k.topics = strings.Split(topicStr, kafkaContextSep)

	groupId := ctx.GetStringOrDefault(kafkaGroupIdTitle, kafkaGroupIdDefault)
	k.groupId = groupId
}

func (k *KafkaExtractor) Extract() (*plugin.Event, error) {
	payload := k.kafkaConsumer.Receive()
	event := plugin.NewEvent(payload)
	return event, nil
}

func init() {
	extractor.Add("kafka_extractor", func() extractor.Plugin {
		return &KafkaExtractor{}
	})
}
