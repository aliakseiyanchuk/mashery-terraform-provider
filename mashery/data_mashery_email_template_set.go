package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschema"
)

var emailTemplateSet *DatasourceTemplate

func init() {
	emailTemplateSet = &DatasourceTemplate{
		V3ObjectTypeName: "email template set",
		QueryExtractor:   DataSourceDefaultQueryExtractor,
		RequireUnique:    true,
		Query: func(ctx context.Context, cl v3client.Client, query map[string]string) ([]interface{}, error) {
			if v3State, err := cl.ListEmailTemplateSetsFiltered(ctx, query, mashschema.EmptyStringArray); err != nil {
				return nil, err
			} else {
				return mashschema.CoerceEmailTemplateSetArrayToInterfaceType(v3State), nil
			}
		},
		Mapper: mashschema.EmailTemplateSetMapper,
	}
}
