package snackbar

import (
	"context"
	"embed"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"github.com/gotrino/fusion/spec/svg"
	"honnef.co/go/js/dom/v2"
	"html/template"
	"time"
)

//go:embed *.gohtml
var tpl embed.FS

type Snackbar struct {
	ID  string
	Msg string
	SVG template.HTML
}

func NewSnackbar(ctx context.Context) *Snackbar {
	return &Snackbar{ID: tree.NextID()}
}

func (c *Snackbar) Render(ctx context.Context) *tree.Component {
	return tree.Template(ctx, tpl, c)
}

func ShowToast(ctx context.Context, svg svg.SVG, msg string) {
	elem := dom.GetWindow().Document().GetElementByID("snackbar-panel")
	bar := NewSnackbar(ctx)
	bar.SVG = template.HTML(svg)
	bar.Msg = msg
	snackElem := bar.Render(ctx)
	snackElem.Unwrap().Class().Add("translate-y-full")
	go func() {
		time.Sleep(10 * time.Millisecond)
		snackElem.Unwrap().Class().Remove("translate-y-full")
	}()
	go func() {
		time.Sleep(2 * time.Second)
		snackElem.Unwrap().Class().Add("opacity-0")
		time.Sleep(500 * time.Millisecond)
		snackElem.DeleteSelf()
	}()
	elem.AppendChild(snackElem.Unwrap())
}
