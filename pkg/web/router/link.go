package router

import (
	"context"
	"strings"
)

// LinkTo inspects the context and the configured router to decide how to encode the given route so that it fits
// to the routers settings.
func LinkTo(ctx context.Context, route string) string {
	if !strings.HasPrefix(route, "/") {
		route = "/" + route
	}
	return "#" + route
}
