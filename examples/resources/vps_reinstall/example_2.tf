# Reinstall a VPS with a raw public SSH key and suppressed root password email.
#
# When do_not_send_password = true, make sure ssh_keys or public_ssh_key is
# also set, otherwise you will not be able to log in to the reinstalled VPS.

resource "ovh_vps_reinstall" "reinstall" {
  service_name         = "vpsXXXXXX.vps.ovh.net"
  template_id          = 1144
  language             = "en"
  public_ssh_key       = "ssh-ed25519 AAAAC3Nz... user@host"
  do_not_send_password = true
}
