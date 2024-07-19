package tests

import (
	"context"
	"testing"

	"github.com/modylegi/service/internal/domain/repository/mocks"
	"github.com/modylegi/service/internal/repository"
	repositoryImpl "github.com/modylegi/service/internal/repository"
)

func TestBlockRepository_FindBlockID(t *testing.T) {

	type args struct {
		ctx       context.Context
		condition *repository.Condition
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "base test",
			args: args{
				ctx:       context.Background(),
				condition: &repositoryImpl.Condition{UserID: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blockRepository := mocks.NewBlockRepository(t)
			blockRepository.
				On("FindBlockID", tt.args.ctx, tt.args.condition).
				Once().
				Return(nil, nil)

			if err := blockRepository.FindBlockID(tt.args.ctx, tt.args.condition); (err != nil) != tt.wantErr {
				t.Errorf("FindBlockID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
