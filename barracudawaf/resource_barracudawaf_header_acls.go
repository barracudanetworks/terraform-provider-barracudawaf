package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceHeaderAclsParams = map[string][]string{}
)

func resourceCudaWAFHeaderAcls() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFHeaderAclsCreate,
		Read:   resourceCudaWAFHeaderAclsRead,
		Update: resourceCudaWAFHeaderAclsUpdate,
		Delete: resourceCudaWAFHeaderAclsDelete,

		Schema: map[string]*schema.Schema{
			"comments":                {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"max_header_value_length": {Type: schema.TypeString, Optional: true, Description: "Max Header Value Length"},
			"header_name":             {Type: schema.TypeString, Required: true, Description: "Header Name"},
			"blocked_attack_types":    {Type: schema.TypeString, Optional: true, Description: "Blocked Attack Types"},
			"denied_metachars":        {Type: schema.TypeString, Optional: true, Description: "Denied Metacharacters"},
			"mode":                    {Type: schema.TypeString, Optional: true, Description: "Mode"},
			"custom_blocked_attack_types": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom Blocked Attack Types",
			},
			"exception_patterns": {Type: schema.TypeString, Optional: true, Description: "Exception Patterns"},
			"status":             {Type: schema.TypeString, Optional: true, Description: "Status"},
			"name":               {Type: schema.TypeString, Required: true, Description: "Header ACL Name"},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},

		Description: "`barracudawaf_header_acls` manages `Header Acls` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFHeaderAclsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/header-acls"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFHeaderAclsResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFHeaderAclsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFHeaderAclsRead(d, m)
}

func resourceCudaWAFHeaderAclsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/header-acls"
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

func resourceCudaWAFHeaderAclsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/header-acls"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFHeaderAclsResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFHeaderAclsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFHeaderAclsRead(d, m)
}

func resourceCudaWAFHeaderAclsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/header-acls"
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

func hydrateBarracudaWAFHeaderAclsResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"comments":                    d.Get("comments").(string),
		"max-header-value-length":     d.Get("max_header_value_length").(string),
		"header-name":                 d.Get("header_name").(string),
		"blocked-attack-types":        d.Get("blocked_attack_types").(string),
		"denied-metachars":            d.Get("denied_metachars").(string),
		"mode":                        d.Get("mode").(string),
		"custom-blocked-attack-types": d.Get("custom_blocked_attack_types").(string),
		"exception-patterns":          d.Get("exception_patterns").(string),
		"status":                      d.Get("status").(string),
		"name":                        d.Get("name").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFHeaderAclsSubResource(d *schema.ResourceData, name string, endpoint string) error {

	for subResource, subResourceParams := range subResourceHeaderAclsParams {
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
