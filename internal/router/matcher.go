package router

import (
	"fmt"
	"github.com/gotrino/fusion-rt-wasmjs/internal/reflect"
	"github.com/gotrino/fusion/spec/app"
	"net/url"
	reflect2 "reflect"
	"strconv"
	"strings"
)

type Matcher struct {
	composer    app.ActivityComposer // describes a concrete activity and has parameter fields
	pattern     []string
	paramParser map[string]func(v string) error
	queryParams []string
}

func NewMatcher(composer app.ActivityComposer) (*Matcher, error) {
	hasRoute := false
	cval := reflect2.ValueOf(composer).Elem()

	m := Matcher{
		composer:    composer,
		paramParser: map[string]func(v string) error{},
	}
	for k, field := range reflect.Fields(composer) {
		idx := k
		if !field.IsExported() {
			continue
		}

		v, ok := field.Tag.Lookup("route")
		if ok {
			if hasRoute {
				return nil, fmt.Errorf("multiple routes, must only have a single one")
			}

			hasRoute = true
			m.pattern = stripNames(v)
			continue
		}

		var param string
		if v, ok = field.Tag.Lookup("query"); ok {
			param = v
			m.queryParams = append(m.queryParams, v)
		}

		if v, ok = field.Tag.Lookup("path"); ok {
			param = v
		}

		if param != "" {
			if _, has := m.paramParser[param]; has {
				return nil, fmt.Errorf("path and query parameter must be unique: '%s'", v)
			}

			switch field.Type.Name() {
			case "int":
				m.paramParser[param] = func(v string) error {
					if v == "" {
						cval.Field(idx).Set(reflect2.ValueOf(0))
						return nil
					}

					i, err := strconv.Atoi(v)
					if err != nil {
						return fmt.Errorf("cannot parse parameter: %w", err)
					}

					cval.Field(idx).Set(reflect2.ValueOf(i))
					return nil
				}
			case "bool":
				m.paramParser[param] = func(v string) error {
					if v == "" {
						cval.Field(idx).Set(reflect2.ValueOf(false))
						return nil
					}

					i, err := strconv.ParseBool(v)
					if err != nil {
						return fmt.Errorf("cannot parse parameter: %w", err)
					}

					cval.Field(idx).Set(reflect2.ValueOf(i))
					return nil
				}
			case "string":
				m.paramParser[param] = func(v string) error {
					cval.Field(idx).Set(reflect2.ValueOf(v))
					return nil
				}
			default:
				return nil, fmt.Errorf("unsupported parameter type: %s", field.Type.Name())
			}
		}
	}

	if !hasRoute {
		return nil, fmt.Errorf("route has not been defined")
	}

	return &m, nil
}

func (r *Matcher) Pattern() string {
	return strings.Join(r.pattern, "/")
}

func (r *Matcher) Composer() app.ActivityComposer {
	return r.composer
}

// Reset populates all registered parameter bindings with their zero values.
func (r *Matcher) Reset() {
	for _, f := range r.paramParser {
		if err := f(""); err != nil {
			panic(err) // cannot happen, implementation failure
		}
	}
}

// Matches returns true if the given url path can be matched to the composer pattern.
func (r *Matcher) Matches(uri *url.URL) bool {
	pathNames := stripNames(uri.Path)
	if len(pathNames) != len(r.pattern) {
		return false
	}

	for i := 0; i < len(pathNames); i++ {
		isPathVar := strings.HasPrefix(r.pattern[i], ":")
		if isPathVar {
			continue
		}

		if pathNames[i] != r.pattern[i] {
			return false
		}
	}

	return true
}

// Apply must only be called if the url Matches. Reset is not required, because by definition all declared
// params are parsed anyway and if not defined the default value is set (empty string).
func (r *Matcher) Apply(uri *url.URL) error {
	pathNames := stripNames(uri.Path)
	if len(pathNames) != len(r.pattern) {
		panic("illegal state") // programming error
	}

	// match and parse path params
	for i := 0; i < len(pathNames); i++ {
		isPathVar := strings.HasPrefix(r.pattern[i], ":")
		if isPathVar {
			varName := r.pattern[i][1:]
			varValue := pathNames[i]
			f, ok := r.paramParser[varName]
			if !ok {
				panic("illegal state") // programming error
			}

			if err := f(varValue); err != nil {
				return fmt.Errorf("cannot parse %s: %w", varName, err)
			}
		} else {
			if pathNames[i] != r.pattern[i] {
				panic("illegal state") // programming error
			}
		}
	}

	// match and parse query params
	for _, param := range r.queryParams {
		v := uri.Query().Get(param)
		f, ok := r.paramParser[param]
		if !ok {
			panic("illegal state") // programming error
		}

		if err := f(v); err != nil {
			return fmt.Errorf("cannot parse query value %s: %w", param, err)
		}
	}

	return nil
}

func stripNames(v string) []string {
	if strings.HasPrefix(v, "/") {
		v = v[1:]
	}

	if strings.HasSuffix(v, "/") {
		v = v[:len(v)-1]
	}

	return strings.Split(v, "/")
}
