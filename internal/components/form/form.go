package form

import (
	"context"
	"embed"
	"github.com/gotrino/fusion-rt-wasmjs/internal/components/ace"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/i18n"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"github.com/gotrino/fusion/spec/app"
	"github.com/gotrino/fusion/spec/form"
	"honnef.co/go/js/dom/v2"
	"log"
	"syscall/js"
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
	BtnCancelText string
	BtnCancelID   string
	Error         error
	Entity        any
	repo          app.RepositoryImplStencil
	lastRendered  *tree.Component
}

func NewForm(ctx context.Context, model form.Form) *Form {
	return &Form{ctx: ctx, Model: model, repo: model.Repository.New(ctx)}
}

func (c *Form) Render(ctx context.Context) *tree.Component {
	c.Entity = c.Model.Repository.GetDefault()

	v, err := c.repo.Load(c.Model.ResourceID)
	if err == nil {
		c.Entity = v
	} else {
		if !app.NotFound(err) {
			c.Error = err
		}
	}

	if c.Model.CanDelete {
		c.BtnDeleteText = i18n.Text(ctx, "Delete")
		c.BtnDeleteID = tree.NextID()
	}

	if c.Model.CanWrite {
		c.BtnSaveText = i18n.Text(ctx, "Save")
		c.BtnSaveID = tree.NextID()
	}

	if c.Model.CanCancel {
		c.BtnCancelText = i18n.Text(ctx, "Cancel")
		c.BtnCancelID = tree.NextID()
	}

	formElem := tree.Template(ctx, tplForm, c)

	for _, field := range c.Model.Fields {
		switch t := field.(type) {
		case toStencilText:
			textModel := t.ToStencil()
			text := NewText(textModel, &c.Entity)

			formElem.AppendChild("content", text.Render(ctx))
		case form.CodeEditor:
			editor := ace.NewACE(ctx, t)
			formElem.AppendChild("content", editor.Render(ctx))
		}
	}

	if c.BtnSaveID != "" {
		btnElem := formElem.FindChild(c.BtnSaveID)
		formElem.Attach(btnElem.AddEventListener("click", true, func(event dom.Event) {
			log.Printf("%+v", c.Entity)
			go c.onSave()
		}))
	}

	if c.BtnDeleteID != "" {
		btnElem := formElem.FindChild(c.BtnDeleteID)
		formElem.Attach(btnElem.AddEventListener("click", true, func(event dom.Event) {
			go c.onDelete()
		}))
	}

	if c.BtnCancelID != "" {
		btnElem := formElem.FindChild(c.BtnCancelID)
		formElem.Attach(btnElem.AddEventListener("click", true, func(event dom.Event) {
			go c.onCancel()
		}))
	}

	c.lastRendered = formElem
	return formElem
}

func (c *Form) onDelete() {
	c.Error = c.repo.Delete(c.Model.ResourceID)
	c.invalidate()
}

func (c *Form) onSave() {
	c.Error = c.repo.Save(c.Entity)
	if c.Error != nil {
		c.invalidate()
		return
	}

	js.Global().Get("history").Call("back")
}

func (c *Form) onCancel() {
	js.Global().Get("history").Call("back")
}

func (c *Form) invalidate() {
	lastRendered := c.lastRendered
	lastRendered.ReplaceSelf(c.Render(c.ctx))
}
