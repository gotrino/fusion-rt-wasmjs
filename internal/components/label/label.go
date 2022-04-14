package label

import (
	"context"
	"embed"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"strings"
)

//go:embed *.gohtml
var tpl embed.FS

type Label interface {
	GetFromModel() func(src any) string
	IsLabel() bool
	GetText() string
}

type View struct {
	label  Label
	entity *any
	Lines  []string
}

func NewView(ctx context.Context, label Label, entity *any) *View {
	return &View{
		label:  label,
		entity: entity,
	}
}

func (v *View) Render(ctx context.Context) *tree.Component {
	lines := v.label.GetText()
	if v.label.GetFromModel() != nil {
		lines += v.label.GetFromModel()(*v.entity)
	}

	v.Lines = strings.Split(lines, "\n")

	return tree.Template(ctx, tpl, v)
}
