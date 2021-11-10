package barracudawaf

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceLetsEncryptCertificateParams = map[string][]string{}
)

func resourceCudaWAFLetsEncryptCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFLetsEncryptCertificateCreate,
		Read:   resourceCudaWAFLetsEncryptCertificateRead,
		Update: resourceCudaWAFLetsEncryptCertificateUpdate,
		Delete: resourceCudaWAFLetsEncryptCertificateDelete,

		Schema: map[string]*schema.Schema{
			"allow_private_key_export":   {Type: schema.TypeString, Optional: true, Description: "If set Yes, Private Key gets downloaded along with the certificate"},
			"auto_renew_cert":            {Type: schema.TypeString, Optional: true, Description: "Auto Renew Certificate"},
			"common_name":                {Type: schema.TypeString, Required: true, Description: "Common Name"},
			"multi_cert_trusted_service": {Type: schema.TypeString, Required: true, Description: "Service Name for LetsEncrypt certificate"},
			"schedule_renewal_day":       {Type: schema.TypeString, Optional: true, Description: "Renew Certificate days"},
			"san": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Subject Alternative Names",
			},
			"name": {Type: schema.TypeString, Required: true, Description: "Policy Name"},
		},

		Description: "`barracudawaf_letsencrypt_certificate` manages `Letsencrypt Certificate` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFLetsEncryptCertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/certificates/letsencrypt"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFLetsEncryptCertificateResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFLetsEncryptCertificateSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFLetsEncryptCertificateRead(d, m)
}

func resourceCudaWAFLetsEncryptCertificateRead(d *schema.ResourceData, m interface{}) error {
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

func resourceCudaWAFLetsEncryptCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceCudaWAFLetsEncryptCertificateRead(d, m)
}

func resourceCudaWAFLetsEncryptCertificateDelete(d *schema.ResourceData, m interface{}) error {
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

func hydrateBarracudaWAFLetsEncryptCertificateResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]interface{}{
		"san":                        d.Get("san"),
		"name":                       d.Get("name").(string),
		"common-name":                d.Get("common_name").(string),
		"auto-renew-cert":            d.Get("auto_renew_cert").(string),
		"schedule-renewal-day":       d.Get("schedule_renewal_day").(string),
		"allow-private-key-export":   d.Get("allow_private_key_export").(string),
		"multi-cert-trusted-service": d.Get("multi_cert_trusted_service").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{""}
		for _, param := range updatePayloadExceptions {
			delete(resourcePayload, param)
		}
	}

	// remove empty parameters from resource payload
	for key, val := range resourcePayload {
		if reflect.ValueOf(val).Len() == 0 {
			delete(resourcePayload, key)
		}
	}

	return &APIRequest{
		URL:  endpoint,
		Body: resourcePayload,
	}
}

func (b *BarracudaWAF) hydrateBarracudaWAFLetsEncryptCertificateSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceLetsEncryptCertificateParams {
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
