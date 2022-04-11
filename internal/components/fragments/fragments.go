package fragments

import (
	"context"
	"fmt"
	"github.com/gotrino/fusion-rt-wasmjs/internal/components/datatable"
	form2 "github.com/gotrino/fusion-rt-wasmjs/internal/components/form"
	"github.com/gotrino/fusion-rt-wasmjs/internal/components/launcher"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"github.com/gotrino/fusion/spec/app"
	"github.com/gotrino/fusion/spec/form"
	"github.com/gotrino/fusion/spec/table"
)

type Stencil interface {
	ToStencil() any
}

func Resolve(ctx context.Context, a any) tree.Renderer {
	switch t := a.(type) {
	case app.Icon:
		return launcher.NewLauncher(ctx, t)
	case Stencil:
		return Resolve(ctx, t.ToStencil())
	case table.DataTableStencil:
		return datatable.NewDataTable(ctx, t)
	case form.Form:
		return form2.NewForm(ctx, t)
	default:
		panic(fmt.Errorf("type not implemented: %T", t))
	}
}
