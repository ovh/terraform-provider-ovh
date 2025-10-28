resource "ovh_iam_resource_tags" "project_tags" {
  urn = "urn:v1:eu:resource:cloudProject:1234567890abcdef"
  
  tags = {
    environment    = "staging"
    cost_center    = "engineering"
    project        = "web-app"
    owner          = "team@example.com"
    backup_policy  = "daily"
    compliance     = "gdpr"
  }
}