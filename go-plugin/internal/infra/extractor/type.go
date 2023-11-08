package extractor

type Creator func() Plugin

var Extractors = map[string]Creator{}

func Add(name string, creator Creator) {
	Extractors[name] = creator
}
