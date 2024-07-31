package repository

import (
	"context"
)

type UserRepository interface {
	FindUserID(context.Context, Condition) error
	FindScenario(context.Context, Condition) error
	Create(context.Context, *User) error
	GetByUsername(context.Context, string) (*User, error)
}
