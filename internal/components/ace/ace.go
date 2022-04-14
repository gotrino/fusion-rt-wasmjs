package ace

import (
	"context"
	"embed"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"github.com/gotrino/fusion/spec/form"
	"honnef.co/go/js/dom/v2"
	"html/template"
	"log"
	"syscall/js"
)

//go:embed *.gohtml
var tpl embed.FS

type ACE struct {
	ID      string
	Model   form.CodeEditor
	RawText template.HTML
}

func NewACE(ctx context.Context, model form.CodeEditor) *ACE {
	return &ACE{Model: model}
}

func (c *ACE) Render(ctx context.Context) *tree.Component {
	c.ID = tree.NextID()
	c.RawText = "// check that\nfunc main(){\nprintln(`hello world`)\n}\n"
	aceElem := tree.Template(ctx, tpl, c)

	aceElem.OnAppended(func() {
		log.Println("ace was append to tree")
		script := dom.GetWindow().Document().CreateElement("script")
		aceElem.Attach(script.AddEventListener("load", false, func(event dom.Event) {
			log.Println("ace loaded")
			editor := js.Global().Get("ace").Call("edit", c.ID)
			editor.Call("setTheme", "ace/theme/monokai")
			editor.Get("session").Call("setMode", "ace/mode/golang")
			editor.Call("setReadOnly", c.Model.ReadOnly)
		}))
		script.SetAttribute("type", "text/javascript")
		script.SetAttribute("charset", "utf-8")
		script.SetAttribute("src", "ace/ace.js")

		dom.GetWindow().Document().GetElementsByTagName("head")[0].AppendChild(script)
	})

	return aceElem
}
