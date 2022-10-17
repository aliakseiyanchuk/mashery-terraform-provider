package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschema"
)

var PackagePlanResource *ResourceTemplate

func init() {
	PackagePlanResource = &ResourceTemplate{
		Mapper: mashschema.PlanMapper,
		DoRead: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			return client.GetPlan(ctx, identifier.(masherytypes.PackagePlanIdentifier))
		},
		DoCreate: func(ctx context.Context, client v3client.Client, upsertable mashschema.Upsertable, identifier mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			return client.CreatePlan(ctx, identifier.(masherytypes.PackageIdentifier), upsertable.(masherytypes.Plan))
		},
		DoUpdate: func(ctx context.Context, client v3client.Client, upsertable mashschema.Upsertable) (mashschema.Upsertable, error) {
			return client.UpdatePlan(ctx, upsertable.(masherytypes.Plan))
		},
		DoCountOffending: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) (int64, error) {
			planId := identifier.(masherytypes.PackagePlanIdentifier)
			return client.CountPlanService(ctx, planId)
		},
		DoDelete: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) error {
			return client.DeletePlan(ctx, identifier.(masherytypes.PackagePlanIdentifier))
		},
	}
}
