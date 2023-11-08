package loader

import (
	"github.com/pkg/errors"

	"go-pulgin/internal/datasource/pulsar"
	"go-pulgin/internal/infra/loader"
	"go-pulgin/internal/infra/plugin"
)

const (
	pulsarTopicTitle             = "topic"
	pulsarSerializeStrategyTitle = "serialize"
	defaultSerializeStrategy     = "json"
)

type PulsarLoader struct {
	topic             string
	serializeStrategy string
	producer          *pulsar.Producer
}

func (p *PulsarLoader) Install() error {
	if p.topic == "" {
		return errors.New("topic not specified")
	}

	producer, err := pulsar.NewProducer(p.topic, p.serializeStrategy)
	if err != nil {
		return errors.Errorf("init pulsar consumer failed, topic: %v", p.topic)
	}
	p.producer = producer
	return nil
}

func (p *PulsarLoader) Uninstall() {
	p.producer.Close()
}

func (p *PulsarLoader) Load(event *plugin.Event) error {
	p.producer.SendAsync(event.Payload())
	return nil
}

func (p *PulsarLoader) SetContext(ctx plugin.Context) {
	topic, ok := ctx.GetString(pulsarTopicTitle)
	if !ok {
		return
	}
	p.topic = topic

	serializeType := ctx.GetStringOrDefault(pulsarSerializeStrategyTitle, defaultSerializeStrategy)
	p.serializeStrategy = serializeType
}

func init() {
	loader.Add("pulsar_loader", func() loader.Plugin {
		return &PulsarLoader{}
	})
}
