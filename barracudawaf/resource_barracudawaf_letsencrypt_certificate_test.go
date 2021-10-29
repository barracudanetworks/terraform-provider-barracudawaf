package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var LETSENCRYPT_CERT_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_services" "demo_app_1" {
    name            = "DemoApp1"
    ip_address      = "172.31.50.142"
    port            = "80"
    type            = "HTTP"
    vsite           = "default"
    address_version = "IPv4"
    status          = "On"
    group           = "default"
    comments        = "Demo Service with Terraform"

    basic_security {
      mode = "Active"
    }
}

resource "barracudawaf_letsencrypt_certificate" "demo_letsencrypt_cert" {
    name                       = "DemoLetsEncryptCert"
    common_name                = "app.cudawaf.net"
    allow_private_key_export   = "Yes"
    auto_renew_cert            = "Yes"
    schedule_renewal_day       = "60"

    multi_cert_trusted_service = barracudawaf_services.demo_app_1.name
    depends_on = [ barracudawaf_services.demo_app_1 ]
}
`

func TestAccBarracudaWAFLetsEncryptCertificate_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: LETSENCRYPT_CERT_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckLetsEncryptCertificateExists("DemoLetsEncryptCert"),
					resource.TestCheckResourceAttr("barracudawaf_letsencrypt_certificate.demo_letsencrypt_cert", "auto_renew_cert", "Yes"),
					resource.TestCheckResourceAttr("barracudawaf_letsencrypt_certificate.demo_letsencrypt_cert", "schedule_renewal_day", "60"),
					resource.TestCheckResourceAttr("barracudawaf_letsencrypt_certificate.demo_letsencrypt_cert", "name", "DemoLetsEncryptCert"),
					resource.TestCheckResourceAttr("barracudawaf_letsencrypt_certificate.demo_letsencrypt_cert", "allow_private_key_export", "Yes"),
				),
			},
		},
	})
}

func testCheckLetsEncryptCertificateExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/signed-certificate/" + name
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		letsEncryptCerts, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if letsEncryptCerts == nil {
			return fmt.Errorf("letsencrypt certificate %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range letsEncryptCerts.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("letsencrypt certificate (%s) not found on the system", name)
		}

		return nil
	}
}
