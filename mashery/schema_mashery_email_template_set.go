package mashery

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const MashEmailTemplateSetType = "type"
const MashEmailTemplateSetId = "set_id"

var EmailTemplateSetSchema = map[string]*schema.Schema{}

var DataSourceEmailTemplateSetSchema = map[string]*schema.Schema{
	MashDataSourceSearch: {
		Type:        schema.TypeMap,
		Required:    true,
		Description: "Search conditions for this email set, typically name = value",
		Elem:        stringElem(),
	},
	MashDataSourceRequired: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "If true (default), then email template set must exist. If an element doesn't exist, the error is generated",
	},
}

func V3EmailTemplateSetIdToResourceData(set *masherytypes.MasheryEmailTemplateSet, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashEmailTemplateSetId:   set.Id,
		MashObjCreated:           set.Created.ToString(),
		MashObjUpdated:           set.Updated.ToString(),
		MashObjName:              set.Name,
		MashEmailTemplateSetType: set.Type,
	}

	return SetResourceFields(data, d)
}

func initEmailTemplateSetSchemaBoilerplate() {
	addComputedString(&EmailTemplateSetSchema, MashObjCreated, "Date/time this email template set was created")
	addComputedString(&EmailTemplateSetSchema, MashObjUpdated, "Date/time this email template set was updated")
	addRequiredString(&EmailTemplateSetSchema, MashObjName, "Email data set name")
	addComputedOptionalString(&EmailTemplateSetSchema, MashEmailTemplateSetType, "Email template set type")
}

func init() {
	initEmailTemplateSetSchemaBoilerplate()

	// Copy email source template.
	appendAsComputedInto(&EmailTemplateSetSchema, &DataSourceServiceCacheScheme)
	addComputedString(&DataSourceServiceCacheScheme, MashEmailTemplateSetId, "Email set Id")

}
