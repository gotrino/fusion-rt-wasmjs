package wasmjs

import (
	"github.com/gotrino/fusion/runtime"
	"github.com/gotrino/fusion/spec/app"
)

type Runtime struct {
}

func (r *Runtime) Start(spec app.ApplicationComposer) error {
	//TODO implement me
	panic("implement me")
}

func init() {
	runtime.Register("wasm/js", func() (runtime.Runtime, error) {
		return &Runtime{}, nil
	})
}
