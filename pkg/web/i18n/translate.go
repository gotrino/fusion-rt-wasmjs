package i18n

import "context"

func Text(ctx context.Context, key string, args ...any) string {
	return key
}

func Quantity(ctx context.Context, key string, quantity int, args ...any) string {
	return key
}
