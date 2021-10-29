package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var SIGNED_CERT_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_signed_certificate" "demo_signed_cert" {
    name                     = "DemoSignedCert"
    signed_certificate       = "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"
    certificate_type         = "PKCS12 Token"
    certificate_password     = "secret@123"
    allow_private_key_export = "Yes"
    key_type                 = "RSA"
}
`

func TestAccBarracudaWAFSignedCertificate_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: SIGNED_CERT_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckSignedCertificateExists("DemoSignedCert"),
					resource.TestCheckResourceAttr("barracudawaf_signed_certificate.demo_signed_cert", "name", "DemoSignedCert"),
					resource.TestCheckResourceAttr("barracudawaf_signed_certificate.demo_signed_cert", "certificate_type", "PKCS12 Token"),
					resource.TestCheckResourceAttr("barracudawaf_signed_certificate.demo_signed_cert", "allow_private_key_export", "Yes"),
				),
			},
		},
	})
}

func testCheckSignedCertificateExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/signed-certificate/" + name
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		signedCerts, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if signedCerts == nil {
			return fmt.Errorf("signed certificate %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range signedCerts.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("signed certificate (%s) not found on the system", name)
		}

		return nil
	}
}
