package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceJsonSecurityPoliciesParams = map[string][]string{}
)

func resourceCudaWAFJsonSecurityPolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFJsonSecurityPoliciesCreate,
		Read:   resourceCudaWAFJsonSecurityPoliciesRead,
		Update: resourceCudaWAFJsonSecurityPoliciesUpdate,
		Delete: resourceCudaWAFJsonSecurityPoliciesDelete,

		Schema: map[string]*schema.Schema{
			"name":               {Type: schema.TypeString, Required: true, Description: "Policy Name"},
			"max_array_elements": {Type: schema.TypeString, Optional: true, Description: "Max Array Elements"},
			"max_siblings":       {Type: schema.TypeString, Optional: true, Description: "Max Siblings"},
			"max_keys":           {Type: schema.TypeString, Required: true, Description: "Max Keys"},
			"max_key_length":     {Type: schema.TypeString, Required: true, Description: "Max Key Length"},
			"max_number_value":   {Type: schema.TypeString, Optional: true, Description: "Max Number Value"},
			"max_object_depth":   {Type: schema.TypeString, Optional: true, Description: "Max Object Depth"},
			"max_value_length":   {Type: schema.TypeString, Required: true, Description: "Max Value Length"},
		},

		Description: "`barracudawaf_json_security_policies` manages `Json Security Policies` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFJsonSecurityPoliciesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/json-security-policies"
	err := client.CreateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFJsonSecurityPoliciesResource(d, "post", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFJsonSecurityPoliciesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFJsonSecurityPoliciesRead(d, m)
}

func resourceCudaWAFJsonSecurityPoliciesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/json-security-policies"
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

func resourceCudaWAFJsonSecurityPoliciesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/json-security-policies"
	err := client.UpdateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFJsonSecurityPoliciesResource(d, "put", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFJsonSecurityPoliciesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFJsonSecurityPoliciesRead(d, m)
}

func resourceCudaWAFJsonSecurityPoliciesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/json-security-policies"
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

func hydrateBarracudaWAFJsonSecurityPoliciesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"name":               d.Get("name").(string),
		"max-array-elements": d.Get("max_array_elements").(string),
		"max-siblings":       d.Get("max_siblings").(string),
		"max-keys":           d.Get("max_keys").(string),
		"max-key-length":     d.Get("max_key_length").(string),
		"max-number-value":   d.Get("max_number_value").(string),
		"max-object-depth":   d.Get("max_object_depth").(string),
		"max-value-length":   d.Get("max_value_length").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFJsonSecurityPoliciesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceJsonSecurityPoliciesParams {
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
