package barracudawaf

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//Provider : Schema definition for barracudawaf provider
func Provider() *schema.Provider {

	// The actual provider
	provider := &schema.Provider{

		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP Address of the WAF to be configured",
			},
			"port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Admin port on the WAF to be configured",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username of the WAF to be configured",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Password of the WAF to be configured",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"barracudawaf_trusted_ca_certificate":     resourceCudaWAFTrustedCaCertificate(),
			"barracudawaf_content_rules":              resourceCudaWAFContentRules(),
			"barracudawaf_trusted_server_certificate": resourceCudaWAFTrustedServerCertificate(),
			"barracudawaf_services":                   resourceCudaWAFServices(),
			"barracudawaf_content_rule_servers":       resourceCudaWAFContentRuleServers(),
			"barracudawaf_security_policies":          resourceCudaWAFSecurityPolicies(),
			"barracudawaf_signed_certificate":         resourceCudaWAFSignedCertificate(),
			"barracudawaf_self_signed_certificate":    resourceCudaWAFSelfSignedCertificate(),
			"barracudawaf_servers":                    resourceCudaWAFServers(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		IPAddress: d.Get("address").(string),
		AdminPort: d.Get("port").(string),
		Username:  d.Get("username").(string),
		Password:  d.Get("password").(string),
	}
	cfg, err := config.Client()
	if err != nil {
		return cfg, err
	}
	cfg.UserAgent = fmt.Sprintf("Terraform/%s", terraformVersion)
	return cfg, err
}
