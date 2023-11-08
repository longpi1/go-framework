package config

// Factory 配置抽象工厂接口
type Factory interface {
	CreateExtractorConfig() Extractor
	CreateTransformerConfig() Transformer
	CreateLoaderConfig() Loader
	CreatePipelineConfig() Pipeline
}
