package extractor

import (
	"nodr-datahub/internal/infra/config"
	"nodr-datahub/internal/infra/plugin"
	"nodr-datahub/logger"
)

// Plugin Extractor插件
type Plugin interface {
	plugin.Plugin
	Extract() (*plugin.Event, error)
}

// NewPlugin 输入插件工厂方法
func NewPlugin(config config.Extractor) (Plugin, error) {
	creator, ok := Extractors[config.PluginType]
	if !ok {
		logger.Errorf("extractor plugin not found: %v", config.PluginType)
		return nil, plugin.ErrUnknownPlugin
	}
	extractorPlugin := creator()
	extractorPlugin.SetContext(config.Ctx)
	return extractorPlugin, nil
}
