package utils

import "context"

type ContextKey string

const (
	FirstNameKey ContextKey = "first_name"
	LastNameKey  ContextKey = "last_name"
	EmailKey     ContextKey = "email"
	RoleKey      ContextKey = "role"
	UserIdKey    ContextKey = "user_id"
)

func WithFirstName(ctx context.Context, firstName string) context.Context {
	return context.WithValue(ctx, FirstNameKey, firstName)
}

func WithLastName(ctx context.Context, lastName string) context.Context {
	return context.WithValue(ctx, LastNameKey, lastName)
}

func WithEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, EmailKey, email)
}

func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, RoleKey, role)
}

func WithUserId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, UserIdKey, userId)
}

func FirstNameFromCtx(ctx context.Context) (string, bool) {
	firstName, ok := ctx.Value(FirstNameKey).(string)
	return firstName, ok
}

func LastNameFromCtx(ctx context.Context) (string, bool) {
	lastName, ok := ctx.Value(LastNameKey).(string)
	return lastName, ok
}

func EmailFromCtx(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(EmailKey).(string)
	return email, ok
}

func RoleFromCtx(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(RoleKey).(string)
	return role, ok
}

func UserIdFromCtx(ctx context.Context) (string, bool) {
	userId, ok := ctx.Value(UserIdKey).(string)
	return userId, ok
}
