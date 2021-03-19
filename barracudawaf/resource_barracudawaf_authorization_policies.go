package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceAuthorizationPoliciesParams = map[string][]string{}
)

func resourceCudaWAFAuthorizationPolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFAuthorizationPoliciesCreate,
		Read:   resourceCudaWAFAuthorizationPoliciesRead,
		Update: resourceCudaWAFAuthorizationPoliciesUpdate,
		Delete: resourceCudaWAFAuthorizationPoliciesDelete,

		Schema: map[string]*schema.Schema{
			"allowed_groups":        {Type: schema.TypeString, Optional: true, Description: "Allowed Groups"},
			"allowed_users":         {Type: schema.TypeString, Optional: true, Description: "Allowed Users"},
			"comments":              {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"auth_context_classref": {Type: schema.TypeString, Optional: true, Description: "AuthnContextClassRef"},
			"name":                  {Type: schema.TypeString, Required: true, Description: "Policy Name"},
			"host":                  {Type: schema.TypeString, Optional: true, Description: "Host Match"},
			"extended_match":        {Type: schema.TypeString, Optional: true, Description: "Extended Match"},
			"extended_match_sequence": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Extended Match Sequence",
			},
			"cookie_timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Persistent Cookie Timeout",
			},
			"access_rules": {Type: schema.TypeString, Optional: true, Description: "Access Rules"},
			"enable_signing_on_authrequest": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable Signing on AuthRequest",
			},
			"digest_algorithm":      {Type: schema.TypeString, Optional: true, Description: "Digest Algorithm"},
			"status":                {Type: schema.TypeString, Optional: true, Description: "Status"},
			"url":                   {Type: schema.TypeString, Required: true, Description: "URL Match"},
			"use_persistent_cookie": {Type: schema.TypeString, Optional: true, Description: "Use Persistent Cookie"},
			"allow_any_authenticated_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Allow any Authenticated User",
			},
			"login_method":      {Type: schema.TypeString, Optional: true, Description: "Login Method"},
			"access_denied_url": {Type: schema.TypeString, Optional: true, Description: "Access Denied URL"},
			"login_url":         {Type: schema.TypeString, Optional: true, Description: "Auth Not Done URL"},
			"send_basic_auth": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Send Basic Authentication",
			},
			"send_domain_in_basic_auth": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Send Domain in Basic Authentication",
			},
			"kerberos_spn": {Type: schema.TypeString, Optional: true, Description: "Kerberos SPN"},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},

		Description: "`barracudawaf_authorization_policies` manages `Authorization Policies` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFAuthorizationPoliciesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/authorization-policies"
	err := client.CreateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFAuthorizationPoliciesResource(d, "post", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFAuthorizationPoliciesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFAuthorizationPoliciesRead(d, m)
}

func resourceCudaWAFAuthorizationPoliciesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/authorization-policies"
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

func resourceCudaWAFAuthorizationPoliciesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/authorization-policies"
	err := client.UpdateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFAuthorizationPoliciesResource(d, "put", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFAuthorizationPoliciesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFAuthorizationPoliciesRead(d, m)
}

func resourceCudaWAFAuthorizationPoliciesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/authorization-policies"
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

func hydrateBarracudaWAFAuthorizationPoliciesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"allowed-groups":                d.Get("allowed_groups").(string),
		"allowed-users":                 d.Get("allowed_users").(string),
		"comments":                      d.Get("comments").(string),
		"auth-context-classref":         d.Get("auth_context_classref").(string),
		"name":                          d.Get("name").(string),
		"host":                          d.Get("host").(string),
		"extended-match":                d.Get("extended_match").(string),
		"extended-match-sequence":       d.Get("extended_match_sequence").(string),
		"cookie-timeout":                d.Get("cookie_timeout").(string),
		"access-rules":                  d.Get("access_rules").(string),
		"enable-signing-on-authrequest": d.Get("enable_signing_on_authrequest").(string),
		"digest-algorithm":              d.Get("digest_algorithm").(string),
		"status":                        d.Get("status").(string),
		"url":                           d.Get("url").(string),
		"use-persistent-cookie":         d.Get("use_persistent_cookie").(string),
		"allow-any-authenticated-user":  d.Get("allow_any_authenticated_user").(string),
		"login-method":                  d.Get("login_method").(string),
		"access-denied-url":             d.Get("access_denied_url").(string),
		"login-url":                     d.Get("login_url").(string),
		"send-basic-auth":               d.Get("send_basic_auth").(string),
		"send-domain-in-basic-auth":     d.Get("send_domain_in_basic_auth").(string),
		"kerberos-spn":                  d.Get("kerberos_spn").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFAuthorizationPoliciesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceAuthorizationPoliciesParams {
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
