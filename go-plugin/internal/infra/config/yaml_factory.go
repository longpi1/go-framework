package config

import "gopkg.in/yaml.v2"

// YamlFactory Yaml配置工厂
type YamlFactory struct {
}

func NewYamlFactory() *YamlFactory {
	return &YamlFactory{}
}

func (y YamlFactory) CreateExtractorConfig() Extractor {
	return Extractor{loadConf: loadYaml}
}

func (y YamlFactory) CreateTransformerConfig() Transformer {
	return Transformer{loadConf: loadYaml}
}

func (y YamlFactory) CreateLoaderConfig() Loader {
	return Loader{loadConf: loadYaml}
}

func (y YamlFactory) CreatePipelineConfig() Pipeline {
	pipeline := Pipeline{}
	pipeline.loadConf = loadYaml
	return pipeline
}

func loadYaml(conf string, item interface{}) error {
	return yaml.Unmarshal([]byte(conf), item)
}
