package config

import (
	"go-plugin/internal/infra/plugin"
)

// Config 插件配置抽象接口
type Config interface {

	// Load 将配置文件加载为插件配置项
	Load(conf string) error
}

// item 配置项结构体, 每个配置项均包含固定属性Name、PluginType、Ctx
type item struct {
	// 插件名
	Name string `json:"name" yaml:"name"`

	// 插件类型(Extractor\Transformer\Loader\Pipeline)
	PluginType string `json:"type" yaml:"type"`

	// 插件上下文, map结构, 用于保存插件执行所需参数
	Ctx plugin.Context `json:"context" yaml:"context"`

	//loadConf方法用于实现多态
	loadConf func(conf string, item interface{}) error
}

type Extractor item

func (e *Extractor) Load(conf string) error {
	return e.loadConf(conf, e)
}

type Transformer item

func (f *Transformer) Load(conf string) error {
	return f.loadConf(conf, f)
}

type Loader item

func (l *Loader) Load(conf string) error {
	return l.loadConf(conf, l)
}

type Pipeline struct {
	item         `yaml:",inline"`
	Extractor    Extractor     `json:"extractor" yaml:"extractor"`
	Transformers []Transformer `json:"transformers" yaml:"transformers,flow"`
	Loader       Loader        `json:"loader" yaml:"loader"`
}

func (p *Pipeline) Load(conf string) error {
	return p.loadConf(conf, p)
}
