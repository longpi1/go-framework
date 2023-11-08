package extractor

import (
	"github.com/pkg/errors"

	"go-pulgin/internal/datasource/pulsar"
	"go-pulgin/internal/infra/extractor"
	"go-pulgin/internal/infra/plugin"
	"go-pulgin/logger"
)

const (
	pulsarTopicTitle             = "topic"
	pulsarSerializeStrategyTitle = "serialize"
	defaultSerializeStrategy     = "avro"
)

type PulsarExtractor struct {
	topic             string
	serializeStrategy string
	consumer          *pulsar.Consumer
}

func (p *PulsarExtractor) Install() error {
	if p.topic == "" {
		return errors.New("topic not specified")
	}

	consumer, err := pulsar.NewConsumer(p.topic, p.serializeStrategy)
	if err != nil {
		return errors.Errorf("init pulsar consumer failed, topic: %v", p.topic)
	}
	p.consumer = consumer
	return nil
}

func (p *PulsarExtractor) SetContext(ctx plugin.Context) {
	topic, ok := ctx.GetString(pulsarTopicTitle)
	if !ok {
		return
	}
	p.topic = topic

	serializeType := ctx.GetStringOrDefault(pulsarSerializeStrategyTitle, defaultSerializeStrategy)
	p.serializeStrategy = serializeType
}

func (p *PulsarExtractor) Uninstall() {
	p.consumer.Close()
}

func (p *PulsarExtractor) Extract() (*plugin.Event, error) {
	payload, err := p.consumer.RecvWithDecode()
	if err != nil {
		logger.Errorf("[PulsarExtractor] Extract error: %v", err)
		return nil, err
	}
	logger.Debugf("[PulsarExtractor] extract payload: %v", payload)
	event := plugin.NewEvent(payload)
	return event, nil
}

func init() {
	extractor.Add("pulsar_extractor", func() extractor.Plugin {
		return &PulsarExtractor{}
	})
}
