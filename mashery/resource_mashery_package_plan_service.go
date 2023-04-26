package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschema"
)

var PackagePlanServiceResource *ResourceTemplate

func init() {
	PackagePlanServiceResource = &ResourceTemplate{
		Mapper: mashschema.PlanServiceMapper,
		DoRead: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			cr, err := client.CheckPlanServiceExists(ctx, identifier.(masherytypes.PackagePlanServiceIdentifier))
			if cr {
				return identifier, err
			} else {
				return nil, err
			}
		},
		DoCreate: func(ctx context.Context, client v3client.Client, upsertable mashschema.Upsertable, identifier mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			return client.CreatePlanService(ctx, upsertable.(masherytypes.PackagePlanServiceIdentifier))
		},
		DoCountOffending: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) (int64, error) {
			return client.CountPlanEndpoints(ctx, identifier.(masherytypes.PackagePlanServiceIdentifier))
		},
		DoDelete: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) error {
			return client.DeletePlanService(ctx, identifier.(masherytypes.PackagePlanServiceIdentifier))
		},
	}
}
