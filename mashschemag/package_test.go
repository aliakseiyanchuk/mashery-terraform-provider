package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"reflect"
	"terraform-provider-mashery/mashschema"
	"testing"
	"time"
)

var timeNow = masherytypes.MasheryJSONTime(time.Now())

func TestPackageBuilderWillProduceSchema(t *testing.T) {
	schema := PackageResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestPackageMapperIdentity(t *testing.T) {
	mapper := PackageResourceSchemaBuilder.Mapper()
	testState := PackageResourceSchemaBuilder.TestResourceData()

	testPack := masherytypes.PackageIdentifier{
		PackageId: "123",
	}

	err := mapper.AssignIdentity(testPack, testState)
	assert.Nil(t, err)

	readBack, idErr := mapper.Identity(testState)
	assert.Nil(t, idErr)

	assert.True(t, reflect.DeepEqual(testPack, readBack))
}

func TestPackageMapperUpsertion(t *testing.T) {
	mapper := PackageResourceSchemaBuilder.Mapper()
	testState := PackageResourceSchemaBuilder.TestResourceData()

	threshold := 70
	keyLength := 24
	secretLength := 12

	testPack := masherytypes.Package{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:      "packageId",
			Name:    "packageName",
			Created: &timeNow,
			Updated: &timeNow,
		},
		Description:                 "Desc",
		NotifyDeveloperPeriod:       "hour",
		NotifyDeveloperNearQuota:    false,
		NotifyDeveloperOverQuota:    true,
		NotifyDeveloperOverThrottle: false,
		NotifyAdminPeriod:           "day",
		NotifyAdminNearQuota:        false,
		NotifyAdminOverQuota:        true,
		NotifyAdminOverThrottle:     false,
		NotifyAdminEmails:           "a@b.com,c@d.com",
		NearQuotaThreshold:          &threshold,
		Eav: map[string]string{
			"A": "B",
			"C": "B",
		},
		KeyAdapter:         "adapter",
		KeyLength:          &keyLength,
		SharedSecretLength: &secretLength,
	}

	diags := mapper.RemoteToSchema(&testPack, testState)
	assert.Equal(t, 0, len(diags))

	assert.Equal(t, "packageId", testState.Get(mashschema.MashPackageId))

	upserted := masherytypes.Package{}
	mapper.SchemaToRemote(testState, &upserted)

	assert.Equal(t, testPack.Name, upserted.Name)
	assert.Equal(t, testPack.Description, upserted.Description)

	assert.Equal(t, testPack.NotifyDeveloperPeriod, upserted.NotifyDeveloperPeriod)
	assert.Equal(t, testPack.NotifyDeveloperNearQuota, upserted.NotifyDeveloperNearQuota)
	assert.Equal(t, testPack.NotifyDeveloperOverQuota, upserted.NotifyDeveloperOverQuota)
	assert.Equal(t, testPack.NotifyDeveloperOverThrottle, upserted.NotifyDeveloperOverThrottle)

	assert.Equal(t, testPack.NotifyAdminPeriod, upserted.NotifyAdminPeriod)
	assert.Equal(t, testPack.NotifyAdminNearQuota, upserted.NotifyAdminNearQuota)
	assert.Equal(t, testPack.NotifyAdminOverQuota, upserted.NotifyAdminOverQuota)
	assert.Equal(t, testPack.NotifyAdminOverThrottle, upserted.NotifyAdminOverThrottle)

	assert.Equal(t, *testPack.NearQuotaThreshold, *upserted.NearQuotaThreshold)

	assert.Equal(t, testPack.Eav, upserted.Eav)

	assert.Equal(t, testPack.KeyAdapter, upserted.KeyAdapter)
	assert.Equal(t, *testPack.KeyLength, *upserted.KeyLength)
	assert.Equal(t, *testPack.SharedSecretLength, *upserted.SharedSecretLength)

}
