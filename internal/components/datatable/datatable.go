package datatable

import (
	"context"
	"embed"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"github.com/gotrino/fusion/spec/table"
)

//go:embed *.gohtml
var tpl embed.FS

type DataTable struct {
	model table.DataTableStencil
}

func NewDataTable(ctx context.Context, model table.DataTableStencil) *DataTable {
	return &DataTable{model: model}
}

func (c *DataTable) Render(ctx context.Context) *tree.Component {
	table := tree.Template(ctx, tpl, nil)
	//c.model.Repository.(rest.Repository).New()
	return table
}
