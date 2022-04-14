package ace

import (
	"context"
	"embed"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"honnef.co/go/js/dom/v2"
	"html/template"
	"log"
	"strings"
	"syscall/js"
)

//go:embed *.gohtml
var tpl embed.FS

type CodeEditor interface {
	GetToModel() func(src string, dst any) (any, error)
	GetFromModel() func(src any) string
	GetLang() string
	IsReadOnly() bool
}

type ACE struct {
	ID      string
	Model   CodeEditor
	RawText template.HTML
	Error   error
	entity  *any
}

func NewACE(ctx context.Context, model CodeEditor, entity *any) *ACE {
	return &ACE{Model: model, entity: entity}
}

func langToMode(lang string) string {
	switch strings.ToLower(lang) {
	case "golang":
		fallthrough
	case "go":
		return "ace/mode/golang"
	case "json":
		return "ace/mode/golang"
	default:
		return "ace/mode/text"
	}
}

func (c *ACE) Render(ctx context.Context) *tree.Component {
	c.ID = tree.NextID()
	if c.Model.GetFromModel() != nil {
		c.RawText = template.HTML(c.Model.GetFromModel()(*c.entity))
	} else {
		log.Println("ace: FromModel is not implemented")
	}
	aceElem := tree.Template(ctx, tpl, c)

	aceElem.OnAppended(func() {
		log.Println("ace was append to tree")
		script := dom.GetWindow().Document().CreateElement("script")
		aceElem.Attach(script.AddEventListener("load", false, func(event dom.Event) {
			log.Println("ace loaded")
			editor := js.Global().Get("ace").Call("edit", c.ID)
			editor.Call("setReadOnly", c.Model.IsReadOnly())
			editor.Call("setTheme", "ace/theme/monokai")
			session := editor.Get("session")
			session.Call("setMode", langToMode(c.Model.GetLang()))
			onChangeFun := js.FuncOf(func(this js.Value, args []js.Value) any {
				text := editor.Call("getValue").String()
				if c.Model.GetToModel() != nil {
					*c.entity, c.Error = c.Model.GetToModel()(text, *c.entity)
				} else {
					log.Println("ace: ToModel is not implemented")
				}
				return nil
			})
			aceElem.Attach(onChangeFun)
			session.Call("on", "change", onChangeFun)

		}))
		script.SetAttribute("type", "text/javascript")
		script.SetAttribute("charset", "utf-8")
		script.SetAttribute("src", "ace/ace.js")

		dom.GetWindow().Document().GetElementsByTagName("head")[0].AppendChild(script)
	})

	return aceElem
}
