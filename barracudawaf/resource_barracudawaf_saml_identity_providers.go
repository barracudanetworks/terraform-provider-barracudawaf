package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceSamlIdentityProvidersParams = map[string][]string{}
)

func resourceCudaWAFSamlIdentityProviders() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFSamlIdentityProvidersCreate,
		Read:   resourceCudaWAFSamlIdentityProvidersRead,
		Update: resourceCudaWAFSamlIdentityProvidersUpdate,
		Delete: resourceCudaWAFSamlIdentityProvidersDelete,

		Schema: map[string]*schema.Schema{
			"metadata_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the metadata file being uploaded.",
			},
			"metadata_type":       {Type: schema.TypeString, Optional: true, Description: "Identity Provider Metadata Type"},
			"autoupdate_metadata": {Type: schema.TypeString, Optional: true, Description: "Auto Update Metadata"},
			"metadata_url":        {Type: schema.TypeString, Optional: true, Description: "Metadata URL"},
			"name":                {Type: schema.TypeString, Required: true, Description: "Identity Provider Name"},
			"metadata_content":    {Type: schema.TypeString, Optional: true, Description: "Must be a Base64 encoded value"},
			"parent":              {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}, Required: true},
		},

		Description: "`barracudawaf_saml_identity_providers` manages `Saml Identity Providers` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFSamlIdentityProvidersCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/saml-services/" + d.Get("parent.0").(string) + "/saml-identity-providers"
	err := client.CreateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFSamlIdentityProvidersResource(d, "post", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFSamlIdentityProvidersSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFSamlIdentityProvidersRead(d, m)
}

func resourceCudaWAFSamlIdentityProvidersRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/saml-services/" + d.Get("parent.0").(string) + "/saml-identity-providers"
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

func resourceCudaWAFSamlIdentityProvidersUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/saml-services/" + d.Get("parent.0").(string) + "/saml-identity-providers"
	err := client.UpdateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFSamlIdentityProvidersResource(d, "put", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFSamlIdentityProvidersSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFSamlIdentityProvidersRead(d, m)
}

func resourceCudaWAFSamlIdentityProvidersDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/saml-services/" + d.Get("parent.0").(string) + "/saml-identity-providers"
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

func hydrateBarracudaWAFSamlIdentityProvidersResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"metadata-file":       d.Get("metadata_file").(string),
		"metadata-type":       d.Get("metadata_type").(string),
		"autoupdate-metadata": d.Get("autoupdate_metadata").(string),
		"metadata-url":        d.Get("metadata_url").(string),
		"name":                d.Get("name").(string),
		"metadata-content":    d.Get("metadata_content").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFSamlIdentityProvidersSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceSamlIdentityProvidersParams {
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
