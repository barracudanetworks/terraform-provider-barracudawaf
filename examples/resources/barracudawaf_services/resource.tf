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
    
    depends_on = [ barracudawaf_signed_certificate.demo_signed_cert ]
}