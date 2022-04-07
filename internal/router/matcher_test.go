package router

import (
	"context"
	"github.com/gotrino/fusion/spec/app"
	"net/url"
	"testing"
)

type Details struct {
	Route string `route:"/devices/:id/details"`
	ID    string `path:"id"`
	Limit int    `query:"limit"`
	Admin bool   `query:"admin"`
}

func (d *Details) Compose(ctx context.Context) app.Activity {
	panic("implement me")
}

func TestRouteFrom(t *testing.T) {
	composer := &Details{}
	var boxed app.ActivityComposer
	boxed = composer
	matcher, err := NewMatcher(boxed)
	if err != nil {
		t.Fatal(err)
	}

	table := []struct {
		url     string
		matches bool
		id      string
		limit   int
		admin   bool
	}{
		{
			url:     "https://localhost:9090/devices/1234/details?limit=10&admin=true",
			matches: true,
			id:      "1234",
			limit:   10,
			admin:   true,
		},

		{
			url:     "devices/1/details",
			matches: true,
			id:      "1",
			limit:   0,
			admin:   false,
		},

		{
			url:     "devices/1/details/",
			matches: true,
			id:      "1",
			limit:   0,
			admin:   false,
		},

		{
			url:     "devices/1/details/nope",
			matches: false,
			id:      "",
			limit:   0,
			admin:   false,
		},
	}

	for _, s := range table {
		parsed, err := url.Parse(s.url)
		if err != nil {
			t.Fatal(err)
		}

		fits := matcher.Matches(parsed)
		if fits != s.matches {
			t.Fatalf("expected to match=%v but is not", fits)
		}

		if !fits {
			matcher.Reset()
		} else {
			if err := matcher.Apply(parsed); err != nil {
				t.Fatal(err)
			}
		}

		if composer.ID != s.id {
			t.Fatalf("expected id to be '%v' but is '%v'", s.id, composer.ID)
		}

		if composer.Admin != s.admin {
			t.Fatalf("expected admin to be %v", s.admin)
		}

		if composer.Limit != s.limit {
			t.Fatalf("expected limit to be %v", s.limit)
		}
	}

}
