package pipeline

import (
	"sync/atomic"

	"go-pulgin/internal/infra/config"
	"go-pulgin/internal/infra/extractor"
	"go-pulgin/internal/infra/loader"
	"go-pulgin/internal/infra/plugin"
	"go-pulgin/internal/infra/transformer"
	"go-pulgin/logger"
)

type Plugin interface {
	plugin.Plugin
	SetExtractor(extractor extractor.Plugin)
	SetTransformer(transformer transformer.Plugin)
	SetLoader(loader loader.Plugin)
	IsClosed() bool
}

func NewPlugin(config config.Pipeline) (Plugin, error) {
	pipelineCreator, ok := Pipelines[config.PluginType]
	if !ok {
		logger.Errorf("pipeline plugin not found: %v", config.PluginType)
		return nil, plugin.ErrUnknownPlugin
	}
	pipelinePlugin := pipelineCreator()
	pipelinePlugin.SetContext(config.Ctx)

	// 1.设置extractor
	extractorPlugin, err := extractor.NewPlugin(config.Extractor)
	if err != nil {
		return nil, err
	}
	pipelinePlugin.SetExtractor(extractorPlugin)

	// 2.设置transformer
	var transformerPlugins []transformer.Plugin
	for _, transformerConf := range config.Transformers {
		transformerPlugin, err := transformer.NewPlugin(transformerConf)
		if err != nil {
			return nil, err
		}
		transformerPlugins = append(transformerPlugins, transformerPlugin)
	}
	transformerChain := transformer.NewChain(transformerPlugins)
	pipelinePlugin.SetTransformer(transformerChain)

	// 3.设置loader
	loaderPlugin, err := loader.NewPlugin(config.Loader)
	if err != nil {
		return nil, err
	}
	pipelinePlugin.SetLoader(loaderPlugin)

	return pipelinePlugin, nil
}

type Template struct {
	extractor   extractor.Plugin
	transformer transformer.Plugin
	loader      loader.Plugin
	isRunning   uint32
	Run         func()
}

func (t *Template) Install() error {
	if err := t.loader.Install(); err != nil {
		return err
	}
	if err := t.transformer.Install(); err != nil {
		return err
	}
	if err := t.extractor.Install(); err != nil {
		return err
	}
	atomic.StoreUint32(&t.isRunning, 1)

	t.Run()
	return nil
}

func (t *Template) Uninstall() {
	atomic.StoreUint32(&t.isRunning, 0)
	t.extractor.Uninstall()
	t.transformer.Uninstall()
	t.loader.Uninstall()
}

func (t *Template) SetExtractor(extractor extractor.Plugin) {
	t.extractor = extractor
}

func (t *Template) SetTransformer(transformer transformer.Plugin) {
	t.transformer = transformer
}

func (t *Template) SetLoader(loader loader.Plugin) {
	t.loader = loader
}

// DoRun pipeline默认执行流
func (t *Template) DoRun() {
	for atomic.LoadUint32(&t.isRunning) == 1 {

		// step1: extract
		event, err := t.extractor.Extract()
		if err != nil {
			logger.Errorf("[pipeline DoRun Extract] err: %v", err)
			continue
		}
		if event == nil {
			continue
		}

		// step2: transform chain
		event = t.transformer.Transform(event)
		if event == nil {
			continue
		}

		// step3: load
		if err = t.loader.Load(event); err != nil {
			logger.Errorf("[pipeline DoRun Load] err: %v", err)
		}
	}
}

func (t *Template) IsClosed() bool {
	return atomic.LoadUint32(&t.isRunning) != 1
}
