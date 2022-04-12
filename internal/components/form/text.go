package form

import (
	"context"
	"embed"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"github.com/gotrino/fusion/spec/form"
	"honnef.co/go/js/dom/v2"
)

//go:embed text.gohtml
var tplText embed.FS

type Text struct {
	ID          string
	Label       string
	Description string
	Disabled    bool
	Placeholder string
	Lines       int
	Error       error
	Type        string
	def         *any
	stencil     form.StencilText
	Value       string
}

func NewText(textModel form.StencilText, def *any) *Text {
	modelValue := textModel.FromModel(*def)
	return &Text{
		def:         def,
		stencil:     textModel,
		ID:          tree.NextID(),
		Type:        "text",
		Label:       textModel.Label,
		Description: textModel.Description,
		Disabled:    textModel.Disabled,
		Placeholder: textModel.Placeholder,
		Lines:       textModel.Lines,
		Value:       modelValue,
	}
}

func (c *Text) Render(ctx context.Context) *tree.Component {
	textElem := tree.Template(ctx, tplText, c)
	inputElem := textElem.FindChild(c.ID)
	textElem.Attach(inputElem.AddEventListener("input", true, func(event dom.Event) {
		value := inputElem.Underlying().Get("value").String()
		c.Value = value
		*c.def, c.Error = c.stencil.ToModel(value, *c.def)

		// this is a bit ugly, because we do not have a delta-algorithm, so we have to fix
		priorInputElem := textElem.FindChild(c.ID)
		cursorPos := priorInputElem.Underlying().Get("selectionStart").Int()
		textElem.ReplaceSelf(c.Render(ctx))
		inputElem := textElem.FindChild(c.ID)
		inputElem.Underlying().Call("focus")
		inputElem.Underlying().Call("setSelectionRange", cursorPos, cursorPos)
	}))

	return textElem
}
