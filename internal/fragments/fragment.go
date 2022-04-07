package fragments

import (
	"github.com/gotrino/fusion/spec/app"
	"github.com/gotrino/fusion/spec/table"
)

type FragmentModel struct {
	Type  string
	Model any
}

func NewFragmentModel(a any) FragmentModel {
	var m FragmentModel
	switch t := a.(type) {
	case Stencil:
		return NewFragmentModel(t.ToStencil())
	case app.Icon:
		m.Type = Typename(t)
		m.Model = NewLauncherIconModel(t)
	case []app.Activity:
		m.Type = Typename(t)
		m.Model = NewDashboardModel(t)
	case table.DataTableStencil:
		m.Type = Typename(t)
		m.Model = NewDataTableModel(t)
	default:
		panic("unsupported fragment type: " + Typename(a))
	}

	return m
}

func (f FragmentModel) IsLauncherIcon() bool {
	return f.Type == Typename(app.Icon{})
}

func (f FragmentModel) IsDashboardModel() bool {
	return f.Type == Typename([]app.Activity{})
}

func (f FragmentModel) IsDataTable() bool {
	return f.Type == Typename(table.DataTableStencil{})
}
