package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceClientCertificateCrlsParams = map[string][]string{}
)

func resourceCudaWAFClientCertificateCrls() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFClientCertificateCrlsCreate,
		Read:   resourceCudaWAFClientCertificateCrlsRead,
		Update: resourceCudaWAFClientCertificateCrlsUpdate,
		Delete: resourceCudaWAFClientCertificateCrlsDelete,

		Schema: map[string]*schema.Schema{
			"auto_update_type":  {Type: schema.TypeString, Optional: true, Description: "Auto Update Type"},
			"date_of_month":     {Type: schema.TypeString, Optional: true, Description: "Date Of Month"},
			"day_of_week":       {Type: schema.TypeString, Optional: true, Description: "Day Of Week"},
			"time_of_day":       {Type: schema.TypeString, Optional: true, Description: "Time Of Day"},
			"auto_update":       {Type: schema.TypeString, Optional: true, Description: "CRL Auto Update"},
			"name":              {Type: schema.TypeString, Required: true, Description: "CRL Name"},
			"number_of_retries": {Type: schema.TypeString, Optional: true, Description: "Number of Retries"},
			"url":               {Type: schema.TypeString, Required: true, Description: "CRL URL"},
			"enable":            {Type: schema.TypeString, Optional: true, Description: "Enable CRL"},
			"parent":            {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}, Required: true},
		},

		Description: "`barracudawaf_client_certificate_crls` manages `Client Certificate Crls` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFClientCertificateCrlsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/client-certificate-crls"
	err := client.CreateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFClientCertificateCrlsResource(d, "post", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFClientCertificateCrlsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFClientCertificateCrlsRead(d, m)
}

func resourceCudaWAFClientCertificateCrlsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/client-certificate-crls"
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

func resourceCudaWAFClientCertificateCrlsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/client-certificate-crls"
	err := client.UpdateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFClientCertificateCrlsResource(d, "put", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFClientCertificateCrlsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFClientCertificateCrlsRead(d, m)
}

func resourceCudaWAFClientCertificateCrlsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/client-certificate-crls"
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

func hydrateBarracudaWAFClientCertificateCrlsResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"auto-update-type":  d.Get("auto_update_type").(string),
		"date-of-month":     d.Get("date_of_month").(string),
		"day-of-week":       d.Get("day_of_week").(string),
		"time-of-day":       d.Get("time_of_day").(string),
		"auto-update":       d.Get("auto_update").(string),
		"name":              d.Get("name").(string),
		"number-of-retries": d.Get("number_of_retries").(string),
		"url":               d.Get("url").(string),
		"enable":            d.Get("enable").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{}
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

func (b *BarracudaWAF) hydrateBarracudaWAFClientCertificateCrlsSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceClientCertificateCrlsParams {
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
