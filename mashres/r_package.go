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

		UpsertableFunc: func() masherytypes.Package {
			return masherytypes.Package{}
		},

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackageIdentifier) (masherytypes.Package, bool, error) {
			return client.GetPackage(ctx, identifier)
		},

		DoCreate: func(ctx context.Context, client v3client.Client, orphan tfmapper.Orphan, m masherytypes.Package) (masherytypes.Package, masherytypes.PackageIdentifier, error) {
			if createdPackage, err := client.CreatePackage(ctx, m); err != nil {
				return masherytypes.Package{}, masherytypes.PackageIdentifier{}, err
			} else {
				rvIdent := createdPackage.Identifier()
				return createdPackage, rvIdent, nil
			}
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackageIdentifier, m masherytypes.Package) (masherytypes.Package, error) {
			m.Id = identifier.PackageId

			// The reason why this code is organized this way is that:
			// - We may need to set organization to the area level if it was previously set; however
			// - testing shows that Mashery doesn't treat organization attribute as a purely CRUD logic; e.g. it
			//   will not be returned during the update.
			//
			// So the working process here is:
			// - Send the update to the desired state.
			// - Read the actual state in the separate call
			// - In case the desired state does NOT have an organization, and actual DOES, call a
			//   specialized ownership reset method to drop the package ownership to the area level.

			if updatedPack, updateErr := client.UpdatePackage(ctx, m); updateErr != nil {
				return updatedPack, updateErr
			} else {
				if readBackPack, _, readBackErr := client.GetPackage(ctx, m.Identifier()); readBackErr != nil {
					return readBackPack, readBackErr
				} else if m.Organization == nil && readBackPack.Organization != nil {
					return client.ResetPackageOwnership(ctx, m.Identifier())
				} else {
					return updatedPack, updateErr
				}
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
