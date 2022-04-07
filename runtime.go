package wasmjs

import (
	"context"
	"fmt"
	"github.com/gotrino/fusion-rt-wasmjs/internal/core"
	"github.com/gotrino/fusion-rt-wasmjs/internal/fragments"
	"github.com/gotrino/fusion-rt-wasmjs/internal/router"
	"github.com/gotrino/fusion/runtime"
	"github.com/gotrino/fusion/spec/app"
	"honnef.co/go/js/dom/v2"
	"log"
	"net/url"
	"reflect"
	"strings"
)

type Runtime struct {
	ctx         context.Context
	app         app.Application
	matchers    []*router.Matcher
	releasables []core.Releasable
}

func NewRuntime(ctx context.Context) *Runtime {
	return &Runtime{ctx: ctx}
}

func (r *Runtime) Start(spec app.ApplicationComposer) error {
	r.app = spec.Compose(r.ctx)
	for _, activity := range r.app.Activities {
		matcher, err := router.NewMatcher(activity)
		if err != nil {
			return fmt.Errorf("invalid ActivityComposer '%s': %w ", reflect.TypeOf(activity), err)
		}

		r.matchers = append(r.matchers, matcher)

		log.Printf("registered %v for %v\n", reflect.TypeOf(activity), matcher.Pattern())
	}

	log.Println("hello world")

	r.initRouting()
	select {}
}

func (r *Runtime) initRouting() {
	fun := dom.GetWindow().AddEventListener("hashchange", true, func(event dom.Event) {
		hash := dom.GetWindow().Location().Hash()
		r.dispatchRoute(hash)

	})

	r.dispatchRoute(dom.GetWindow().Location().Hash())
	r.releasables = append(r.releasables, fun)
}

func (r *Runtime) elem() dom.Element {
	return dom.GetWindow().Document().GetElementByID("app")
}

func (r *Runtime) applyActivity(composer app.ActivityComposer) {
	var acs []app.Activity
	for _, cmp := range r.app.Activities {
		ac := cmp.Compose(r.ctx)
		acs = append(acs, ac)
	}

	ac := composer.Compose(r.ctx)
	dom.GetWindow().Document().Underlying().Set("title", ac.Title)
	r.elem().SetInnerHTML("<h1>" + ac.Title + "</h1>")

	r.elem().SetInnerHTML(fragments.Render(r.app, acs, ac))
}

func (r *Runtime) applyHome() {
	var acs []app.Activity
	for _, cmp := range r.app.Activities {
		ac := cmp.Compose(r.ctx)
		acs = append(acs, ac)
	}

	dom.GetWindow().Document().Underlying().Set("title", r.app.Title)
	r.elem().SetInnerHTML(fragments.Render(r.app, acs, acs))
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
			log.Printf("route matches: %s\n", reflect.TypeOf(matcher.Composer()).String())
			r.applyActivity(matcher.Composer())
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
