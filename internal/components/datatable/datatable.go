package datatable

import (
	"context"
	"embed"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"github.com/gotrino/fusion/spec/table"
	"honnef.co/go/js/dom/v2"
	"log"
	"syscall/js"
)

//go:embed *.gohtml
var tpl embed.FS

type DataTable struct {
	model table.DataTableStencil
}

const bla = `
{
  "header": [
    {
      "name": "id",
      "title": "ID",
      "size": 50,
      "sortable": true,
      "sortDir": "asc",
      "format": "number"
    },
    {
      "name": "name",
      "title": "Name",
      "sortable": true
    },
    {
      "name": "start",
      "title": "Start",
      "sortable": true,
      "size": 150,
      "format": "date",
      "formatMask": "dd-mm-yyyy"
    },
    {
      "name": "age",
      "title": "Age",
      "sortable": true,
      "size": 80
    },
    {
      "name": "salary",
      "title": "Salary",
      "sortable": true,
      "size": 150,
      "format": "money",
      "show": true
    }
  ],
  "data": [
    [
      1,
      "Aidan",
      "31-12-2017",
      63,
      "$7,843"
    ]
]
}
`

func NewDataTable(ctx context.Context, model table.DataTableStencil) *DataTable {
	return &DataTable{model: model}
}

func (c *DataTable) Render(ctx context.Context) *tree.Component {
	table := tree.Template(ctx, tpl, nil)

	// loadData see also https://github.com/olton/Metro-UI-CSS/issues/1404
	// see also https://metroui.org.ua/components-api.html
	// see also https://metroui.org.ua/table.html

	table.OnAppended(func() {
		log.Println("table appended")
		elem := dom.GetWindow().Document().GetElementByID("t1")
		_ = elem

		log.Println("--", js.Global().Call("$", "#t1").Call("data", "table"))
		//$('#t1').data('table').toggleInspector()

		metro := js.Global().Get("Metro")
		metro.Call("makePlugin", elem.Underlying(), "table") //  data-role="table" removed
		table := js.Global().Call("$", "#t1").Call("data", "table")
		table.Call("clear")
		table.Call("setData", bla)
		//	table.Call("draw")
		log.Println("are we cool yet?")
	})

	return table
}
