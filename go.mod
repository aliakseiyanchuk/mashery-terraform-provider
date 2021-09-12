module terraform-provider-mashery

go 1.15

require (
	github.com/aliakseiyanchuk/mashery-v3-go-client v0.0.0-20210110193017-ba218ef21d7e
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/terraform-plugin-sdk v1.16.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.4.0
	github.com/stretchr/testify v1.7.0
)

replace github.com/aliakseiyanchuk/mashery-v3-go-client => ../mashery-v3-go-client
