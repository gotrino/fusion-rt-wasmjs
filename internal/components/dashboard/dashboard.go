package dashboard

import (
	"context"
	"github.com/gotrino/fusion-rt-wasmjs/internal/components/launcher"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"github.com/gotrino/fusion/spec/app"
)

type Dashboard struct {
	maxPerRow  int
	activities []app.Activity
}

func NewDashboard(ctx context.Context, activities []app.Activity) *Dashboard {
	return &Dashboard{activities: activities, maxPerRow: 4}
}

func (c *Dashboard) Render(ctx context.Context) *tree.Component {
	rowCount := 0
	root := tree.Elem("div")
	row := tree.Elem("div")
	row.Unwrap().Class().Add("row")

	for _, activity := range c.activities {
		if activity.Launcher != nil {
			ico := launcher.NewLauncher(ctx, activity.Launcher).Render(ctx)
			row.Add(ico)
			rowCount++

			if rowCount == c.maxPerRow {
				row = tree.Elem("div")
				row.Unwrap().Class().Add("row")
				rowCount = 0
				root.Add(row)
			}
		}
	}

	if rowCount != 0 {
		root.Add(row)
	}

	return root
}
