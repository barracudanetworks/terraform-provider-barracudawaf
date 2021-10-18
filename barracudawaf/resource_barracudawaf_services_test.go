package barracudawaf

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var BARRACUDA_WAF_PROVIDER = `
provider "barracudawaf" {
    address  = "` + os.Getenv("BARRACUDA_WAF_IP") + `"
    username = "` + os.Getenv("BARRACUDA_WAF_USERNAME") + `"
    port     = "8443"
    password = "` + os.Getenv("BARRACUDA_WAF_PASSWORD") + `"
}
`

var SERVICE_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
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
`

var SERVICE_RESOURCE_EDIT = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_services" "demo_app_1" {
	name			= "DemoApp1"
    port            = "80"
    comments        = "Edit Service with Terraform"

    basic_security {
      mode = "Active"
    }
}`

func TestAccBarracudaWAFService_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: SERVICE_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckServiceExists("DemoApp1"),
					resource.TestCheckResourceAttr("barracudawaf_services.demo_app_1", "name", "DemoApp1"),
					resource.TestCheckResourceAttr("barracudawaf_services.demo_app_1", "port", "90"),
					resource.TestCheckResourceAttr("barracudawaf_services.demo_app_1", "type", "HTTP"),
				),
			},
			{
				Config: SERVICE_RESOURCE_EDIT,
				Check: resource.ComposeTestCheckFunc(
					testCheckServiceExists("DemoApp1"),
					resource.TestCheckResourceAttr("barracudawaf_services.demo_app_1", "port", "80"),
					resource.TestCheckResourceAttr("barracudawaf_services.demo_app_1", "comments", "Edit Service with Terraform"),
				),
			},
		},
	})
}

func testCheckServiceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/services/" + name
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		service, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if service == nil {
			return fmt.Errorf("service %s was not created.", name)
		}

		return nil
	}
}
