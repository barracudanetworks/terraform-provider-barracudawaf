package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TRUSTED_SERVER_CERT_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_trusted_server_certificate" "demo_trusted_server_cert_1" {
	name        = "DemoTrustedServerCert1"
	certificate = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUdFekNDQS91Z0F3SUJBZ0lRZlZ0UkpyUjJ1\naEhiZEJZTHZGTU5wekFOQmdrcWhraUc5dzBCQVF3RkFEQ0IKaURFTE1Ba0dBMVVFQmhNQ1ZWTXhF\nekFSQmdOVkJBZ1RDazVsZHlCS1pYSnpaWGt4RkRBU0JnTlZCQWNUQzBwbApjbk5sZVNCRGFYUjVN\nUjR3SEFZRFZRUUtFeFZVYUdVZ1ZWTkZVbFJTVlZOVUlFNWxkSGR2Y21zeExqQXNCZ05WCkJBTVRK\nVlZUUlZKVWNuVnpkQ0JTVTBFZ1EyVnlkR2xtYVdOaGRHbHZiaUJCZFhSb2IzSnBkSGt3SGhjTk1U\nZ3gKTVRBeU1EQXdNREF3V2hjTk16QXhNak14TWpNMU9UVTVXakNCanpFTE1Ba0dBMVVFQmhNQ1Iw\nSXhHekFaQmdOVgpCQWdURWtkeVpXRjBaWElnVFdGdVkyaGxjM1JsY2pFUU1BNEdBMVVFQnhNSFUy\nRnNabTl5WkRFWU1CWUdBMVVFCkNoTVBVMlZqZEdsbmJ5Qk1hVzFwZEdWa01UY3dOUVlEVlFRREV5\nNVRaV04wYVdkdklGSlRRU0JFYjIxaGFXNGcKVm1Gc2FXUmhkR2x2YmlCVFpXTjFjbVVnVTJWeWRt\nVnlJRU5CTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQwpBUThBTUlJQkNnS0NBUUVBMW5NejF0\nYzhJTkFBMGhkRnVOWStCNkkveDBIdU1qREpzR3o5OUovTEVwZ1BMVCtOClRRRU1nZzhYZjJJdTZi\naEllZnNXZzA2dDF6SWxrN2NIdjdsUVA2bE13MEFxNlRuLzJZSEtIeFl5UWRxQUpya2oKZW9jZ0h1\nUC9JSm84bFVSdmgzVUdrRUMwTXBNV0NSQUlJejdTM1ljUGIxMVJGR29LYWNWUEFYSnB6OU9UVEcw\nRQpvS01iZ242eG1ybnR4WjdGTjNpZm1nZzArMVl1V01RSkRnWmtXN3czM1BHZktHaW9WckNTbzF5\nZnU0aVlDQnNrCkhhc3doYTZ2c0M2ZWVwM0J3RUljNGdMdzZ1QkswdStRRHJUQlFCYndiNFZDU21U\nM3BEQ2cvcjh1b3lkYWpvdFkKdUszREdSZUVZKzF2VnYyRHkyQTB4SFMrNXAzYjRlVGx5Z3hmRlFJ\nREFRQUJvNElCYmpDQ0FXb3dId1lEVlIwagpCQmd3Rm9BVVUzbS9XcW9yU3M5VWdPSFltOENkOHJJ\nRFpzc3dIUVlEVlIwT0JCWUVGSTJNWHNSVXJZcmhkK21iCitac0Y0YmdCaldIaE1BNEdBMVVkRHdF\nQi93UUVBd0lCaGpBU0JnTlZIUk1CQWY4RUNEQUdBUUgvQWdFQU1CMEcKQTFVZEpRUVdNQlFHQ0Nz\nR0FRVUZCd01CQmdnckJnRUZCUWNEQWpBYkJnTlZIU0FFRkRBU01BWUdCRlVkSUFBdwpDQVlHWjRF\nTUFRSUJNRkFHQTFVZEh3UkpNRWN3UmFCRG9FR0dQMmgwZEhBNkx5OWpjbXd1ZFhObGNuUnlkWE4w\nCkxtTnZiUzlWVTBWU1ZISjFjM1JTVTBGRFpYSjBhV1pwWTJGMGFXOXVRWFYwYUc5eWFYUjVMbU55\nYkRCMkJnZ3IKQmdFRkJRY0JBUVJxTUdnd1B3WUlLd1lCQlFVSE1BS0dNMmgwZEhBNkx5OWpjblF1\nZFhObGNuUnlkWE4wTG1OdgpiUzlWVTBWU1ZISjFjM1JTVTBGQlpHUlVjblZ6ZEVOQkxtTnlkREFs\nQmdnckJnRUZCUWN3QVlZWmFIUjBjRG92CkwyOWpjM0F1ZFhObGNuUnlkWE4wTG1OdmJUQU5CZ2tx\naGtpRzl3MEJBUXdGQUFPQ0FnRUFNcjlodlE1SXcwL0gKdWtkTitKeDRHUUhjRXgyQWIvekRjTFJT\nbWpFem1sZFMrekdlYTZUdlZLcUpqVUFYYVBnUkVIelN5ckh4VlliSAo3ck0ya1liMk9WRy9ScjhQ\nb0xxMDkzNUp4Q28yRjU3a2FEbDZyNVJPVm0reWV6dS9Db2E5emNWM0hBTzRPTEdpCkgxOSsyNHJj\nUmtpMmFBclBzclcwNGpUa1o2azRaZ2xlMHJqOG5TZzZGMEFud25KT0tmMGhQSHpQRS91V0xNVXgK\nUlAwVDdkV2JxV2xvZDN6dTRmK2srVFk0Q0ZNNW9vUTBuQm56dmc2czFTUTM2eU9vZU5EVDUrK1NS\nMlJpT1NMdgp4dmNSdmlLRnhtWkVKQ2FPRURLTnlKT3VCNTZEUGkvWitmVkdqbU8rd2VhMDNLYk5J\nYWlHQ3BYWkxvVW1HdjM4CnNiWlhRbTJWMFRQMk9SUUdna0U0OVk5WTNJQmJwTlY5bFhqOXA1di8v\nY1dvYWFzbTU2ZWtCWWRicWJlNG95QUwKbDZsRmhkMnppK1dKTjQ0cERmd0dGL1k0UUE1QzVCSUcr\nM3Z6eGhGb1l0L2ptUFFUMkJWUGk3RnAyUkJndkdRcQo2akczNUxXak9oU2JKdU1MZS8wQ2pyYVp3\nVGlYV1RiMnFIU2loclplNjhaazZzK2dvL2x1bnJvdEViYUdtQWhZCkxjbXNKV1R5WG5XME9NR3Vm\nMXBHZytwUnlyYnhtUkUxYTZWcWU4WUFzT2Y0dm1TeXJjakM4YXpqVWVxa2srQjUKeU9HQlFNa0tX\nK0VTUE1GZ0t1T1h3SWxDeXBUUFJwZ1NhYnVZME1MVERYSkxSMjdsazhReUtHT0hRK1N3TWo0Swow\nMHUvSTVzVUtVRXJtZ1Fma3kzeHh6bElQSzFhRW44PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t\n"
}
`

func TestAccBarracudaWAFTrustedServerCertificate_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TRUSTED_SERVER_CERT_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckTrustedServerCertificateExists("DemoTrustedServerCert1"),
					resource.TestCheckResourceAttr("barracudawaf_trusted_server_certificate.demo_trusted_server_cert_1", "name", "DemoTrustedServerCert1"),
				),
			},
		},
	})
}

func testCheckTrustedServerCertificateExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/trusted-server-certificate/" + name
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		trustedServerCerts, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if trustedServerCerts == nil {
			return fmt.Errorf("trusted server certificate %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range trustedServerCerts.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("trusted server certificate (%s) not found on the system", name)
		}

		return nil
	}
}
