resource "ovh_okms" "new_kms" {
  ovh_subsidiary = "FR"
  region         = "eu-west-rbx"
  display_name   = "terraformed KMS"
}
