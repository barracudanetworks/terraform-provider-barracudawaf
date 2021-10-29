package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var SECPOLICY_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_security_policies" "demo_security_policy_1" {
    name       = "DemoPolicy1"
    based_on   = "Create New"
}
`

func TestAccBarracudaWAFSecurityPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: SECPOLICY_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityPolicyExists("DemoPolicy1"),
					resource.TestCheckResourceAttr("barracudawaf_security_policies.demo_security_policy_1", "name", "DemoPolicy1"),
				),
			},
		},
	})
}

func testCheckSecurityPolicyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/security-policies/" + name
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		securityPolicies, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if securityPolicies == nil {
			return fmt.Errorf("security policy %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range securityPolicies.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("security policy (%s) not found on the system", name)
		}

		return nil
	}
}
