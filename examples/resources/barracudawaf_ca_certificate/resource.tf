resource "barracudawaf_trusted_ca_certificate" "demo_trusted_ca_cert_1" {
  name        = "DemoTrustedCACert1"
  certificate = "<base_64_encoded_content>"
  
  depends_on  = [ barracudawaf_xxxx.xxxx ]
}