package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceExternalLdapServicesParams = map[string][]string{}
)

func resourceCudaWAFExternalLdapServices() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFExternalLdapServicesCreate,
		Read:   resourceCudaWAFExternalLdapServicesRead,
		Update: resourceCudaWAFExternalLdapServicesUpdate,
		Delete: resourceCudaWAFExternalLdapServicesDelete,

		Schema: map[string]*schema.Schema{
			"ip_address":           {Type: schema.TypeString, Required: true, Description: "IP Address"},
			"search_base":          {Type: schema.TypeString, Optional: true, Description: "LDAP Search Base"},
			"bind_dn":              {Type: schema.TypeString, Required: true, Description: "Bind DN"},
			"bind_password":        {Type: schema.TypeString, Optional: true, Description: "Bind Password"},
			"retype_bind_password": {Type: schema.TypeString, Optional: true, Description: "Re-enter Bind Password"},
			"default_role":         {Type: schema.TypeString, Required: true, Description: "Default Role"},
			"role_map": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "WAF-LDAP Role group Mapping",
			},
			"role_order": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "WAF-LDAP Role group priority",
			},
			"encryption":                 {Type: schema.TypeString, Required: true, Description: "Encryption"},
			"group_filter":               {Type: schema.TypeString, Optional: true, Description: "Group Filter"},
			"group_member_uid_attribute": {Type: schema.TypeString, Required: true, Description: "None"},
			"group_membership_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group Membership format description",
			},
			"group_name_attribute": {Type: schema.TypeString, Optional: true, Description: "Query For Group"},
			"name":                 {Type: schema.TypeString, Required: true, Description: "Realm Name"},
			"port":                 {Type: schema.TypeString, Optional: true, Description: "Port"},
			"uid_attribute":        {Type: schema.TypeString, Optional: true, Description: "UID Attribute"},
			"allow_nested_groups":  {Type: schema.TypeString, Required: true, Description: "Allow Nested Groups"},
			"ldap_server_type":     {Type: schema.TypeString, Optional: true, Description: "LDAP Server Type"},
			"validate_server_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Validate Server Certificate",
			},
		},

		Description: "`barracudawaf_external_ldap_services` manages `External Ldap Services` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFExternalLdapServicesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/external-ldap-services"
	err := client.CreateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFExternalLdapServicesResource(d, "post", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFExternalLdapServicesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFExternalLdapServicesRead(d, m)
}

func resourceCudaWAFExternalLdapServicesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/external-ldap-services"
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

func resourceCudaWAFExternalLdapServicesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/external-ldap-services"
	err := client.UpdateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFExternalLdapServicesResource(d, "put", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFExternalLdapServicesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFExternalLdapServicesRead(d, m)
}

func resourceCudaWAFExternalLdapServicesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/external-ldap-services"
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

func hydrateBarracudaWAFExternalLdapServicesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"ip-address":                  d.Get("ip_address").(string),
		"search-base":                 d.Get("search_base").(string),
		"bind-dn":                     d.Get("bind_dn").(string),
		"bind-password":               d.Get("bind_password").(string),
		"retype-bind-password":        d.Get("retype_bind_password").(string),
		"default-role":                d.Get("default_role").(string),
		"role-map":                    d.Get("role_map").(string),
		"role-order":                  d.Get("role_order").(string),
		"encryption":                  d.Get("encryption").(string),
		"group-filter":                d.Get("group_filter").(string),
		"group-member-uid-attribute":  d.Get("group_member_uid_attribute").(string),
		"group-membership-format":     d.Get("group_membership_format").(string),
		"group-name-attribute":        d.Get("group_name_attribute").(string),
		"name":                        d.Get("name").(string),
		"port":                        d.Get("port").(string),
		"uid-attribute":               d.Get("uid_attribute").(string),
		"allow-nested-groups":         d.Get("allow_nested_groups").(string),
		"ldap_server_type":            d.Get("ldap_server_type").(string),
		"validate-server-certificate": d.Get("validate_server_certificate").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFExternalLdapServicesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceExternalLdapServicesParams {
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
