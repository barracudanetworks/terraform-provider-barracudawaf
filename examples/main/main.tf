provider "barracudawaf" {
    address  = "x.x.x.x"
    username = "xxxxxxx"
    port     = "8443"
    password = "xxxxxxx"
}

resource "barracudawaf_trusted_server_certificate" "demo_trusted_server_cert_1" {
  name        = "DemoTrustedServerCert1"
  certificate = "<base_64_encoded_content>"
}

resource "barracudawaf_trusted_ca_certificate" "demo_trusted_ca_cert_1" {
  name        = "DemoTrustedCACert1"
  certificate = "<base_64_encoded_content>"
  depends_on  = [ barracudawaf_trusted_server_certificate.demo_trusted_server_cert_1 ]
}

resource "barracudawaf_self_signed_certificate" "demo_self_signed_cert_1" {
    name                     = "DemoSelfSignedCert1"
    allow_private_key_export = "Yes"
    city                     = "Bangalore"
    common_name              = "waf.test.local"
    country_code             = "IN"
    key_size                 = "1024"
    key_type                 = "rsa"
    organization_name        = "Barracuda Networks"
    organizational_unit      = "Engineering"
    state                    = "Karnataka"
    depends_on               = [barracudawaf_trusted_ca_certificate.demo_trusted_ca_cert_1]
}

resource "barracudawaf_services" "demo_app_1" {
    name            = "DemoApp1"
    ip_address      = "x.x.x.x"
    port            = "80"
    type            = "HTTP"
    vsite           = "default"
    address_version = "IPv4"
    status          = "On"
    group           = "default"
    comments        = "Demo Service with Terraform"

    basic_security {
      mode = "Active"
    }

    depends_on = [ barracudawaf_self_signed_certificate.demo_self_signed_cert_1 ]
}

resource "barracudawaf_servers" "demo_server_1" {
    name            = "DemoServer1"
    identifier      = "IP Address"
    address_version = "IPv4"
    status          = "In Service"
    ip_address      = "x.x.x.x"
    port            = "80"
    comments        = "Creating the Demo Server"
    parent          = [ barracudawaf_services.demo_app_1.name ]
    
    out_of_band_health_checks {
      enable_oob_health_checks = "Yes"
      interval                 = "900"
    }

    depends_on      = [ barracudawaf_services.demo_app_1 ]
}

resource "barracudawaf_services" "demo_app_2" {
    name            = "DemoApp2"
    ip_address      = "x.x.x.x"
    port            = "443"
    type            = "HTTPS"
    vsite           = "default"
    address_version = "IPv4"
    status          = "On"
    group           = "default"
    comments        = "Demo Service with Terraform"
    certificate     = barracudawaf_self_signed_certificate.demo_self_signed_cert_1.name

    basic_security {
      mode = "Active"
    }

    depends_on = [ barracudawaf_servers.demo_server_1 ]
}

resource "barracudawaf_servers" "demo_server_2" {
    name            = "TestServer2"
    identifier      = "IP Address"
    address_version = "IPv4"
    status          = "In Service"
    ip_address      = "x.x.x.x"
    port            = "80"
    comments        = "Creating the Demo Server"
    parent          = [ barracudawaf_services.demo_app_2.name ]

    out_of_band_health_checks {
      enable_oob_health_checks = "Yes"
      interval                 = "900"
    }

    depends_on = [ barracudawaf_services.demo_app_2 ]
}

resource "barracudawaf_security_policies" "demo_security_policy_1" {
    name       = "DemoPolicy1"
    based_on   = "Create New"
    
    depends_on = [ barracudawaf_servers.demo_server_2 ]
}

resource "barracudawaf_content_rules" "demo_rule_group_1" {
    name                = "DemoRuleGroup1"
    url_match           = "/index.html"
    host_match          = "www.example.com"
    web_firewall_policy = "DemoPolicy1"
    mode                = "Active"
    parent              = [ barracudawaf_services.demo_app_1.name ]
    
    depends_on          = [ barracudawaf_security_policies.demo_security_policy_1 ]
}
 
resource "barracudawaf_content_rule_servers" "demo_rule_group_server_1" {
    name        = "DemoRuleGroupServer1"
    identifier  = "Hostname"
    hostname    = "barracuda.com"
    parent      = [ barracudawaf_services.demo_app_1.name, barracudawaf_content_rules.demo_rule_group_1.name ]
    

    application_layer_health_checks {
        method               = "POST"
        match_content_string = "index"
        domain               = "example.com"
    }

    depends_on = [ barracudawaf_content_rules.demo_rule_group_1 ]
}

resource "barracudawaf_url_acls" "demo_url_acl_1" {
    name         = "DemoUrlAcl1"
    redirect_url = "http://www.example.com/index.html"
    action       = "Allow and Log"
    parent       = [ barracudawaf_services.demo_app_1.name ]
    
    depends_on   = [ barracudawaf_content_rule_servers.demo_rule_group_server_1 ]
}
