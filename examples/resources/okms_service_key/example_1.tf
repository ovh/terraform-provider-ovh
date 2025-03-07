resource "ovh_okms_service_key" "key_symetric" {
  okms_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name       = "key_oct"
  type       = "oct"
  size       = 256
  operations = ["encrypt", "decrypt"]
}

resource "ovh_okms_service_key" "key_rsa" {
  okms_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name       = "key_rsa"
  type       = "RSA"
  size       = 2048
  operations = ["sign", "verify"]
}

resource "ovh_okms_service_key" "key_ecdsa" {
  okms_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name       = "key_ecdsa"
  type       = "EC"
  curve      = "P-256"
  operations = ["sign", "verify"]
}
