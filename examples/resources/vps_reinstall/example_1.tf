# Reinstall a VPS using the legacy template-based path (numeric templateId).
#
# Use the ovh_vps_templates / ovh_vps_template data sources to discover a
# template id that is compatible with your VPS plan.

resource "ovh_vps_reinstall" "reinstall" {
  service_name = "vpsXXXXXX.vps.ovh.net"
  template_id  = 1144 # e.g. Debian 12
  language     = "en"

  # Provide either ssh_keys (by name from /me/sshKey) and/or a raw public_ssh_key.
  ssh_keys = ["my-key"]

  # Optional: re-run the reinstall when one of these values changes.
  triggers = {
    rebuild_at = "2025-01-01T00:00:00Z"
  }
}
