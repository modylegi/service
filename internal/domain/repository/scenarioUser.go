package repository

import (
	"context"
)

type ScenarioUserRepository interface {
	GetScenarioByUserID(context.Context, Condition) error
}
