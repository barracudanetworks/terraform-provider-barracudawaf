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
