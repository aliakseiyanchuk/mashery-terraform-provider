package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var PackagePlanServiceResource *ResourceTemplate[masherytypes.PackagePlanIdentifier, masherytypes.PackagePlanServiceIdentifier, mashschemag.PackagePlanServiceParam]

func init() {
	PackagePlanServiceResource = &ResourceTemplate[masherytypes.PackagePlanIdentifier, masherytypes.PackagePlanServiceIdentifier, mashschemag.PackagePlanServiceParam]{
		Schema: mashschemag.PackagePlanServiceResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.PackagePlanServiceResourceSchemaBuilder.Mapper(),

		ImportIdentityParser: regexpImportIdentityParser("/packages/(.+)/plans/(.+)/services/(.+)",
			masherytypes.PackagePlanServiceIdentifier{},
			func(items []string) masherytypes.PackagePlanServiceIdentifier {
				rv := masherytypes.PackagePlanServiceIdentifier{}
				rv.PackageId = items[1]
				rv.PlanId = items[2]
				rv.ServiceId = items[3]

				return rv
			}),

		UpsertableFunc: func() mashschemag.PackagePlanServiceParam {
			return mashschemag.PackagePlanServiceParam{}
		},

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanServiceIdentifier) (mashschemag.PackagePlanServiceParam, bool, error) {
			rvObj := mashschemag.PackagePlanServiceParam{
				ServiceIdentifier: identifier.ServiceIdentifier,
			}

			serviceExists, err := client.CheckPlanServiceExists(ctx, identifier)
			return rvObj, serviceExists, err
		},

		DoCreate: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanIdentifier, m mashschemag.PackagePlanServiceParam) (mashschemag.PackagePlanServiceParam, masherytypes.PackagePlanServiceIdentifier, error) {

			ident := masherytypes.PackagePlanServiceIdentifier{
				ServiceIdentifier:     m.ServiceIdentifier,
				PackagePlanIdentifier: identifier,
			}

			if _, err := client.CreatePlanService(ctx, ident); err != nil {
				return mashschemag.PackagePlanServiceParam{}, masherytypes.PackagePlanServiceIdentifier{}, err
			} else {
				return m, ident, nil
			}
		},

		// Update is not required: it will be delete-only

		DoDelete: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanServiceIdentifier) error {
			return client.DeletePlanService(ctx, identifier)
		},

		DoCountOffending: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanServiceIdentifier) (int64, error) {
			return client.CountPlanEndpoints(ctx, identifier)
		},
	}
}
