package mashery_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashery"
	"testing"
	"time"
)

func TestV3EmailTemplateSetIdToResourceData_DataSource(t *testing.T) {
	res := schema.Resource{
		Schema: mashery.DataSourceServiceCacheScheme,
	}
	d := res.TestResourceData()

	now := v3client.MasheryJSONTime(time.Now())

	orig := v3client.MasheryEmailTemplateSet{
		AddressableV3Object: v3client.AddressableV3Object{
			Id:      "setId",
			Name:    "name",
			Created: &now,
			Updated: &now,
		},
		Type:           "Type",
		EmailTemplates: nil,
	}

	mashery.V3EmailTemplateSetIdToResourceData(&orig, d)

	name := d.Get(mashery.MashObjName).(string)
	setType := d.Get(mashery.MashEmailTemplateSetType).(string)
	setId := d.Get(mashery.MashEmailTemplateSetId).(string)

	assertSameString(t, "id", &orig.Id, &setId)
	assertSameString(t, "name", &orig.Name, &name)
	assertSameString(t, "type", &orig.Type, &setType)
}
