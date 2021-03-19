package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceSelfSignedCertificateParams = map[string][]string{}
)

func resourceCudaWAFSelfSignedCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFSelfSignedCertificateCreate,
		Read:   resourceCudaWAFSelfSignedCertificateRead,
		Update: resourceCudaWAFSelfSignedCertificateUpdate,
		Delete: resourceCudaWAFSelfSignedCertificateDelete,

		Schema: map[string]*schema.Schema{
			"download_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A Certificate Signing Request (CSR) and/or Certificate can be downloaded.",
			},
			"encrypt_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Encryption Password is used to extract the private key from PKCS #12 token.",
			},
			"city":                {Type: schema.TypeString, Optional: true, Description: "Locality Name"},
			"common_name":         {Type: schema.TypeString, Required: true, Description: "Common Name"},
			"country_code":        {Type: schema.TypeString, Required: true, Description: "Country"},
			"elliptic_curve_name": {Type: schema.TypeString, Optional: true, Description: "Elliptic Curve Name"},
			"expiry":              {Type: schema.TypeString, Optional: true},
			"key_size":            {Type: schema.TypeString, Optional: true, Description: "Key Size"},
			"key_type":            {Type: schema.TypeString, Optional: true, Description: "Select Key Type:"},
			"allow_private_key_export": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If set to <b>Yes</b>, the Private Key gets downloaded along with the certificate.",
			},
			"name":                {Type: schema.TypeString, Required: true, Description: "None"},
			"organization_name":   {Type: schema.TypeString, Optional: true, Description: "Organization Name"},
			"organizational_unit": {Type: schema.TypeString, Optional: true, Description: "Organizational Unit Name"},
			"san_certificate":     {Type: schema.TypeString, Optional: true, Description: "None"},
			"serial":              {Type: schema.TypeString, Optional: true},
			"state":               {Type: schema.TypeString, Optional: true, Description: "State or Province"},
		},

		Description: "`barracudawaf_self_signed_certificate` manages `Self Signed Certificate` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFSelfSignedCertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/self-signed-certificate"
	err := client.CreateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFSelfSignedCertificateResource(d, "post", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFSelfSignedCertificateSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFSelfSignedCertificateRead(d, m)
}

func resourceCudaWAFSelfSignedCertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/self-signed-certificate"
	request := &APIRequest{
		Method: "get",
		URL:    resourceEndpoint,
	}

	var dataItems map[string]interface{}
	resources, err := client.GetBarracudaWAFResource(name, request)

	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	if resources.Data == nil {
		log.Printf("[WARN] Barracuda WAF resource (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	for _, dataItems = range resources.Data {
		if dataItems["name"] == name {
			break
		}
	}

	if dataItems["name"] != name {
		return fmt.Errorf("Barracuda WAF resource (%s) not found on the system", name)
	}

	d.Set("name", name)
	return nil
}

func resourceCudaWAFSelfSignedCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/self-signed-certificate"
	err := client.UpdateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFSelfSignedCertificateResource(d, "put", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFSelfSignedCertificateSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFSelfSignedCertificateRead(d, m)
}

func resourceCudaWAFSelfSignedCertificateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/self-signed-certificate"
	request := &APIRequest{
		Method: "delete",
		URL:    resourceEndpoint,
	}

	err := client.DeleteBarracudaWAFResource(name, request)

	if err != nil {
		return fmt.Errorf("Unable to delete the Barracuda WAF resource (%s) (%v)", name, err)
	}

	return nil
}

func hydrateBarracudaWAFSelfSignedCertificateResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"download-type":            d.Get("download_type").(string),
		"encrypt-password":         d.Get("encrypt_password").(string),
		"city":                     d.Get("city").(string),
		"common-name":              d.Get("common_name").(string),
		"country-code":             d.Get("country_code").(string),
		"elliptic-curve-name":      d.Get("elliptic_curve_name").(string),
		"expiry":                   d.Get("expiry").(string),
		"key-size":                 d.Get("key_size").(string),
		"key-type":                 d.Get("key_type").(string),
		"allow-private-key-export": d.Get("allow_private_key_export").(string),
		"name":                     d.Get("name").(string),
		"organization-name":        d.Get("organization_name").(string),
		"organizational-unit":      d.Get("organizational_unit").(string),
		"san-certificate":          d.Get("san_certificate").(string),
		"serial":                   d.Get("serial").(string),
		"state":                    d.Get("state").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{
			"city",
			"common-name",
			"country-code",
			"elliptic-curve-name",
			"expiry",
			"key-size",
			"key-type",
			"allow-private-key-export",
			"name",
			"organization-name",
			"organizational-unit",
			"san-certificate",
			"serial",
			"state",
		}
		for _, param := range updatePayloadExceptions {
			delete(resourcePayload, param)
		}
	}

	// remove empty parameters from resource payload
	for key, val := range resourcePayload {
		if len(val) == 0 {
			delete(resourcePayload, key)
		}
	}

	return &APIRequest{
		URL:  endpoint,
		Body: resourcePayload,
	}
}

func (b *BarracudaWAF) hydrateBarracudaWAFSelfSignedCertificateSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceSelfSignedCertificateParams {
		subResourceParamsLength := d.Get(subResource + ".#").(int)

		log.Printf("[INFO] Updating Barracuda WAF sub resource (%s) (%s)", name, subResource)

		for i := 0; i < subResourceParamsLength; i++ {
			subResourcePayload := map[string]string{}
			suffix := fmt.Sprintf(".%d", i)

			for _, param := range subResourceParams {
				paramSuffix := fmt.Sprintf(".%s", param)
				paramVaule := d.Get(subResource + suffix + paramSuffix).(string)

				if len(paramVaule) > 0 {
					param = strings.Replace(param, "_", "-", -1)
					subResourcePayload[param] = paramVaule
				}
			}

			err := b.UpdateBarracudaWAFSubResource(name, endpoint, &APIRequest{
				URL:  strings.Replace(subResource, "_", "-", -1),
				Body: subResourcePayload,
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}
