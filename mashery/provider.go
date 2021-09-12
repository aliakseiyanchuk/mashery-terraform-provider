package mashery

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"os"
	"strconv"
	"time"
)

// Provider schema configuration settings
//
// The provide accepts three types of access credentials:
// - the token value passed to either from Vault secret engine, or by reading  <code>TF_MASHERY_V3_ACCESS_TOKEN</code>
//   environment variable;
// - by specifying a path to the token file created previously with the <code>mash-connect</code> command. The token
//   file should be valid and allow at least 5 minutes of validity.
// - by specifying the path to the file where the long-term credentials are saved.
//
// Due to the nature of the credentials required for Mashery V3, the long-term credentials cannot appear in the
// Terraform state in the clear. V3 credentials allow full access to the API configuration as well as it allow
// retrieving all package key/secret combination.

const (
	envCredentialPass = "TF_MASHERY_CREDS_PASS"
)

var providerConfigSchema = map[string]*schema.Schema{
	"log_file": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Log file where detailed Mashery session information will be saved",
	},
	"qps": {
		Type:        schema.TypeInt,
		Optional:    true,
		DefaultFunc: intValueFromVariable("TF_MASHERY_V3_QPS", 2),
		Description: "Queries per second to observe. Default to 2 queries per second",
	},
	"deploy_duration": {
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("TF_MASHERY_DEPLOY_DURATION", "3m"),
		Description: "Typical duration required to perform deployment. Defaults to 3 minutes",
	},
	"network_latency": {
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("TF_MASHERY_V3_NETWORK_LATENCY", "173ms"),
		Description: "Mean travel time between machine where the Terraform is running and Mashery API. Defaults to 173 (milliseconds).",
	},
	"token": {
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("TF_MASHERY_V3_ACCESS_TOKEN", ""),
		Description: "Actual access token to be used. For best use, obtain this from Vault",
	},
	"token_file": {
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("TF_MASHERY_V3_TOKEN_FILE", v3client.DefaultSavedAccessTokenFilePath()),
		Description: "Read access token from the local file system",
	},
	"creds_file": {
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("TF_MASHERY_V3_CREDS", v3client.DefaultCredentialsFile()),
		Description: "Load long-term V3 credentials from file to obtain V3 access tokens",
	},
}

func intValueFromVariable(envVar string, defaultVal int) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(envVar); v != "" {
			return strconv.ParseInt(v, 10, 0)
		} else {
			return defaultVal, nil
		}
	}
}

var logger *log.Logger
var encoder *json.Encoder

// Send a message to the log file if it exists
func doLogf(format string, params ...interface{}) {
	if logger != nil {
		logger.Printf(format, params...)
	}
}

func doLogJson(msg string, obj interface{}) {
	if logger != nil {
		logger.Println(msg)
		if obj != nil {
			if b, err := json.Marshal(obj); err != nil {
				logger.Println(err.Error())
			} else {
				logger.Println(string(b))
			}
		} else {
			logger.Print("NULL JSON")
		}
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	var diags diag.Diagnostics
	logFile := d.Get("log_file").(string)
	if len(logFile) > 0 {
		encoder = new(json.Encoder)
		encoder.SetIndent("", "  ")

		now := time.Now()
		f, _ := os.Create(fmt.Sprintf("%s_%d%d%d_%d%d%d.log", logFile, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()))
		logger = log.New(f, "TF_MASHERY :", log.LstdFlags)
	}

	var tokenProvider v3client.V3AccessTokenProvider
	qps := d.Get("qps").(int)

	cfgNetLatency := d.Get("network_latency").(string)
	cfgDeployDuration := d.Get("deploy_duration").(string)

	travelComp, latErr := time.ParseDuration(cfgNetLatency)
	deployDur, deplErr := time.ParseDuration(cfgDeployDuration)

	if latErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Not a valid network latency: %s (%s)", cfgNetLatency, latErr),
		})
	}

	if deplErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Not a valid deployment duration: %s (%s)", deployDur, deplErr),
		})
	}

	// 1. Give preference to explicitly specified token
	tkn := d.Get("token").(string)
	// Mashery tokens is 24 characters long. This condition is a built-in protection for the
	// developer passing invalid tokens
	if len(tkn) > 20 {
		tokenProvider = v3client.NewFixedTokenProvider(tkn)
		doLogf("Initialized static token supplied vy the configuration")
	} else {
		// 2. Explicitly specified token is not available. Second preference is to read the
		// token data from file, if its available.
		savedTknFile := d.Get("token_file").(string)
		doLogf("Saved token file is %s", savedTknFile)

		savedTknDat, err := v3client.ReadSavedV3TokenData(savedTknFile)
		if savedTknDat != nil {
			doLogf("Time left in this token: %d", savedTknDat.TimeLeft())
		}
		if err != nil {
			doLogf("Error reading saved token file: %s", err)
		}

		if err == nil && savedTknDat != nil && savedTknDat.TimeLeft() > int(deployDur.Seconds()) {
			tokenProvider = v3client.NewFileSystemTokenProviderFrom(savedTknFile)
			doLogf("Initialized file system token from provider the saved token file")
		} else {
			// 3. The pre-fetched token is not found. At this point, the last remaining option for the
			// developer is to supply the credentials, either within the file, or by supplying environmental
			// variables/
			credentialsFile := d.Get("creds_file").(string)
			mashCredentials := v3client.DeriveAccessCredentials(credentialsFile, os.Getenv(envCredentialPass), nil)
			if !mashCredentials.FullySpecified() {
				doLogf("insufficient token credentials loaded from credentials file; most likely because of incorrect decryption password")
			} else {
				tokenProvider = v3client.NewClientCredentialsProvider(mashCredentials)
				doLogf("Initialized token from mashery credentials")

				// Default max-QPS to 2 queries per second, which is Mashery's default QPS assigned when
				// the Mashery credentials are obtained. If the developer has different (larger) credentials,
				// the agreed QPS value need to be supplied
				if mashCredentials.MaxQPS > 0 {
					qps = mashCredentials.MaxQPS
				}
			}
		}
	}

	if tokenProvider == nil {
		return nil, diag.Errorf("no mashery authentication method supplied")
	}

	cl := v3client.NewHttpClient(tokenProvider, int64(qps), travelComp)
	return cl, diags
}

// Mashery Terraform Provider schema definition
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: providerConfigSchema,
		ResourcesMap: map[string]*schema.Resource{
			"mashery_service":                      resourceMasheryService(),
			"mashery_service_error_set":            resourceMasheryErrorSet(),
			"mashery_processor_chain":              resourceMasheryProcessorChain(),
			"mashery_endpoint":                     resourceMasheryEndpoint(),
			"mashery_endpoint_method":              resourceMasheryEndpointMethod(),
			"mashery_endpoint_method_filter":       resourceMasheryEndpointMethodFilter(),
			"mashery_package":                      resourceMasheryPackage(),
			"mashery_package_plan":                 resourceMasheryPlan(),
			"mashery_package_plan_service":         resourceMasheryPlanService(),
			"mashery_package_plan_endpoint":        resourceMasheryPackagePlanEndpoint(),
			"mashery_package_plan_endpoint_method": resourceMasheryPlanMethod(),
			"mashery_member":                       resourceMasheryMember(),
			"mashery_application":                  resourceMasheryApplication(),
			"mashery_package_key":                  resourceMasheryPackageKey(),
			"mashery_unique_path":                  resourceMasheryUniquePath(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"mashery_public_domains":     dataSourceMasheryPublicDomains(),
			"mashery_system_domains":     dataSourceMasherySystemDomains(),
			"mashery_service":            dataSourceMasheryService(),
			"mashery_service_endpoints":  dataSourceMasheryServiceEndpoints(),
			"mashery_email_template_set": dataSourceMasheryEmailTemplateSet(),
			"mashery_role":               dataSourceMasheryRole(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}
