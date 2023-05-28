package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	v3client.Client
	mock.Mock
}

func (mc *MockClient) ListRolesFiltered(ctx context.Context, params map[string]string, fields []string) ([]masherytypes.Role, error) {
	args := mc.Called(ctx, params, fields)
	return args.Get(0).([]masherytypes.Role), args.Error(1)
}

func (mc *MockClient) ListEmailTemplateSetsFiltered(ctx context.Context, params map[string]string, fields []string) ([]masherytypes.EmailTemplateSet, error) {
	args := mc.Called(ctx, params, fields)
	return args.Get(0).([]masherytypes.EmailTemplateSet), args.Error(1)
}
