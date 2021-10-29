package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var CONTENTRULES_SERVER_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
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

resource "barracudawaf_content_rule_servers" "demo_rule_group_server_1" {
    name        = "DemoRuleGroupServer1"
    identifier  = "Hostname"
    hostname    = "barracuda.com"
    parent      = [ barracudawaf_services.demo_app_1.name, barracudawaf_content_rules.demo_rule_group_1.name ]

    depends_on = [ barracudawaf_content_rules.demo_rule_group_1 ]
}
`

func TestAccBarracudaWAFContentRulesServer_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CONTENTRULES_SERVER_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckContentRulesServerExists("DemoRuleGroupServer1"),
					resource.TestCheckResourceAttr("barracudawaf_content_rule_servers.demo_rule_group_server_1", "identifier", "Hostname"),
					resource.TestCheckResourceAttr("barracudawaf_content_rule_servers.demo_rule_group_server_1", "hostname", "barracuda.com"),
					resource.TestCheckResourceAttr("barracudawaf_content_rule_servers.demo_rule_group_server_1", "name", "DemoRuleGroupServer1"),
				),
			},
		},
	})
}

func testCheckContentRulesServerExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/services/DemoApp1/content-rules/DemoRuleGroup1/content-rule-servers/" + name
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		contentRuleServers, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if contentRuleServers == nil {
			return fmt.Errorf("content rule server %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range contentRuleServers.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("content rule server (%s) not found on the system", name)
		}

		return nil
	}
}
