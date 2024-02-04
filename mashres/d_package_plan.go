package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/mashschemag"
)

var PackagePlanDataSource *SingularDatasourceTemplate[masherytypes.PackageIdentifier, masherytypes.PackagePlanIdentifier, masherytypes.Plan] = CreateSingularParentScopedDataSource(
	mashschemag.PackagePlanResourceSchemaBuilder,
	mashschema.MashPackageRef,
	queryPackagePlans,
)

func queryPackagePlans(ctx context.Context, client v3client.Client, ident masherytypes.PackageIdentifier, m map[string]string) (masherytypes.PackagePlanIdentifier, *masherytypes.Plan, error) {
	if sets, err := client.ListPlansFiltered(ctx, ident, m); err != nil {
		return masherytypes.PackagePlanIdentifier{}, nil, err
	} else if len(sets) == 1 {
		return sets[0].Identifier(), &sets[0], nil
	} else {
		return masherytypes.PackagePlanIdentifier{}, nil, nil
	}
}
