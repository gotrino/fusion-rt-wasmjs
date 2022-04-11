package form

import (
	"context"
	"embed"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"github.com/gotrino/fusion/spec/form"
)

//go:embed text.gohtml
var tplText embed.FS

//go:embed form.gohtml
var tplForm embed.FS

type toStencilText interface {
	ToStencil() form.StencilText
}

type Section struct {
}

type Form struct {
	ctx   context.Context
	Model form.Form
}

func NewForm(ctx context.Context, model form.Form) *Form {
	return &Form{ctx: ctx, Model: model}
}

func (c *Form) Render(ctx context.Context) *tree.Component {
	formElem := tree.Template(ctx, tplForm, c)

	for _, field := range c.Model.Fields {
		switch t := field.(type) {
		case toStencilText:
			textModel := t.ToStencil()
			textElem := tree.Template(ctx, tplText, textModel)
			formElem.AppendChild("content", textElem)
		}
	}

	return formElem
}
