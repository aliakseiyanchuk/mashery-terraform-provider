package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	mashschema "terraform-provider-mashery/mashschema"
)

var roleDataSource *DatasourceTemplate

func init() {
	roleDataSource = &DatasourceTemplate{
		QueryExtractor: DataSourceDefaultQueryExtractor,
		RequireUnique:  true,
		Query: func(ctx context.Context, cl v3client.Client, query map[string]string) ([]interface{}, error) {
			if v3Objects, err := cl.ListRolesFiltered(ctx, query, mashschema.EmptyStringArray); err != nil {
				return nil, err
			} else {
				return mashschema.CoerceRolesArrayToInterfaceType(v3Objects), nil
			}

		},
		Mapper: mashschema.RoleMapper,
	}
}
