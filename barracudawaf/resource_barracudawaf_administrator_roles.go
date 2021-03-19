package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceAdministratorRolesParams = map[string][]string{}
)

func resourceCudaWAFAdministratorRoles() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFAdministratorRolesCreate,
		Read:   resourceCudaWAFAdministratorRolesRead,
		Update: resourceCudaWAFAdministratorRolesUpdate,
		Delete: resourceCudaWAFAdministratorRolesDelete,

		Schema: map[string]*schema.Schema{
			"api_privilege":           {Type: schema.TypeString, Optional: true, Description: "API Privilege"},
			"authentication_services": {Type: schema.TypeString, Optional: true, Description: "Auth Services"},
			"name":                    {Type: schema.TypeString, Required: true, Description: "Role Name"},
			"objects":                 {Type: schema.TypeString, Optional: true, Description: "Object access permissions"},
			"operations":              {Type: schema.TypeString, Optional: true, Description: "Specify Allowed Operations"},
			"security_policies":       {Type: schema.TypeString, Optional: true, Description: "Security Policies"},
			"service_groups":          {Type: schema.TypeString, Optional: true, Description: "Service Group"},
			"services":                {Type: schema.TypeString, Optional: true, Description: "Services"},
			"role_type":               {Type: schema.TypeString, Optional: true},
			"vsites":                  {Type: schema.TypeString, Optional: true, Description: "Vsites"},
		},

		Description: "`barracudawaf_administrator_roles` manages `Administrator Roles` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFAdministratorRolesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/administrator-roles"
	err := client.CreateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFAdministratorRolesResource(d, "post", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFAdministratorRolesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFAdministratorRolesRead(d, m)
}

func resourceCudaWAFAdministratorRolesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/administrator-roles"
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

func resourceCudaWAFAdministratorRolesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/administrator-roles"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFAdministratorRolesResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFAdministratorRolesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFAdministratorRolesRead(d, m)
}

func resourceCudaWAFAdministratorRolesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/administrator-roles"
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

func hydrateBarracudaWAFAdministratorRolesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"api-privilege":           d.Get("api_privilege").(string),
		"authentication-services": d.Get("authentication_services").(string),
		"name":                    d.Get("name").(string),
		"objects":                 d.Get("objects").(string),
		"operations":              d.Get("operations").(string),
		"security-policies":       d.Get("security_policies").(string),
		"service-groups":          d.Get("service_groups").(string),
		"services":                d.Get("services").(string),
		"role-type":               d.Get("role_type").(string),
		"vsites":                  d.Get("vsites").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"role-type"}
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

func (b *BarracudaWAF) hydrateBarracudaWAFAdministratorRolesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceAdministratorRolesParams {
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
