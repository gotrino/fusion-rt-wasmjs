package datatable

import (
	"context"
	"embed"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"github.com/gotrino/fusion/spec/table"
	"honnef.co/go/js/dom/v2"
	"html/template"
	"log"
	"strings"
)

//go:embed *.gohtml
var tpl embed.FS

type Column struct {
	Title string
}

type Row struct {
	ID    string
	Hover bool
	Cells []Cell
}

// Cell represents the following ugly state machine:
//  * type == "text-1"
//    only Values[0] is shown as is
//  * type == "text-2"
//    Values[0] is shown in a first line and Values[1] in a second line
//  * type == "img-text-2"
//    Values[0] is shown in a first line and Values[1] in a second line and Values[2] is the image source
//  * type == "svg-text-2"
//    Values[0] is shown in a first line and Values[1] in a second line and Values[2] hopefully contains an svg
//    image which will not be sanitized. If this content is checked by the developer itself or if it comes
//    from a third-party source (like a db) this is open to inject any client-side javascripts.
//  * type == "link"
//    Values[0] contains the link and Values[1] contains the link-name
//  * type == "badges"
//    Values contains each badge in the combo: <color>:<text> where color is one of red|blue|green
type Cell struct {
	Values []string
	Type   string
}

func (c *Cell) String() string {
	return strings.Join(c.Values, ", ")
}

func (c *Cell) UnsafeHTML(idx int) template.HTML {
	return template.HTML(c.Values[idx])
}

func (c *Cell) Badges() []Badge {
	res := make([]Badge, 0, len(c.Values))
	for _, value := range c.Values {
		tokens := strings.Split(value, ":")
		switch len(tokens) {
		case 2:
			res = append(res, Badge{
				Color: strings.TrimSpace(tokens[0]),
				Text:  strings.TrimSpace(tokens[1]),
			})
		default:
			log.Printf("invalid badge format: %s\n", value)
		}
	}

	return res
}

type Badge struct {
	Color string
	Text  string
}

type DataTable struct {
	model   table.DataTableStencil
	ErrMsg  string
	Columns []Column
	Rows    []Row
	ctx     context.Context
}

func NewDataTable(ctx context.Context, model table.DataTableStencil) *DataTable {
	return &DataTable{ctx: ctx, model: model}
}

func (c *DataTable) Render(ctx context.Context) *tree.Component {
	c.Rows = nil
	c.Columns = nil

	repo := c.model.Repository.New(ctx)
	entities, err := repo.List()
	if err != nil {
		c.ErrMsg = err.Error()
	}

	for _, column := range c.model.Columns {
		c.Columns = append(c.Columns, Column{Title: column.Name})
	}

	for _, entity := range entities {
		id := tree.NextID()
		row := Row{
			ID:    id,
			Hover: c.model.OnClick != nil,
		}

		for i := 0; i < len(c.Columns); i++ {
			cell := c.model.OnRender(ctx, entity, i)
			row.Cells = append(row.Cells, Cell{
				Values: cell.Values,
				Type:   cell.RenderHint,
			})
		}

		c.Rows = append(c.Rows, row)
	}

	elem := tree.Template(ctx, tpl, c)

	if c.model.OnClick != nil {
		for i, row := range c.Rows {
			dom.GetWindow().Document().GetElementByID(row.ID)
			elem.Attach(elem.Unwrap().AddEventListener("click", true, func(event dom.Event) {
				entity := entities[i]
				log.Printf("clicked %v -> %v\n", i, entity)
				c.model.OnClick(ctx, entity)
			}))
		}
	}

	return elem
}
