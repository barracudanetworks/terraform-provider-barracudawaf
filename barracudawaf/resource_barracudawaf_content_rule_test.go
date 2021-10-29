package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var CONTENTRULES_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_services" "demo_app_1" {
    name            = "DemoApp1"
    ip_address      = "172.30.1.4"
    port            = "90"
    type            = "HTTP"
    vsite           = "default"
    address_version = "IPv4"
    status          = "On"
    group           = "default"
    comments        = "Demo Service with Terraform"
}

resource "barracudawaf_security_policies" "demo_security_policy_1" {
    name       = "DemoPolicy1"
    based_on   = "Create New"
    
    depends_on = [ barracudawaf_services.demo_app_1 ]
}

resource "barracudawaf_content_rules" "demo_rule_group_1" {
    name                = "DemoRuleGroup1"
    url_match           = "/index.html"
    host_match          = "www.example.com"
    web_firewall_policy = "DemoPolicy1"
    mode                = "Active"
    parent              = [ barracudawaf_services.demo_app_1.name ]
    
    depends_on          = [ barracudawaf_security_policies.demo_security_policy_1 ]
}
`

func TestAccBarracudaWAFContentRules_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CONTENTRULES_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckContentRulesExists("DemoRuleGroup1"),
					resource.TestCheckResourceAttr("barracudawaf_content_rules.demo_rule_group_1", "mode", "Active"),
					resource.TestCheckResourceAttr("barracudawaf_content_rules.demo_rule_group_1", "name", "DemoRuleGroup1"),
					resource.TestCheckResourceAttr("barracudawaf_content_rules.demo_rule_group_1", "url_match", "/index.html"),
					resource.TestCheckResourceAttr("barracudawaf_content_rules.demo_rule_group_1", "host_match", "www.example.com"),
					resource.TestCheckResourceAttr("barracudawaf_content_rules.demo_rule_group_1", "web_firewall_policy", "DemoPolicy1"),
				),
			},
		},
	})
}

func testCheckContentRulesExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/services/DemoApp1/content-rules/" + name
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		contentRules, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if contentRules == nil {
			return fmt.Errorf("content rule %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range contentRules.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("content rule (%s) not found on the system", name)
		}

		return nil
	}
}
