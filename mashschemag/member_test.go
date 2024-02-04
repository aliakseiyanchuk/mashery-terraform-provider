package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemberBuilderWillProduceSchema(t *testing.T) {
	schema := MemberResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestMemberBuilderMapperIdentity(t *testing.T) {
	ref := masherytypes.MemberIdentifier{
		MemberId: "mid",
		Username: "un",
	}
	autoTestIdentity(t, MemberResourceSchemaBuilder, ref)
}

func TestMemberBuilderMapper(t *testing.T) {
	autoTestMappings(t, MemberResourceSchemaBuilder, func() masherytypes.Member {
		return masherytypes.Member{}
	})
}
