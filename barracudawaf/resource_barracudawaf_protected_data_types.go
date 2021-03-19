package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceProtectedDataTypesParams = map[string][]string{}
)

func resourceCudaWAFProtectedDataTypes() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFProtectedDataTypesCreate,
		Read:   resourceCudaWAFProtectedDataTypesRead,
		Update: resourceCudaWAFProtectedDataTypesUpdate,
		Delete: resourceCudaWAFProtectedDataTypesDelete,

		Schema: map[string]*schema.Schema{
			"action": {Type: schema.TypeString, Optional: true, Description: "Action"},
			"initial_characters_to_keep": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Initial Characters to Keep",
			},
			"trailing_characters_to_keep": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Trailing Characters to Keep",
			},
			"name": {Type: schema.TypeString, Required: true, Description: "Data Theft Element Name"},
			"custom_identity_theft_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom Identity Theft Type",
			},
			"enable":              {Type: schema.TypeString, Optional: true, Description: "Enabled"},
			"identity_theft_type": {Type: schema.TypeString, Optional: true, Description: "Identity Theft Type"},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},

		Description: "`barracudawaf_protected_data_types` manages `Protected Data Types` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFProtectedDataTypesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/security-policies/" + d.Get("parent.0").(string) + "/protected-data-types"
	err := client.CreateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFProtectedDataTypesResource(d, "post", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFProtectedDataTypesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFProtectedDataTypesRead(d, m)
}

func resourceCudaWAFProtectedDataTypesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/security-policies/" + d.Get("parent.0").(string) + "/protected-data-types"
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

func resourceCudaWAFProtectedDataTypesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/security-policies/" + d.Get("parent.0").(string) + "/protected-data-types"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFProtectedDataTypesResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFProtectedDataTypesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFProtectedDataTypesRead(d, m)
}

func resourceCudaWAFProtectedDataTypesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/security-policies/" + d.Get("parent.0").(string) + "/protected-data-types"
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

func hydrateBarracudaWAFProtectedDataTypesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"action":                      d.Get("action").(string),
		"initial-characters-to-keep":  d.Get("initial_characters_to_keep").(string),
		"trailing-characters-to-keep": d.Get("trailing_characters_to_keep").(string),
		"name":                        d.Get("name").(string),
		"custom-identity-theft-type":  d.Get("custom_identity_theft_type").(string),
		"enable":                      d.Get("enable").(string),
		"identity-theft-type":         d.Get("identity_theft_type").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFProtectedDataTypesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceProtectedDataTypesParams {
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
