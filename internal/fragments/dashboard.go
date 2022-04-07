package fragments

import (
	"github.com/gotrino/fusion/spec/app"
)

type DashboardRow []LauncherIconModel

type DashboardModel struct {
	Rows []DashboardRow
}

func NewDashboardModel(activities []app.Activity) DashboardModel {
	var launchers []app.Launcher
	for _, activity := range activities {
		if activity.Launcher != nil {
			launchers = append(launchers, activity.Launcher)
		}
	}

	return convertLauncher(launchers)
}

func convertLauncher(icons []app.Launcher) DashboardModel {
	var model DashboardModel
	if len(icons) == 0 {
		return model
	}

	const maxRowLen = 4
	for _, icon := range icons {
		if len(model.Rows) == 0 {
			model.Rows = append(model.Rows, []LauncherIconModel{})
		}

		if len(model.Rows[len(model.Rows)-1]) == maxRowLen {
			model.Rows = append(model.Rows, []LauncherIconModel{})
		}

		switch t := icon.(type) {
		case app.Icon:
			model.Rows[len(model.Rows)-1] = append(model.Rows[len(model.Rows)-1], NewLauncherIconModel(t))
		default:
			panic("not implemented: " + Typename(t))
		}

	}

	return model
}
