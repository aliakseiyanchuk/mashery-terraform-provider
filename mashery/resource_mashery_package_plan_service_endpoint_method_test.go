package mashery_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	mock "github.com/stretchr/testify/mock"
)

type MasheryPlanMethodMockClient struct {
	v3client.PluggableClient
	mock.Mock
}
