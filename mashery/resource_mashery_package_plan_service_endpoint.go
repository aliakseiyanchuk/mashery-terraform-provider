package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	mashschema "terraform-provider-mashery/mashschema"
)

var PackagePlanServiceEndpointResource *ResourceTemplate

func init() {
	PackagePlanServiceEndpointResource = &ResourceTemplate{
		Mapper: mashschema.PlanServiceEndpointMapper,
		DoCreate: func(ctx context.Context, client v3client.Client, upsertable mashschema.Upsertable, identifier mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			return client.CreatePlanEndpoint(ctx, upsertable.(masherytypes.PackagePlanServiceEndpointIdentifier))
		},
		DoDelete: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) error {
			return client.DeletePlanEndpoint(ctx, identifier.(masherytypes.PackagePlanServiceEndpointIdentifier))
		},
	}
}
