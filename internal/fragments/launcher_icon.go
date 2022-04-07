package fragments

import (
	"github.com/gotrino/fusion/spec/app"
)

type LauncherIconModel struct {
	Title    string
	SubTitle string
	Icon     string
	Link     string
	More     string
}

func NewLauncherIconModel(t app.Icon) LauncherIconModel {
	return LauncherIconModel{
		Title:    t.Title,
		SubTitle: t.Hint,
		Icon:     t.Icon,
		Link:     "#" + t.Link,
		More:     "Mehr",
	}
}
