package pipeline

type Creator func() Plugin

var Pipelines = map[string]Creator{}

func Add(name string, creator Creator) {
	Pipelines[name] = creator
}
