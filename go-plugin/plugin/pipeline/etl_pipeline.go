package pipeline

import (
	"go-pulgin/internal/infra/pipeline"
	"go-pulgin/internal/infra/plugin"
)

// SimplePipeline 简单Pipeline实现，每次运行时新启一个goroutine
type SimplePipeline struct {
	pipeline.Template
}

func (s *SimplePipeline) SetContext(ctx plugin.Context) {
	s.Run = func() {
		go func() {
			s.DoRun()
		}()
	}
}

func init() {
	pipeline.Add("simple", func() pipeline.Plugin {
		return &SimplePipeline{}
	})
}
