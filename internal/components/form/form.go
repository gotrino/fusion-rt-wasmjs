package form

import (
	"context"
	"embed"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/i18n"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"github.com/gotrino/fusion/spec/form"
	"honnef.co/go/js/dom/v2"
	"log"
)

//go:embed form.gohtml
var tplForm embed.FS

type toStencilText interface {
	ToStencil() form.StencilText
}

type Form struct {
	ctx           context.Context
	Model         form.Form
	BtnSaveText   string
	BtnSaveID     string
	BtnDeleteText string
	BtnDeleteID   string
}

func NewForm(ctx context.Context, model form.Form) *Form {
	return &Form{ctx: ctx, Model: model}
}

func (c *Form) Render(ctx context.Context) *tree.Component {

	var myDefault any
	switch t := c.Model.Resource.(type) {
	case interface{ GetDefault() any }:
		myDefault = t.GetDefault()

	}

	c.BtnDeleteText = i18n.Text(ctx, "Delete")
	c.BtnSaveText = i18n.Text(ctx, "Save")

	c.BtnSaveID = tree.NextID()
	c.BtnDeleteID = tree.NextID()
	formElem := tree.Template(ctx, tplForm, c)

	for _, field := range c.Model.Fields {
		switch t := field.(type) {
		case toStencilText:
			textModel := t.ToStencil()
			text := NewText(textModel, myDefault)

			formElem.AppendChild("content", text.Render(ctx))
		}
	}

	if c.BtnSaveID != "" {
		btnElem := formElem.FindChild(c.BtnSaveID)
		formElem.Attach(btnElem.AddEventListener("click", true, func(event dom.Event) {
			c.onSave()
		}))
	}

	if c.BtnDeleteID != "" {
		btnElem := formElem.FindChild(c.BtnDeleteID)
		formElem.Attach(btnElem.AddEventListener("click", true, func(event dom.Event) {
			c.onDelete()
		}))
	}

	return formElem
}

func (c *Form) onDelete() {
	log.Println("yo should delete")
}

func (c *Form) onSave() {
	log.Println("yo should save")
}
