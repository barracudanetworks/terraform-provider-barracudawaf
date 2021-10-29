package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var SERVER_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
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

resource "barracudawaf_servers" "demo_server_1" {
    name            = "DemoServer1"
    ip_address      = "99.86.47.44"
    identifier      = "IP Address"
    address_version = "IPv4"
    status          = "In Service"
    port            = "80"
    comments        = "Creating the Demo Server"
    parent          = [ "DemoApp1" ]

    depends_on = [ barracudawaf_services.demo_app_1 ]
}
`

func TestAccBarracudaWAFServer_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: SERVER_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckServerExists("DemoServer1"),
					resource.TestCheckResourceAttr("barracudawaf_servers.demo_server_1", "port", "80"),
					resource.TestCheckResourceAttr("barracudawaf_servers.demo_server_1", "name", "DemoServer1"),
					resource.TestCheckResourceAttr("barracudawaf_servers.demo_server_1", "status", "In Service"),
					resource.TestCheckResourceAttr("barracudawaf_servers.demo_server_1", "address_version", "IPv4"),
					resource.TestCheckResourceAttr("barracudawaf_servers.demo_server_1", "identifier", "IP Address"),
				),
			},
		},
	})
}

func testCheckServerExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/services/DemoApp1/servers/" + name
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		servers, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if servers == nil {
			return fmt.Errorf("server %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range servers.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("server (%s) not found on the system", name)
		}

		return nil
	}
}
