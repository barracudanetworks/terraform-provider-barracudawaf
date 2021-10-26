resource "barracudawaf_trusted_server_certificate" "demo_trusted_server_cert_1" {
  name        = "DemoTrustedServerCert1"
  certificate = "<base_64_encoded_content>"

  depends_on  = [ barracudawaf_xxxx.xxxx ]
}