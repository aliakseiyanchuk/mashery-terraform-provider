package mashschema

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const MashEmailTemplateSetType = "type"

const MashEmailTemplateSetId = "set_id"

var EmailTemplateSetMapper *emailTemplateSetMapperImpl

type emailTemplateSetMapperImpl struct {
	DataSourceMapperImpl
}

func (etsm *emailTemplateSetMapperImpl) PersistTyped(set masherytypes.EmailTemplateSet, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashObjCreated:           set.Created.ToString(),
		MashObjUpdated:           set.Updated.ToString(),
		MashObjName:              set.Name,
		MashEmailTemplateSetType: set.Type,
	}

	return SetResourceFields(data, d)
}

func (etsm *emailTemplateSetMapperImpl) initEmailTemplateSetSchemaBoilerplate() {
	etsm.SchemaBuilder().
		AddComputedString(MashObjCreated, "Date/time this email template set was created").
		AddComputedString(MashObjUpdated, "Date/time this email template set was updated").
		AddRequiredString(MashObjName, "Email data set name").
		AddComputedOptionalString(MashEmailTemplateSetType, "Email template set type").
		AddComputedString(MashEmailTemplateSetId, "Email set Id")
}

func init() {
	EmailTemplateSetMapper = &emailTemplateSetMapperImpl{
		DataSourceMapperImpl{
			v3ObjectName: "email template set",
			schema: map[string]*schema.Schema{
				MashDataSourceSearch: {
					Type:        schema.TypeMap,
					Required:    true,
					Description: "Search conditions for this email set, typically name = value",
					Elem:        StringElem(),
				},
				MashDataSourceRequired: {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "If true (default), then email template set must exist. If an element doesn't exist, the error is generated",
				},
			},

			persistOne: func(rv interface{}, d *schema.ResourceData) diag.Diagnostics {
				return EmailTemplateSetMapper.PersistTyped(rv.(masherytypes.EmailTemplateSet), d)
			},
		},
	}

	EmailTemplateSetMapper.initEmailTemplateSetSchemaBoilerplate()
}
