package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"reflect"
	"terraform-provider-mashery/mashschema"
	"testing"
)

func TestServiceEndpointBuilderWillProduceSchema(t *testing.T) {
	schema := ServiceEndpointResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestServiceEndpointIdentityMapping(t *testing.T) {
	autoTestIdentity(t, ServiceEndpointResourceSchemaBuilder, masherytypes.ServiceEndpointIdentifier{
		EndpointId: "defg",
		ServiceIdentifier: masherytypes.ServiceIdentifier{
			ServiceId: "abc",
		},
	})

	autoTestParentIdentity(t, ServiceCacheResourceSchemaBuilder, masherytypes.ServiceIdentifier{
		ServiceId: "defg",
	})
}

func TestServiceEndpointBuilderMappings(t *testing.T) {
	autoTestMappings(t, ServiceEndpointResourceSchemaBuilder, func() masherytypes.Endpoint {
		return masherytypes.Endpoint{}
	})
}

func TestServiceEndpointMappingProcessorReadWrite(t *testing.T) {
	mapper, state := ServiceEndpointResourceSchemaBuilder.MapperAndTestData()

	proc := masherytypes.Processor{
		PreProcessEnabled:  true,
		PostProcessEnabled: false,
		PostInputs:         map[string]string{},
		PreInputs: map[string]string{
			"A": "B",
		},
		Adapter: "abc",
	}
	remote := masherytypes.Endpoint{
		Processor: &proc,
	}

	mapper.RemoteToSchema(&remote, state)

	rv := masherytypes.Endpoint{}
	mapper.SchemaToRemote(state, &rv)

	assert.NotNil(t, rv.Processor)
	assert.True(t, reflect.DeepEqual(proc, *rv.Processor))
}

func TestServiceEndpointMappingSystemAuthReadWrite(t *testing.T) {
	mapper, state := ServiceEndpointResourceSchemaBuilder.MapperAndTestData()

	userName := "aa"
	password := "bb"
	cert := "cc"
	proc := masherytypes.SystemDomainAuthentication{
		Type:        "dd",
		Username:    &userName,
		Certificate: &cert,
		Password:    &password,
	}
	remote := masherytypes.Endpoint{
		SystemDomainAuthentication: &proc,
	}

	mapper.RemoteToSchema(&remote, state)

	rv := masherytypes.Endpoint{}
	mapper.SchemaToRemote(state, &rv)

	assert.NotNil(t, rv.SystemDomainAuthentication)
	assert.True(t, reflect.DeepEqual(proc, *rv.SystemDomainAuthentication))
}

func TestServiceEndpointErrorSetIdMappingReadWrite(t *testing.T) {
	mapper, state := ServiceEndpointResourceSchemaBuilder.MapperAndTestData()
	remote := masherytypes.Endpoint{
		ErrorSet: &masherytypes.AddressableV3Object{
			Id: "ErrorSetId",
		},
	}

	mapper.RemoteToSchema(&remote, state)

	rv := masherytypes.Endpoint{}
	mapper.SchemaToRemote(state, &rv)

	assert.NotNil(t, rv.ErrorSet)
	assert.Equal(t, "ErrorSetId", rv.ErrorSet.Id)
}

func TestServiceEndpointCorsMapping(t *testing.T) {
	autoTestNestedObjectMappings(t, ServiceEndpointResourceSchemaBuilder, func() (masherytypes.Endpoint, *masherytypes.Cors) {
		rv := masherytypes.Endpoint{
			Cors: &masherytypes.Cors{},
		}
		return rv, rv.Cors
	})
}

func TestServiceEndpointProcessorMapping(t *testing.T) {
	autoTestNestedObjectMappings(t, ServiceEndpointResourceSchemaBuilder, func() (masherytypes.Endpoint, *masherytypes.Processor) {
		rv := masherytypes.Endpoint{
			Processor: &masherytypes.Processor{},
		}
		return rv, rv.Processor
	})
}

func TestServiceEndpointCacheMapping(t *testing.T) {
	autoTestNestedObjectMappings(t, ServiceEndpointResourceSchemaBuilder, func() (masherytypes.Endpoint, *masherytypes.Cache) {
		rv := masherytypes.Endpoint{
			Cache: &masherytypes.Cache{},
		}
		return rv, rv.Cache
	})
}

func TestEndpointWillContainEmptyNonNilProcessor(t *testing.T) {
	mapper, state := ServiceEndpointResourceSchemaBuilder.MapperAndTestData()
	rv := masherytypes.Endpoint{}

	mapper.SchemaToRemote(state, &rv)

	assert.NotNil(t, rv.Processor)
	assert.Equal(t, "", rv.Processor.Adapter)
	assert.False(t, rv.Processor.PreProcessEnabled)
	assert.False(t, rv.Processor.PostProcessEnabled)

	assert.NotNil(t, rv.Processor.PreInputs)
	assert.NotNil(t, rv.Processor.PostInputs)
}

func TestServiceEndpointDeveloperKeyLocation_WillAcceptExpectedValues(t *testing.T) {
	autoTestFieldEnumValuesValidation(t,
		ServiceEndpointResourceSchemaBuilder,
		mashschema.MashEndpointApiKeyValueLocations,
		mashschema.DeveloperAPIKeyLocationsEnum)
}

func TestServiceEndpointDeveloperKeyLocation_WillAcceptCombination(t *testing.T) {
	mapper, state := ServiceEndpointResourceSchemaBuilder.MapperAndTestData()

	err := state.Set(mashschema.MashEndpointApiKeyValueLocations, []string{"request-body", "request-parameters"})
	assert.Nil(t, err)
	dg := mapper.IsStateValid(state)
	assert.Equal(t, 0, len(dg))

	err = state.Set(mashschema.MashEndpointApiKeyValueLocations, []string{"request-parameters", "request-body"})
	assert.Nil(t, err)
	dg = mapper.IsStateValid(state)
	assert.Equal(t, 0, len(dg))
}

func TestServiceEndpointDeveloperKeyLocation_WillRejectIllegalCombination(t *testing.T) {
	mapper, state := ServiceEndpointResourceSchemaBuilder.MapperAndTestData()

	err := state.Set(mashschema.MashEndpointApiKeyValueLocations, []string{"request-body", "request-header"})
	assert.Nil(t, err)
	dg := mapper.IsStateValid(state)

	assert.Equal(t, 1, len(dg))
	assert.Equal(t,
		"field validation has returned the following error: request-header and request-body is an illegal combination, only request-body and request-parameters can be specified together",
		dg[0].Detail)

	err = state.Set(mashschema.MashEndpointApiKeyValueLocations, []string{"request-body", "request-parameters", "request-header"})
	assert.Nil(t, err)
	dg = mapper.IsStateValid(state)

	assert.Equal(t, 1, len(dg))
	assert.Equal(t,
		"field validation has returned the following error: too many options specified (3), maximum possible is 2",
		dg[0].Detail)
}

func TestServiceEndpointMethodLocation_WillAcceptExpectedValues(t *testing.T) {
	autoTestFieldEnumValuesValidation(t,
		ServiceEndpointResourceSchemaBuilder,
		mashschema.MashEndpointApiMethodDetectionLocations,
		mashschema.MethodLocationsEnum)
}

func TestServiceEndpointMethod_WillAcceptCombination(t *testing.T) {
	mapper, state := ServiceEndpointResourceSchemaBuilder.MapperAndTestData()

	err := state.Set(mashschema.MashEndpointApiMethodDetectionLocations, []string{"request-body", "request-parameters"})
	assert.Nil(t, err)
	dg := mapper.IsStateValid(state)
	assert.Equal(t, 0, len(dg))

	err = state.Set(mashschema.MashEndpointApiMethodDetectionLocations, []string{"request-parameters", "request-body"})
	assert.Nil(t, err)
	dg = mapper.IsStateValid(state)
	assert.Equal(t, 0, len(dg))
}

func TestServiceEndpointMethod_WillRejectIllegalCombination(t *testing.T) {
	mapper, state := ServiceEndpointResourceSchemaBuilder.MapperAndTestData()

	err := state.Set(mashschema.MashEndpointApiMethodDetectionLocations, []string{"request-body", "request-header"})
	assert.Nil(t, err)
	dg := mapper.IsStateValid(state)

	assert.Equal(t, 1, len(dg))
	assert.Equal(t,
		"field validation has returned the following error: request-header and request-body is an illegal combination, only request-body and request-parameters can be specified together",
		dg[0].Detail)

	err = state.Set(mashschema.MashEndpointApiMethodDetectionLocations, []string{"request-body", "request-parameters", "request-header"})
	assert.Nil(t, err)
	dg = mapper.IsStateValid(state)

	assert.Equal(t, 1, len(dg))
	assert.Equal(t,
		"field validation has returned the following error: too many options specified (3), maximum possible is 2",
		dg[0].Detail)
}

func TestServiceEndpointForwardedHeaders(t *testing.T) {
	autoTestFieldEnumValuesValidation(t, ServiceEndpointResourceSchemaBuilder, mashschema.MashEndpointForwardedHeaders, mashschema.ForwardedHeadersEnum)
}

func TestServiceEndpointReturnedHeaders(t *testing.T) {
	autoTestFieldEnumValuesValidation(t, ServiceEndpointResourceSchemaBuilder, mashschema.MashEndpointReturnedHeaders, mashschema.ReturnedHeadersEnum)
}

func TestServiceEndpointHttpMethods(t *testing.T) {
	autoTestFieldEnumValuesValidation(t, ServiceEndpointResourceSchemaBuilder, mashschema.MashEndpointSupportedHttpMethods, mashschema.HttoMethodsEnum)
}
