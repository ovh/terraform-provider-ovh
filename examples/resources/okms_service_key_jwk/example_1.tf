resource "ovh_okms_service_key_jwk" "key_symetric" {
  okms_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name       = "key_oct"
  keys       = [
    jsondecode(<<EOT
      {
        "k": "Wc7IgEZzWicZf1LTJUtA0w",
        "key_ops": ["encrypt", "decrypt"],
        "kty": "oct"
      }
    EOT
  )]
}

resource "ovh_okms_service_key_jwk" "key_rsa" {
  okms_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name       = "key_rsa"
  keys       = [
    jsondecode(<<EOT
      {
        "key_ops": ["sign", "verify"],
        "kty": "RSA",
        "d": "UgaHfwn0kjl...",
        "dp": "SIrzLAa0Ll...",
        "dq": "aQv6Kg0Lw1...",
        "e": "AQAB",
        "n": "qrFKVDudlle...",
        "p": "7O4PCo_cWzu...",
        "q": "uG5pDYvV-eu...",
        "qi": "B5z00bGrZO..."
      }
    EOT
  )]
}

resource "ovh_okms_service_key_jwk" "key_ecdsa" {
  okms_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name       = "key_ecdsa"
  keys = [
    jsondecode(<<EOT
      {
        "key_ops": ["sign", "verify"],
        "kty": "EC",
        "crv": "P-256",
        "alg": "ES256",
        "d": "SIy6AYrv5nGBLQsM7bg7WCbAPxHyUIVTaDyTxrCWPks",
        "x": "V7a79Iv0RdykDIzhJhu5OvkCFJ8rFkFm5r11qwR9QeY",
        "y": "RSUYb-RPSkF5al1D2fnxerahFpHCHtmJRAlUeS1Ehtw"
      }
    EOT
  )]
}
