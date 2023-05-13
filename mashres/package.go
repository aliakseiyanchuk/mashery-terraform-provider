package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
	"terraform-provider-mashery/tfmapper"
)

var PackageResource *ResourceTemplate[tfmapper.Orphan, masherytypes.PackageIdentifier, masherytypes.Package]

func init() {
	PackageResource = &ResourceTemplate[tfmapper.Orphan, masherytypes.PackageIdentifier, masherytypes.Package]{
		Schema: mashschemag.PackageResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.PackageResourceSchemaBuilder.Mapper(),

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackageIdentifier) (*masherytypes.Package, error) {
			return client.GetPackage(ctx, identifier)
		},

		DoCreate: func(ctx context.Context, client v3client.Client, orphan tfmapper.Orphan, m masherytypes.Package) (*masherytypes.Package, *masherytypes.PackageIdentifier, error) {
			if createdPackage, err := client.CreatePackage(ctx, m); err != nil {
				return nil, nil, err
			} else {
				rvIdent := createdPackage.Identifier()
				return createdPackage, &rvIdent, nil
			}
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackageIdentifier, m masherytypes.Package) (*masherytypes.Package, error) {
			m.Id = identifier.PackageId

			if updatedPack, err := client.UpdatePackage(ctx, m); err != nil {
				return nil, err
			} else {
				return updatedPack, err
			}
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackageIdentifier) error {
			return client.DeletePackage(ctx, identifier)
		},

		DoCountOffending: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackageIdentifier) (int64, error) {
			return client.CountPlans(ctx, identifier)
		},
	}
}
