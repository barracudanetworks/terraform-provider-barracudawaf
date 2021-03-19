resource "barracudawaf_url_acls" "url_acl_1" {
    redirect_url = "http://www.example.com/index.html"
    action       = "Allow and Log"
    name         = "DemoUrlAcl1"
    parent       = [ "DemoApp1" ]
    depends_on   = [ barracudawaf_content_rule_servers.rule_group_server_1 ]
}