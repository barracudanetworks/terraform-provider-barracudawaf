resource "barracudawaf_letsencrypt_certificate" "demo_letsencrypt_cert" {
    name                       = "DemoLetsEncryptCert"
    common_name                = "xxxxxx"
    allow_private_key_export   = "Yes"
    auto_renew_cert            = "Yes"
    schedule_renewal_day       = "60"

    multi_cert_trusted_service = barracudawaf_services.application_1.name
    depends_on = [ barracudawaf_xxxx.xxxx ]
}