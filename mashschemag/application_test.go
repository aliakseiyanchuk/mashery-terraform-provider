package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplicationBuilderWillProduceSchema(t *testing.T) {
	schema := ApplicationResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestApplicationBuilderMapperIdentity(t *testing.T) {
	ref := ApplicationOfMemberIdentifier{
		MemberIdentifier: masherytypes.MemberIdentifier{
			MemberId: "member-id",
			Username: "member-user-name",
		},
		ApplicationIdentifier: masherytypes.ApplicationIdentifier{
			ApplicationId: "app-id",
		},
	}

	autoTestIdentity(t, ApplicationResourceSchemaBuilder, ref)
}

func TestApplicationBuilderMapper(t *testing.T) {
	autoTestMappings(t, ApplicationResourceSchemaBuilder, func() masherytypes.Application {
		return masherytypes.Application{}
	}, "Eav")
}
