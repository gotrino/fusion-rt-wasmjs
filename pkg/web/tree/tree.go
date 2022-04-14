package tree

import (
	"bytes"
	"context"
	"honnef.co/go/js/dom/v2"
	"html/template"
	"io/fs"
	"strconv"
	"sync"
	"sync/atomic"
)

var tplCache = map[fs.FS]*template.Template{}
var mutex sync.Mutex
var id int32

func NextID() string {
	r := atomic.AddInt32(&id, 1)
	return "gtr-" + strconv.Itoa(int(r))
}

type Releasable interface {
	Release()
}

// A Renderer creates a dom.Element ready to display and interact with.
type Renderer interface {
	Render(ctx context.Context) *Component
}

// A Component consists of a dom.Element and potentially some registered resources which can be released.
type Component struct {
	wrapped   dom.Element
	resources []Releasable

	appendListeners []func()
	wasAppended     bool
}

func Wrap(elem dom.Element) *Component {
	return &Component{
		wrapped: elem,
	}
}

func Elem(name string) *Component {
	return Wrap(dom.GetWindow().Document().CreateElement(name))
}

func parse(fs fs.FS) *template.Template {
	mutex.Lock()
	defer mutex.Unlock()
	tpl := tplCache[fs]
	if tpl == nil {

		t, err := template.ParseFS(fs, "*.gohtml")
		if err != nil {
			panic(err) // programming error, templates are usually compile-time only
		}
		tplCache[fs] = t
		tpl = t
	}

	return tpl
}

func Template(ctx context.Context, fs fs.FS, state any) *Component {
	tpl := parse(fs)

	// allow concurrent execution, even though this does not matter today
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, state); err != nil {
		panic(err) // programming error in template, in 99% of the cases this is not recoverable anyway
	}

	elem := dom.GetWindow().Document().CreateElement("div")
	elem.SetInnerHTML(buf.String())
	return Wrap(elem.FirstChild().(dom.Element))
}

// Replace replaces the element denoted by the data-id entirely using the other components wrapped element.
// It also attaches other to this, so that it gets released.
func (c *Component) Replace(dataId string, other *Component) *Component {
	elem := c.wrapped.QuerySelector(`[data-id="` + dataId + `"]`)
	if elem == nil {
		panic("there is no such element '" + dataId + "'")
	}

	p := elem.ParentElement()
	if p == nil {
		panic("no parent elem for data node")
	}
	p.ReplaceChild(other.wrapped, elem)

	c.Attach(other)

	return c
}

func (c *Component) DeleteSelf() {
	p := c.wrapped.ParentElement()
	if p == nil {
		panic("no parent elem for data node")
	}

	p.RemoveChild(c.wrapped)
	c.Release()
}

func (c *Component) ReplaceSelf(other *Component) *Component {
	p := c.wrapped.ParentElement()
	if p == nil {
		panic("no parent elem for data node")
	}
	p.ReplaceChild(other.wrapped, c.wrapped)

	// clean up ourself
	c.Release()

	// incorporate everything else
	c.wrapped = other.wrapped
	c.resources = other.resources
	c.appendListeners = other.appendListeners

	return c
}

func (c *Component) AppendChild(dataId string, other *Component) *Component {
	elem := c.wrapped.QuerySelector(`[data-id="` + dataId + `"]`)
	if elem == nil {
		panic("there is no such element '" + dataId + "'")
	}

	elem.AppendChild(other.wrapped)

	c.Attach(other)

	return c
}

func (c *Component) FindChild(dataId string) dom.Element {
	return c.wrapped.QuerySelector(`[data-id="` + dataId + `"]`)
}

func (c *Component) Add(other *Component) *Component {
	c.wrapped.AppendChild(other.wrapped)
	c.Attach(other)

	return c
}

// Attach just adds the releasable to the receivers lifetime.
func (c *Component) Attach(r Releasable) {
	c.resources = append(c.resources, r)
	if c.wasAppended {
		if cmp, ok := r.(*Component); ok {
			cmp.Appended()
		}
	}
}

func (c *Component) Unwrap() dom.Element {
	return c.wrapped
}

func (c *Component) OnAppended(f func()) {
	c.appendListeners = append(c.appendListeners, f)
}

func (c *Component) Appended() {
	for _, listener := range c.appendListeners {
		listener()
	}

	for _, resource := range c.resources {
		if other, ok := resource.(*Component); ok {
			other.Appended()
		}
	}

	c.wasAppended = true
}

func (c *Component) Release() {
	if c == nil {
		return
	}

	for _, resource := range c.resources {
		resource.Release()
	}

	c.resources = nil
	c.wrapped = nil
	c.appendListeners = nil
}
