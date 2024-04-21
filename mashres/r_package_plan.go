package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var PackagePlanResource *ResourceTemplate[masherytypes.PackageIdentifier, masherytypes.PackagePlanIdentifier, masherytypes.Plan]

func init() {
	PackagePlanResource = &ResourceTemplate[masherytypes.PackageIdentifier, masherytypes.PackagePlanIdentifier, masherytypes.Plan]{
		Schema: mashschemag.PackagePlanResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.PackagePlanResourceSchemaBuilder.Mapper(),

		ImportIdentityParser: regexpImportIdentityParser("/packages/(.)/plans/(.+)",
			masherytypes.PackagePlanIdentifier{},
			func(items []string) masherytypes.PackagePlanIdentifier {
				rv := masherytypes.PackagePlanIdentifier{}
				rv.PackageId = items[1]
				rv.PlanId = items[2]

				return rv
			}),

		UpsertableFunc: func() masherytypes.Plan {
			return masherytypes.Plan{}
		},

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanIdentifier) (masherytypes.Plan, bool, error) {
			return client.GetPlan(ctx, identifier)
		},

		DoCreate: func(ctx context.Context, client v3client.Client, packageId masherytypes.PackageIdentifier, m masherytypes.Plan) (masherytypes.Plan, masherytypes.PackagePlanIdentifier, error) {
			if createdPackage, err := client.CreatePlan(ctx, packageId, m); err != nil {
				return masherytypes.Plan{}, masherytypes.PackagePlanIdentifier{}, err
			} else {
				rvIdent := createdPackage.Identifier()
				return createdPackage, rvIdent, nil
			}
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanIdentifier, m masherytypes.Plan) (masherytypes.Plan, error) {
			m.Id = identifier.PlanId
			m.ParentPackageId = identifier.PackageIdentifier

			if updatedPack, err := client.UpdatePlan(ctx, m); err != nil {
				return masherytypes.Plan{}, err
			} else {
				return updatedPack, err
			}
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanIdentifier) error {
			return client.DeletePlan(ctx, identifier)
		},

		DoCountOffending: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanIdentifier) (int64, error) {
			return client.CountPlanService(ctx, identifier)
		},
	}
}
