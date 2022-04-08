package page

import (
	"context"
	"embed"
	"fmt"
	"github.com/gotrino/fusion-rt-wasmjs/internal/components/dashboard"
	"github.com/gotrino/fusion-rt-wasmjs/internal/components/fragments"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/i18n"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/router"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"github.com/gotrino/fusion/runtime"
	"github.com/gotrino/fusion/spec/app"
	"honnef.co/go/js/dom/v2"
)

//go:embed *.gohtml
var tpl embed.FS

type MenuItem struct {
	Icon    string
	Link    string
	Caption string
	Active  bool
}

type Model struct {
	Navigation     []MenuItem
	AppBar         []MenuItem
	AppBarTitle    string
	AppBarSubTitle string
	Hamburger      []MenuItem
}

type Page struct {
	state runtime.State
	model Model
}

func NewPage(ctx context.Context, state runtime.State) *Page {
	var m Model
	m.AppBarTitle = state.Application.Title
	m.AppBarSubTitle = i18n.Text(ctx, "Dashboard")
	m.Navigation = append(m.Navigation, MenuItem{
		Icon:    "mif-apps",
		Link:    router.LinkTo(ctx, "/"),
		Caption: "Dashboard",
		Active:  state.Active < 0,
	})

	if state.Active >= 0 {
		m.AppBarSubTitle = state.Activities[state.Active].Title
	}

	for i, activity := range state.Activities {
		if activity.Launcher != nil {
			switch t := activity.Launcher.(type) {
			case app.Icon:
				m.Navigation = append(m.Navigation, MenuItem{
					Icon:    t.Icon,
					Link:    router.LinkTo(ctx, t.Link),
					Caption: t.Title,
					Active:  i == state.Active,
				})
			default:
				panic(fmt.Errorf("unsupported type: %T", t))
			}
		}
	}

	return &Page{model: m, state: state}
}

func (c *Page) Render(ctx context.Context) *tree.Component {
	dom.GetWindow().Document().Underlying().Set("title", c.model.AppBarTitle+" | "+c.model.AppBarSubTitle)
	cmp := tree.Template(ctx, tpl, c.model)

	if c.state.Active < 0 {
		db := dashboard.NewDashboard(ctx, c.state.Activities).Render(ctx)
		cmp.Replace("content", db)
	} else {
		activity := c.state.Activities[c.state.Active]
		frags := tree.Elem("div")
		for _, fragment := range activity.Fragments {
			renderer := fragments.Resolve(ctx, fragment)
			frags.Add(renderer.Render(ctx))
		}

		cmp.Replace("content", frags)
	}

	return cmp
}
