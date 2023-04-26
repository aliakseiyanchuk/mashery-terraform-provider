package mashery_test

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"terraform-provider-mashery/mashery"
	"terraform-provider-mashery/mashschema"
	"testing"
	"time"
)

func masheryTimeNow() *masherytypes.MasheryJSONTime {
	rv := masherytypes.MasheryJSONTime(time.Now())
	return &rv
}

// TestServiceWillCreateServiceWithoutRoles
// Perform a simplest possible test from the following HCL configuration
//
//	resource "mashery_service" "lspwd2-first" {
//	 name_prefix = "x002924"
//	 version = "0.0.1a"
//	}
func TestServiceWillCreateServiceWithoutRoles(t *testing.T) {

	tf := map[string]interface{}{
		mashschema.MashSvcName:    "mock-service",
		mashschema.MashSvcVersion: "0.0.1a",
	}

	d, dg := mashschema.ServiceMapper.TestResourceDataWith(tf)
	assert.Equal(t, 0, len(dg))

	mockCl := &MasheryPlanMethodMockClient{}
	mockCl.On("CreateService", mock.Anything, mock.Anything).Return(&masherytypes.Service{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:      "svcId",
			Name:    "mock-service",
			Created: masheryTimeNow(),
			Updated: masheryTimeNow(),
		},
		Version:      "0.0.1a",
		RobotsPolicy: "robots",
	}, nil).Once()

	dg = mashery.ServiceCreate(context.TODO(), d, mockCl)
	assert.Equal(t, 0, len(dg))

	mockCl.Mock.AssertExpectations(t)

	identRaw, dg := mashschema.ServiceMapper.V3Identity(d)
	ident := identRaw.(masherytypes.ServiceIdentifier)
	assert.Equal(t, 0, len(dg))
	assert.Equal(t, "svcId", ident.ServiceId)
}
