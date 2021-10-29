package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var SELF_SIGNED_CERT_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_self_signed_certificate" "demo_self_signed_cert_1" {
    name                     = "DemoSelfSignedCert1"
    allow_private_key_export = "Yes"
    city                     = "Bangalore"
    common_name              = "barracuda.com"
    country_code             = "IN"
    key_size                 = "1024"
    key_type                 = "rsa"
    organization_name        = "Barracuda Networks"
    organizational_unit      = "Engineering"
    state                    = "Karnataka"
}
`

func TestAccBarracudaWAFSelfSignedCertificate_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: SELF_SIGNED_CERT_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckSelfSignedCertificateExists("DemoSelfSignedCert1"),
					resource.TestCheckResourceAttr("barracudawaf_self_signed_certificate.demo_self_signed_cert_1", "key_type", "rsa"),
					resource.TestCheckResourceAttr("barracudawaf_self_signed_certificate.demo_self_signed_cert_1", "key_size", "1024"),
					resource.TestCheckResourceAttr("barracudawaf_self_signed_certificate.demo_self_signed_cert_1", "name", "DemoSelfSignedCert1"),
					resource.TestCheckResourceAttr("barracudawaf_self_signed_certificate.demo_self_signed_cert_1", "allow_private_key_export", "Yes"),
				),
			},
		},
	})
}

func testCheckSelfSignedCertificateExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/self-signed-certificate/" + name
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		selfSignedCerts, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if selfSignedCerts == nil {
			return fmt.Errorf("self signed certificate %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range selfSignedCerts.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("self signed certificate (%s) not found on the system", name)
		}

		return nil
	}
}
