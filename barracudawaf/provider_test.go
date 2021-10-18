package barracudawaf

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"barracudawaf": testAccProvider,
	}
}

func TestAccProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAcctPreCheck(t *testing.T) {
	if os.Getenv("BARRACUDA_WAF_IP") != "" && (os.Getenv("BARRACUDA_WAF_USERNAME") != "" && os.Getenv("BARRACUDA_WAF_PASSWORD") != "") {
		return
	} else {
		t.Fatal("BARRACUDA_WAF_IP, BARRACUDA_WAF_USERNAME and BARRACUDA_WAF_PASSWORD are required for tests.")
		return
	}
}
