package wasmjs

import (
	"context"
	"fmt"
	"github.com/gotrino/fusion-rt-wasmjs/internal/components/dialog"
	"github.com/gotrino/fusion-rt-wasmjs/internal/components/page"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/router"
	"github.com/gotrino/fusion-rt-wasmjs/pkg/web/tree"
	"github.com/gotrino/fusion/runtime"
	"github.com/gotrino/fusion/spec/app"
	"github.com/gotrino/fusion/spec/svg"
	"honnef.co/go/js/dom/v2"
	"html/template"
	"log"
	"net/url"
	"reflect"
	"runtime/debug"
	"strings"
)

type Runtime struct {
	ctx         context.Context
	app         app.Application
	matchers    []*router.Matcher[app.ActivityComposer]
	releasables []tree.Releasable
	component   *tree.Component
}

func NewRuntime(ctx context.Context) *Runtime {
	return &Runtime{ctx: ctx}
}

func (r *Runtime) Navigate(ac app.ActivityComposer) {
	matcher, err := router.NewMatcher[app.ActivityComposer](ac)
	if err != nil {
		panic(fmt.Errorf("invalid matcher: %w", err))
	}

	route := matcher.Render()
	route = router.LinkTo(r.ctx, route)
	dom.GetWindow().Location().SetHref(route)
}

func (r *Runtime) Start(spec app.ApplicationComposer) error {
	r.ctx = context.Background()
	r.ctx = app.WithContext[app.RT](r.ctx, app.RT{Delegate: r})

	r.app = spec.Compose(r.ctx)
	r.ctx = app.WithContext[app.Application](r.ctx, r.app)

	for _, activity := range r.app.Activities {
		matcher, err := router.NewMatcher(activity)
		if err != nil {
			return fmt.Errorf("invalid ActivityComposer '%s': %w ", reflect.TypeOf(activity), err)
		}

		r.matchers = append(r.matchers, matcher)

		log.Printf("registered %v for %v\n", reflect.TypeOf(activity), matcher.Pattern())
	}

	r.initRouting()
	select {}
}

func (r *Runtime) Spawn(f func()) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.Println("Recovered in f", r)
				stacktrace := string(debug.Stack())
				log.Println(string(debug.Stack()))
				dlg := dialog.NewActionDialog(r.ctx,
					svg.OutlineExclamation,
					"oops, something went wrong",
					"",
					dialog.Button{
						Caption: "not ok",
					},
				)
				dlg.RawMessage = `<div class="overscroll-auto"><pre>` + template.HTML(strings.ReplaceAll(stacktrace, "\n", "<br>")) + "</pre></div>"
				dlgElem := dlg.Render(r.ctx)
				if r.component != nil {
					r.component.Add(dlgElem)
				}
			}
		}()

		f()
	}()
}

func (r *Runtime) Refresh() {
	hash := dom.GetWindow().Location().Hash()
	r.dispatchRoute(hash)
}

func (r *Runtime) initRouting() {
	fun := dom.GetWindow().AddEventListener("hashchange", true, func(event dom.Event) {
		hash := dom.GetWindow().Location().Hash()
		r.dispatchRoute(hash)

	})

	r.dispatchRoute(dom.GetWindow().Location().Hash())
	r.releasables = append(r.releasables, fun)
}

func (r *Runtime) Release() {
	for _, releasable := range r.releasables {
		releasable.Release()
	}

	r.releasables = nil
}

func (r *Runtime) elem() dom.Element {
	elem := dom.GetWindow().Document().GetElementByID("app")
	if elem == nil {
		panic("element with id='app' not found")
	}

	return elem
}

func (r *Runtime) applyActivity(composer app.ActivityComposer) {
	r.component.Release()

	var acs []app.Activity
	idx := -1
	for i, cmp := range r.app.Activities {
		ac := cmp.Compose(r.ctx)
		acs = append(acs, ac)
		if cmp == composer { // composer must be pointer-type anyway
			idx = i
		}
	}

	state := runtime.State{
		Context:     r.ctx,
		Application: r.app,
		Activities:  acs,
		Active:      idx,
	}

	myPage := page.NewPage(r.ctx, state)
	r.component = myPage.Render(r.ctx)

	root := r.elem()
	root.SetTextContent("")
	root.AppendChild(r.component.Unwrap())
	r.component.Appended()
}

func (r *Runtime) applyHome() {
	r.applyActivity(nil)
}

func (r *Runtime) dispatchRoute(route string) {
	if strings.HasPrefix(route, "#") {
		route = route[1:]
	}

	u, err := url.Parse(route)
	if err != nil {
		log.Printf("cannot parse route as url '%s': %v\n", dom.GetWindow().Location().String(), err)
	}

	for _, matcher := range r.matchers {
		if matcher.Matches(u) {
			if err := matcher.Apply(u); err != nil {
				log.Printf("cannot apply matching match: %v\n", err)
			}
			log.Printf("route matches: %T\n", matcher.Unwrap())
			r.applyActivity(matcher.Unwrap())
			return
		}
	}

	log.Printf("no match found for route %s\n", route)
	r.applyHome()
}

func init() {
	runtime.Register("wasm/js", func() (runtime.Runtime, error) {
		return NewRuntime(context.Background()), nil
	})
}
