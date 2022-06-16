package collector

import (
	"errors"
	"fmt"
)

var Registry = &registry{
	metrics: make(map[string]initialize),
}

type initialize func() Collector

type registry struct {
	metrics map[string]initialize
}

func (r *registry) Add(name string, init initialize) {
	if _, exists := r.metrics[name]; exists {
		panic(fmt.Sprintf("already registered of %s", name))
	}

	r.metrics[name] = init
}

func (r *registry) Load(feats ...string) ([]Collector, error) {
	var cs []Collector

	for _, feat := range feats {
		if init, exists := r.metrics[feat]; exists {
			cs = append(cs, init())
		} else {
			return nil, errors.New(fmt.Sprintf("no collector for %s", feat))
		}
	}

	return cs, nil
}
