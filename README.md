# Overview #

A [Terraform](https://www.terraform.io/) provider for Barracuda Web Application Firewall.

&nbsp;
## Requirements ##
-	[Terraform](https://www.terraform.io/downloads.html) v0.14.x
-	[Go](https://golang.org/doc/install) 1.15 (to build the provider plugin)

&nbsp;
## Usage ##

**Use provider**
```hcl
variable address {}
variable username {}
variable password {}
variable port {}

provider "barracudawaf" {
    address  = "x.x.x.x"
    username = "xxxxxxx"
    port     = "8443"
    password = "xxxxxxx"
}

```
**Create Self Signed Certificates**
```hcl
resource "barracudawaf_self_signed_certificate" "demo_self_signed_cert_1" {
    name                     = "DemoSelfSignedCert1"
    allow_private_key_export = "Yes"
    city                     = "xxxxxx"
    common_name              = "xxxxxx"
    country_code             = "IN"
    key_size                 = "1024"
    key_type                 = "rsa"
    organization_name        = "xxxxxx"
    organizational_unit      = "xxxxxx"
    state                    = "xxxxxx"
}
```

**Create Service**
```hcl
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

    depends_on = [ barracudawaf_services.demo_app_1 ]
}
```
**Create Servers**
```hcl
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

    depends_on      = [ barracudawaf_services.demo_app_2 ]
}
```
**Create Security Policies**
```hcl
resource "barracudawaf_security_policies" "demo_security_policy_1" {
    name       = "DemoPolicy1"
    based_on   = "Create New"
    
    depends_on = [ barracudawaf_servers.demo_server_1 ]
}
```
**Create Rule Groups**
```hcl
resource "barracudawaf_content_rules" "demo_rule_group_1" {
    name                = "DemoRuleGroup1"
    url_match           = "/index.html"
    host_match          = "www.example.com"
    web_firewall_policy = "DemoPolicy1"
    mode                = "Active"
    parent              = [ barracudawaf_services.demo_app_1.name ]
    
    depends_on          = [ barracudawaf_security_policies.demo_security_policy_1 ]
}
```
**Create Rule Group Servers**
```hcl
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
```

&nbsp;&nbsp;
## Building The Provider ##

### Dependencies for building from source ###
If you need to build from source, you should have a working Go environment setup.  If not check out the Go [getting started](http://golang.org/doc/install) guide.

This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) for dependency management.  To fetch all dependencies run `go get` inside this repository.

&nbsp;&nbsp;
### Build ###

Clone repository to: $GOPATH/src/github.com/hhakkaev/terraform-provider-barracudawaf
```shell
$ mkdir -p $GOPATH/src/github.com/hhakkaev; cd $GOPATH/src/github.com/hhakkaev
$ git clone https://github.com/hhakkaev/terraform-provider-barracudawaf.git
```

Enter the provider directory and build the provider
```shell
cd $GOPATH/src/github.com/hhakkaev/terraform-provider-barracudawaf
make build
```

&nbsp;&nbsp;
### Install ###

```shell
$ cd $GOPATH/src/github.com/hhakkaev/terraform-provider-barracudawaf
$ make install

```

&nbsp;&nbsp;
# Using the Provider

If you're building the provider, follow the instructions to install it as a plugin. After placing it into your plugins directory, run terraform init to initialize it.

&nbsp;&nbsp;
# Developing the Provider

If you wish to work on the provider, you'll first need Go installed on your machine (version 1.15 is required). You'll also need to correctly setup a GOPATH, as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run make build. This will create a binary with name `terraform-provider-barracudawaf` in `$GOPATH/src/github.com/hhakkaev/terraform-provider-barracudawaf` directory.

```shell
$ make build
...
$ $GOPATH/src/github.com/hhakkaev/terraform-provider-barracudawaf
...

```

&nbsp;
# Using the binary instead of building it from source #

Download the binary added under [releases](https://github.com/hhakkaev/terraform-provider-barracudawaf/releases), and follow below :

Copy the downloaded binary `terraform-provider-barracudawaf_v<tag>` into `plugins` Terraform directory.

