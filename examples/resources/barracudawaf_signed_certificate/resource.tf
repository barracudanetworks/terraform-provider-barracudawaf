resource "barracudawaf_self_signed_certificate" "demo_self_signed_cert_1" {
    name                     = "DemoSelfSignedCert1"
    allow_private_key_export = "Yes"
    city                     = "xxxxxx"
    common_name              = "xxxxxx"
    country_code             = "xx"
    key_size                 = "1024"
    key_type                 = "rsa"
    organization_name        = "xxxxxx"
    organizational_unit      = "xxxxxx"
    state                    = "xxxxxx"
    
    depends_on               = [ barracudawaf_xxxx.xxxx ]
}