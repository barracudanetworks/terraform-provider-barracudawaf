package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceTrustedServerCertificateParams = map[string][]string{}
)

func resourceCudaWAFTrustedServerCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFTrustedServerCertificateCreate,
		Read:   resourceCudaWAFTrustedServerCertificateRead,
		Update: resourceCudaWAFTrustedServerCertificateUpdate,
		Delete: resourceCudaWAFTrustedServerCertificateDelete,

		Schema: map[string]*schema.Schema{
			"common_name": {Type: schema.TypeString, Optional: true, Description: "Common Name"},
			"expiry":      {Type: schema.TypeString, Optional: true},
			"serial":      {Type: schema.TypeString, Optional: true},
			"name":        {Type: schema.TypeString, Optional: true, Description: "Certificate Name"},
			"certificate": {Type: schema.TypeString, Optional: true},
			"download_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A Certificate Signing Request (CSR) and/or Certificate can be downloaded.",
			},
		},

		Description: "`barracudawaf_trusted_server_certificate` manages `Trusted Server Certificate` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFTrustedServerCertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/trusted-server-certificate"
	err := client.CreateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFTrustedServerCertificateResource(d, "post", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFTrustedServerCertificateSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFTrustedServerCertificateRead(d, m)
}

func resourceCudaWAFTrustedServerCertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/trusted-server-certificate"
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

func resourceCudaWAFTrustedServerCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/trusted-server-certificate"
	err := client.UpdateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFTrustedServerCertificateResource(d, "put", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFTrustedServerCertificateSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFTrustedServerCertificateRead(d, m)
}

func resourceCudaWAFTrustedServerCertificateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/trusted-server-certificate"
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

func hydrateBarracudaWAFTrustedServerCertificateResource(
	d *schema.ResourceData,
	method string,
	endpoint string,
) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"common-name":   d.Get("common_name").(string),
		"expiry":        d.Get("expiry").(string),
		"serial":        d.Get("serial").(string),
		"name":          d.Get("name").(string),
		"certificate":   d.Get("certificate").(string),
		"download-type": d.Get("download_type").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"common-name", "expiry", "serial", "name", "certificate"}
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

func (b *BarracudaWAF) hydrateBarracudaWAFTrustedServerCertificateSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceTrustedServerCertificateParams {
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
