package fragments

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/gotrino/fusion/spec/app"
	"html/template"
	"reflect"
)

//go:embed *.gohtml
var files embed.FS
var Templates *template.Template

func init() {
	t, err := template.ParseFS(files, "*.gohtml")
	if err != nil {
		panic(err)
	}

	Templates = t
}

func Render(a app.Application, activities []app.Activity, content any) string {
	var pageModel PageModel
	pageModel.AppBarCaption = a.Title
	pageModel.Navigation = append(pageModel.Navigation, MenuItem{
		Icon:    "mif-apps",
		Link:    "/",
		Caption: "Dashboard",
		Active:  reflect.TypeOf(activities) == reflect.TypeOf(content),
	})

	var activeLauncher string
	switch t := content.(type) {
	case app.Activity:
		pageModel.AppBarCaption = t.Title
		for _, fragment := range t.Fragments {
			pageModel.Fragments = append(pageModel.Fragments, NewFragmentModel(fragment))
		}
		if ico, ok := t.Launcher.(app.Icon); ok {
			activeLauncher = ico.Link
		}
	default:
		pageModel.Fragments = append(pageModel.Fragments, NewFragmentModel(t))
	}

	for _, activity := range activities {
		if activity.Launcher != nil {
			switch t := activity.Launcher.(type) {
			case app.Icon:
				pageModel.Navigation = append(pageModel.Navigation, MenuItem{
					Icon:    t.Icon,
					Link:    makeLink(t.Link),
					Caption: t.Title,
					Active:  t.Link == activeLauncher && activeLauncher != "",
				})
			default:
				panic("unsupported type: " + Typename(t))
			}
		}
	}

	var buf bytes.Buffer
	if err := Templates.ExecuteTemplate(&buf, "page.gohtml", pageModel); err != nil {
		panic(fmt.Errorf("cannot apply template: %w", err))
	}

	return buf.String()
}

func makeLink(s string) string {
	return "#" + s
}

func Typename(a any) string {
	return reflect.TypeOf(a).String()
}

type Stencil interface {
	ToStencil() any
}
