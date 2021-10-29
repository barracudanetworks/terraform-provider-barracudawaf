package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TRUSTED_CA_CERT_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_trusted_ca_certificate" "demo_trusted_ca_cert_1" {
	name        = "DemoTrustedCACert1"
	certificate = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUYzakNDQThhZ0F3SUJBZ0lRQWYxdE1QeWp5\nbEdvRzd4a0RqVURMVEFOQmdrcWhraUc5dzBCQVF3RkFEQ0IKaURFTE1Ba0dBMVVFQmhNQ1ZWTXhF\nekFSQmdOVkJBZ1RDazVsZHlCS1pYSnpaWGt4RkRBU0JnTlZCQWNUQzBwbApjbk5sZVNCRGFYUjVN\nUjR3SEFZRFZRUUtFeFZVYUdVZ1ZWTkZVbFJTVlZOVUlFNWxkSGR2Y21zeExqQXNCZ05WCkJBTVRK\nVlZUUlZKVWNuVnpkQ0JTVTBFZ1EyVnlkR2xtYVdOaGRHbHZiaUJCZFhSb2IzSnBkSGt3SGhjTk1U\nQXcKTWpBeE1EQXdNREF3V2hjTk16Z3dNVEU0TWpNMU9UVTVXakNCaURFTE1Ba0dBMVVFQmhNQ1ZW\nTXhFekFSQmdOVgpCQWdUQ2s1bGR5QktaWEp6WlhreEZEQVNCZ05WQkFjVEMwcGxjbk5sZVNCRGFY\nUjVNUjR3SEFZRFZRUUtFeFZVCmFHVWdWVk5GVWxSU1ZWTlVJRTVsZEhkdmNtc3hMakFzQmdOVkJB\nTVRKVlZUUlZKVWNuVnpkQ0JTVTBFZ1EyVnkKZEdsbWFXTmhkR2x2YmlCQmRYUm9iM0pwZEhrd2dn\nSWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUNEd0F3Z2dJSwpBb0lDQVFDQUVtVVhOZzdEMndpejBL\neFhEWGJ0elNmVFRLMVFnMkhpcWlCTkNTMWtDZHpPaVovTVBhbnM5cy9CCjNQSFRzZFo3TnlnUksw\nZmFPY2E4T2htMFg2YTlmWjJqWTBLMmR2S3BPeXVSK09KdjBPd1dJSkFKUHVMb2RNa1kKdEpIVVlt\nVGJmNk1HOFlnWWFwQWlQTHorRS9DSEZIdjI1QitPMU9SUnhoRm5SZ2hSeTRZVVZEKzhNLzUrYkp6\nLwpGcDBZdlZHT05hYW5ac2h5WjlzaFpySFVtM2dEd0ZBNjZNenczTHllVFA2dkJaWTFIMWRhdC8v\nTytUMjNMTGIyClZOM0k1eEk2VGE1TWlyZGNtclMzSUQzS2Z5STBybjQ3YUdZQlJPY0JUa1pUbXpO\nZzk1UytVemVRYzBQek1zTlQKNzl1cS9uUk9hY2RyakdDVDNzVEhETi9oTXE3TWt6dFJlSlZuaSs0\nOVZ2NE0wR2tQR3cvekpTWnJNMjMzYmtmNgpjMFBsZmc2bFpyRXBmREtFWTFXSnhBM0JrMVF3R1JP\nczAzMDNwK3RkT213MVhOdEIxeExhcVVrTDM5aUFpZ21UCllvNjFaczhsaU0yRXVMRS9wRGtQMlFL\nZTZ4Sk1sWHp6YXdXcFhoYUR6TGhuNHVnVG5jeGJndE5NcysxYi85N2wKYzZ3ak95MEF2elZWZEFs\nSjJFbFlHbitTTnVaUmtnN3pKbjBjVFJlOHlleERKdEMvUVY5QXFVUkU5Sm5uVjRlZQpVQjlYVktn\nKy9YUmpMN0ZRWlFubVdFSXVReHBNdFBBbFIxbjZCQjZUMUNaR1NsQ0JzdDYrZUxmOFp4WGh5VmVF\nCkhnOWoxdWxpdXRaZlZTN3FYTVlvQ0FRbE9iZ09LNm55VEpjY0J6OE5Vdlh0N3krQ0R3SURBUUFC\nbzBJd1FEQWQKQmdOVkhRNEVGZ1FVVTNtL1dxb3JTczlVZ09IWW04Q2Q4cklEWnNzd0RnWURWUjBQ\nQVFIL0JBUURBZ0VHTUE4RwpBMVVkRXdFQi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRU1CUUFE\nZ2dJQkFGelVmQTNQOXdGOVFabGxESFBGClVwL0wrTStaQm44YjJrTVZuNTRDVlZlV0ZQRlNQQ2VI\nbENqdEh6b0JONkoyL0ZOUXdJU2J4bXRPdW93aFQ2S08KVldLUjgya1YyTHlJNDhTcUMvM3ZxT2xM\nVlNvR0lHMVZlQ2taN2w4d1hFc2tFVlgvSkpwdVhpb3I3Z3RObjMvMwpBVGlVRkpWREJ3bjdZS251\nSEtzU2pLQ2FYcWVZYWxsdGl6OEkrOGpSUmE4WUZXU1FFZzl6S0M3RjRpUk8vRmpzCjhQUkYvaUt6\nNnkrTzB0bEZZUVhCbDIrb2RuS1BpNHcycjc4TkJjNXhqZWFtYng5c3BuRml4ZGpRZzNJTThXY1IK\naVF5Y0UweHlOTis4MVhIZnFuSGQ0YmxzakR3U1hXWGF2VmNTdGtOci8rWGVUV1lSVWMrWnJ1d1h0\ndWh4a1l6ZQpTZjdkTlhHaUZTZVVITTloNHlhN2I2Tm5KU0ZkNXQwZEN5NW9HenVDcit5RFo0WFVt\nRkYwc2JtWmdJbi9mM2daClhIbEtZQzZTUUs1TU55b3N5Y2RpeUE1ZDl6WmJ5dUFsSlFHMDNSb0hu\nSGNBUDlEYzFldzkxUHE3UDh5RjFtOS8KcVMzZnVRTDM5WmVhdFRYYXcyZXdoMHFwS0o0amp2OWNK\nMnZoc0UvekIrNEFMdFJaaDh0U1FaWHE5RWZYN21SQgpWWHlOV1FLVjNXS2R3cm51V2loMGhLV2J0\nNURIREFmZjlZazJkRExXS01Hd3NBdmduRXpESE5iODQybTFSMGFCCkw2S0NxOU5qUkhERWpmOHRN\nN3F0ajN1MWNJaXVQaG5QUUNqWS9NaVF1MTJaSXZWUzVsakZINGd4USs2SUhkZkcKamp4RGFoMm5H\nTjU5UFJieFl2bktrS2o5Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K\n"
}
`

func TestAccBarracudaWAFTrustedCACertificate_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TRUSTED_CA_CERT_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckTrustedCACertificateExists("DemoTrustedCACert1"),
					resource.TestCheckResourceAttr("barracudawaf_trusted_ca_certificate.demo_trusted_ca_cert_1", "name", "DemoTrustedCACert1"),
				),
			},
		},
	})
}

func testCheckTrustedCACertificateExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/trusted-ca-certificate/" + name
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		trustedCACerts, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if trustedCACerts == nil {
			return fmt.Errorf("trusted ca certificate %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range trustedCACerts.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("trusted ca certificate (%s) not found on the system", name)
		}

		return nil
	}
}
