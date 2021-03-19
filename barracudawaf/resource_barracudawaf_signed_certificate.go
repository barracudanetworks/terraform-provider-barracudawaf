package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceSignedCertificateParams = map[string][]string{
		"ocsp_stapling": {
			"cache_timeout",
			"clock_skew",
			"error_timeout",
			"issuer_certificate",
			"ocsp_stapling",
			"override_ocsp_responder",
			"ocsp_issuer_certificate",
		},
	}
)

func resourceCudaWAFSignedCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFSignedCertificateCreate,
		Read:   resourceCudaWAFSignedCertificateRead,
		Update: resourceCudaWAFSignedCertificateUpdate,
		Delete: resourceCudaWAFSignedCertificateDelete,

		Schema: map[string]*schema.Schema{
			"assign_associated_key": {Type: schema.TypeString, Optional: true},
			"signed_certificate":    {Type: schema.TypeString, Optional: true},
			"certificate_key":       {Type: schema.TypeString, Optional: true},
			"certificate_password":  {Type: schema.TypeString, Optional: true},
			"certificate_type":      {Type: schema.TypeString, Optional: true},
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
			"intermediary_certificates": {Type: schema.TypeString, Optional: true},
			"name":                      {Type: schema.TypeString, Optional: true, Description: "Certificate Name"},
			"auto_renew_cert":           {Type: schema.TypeString, Optional: true, Description: "None"},
			"common_name":               {Type: schema.TypeString, Optional: true, Description: "Common Name"},
			"expiry":                    {Type: schema.TypeString, Optional: true},
			"key_type":                  {Type: schema.TypeString, Optional: true, Description: "Select Key Type:"},
			"allow_private_key_export":  {Type: schema.TypeString, Optional: true},
			"schedule_renewal_day":      {Type: schema.TypeString, Optional: true, Description: "None"},
			"serial":                    {Type: schema.TypeString, Optional: true},
			"ocsp_stapling": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cache_timeout": {Type: schema.TypeString, Required: true, Description: "Cache timeout"},
						"clock_skew":    {Type: schema.TypeString, Optional: true, Description: "Clock Skew"},
						"error_timeout": {Type: schema.TypeString, Optional: true, Description: "Error Timeout"},
						"issuer_certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "OCSP Issuer Cetificate",
						},
						"ocsp_stapling": {Type: schema.TypeString, Optional: true, Description: "OCSP Stapling"},
						"override_ocsp_responder": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "OCSP Responder URL",
						},
						"ocsp_issuer_certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Ocsp Issuer Certificate content as a Base64 encoded string.",
						},
					},
				},
			},
		},

		Description: "`barracudawaf_signed_certificate` manages `Signed Certificate` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFSignedCertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/signed-certificate"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFSignedCertificateResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFSignedCertificateSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFSignedCertificateRead(d, m)
}

func resourceCudaWAFSignedCertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/signed-certificate"
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

func resourceCudaWAFSignedCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/signed-certificate"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFSignedCertificateResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFSignedCertificateSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFSignedCertificateRead(d, m)
}

func resourceCudaWAFSignedCertificateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/signed-certificate"
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

func hydrateBarracudaWAFSignedCertificateResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"assign-associated-key":     d.Get("assign_associated_key").(string),
		"signed-certificate":        d.Get("signed_certificate").(string),
		"certificate-key":           d.Get("certificate_key").(string),
		"certificate-password":      d.Get("certificate_password").(string),
		"certificate-type":          d.Get("certificate_type").(string),
		"download-type":             d.Get("download_type").(string),
		"encrypt-password":          d.Get("encrypt_password").(string),
		"intermediary-certificates": d.Get("intermediary_certificates").(string),
		"name":                      d.Get("name").(string),
		"auto-renew-cert":           d.Get("auto_renew_cert").(string),
		"common-name":               d.Get("common_name").(string),
		"expiry":                    d.Get("expiry").(string),
		"key-type":                  d.Get("key_type").(string),
		"allow-private-key-export":  d.Get("allow_private_key_export").(string),
		"schedule-renewal-day":      d.Get("schedule_renewal_day").(string),
		"serial":                    d.Get("serial").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{
			"assign-associated-key",
			"signed-certificate",
			"certificate-key",
			"certificate-password",
			"certificate-type",
			"intermediary-certificates",
			"name",
			"common-name",
			"expiry",
			"key-type",
			"allow-private-key-export",
			"serial",
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

func (b *BarracudaWAF) hydrateBarracudaWAFSignedCertificateSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceSignedCertificateParams {
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
