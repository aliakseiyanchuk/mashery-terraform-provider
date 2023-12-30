package tfmapper

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"terraform-provider-mashery/mashschema"
	"testing"
	"time"
)

func TestSuppressSameDuration(t *testing.T) {
	assert.True(t, ShouldSuppressSameDuration("1h", "1h0m0s"))
	assert.False(t, ShouldSuppressSameDuration("1h", "1h0m1s"))
}

func TestDurationMapperSettingLeniency(t *testing.T) {
	builder := durationMapperBuilder()

	resourceData := builder.TestResourceData()
	err := resourceData.Set(mashschema.MashSvcOAuthAccessTokenTtl, "1h")
	assert.Nil(t, err)

	remote := masherytypes.MasheryOAuth{
		AccessTokenTtl: 3600,
	}

	mapper := builder.Mapper()
	mapper.RemoteToSchema(&remote, resourceData)

	readBack := resourceData.Get(mashschema.MashSvcOAuthAccessTokenTtl).(string)
	assert.Equal(t, "1h", readBack)
}

func TestDurationMapperConvertingRemoteTime(t *testing.T) {
	builder := durationMapperBuilder()

	resourceData := builder.TestResourceData()
	err := resourceData.Set(mashschema.MashSvcOAuthAccessTokenTtl, "1h")
	assert.Nil(t, err)

	remote := masherytypes.MasheryOAuth{
		AccessTokenTtl: 3601,
	}

	mapper := builder.Mapper()
	mapper.RemoteToSchema(&remote, resourceData)

	readBack := resourceData.Get(mashschema.MashSvcOAuthAccessTokenTtl).(string)
	assert.Equal(t, "1h0m1s", readBack)
}

func durationMapperBuilder() *SchemaBuilder[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.MasheryOAuth] {
	builder := NewSchemaBuilder[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.MasheryOAuth]()

	builder.Add(&DurationFieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *int64 {
			return &in.AccessTokenTtl
		},
		Unit: time.Second,
		FieldMapperBase: FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthAccessTokenTtl,
			Schema: &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Access token expires after the specified time has passed. TTL time is specified in seconds",
				Default:          "1h",
				ValidateDiagFunc: mashschema.ValidateDuration,
			},
		},
	})
	return builder
}
