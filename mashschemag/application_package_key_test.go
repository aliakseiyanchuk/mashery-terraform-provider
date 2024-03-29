package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"terraform-provider-mashery/mashschema"
	"testing"
)

func TestApplicationPackageKeyBuilderWillProduceSchema(t *testing.T) {
	schema := ApplicationPackageKeyResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestApplicationPackageKeyBuilderMapperIdentity(t *testing.T) {
	ref := ApplicationPackageKeyIdentifier{
		PackageKeyIdentifier: masherytypes.PackageKeyIdentifier{
			PackageKeyId: "key-id",
		},
		ApplicationIdentifier: masherytypes.ApplicationIdentifier{
			ApplicationId: "app-id",
		},
	}

	autoTestIdentity(t, ApplicationPackageKeyResourceSchemaBuilder, ref)
}

func TestApplicationPackageKeyBuilderMapper(t *testing.T) {
	autoTestMappings(t, ApplicationPackageKeyResourceSchemaBuilder, func() masherytypes.PackageKey {
		return masherytypes.PackageKey{}
	})
}

func TestApplicationPackageKeyBuilderMapperLimitsField(t *testing.T) {
	mapper, testData := ApplicationPackageKeyResourceSchemaBuilder.MapperAndTestData()

	packageKey := masherytypes.PackageKey{
		Limits: &[]masherytypes.Limit{
			{
				Period:  "second",
				Source:  "plan",
				Ceiling: 2,
			},
			{
				Period:  "day",
				Source:  "plan",
				Ceiling: 5000,
			},
		},
	}

	dg := mapper.RemoteToSchema(&packageKey, testData)
	setData := testData.Get(mashschema.ApplicationPackageKeyLimits)
	assert.NotNil(t, setData)
	assert.Equal(t, 0, len(dg))

	readBackKey := masherytypes.PackageKey{}
	mapper.SchemaToRemote(testData, &readBackKey)

	assert.Equal(t, 2, len(*readBackKey.Limits))

	assert.Equal(t, "second", (*readBackKey.Limits)[0].Period)
	assert.Equal(t, "plan", (*readBackKey.Limits)[0].Source)
	assert.Equal(t, int64(2), (*readBackKey.Limits)[0].Ceiling)

	assert.Equal(t, "day", (*readBackKey.Limits)[1].Period)
	assert.Equal(t, "plan", (*readBackKey.Limits)[1].Source)
	assert.Equal(t, int64(5000), (*readBackKey.Limits)[1].Ceiling)

}

func TestApplicationPackageKeyBuilderMapperPackagePlan(t *testing.T) {
	mapper, testData := ApplicationPackageKeyResourceSchemaBuilder.MapperAndTestData()

	packageKey := masherytypes.PackageKey{
		Package: &masherytypes.Package{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id: "package-id",
			},
		},
		Plan: &masherytypes.Plan{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id: "plan-id",
			},
		},
	}

	dg := mapper.RemoteToSchema(&packageKey, testData)
	assert.Equal(t, 0, len(dg))

	readBackKey := masherytypes.PackageKey{}
	mapper.SchemaToRemote(testData, &readBackKey)

	assert.NotNil(t, readBackKey.Package)
	assert.Equal(t, "package-id", readBackKey.Package.Id)

	assert.NotNil(t, readBackKey.Plan)
	assert.Equal(t, "plan-id", readBackKey.Plan.Id)
}

func TestApplicationPackageKeyBuilderMapperCeilings(t *testing.T) {
	mapper, testData := ApplicationPackageKeyResourceSchemaBuilder.MapperAndTestData()

	testData.Set(mashschema.ApplicationPackageKeyQpsLimitCeiling, 1)
	testData.Set(mashschema.ApplicationPackageKeyRateLimitCeiling, 2)

	v := masherytypes.PackageKey{}
	mapper.SchemaToRemote(testData, &v)

	assert.NotNil(t, v.QpsLimitCeiling)
	assert.NotNil(t, v.RateLimitCeiling)

	assert.Equal(t, int64(1), *v.QpsLimitCeiling)
	assert.Equal(t, int64(2), *v.RateLimitCeiling)
}

func TestApplicationPackageKeyBuilderMapperCeilingsWithZero(t *testing.T) {
	mapper, testData := ApplicationPackageKeyResourceSchemaBuilder.MapperAndTestData()

	var err error
	err = testData.Set(mashschema.ApplicationPackageKeyQpsLimitCeiling, 0)
	assert.Nil(t, err)

	err = testData.Set(mashschema.ApplicationPackageKeyRateLimitCeiling, 0)
	assert.Nil(t, err)

	v := masherytypes.PackageKey{}
	mapper.SchemaToRemote(testData, &v)

	assert.NotNil(t, v.QpsLimitCeiling)
	assert.NotNil(t, v.RateLimitCeiling)

	assert.Equal(t, int64(0), *v.QpsLimitCeiling)
	assert.Equal(t, int64(0), *v.RateLimitCeiling)
}
