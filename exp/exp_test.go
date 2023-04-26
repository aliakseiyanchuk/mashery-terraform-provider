package exp_test

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestLoadingHCLFile(t *testing.T) {

	data, _ := ioutil.ReadFile("example.tf")

	p := hclparse.NewParser()

	file, diags := p.ParseHCL(data, "example.tf")
	assert.Equal(t, 0, len(diags))

	block := file.BlocksAtPos(hcl.Pos{Line: 1})

	ctx := hcl.EvalContext{}

	attrs, _ := block[0].Body.JustAttributes()
	for k, v := range attrs {
		value, _ := v.Expr.Value(&ctx)
		fmt.Printf("%s = %s\n", k, value.AsValueMap())
	}
}
