package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var PackageDataSource = CreateSingularDataSource(
	mashschemag.PackageResourceSchemaBuilder,
	queryPackage,
)

func queryPackage(ctx context.Context, client v3client.Client, m map[string]string) (masherytypes.PackageIdentifier, *masherytypes.Package, error) {
	if sets, err := client.ListPackagesFiltered(ctx, m); err != nil {
		return masherytypes.PackageIdentifier{}, nil, err
	} else if len(sets) == 1 {
		return sets[0].Identifier(), &sets[0], nil
	} else {
		return masherytypes.PackageIdentifier{}, nil, nil
	}
}
