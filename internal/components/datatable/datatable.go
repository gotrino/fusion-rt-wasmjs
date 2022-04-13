package datatable

import (
	"context"
	"embed"
	"github.com/gotrino/fusion-rt-wasmjs/internal/components/dialog"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/i18n"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"github.com/gotrino/fusion/runtime/rest"
	"github.com/gotrino/fusion/spec/app"
	"github.com/gotrino/fusion/spec/svg"
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
//  * type == "data-id-svg"
//    Values[0] contains the data-id and Values[1] contains the SVG
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

type dataIdEntity struct {
	dataId string
	entity any
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

	// add empty column for actions
	if c.model.Deletable {
		c.Columns = append(c.Columns, Column{})
	}

	var deletables []dataIdEntity

	for _, entity := range entities {
		id := tree.NextID()
		row := Row{
			ID:    id,
			Hover: c.model.OnClick != nil,
		}

		for i := 0; i < len(c.model.Columns); i++ {
			cell := c.model.OnRender(ctx, entity, i)
			row.Cells = append(row.Cells, Cell{
				Values: cell.Values,
				Type:   cell.RenderHint,
			})
		}

		if c.model.Deletable {
			treeId := tree.NextID()
			row.Cells = append(row.Cells, Cell{
				Values: []string{treeId, string(svg.OutlineTrash)},
				Type:   "data-id-svg",
			})

			deletables = append(deletables, dataIdEntity{
				dataId: treeId,
				entity: entity,
			})
		}

		c.Rows = append(c.Rows, row)
	}

	elem := tree.Template(ctx, tpl, c)

	if c.model.OnClick != nil {
		for i, row := range c.Rows {
			rowNum := i
			rowElem := elem.FindChild(row.ID)
			elem.Attach(rowElem.AddEventListener("click", false, func(event dom.Event) {
				entity := entities[rowNum]
				log.Printf("clicked %v -> %v\n", i, entity)
				c.model.OnClick(ctx, entity)
			}))
		}
	}

	if c.model.Deletable {
		for _, deletable := range deletables {
			entity := deletable.entity
			trashBtn := elem.FindChild(deletable.dataId)
			trashBtn.Class().Add("hover:bg-blue-100")
			elem.Attach(trashBtn.AddEventListener("click", false, func(event dom.Event) {
				event.StopPropagation()

				dlgElem := dialog.NewActionDialog(
					ctx,
					svg.OutlineExclamation,
					i18n.Text(ctx, "Bestätigung löschen"),
					i18n.Text(ctx, "Sind Sie sicher, dass Sie den Eintrag löschen möchten?"),
					dialog.Button{
						Caption: i18n.Text(ctx, "cancel"),
					},
					dialog.Button{
						Caption:      "delete",
						CallToAction: true,
						Action: func() {
							rt := app.FromContext[app.RT](ctx)
							log.Printf("%+v delete", entity)
							rt.Delegate.Spawn(func() {
								id, err := rest.GetID(entity)
								if err != nil {
									panic(err)
								}
								if err := repo.Delete(id); err != nil {
									panic(err)
								}

								rt.Delegate.Refresh()
							})
						},
					},
				).Render(ctx)
				elem.Add(dlgElem)
			}))
		}
	}

	return elem
}
