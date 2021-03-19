package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceUrlAclsParams = map[string][]string{}
)

func resourceCudaWAFUrlAcls() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFUrlAclsCreate,
		Read:   resourceCudaWAFUrlAclsRead,
		Update: resourceCudaWAFUrlAclsUpdate,
		Delete: resourceCudaWAFUrlAclsDelete,

		Schema: map[string]*schema.Schema{
			"action":                  {Type: schema.TypeString, Optional: true, Description: "Action"},
			"comments":                {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"deny_response":           {Type: schema.TypeString, Optional: true, Description: "Deny Response"},
			"extended_match":          {Type: schema.TypeString, Optional: true, Description: "Extended Match"},
			"extended_match_sequence": {Type: schema.TypeString, Optional: true, Description: "Extended Match Sequence"},
			"follow_up_action":        {Type: schema.TypeString, Optional: true, Description: "Follow Up Action"},
			"follow_up_action_time":   {Type: schema.TypeString, Optional: true, Description: "Follow Up Action Time"},
			"host":                    {Type: schema.TypeString, Optional: true, Description: "Host Match"},
			"name":                    {Type: schema.TypeString, Required: true, Description: "URL ACL Name"},
			"redirect_url":            {Type: schema.TypeString, Optional: true, Description: "Redirect URL"},
			"response_page":           {Type: schema.TypeString, Optional: true, Description: "Response Page"},
			"enable":                  {Type: schema.TypeString, Optional: true, Description: "Enable URL ACL"},
			"url":                     {Type: schema.TypeString, Optional: true, Description: "URL Match"},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},

		Description: "`barracudawaf_url_acls` manages `Url Acls` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFUrlAclsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/url-acls"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFUrlAclsResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFUrlAclsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFUrlAclsRead(d, m)
}

func resourceCudaWAFUrlAclsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/url-acls"
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

func resourceCudaWAFUrlAclsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/url-acls"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFUrlAclsResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFUrlAclsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFUrlAclsRead(d, m)
}

func resourceCudaWAFUrlAclsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/url-acls"
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

func hydrateBarracudaWAFUrlAclsResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"action":                  d.Get("action").(string),
		"comments":                d.Get("comments").(string),
		"deny-response":           d.Get("deny_response").(string),
		"extended-match":          d.Get("extended_match").(string),
		"extended-match-sequence": d.Get("extended_match_sequence").(string),
		"follow-up-action":        d.Get("follow_up_action").(string),
		"follow-up-action-time":   d.Get("follow_up_action_time").(string),
		"host":                    d.Get("host").(string),
		"name":                    d.Get("name").(string),
		"redirect-url":            d.Get("redirect_url").(string),
		"response-page":           d.Get("response_page").(string),
		"enable":                  d.Get("enable").(string),
		"url":                     d.Get("url").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFUrlAclsSubResource(d *schema.ResourceData, name string, endpoint string) error {

	for subResource, subResourceParams := range subResourceUrlAclsParams {
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
