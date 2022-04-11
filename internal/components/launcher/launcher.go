package launcher

import (
	"context"
	"embed"
	"fmt"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/i18n"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/router"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"github.com/gotrino/fusion/spec/app"
	"github.com/gotrino/fusion/spec/svg"
)

//go:embed *.gohtml
var tpl embed.FS

type Model struct {
	Title    string
	SubTitle string
	Icon     svg.SVG
	Link     string
	More     string
}

type Launcher struct {
	model Model
}

func NewLauncher(ctx context.Context, l app.Launcher) *Launcher {
	switch t := l.(type) {
	case app.Icon:
		m := Model{
			Title:    t.Title,
			SubTitle: t.Hint,
			Icon:     t.Icon,
			Link:     router.LinkTo(ctx, t.Link),
			More:     i18n.Text(ctx, "more"),
		}

		return &Launcher{model: m}
	default:
		panic(fmt.Errorf("unsupported type '%T'", t))
	}
}

func (c *Launcher) Render(ctx context.Context) *tree.Component {
	return tree.Template(ctx, tpl, c.model)
}
