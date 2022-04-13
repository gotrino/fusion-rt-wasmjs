package dialog

import (
	"context"
	"embed"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"github.com/gotrino/fusion/spec/svg"
	"honnef.co/go/js/dom/v2"
	"html/template"
	"reflect"
)

//go:embed actiondialog.gohtml
var tpl embed.FS

type Button struct {
	ID           string
	Caption      string
	CallToAction bool
	Action       func()
}

type ActionDialog struct {
	SVG        template.HTML
	Title      string
	Message    string
	RawMessage template.HTML
	Buttons    []Button
}

func NewActionDialog(ctx context.Context, svg svg.SVG, title string, msg string, buttons ...Button) *ActionDialog {
	return &ActionDialog{
		SVG:     template.HTML(svg),
		Title:   title,
		Message: msg,
		Buttons: buttons,
	}
}

func (c *ActionDialog) Render(ctx context.Context) *tree.Component {
	for i := range c.Buttons {
		if c.Buttons[i].ID == "" {
			c.Buttons[i].ID = tree.NextID()
		}
	}

	// we are using tailwind flex-row-reverse for the button, thus we need to reverse the slice order
	reverse(c.Buttons)

	elem := tree.Template(ctx, tpl, c)

	for _, button := range c.Buttons {
		action := button.Action
		btnElem := elem.FindChild(button.ID)
		elem.Attach(btnElem.AddEventListener("click", false, func(event dom.Event) {
			if action != nil {
				action()
			}

			elem.DeleteSelf()
		}))
	}

	return elem
}

func reverse(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}
