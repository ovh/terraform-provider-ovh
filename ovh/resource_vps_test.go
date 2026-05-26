package ovh

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVpsBasic = `
data "ovh_me" "myaccount" {}

data "ovh_order_cart" "mycart" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
}

data "ovh_order_cart_product_plan" "vps" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "vps"
  plan_code      = "vps-le-2-2-40"
}

resource "ovh_vps" "myvps" {
  display_name = "%s"
  netboot_mode = "rescue"

  image_id = "45b2f222-ab10-44ed-863f-720942762b6f"
  public_ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDSD76EaLUzJjf70W8W2uU9FzEyl68di67Bd20qtYfBLJpFJuX/RJC9StI1y1RnXXqC1Lf/Yo+yJzvNx0iqLxCX1G7g0XYex74HkgC6a2QeNhp9M56ANZtA3TKKAbkZ1xobfhOPWpq3lEFp7dgJctcILBPL3l6OjKf6NIxHo5yF67Vy4D0nWl5utumNdWhhlX7MtVQooszLyIwPlNO+DzD3ZnJFCt2Z1jdRkhm/Oobtx17CZ+5SN23tgHXS6pLOgM6w30M11zkI510z95IAIHhRT7MbiXICkvG/0qHuSftz1j/CcHFbttNB27dH86vByumfSEgRKaoRkCqrn64IWrSsFr3Smsf7gZWLBlYLliGPyn8Tsr9bT5pRul6yTvVbfZ31RREBr1I0Lp4q++d+fIpa3LtMGRaMb9huJYy8cwW/Vfzbxsqfz9xzjIOFNcYl7J9l4cvz3hgSlai2Jgngw5ShNVlxcIKUdiynZWm09nQudlYNHgor9ID+JACzCfPkUZ8"

  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
  plan = [
    {
      duration     = "P1M"
      plan_code    = data.ovh_order_cart_product_plan.vps.plan_code
      pricing_mode = "default"

      configuration = [
        {
          label = "vps_datacenter"
          value = "WAW"
        },
        {
          label = "vps_os"
          value = "Debian 10"
        }
      ]
    }
  ]
}
`

const testAccVpsDoNotSendPassword = `
data "ovh_me" "myaccount" {}

data "ovh_order_cart" "mycart" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
}

data "ovh_order_cart_product_plan" "vps" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "vps"
  plan_code      = "vps-le-2-2-40"
}

resource "ovh_vps" "myvps" {
  display_name = "%s"
  netboot_mode = "rescue"

  image_id = "45b2f222-ab10-44ed-863f-720942762b6f"
  public_ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDSD76EaLUzJjf70W8W2uU9FzEyl68di67Bd20qtYfBLJpFJuX/RJC9StI1y1RnXXqC1Lf/Yo+yJzvNx0iqLxCX1G7g0XYex74HkgC6a2QeNhp9M56ANZtA3TKKAbkZ1xobfhOPWpq3lEFp7dgJctcILBPL3l6OjKf6NIxHo5yF67Vy4D0nWl5utumNdWhhlX7MtVQooszLyIwPlNO+DzD3ZnJFCt2Z1jdRkhm/Oobtx17CZ+5SN23tgHXS6pLOgM6w30M11zkI510z95IAIHhRT7MbiXICkvG/0qHuSftz1j/CcHFbttNB27dH86vByumfSEgRKaoRkCqrn64IWrSsFr3Smsf7gZWLBlYLliGPyn8Tsr9bT5pRul6yTvVbfZ31RREBr1I0Lp4q++d+fIpa3LtMGRaMb9huJYy8cwW/Vfzbxsqfz9xzjIOFNcYl7J9l4cvz3hgSlai2Jgngw5ShNVlxcIKUdiynZWm09nQudlYNHgor9ID+JACzCfPkUZ8"
  do_not_send_password = true

  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
  plan = [
    {
      duration     = "P1M"
      plan_code    = data.ovh_order_cart_product_plan.vps.plan_code
      pricing_mode = "default"

      configuration = [
        {
          label = "vps_datacenter"
          value = "WAW"
        },
        {
          label = "vps_os"
          value = "Debian 10"
        }
      ]
    }
  ]
}
`

const testAccVpsReinstallImageOnly = `
data "ovh_me" "myaccount" {}

data "ovh_order_cart" "mycart" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
}

data "ovh_order_cart_product_plan" "vps" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "vps"
  plan_code      = "vps-le-2-2-40"
}

resource "ovh_vps" "myvps" {
  display_name = "%s"
  netboot_mode = "rescue"

  image_id = "45b2f222-ab10-44ed-863f-720942762b6f"

  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
  plan = [
    {
      duration     = "P1M"
      plan_code    = data.ovh_order_cart_product_plan.vps.plan_code
      pricing_mode = "default"

      configuration = [
        {
          label = "vps_datacenter"
          value = "WAW"
        },
        {
          label = "vps_os"
          value = "Debian 10"
        }
      ]
    }
  ]
}
`

func TestAccResourceVps_basic(t *testing.T) {
	displayName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccVpsBasic,
		displayName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOrderVPS(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_vps.myvps", "netboot_mode", "rescue"),
					resource.TestCheckResourceAttr(
						"ovh_vps.myvps", "display_name", displayName),
					resource.TestCheckResourceAttr(
						"ovh_vps.myvps", "public_ssh_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDSD76EaLUzJjf70W8W2uU9FzEyl68di67Bd20qtYfBLJpFJuX/RJC9StI1y1RnXXqC1Lf/Yo+yJzvNx0iqLxCX1G7g0XYex74HkgC6a2QeNhp9M56ANZtA3TKKAbkZ1xobfhOPWpq3lEFp7dgJctcILBPL3l6OjKf6NIxHo5yF67Vy4D0nWl5utumNdWhhlX7MtVQooszLyIwPlNO+DzD3ZnJFCt2Z1jdRkhm/Oobtx17CZ+5SN23tgHXS6pLOgM6w30M11zkI510z95IAIHhRT7MbiXICkvG/0qHuSftz1j/CcHFbttNB27dH86vByumfSEgRKaoRkCqrn64IWrSsFr3Smsf7gZWLBlYLliGPyn8Tsr9bT5pRul6yTvVbfZ31RREBr1I0Lp4q++d+fIpa3LtMGRaMb9huJYy8cwW/Vfzbxsqfz9xzjIOFNcYl7J9l4cvz3hgSlai2Jgngw5ShNVlxcIKUdiynZWm09nQudlYNHgor9ID+JACzCfPkUZ8"),
					resource.TestCheckResourceAttr(
						"ovh_vps.myvps", "image_id", os.Getenv("OVH_VPS_IMAGE_ID")),
				),
			},
		},
	})
}

func TestAccResourceVps_doNotSendPassword(t *testing.T) {
	displayName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccVpsDoNotSendPassword,
		displayName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOrderVPS(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_vps.myvps", "netboot_mode", "rescue"),
					resource.TestCheckResourceAttr(
						"ovh_vps.myvps", "display_name", displayName),
					resource.TestCheckResourceAttr(
						"ovh_vps.myvps", "public_ssh_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDSD76EaLUzJjf70W8W2uU9FzEyl68di67Bd20qtYfBLJpFJuX/RJC9StI1y1RnXXqC1Lf/Yo+yJzvNx0iqLxCX1G7g0XYex74HkgC6a2QeNhp9M56ANZtA3TKKAbkZ1xobfhOPWpq3lEFp7dgJctcILBPL3l6OjKf6NIxHo5yF67Vy4D0nWl5utumNdWhhlX7MtVQooszLyIwPlNO+DzD3ZnJFCt2Z1jdRkhm/Oobtx17CZ+5SN23tgHXS6pLOgM6w30M11zkI510z95IAIHhRT7MbiXICkvG/0qHuSftz1j/CcHFbttNB27dH86vByumfSEgRKaoRkCqrn64IWrSsFr3Smsf7gZWLBlYLliGPyn8Tsr9bT5pRul6yTvVbfZ31RREBr1I0Lp4q++d+fIpa3LtMGRaMb9huJYy8cwW/Vfzbxsqfz9xzjIOFNcYl7J9l4cvz3hgSlai2Jgngw5ShNVlxcIKUdiynZWm09nQudlYNHgor9ID+JACzCfPkUZ8"),
					resource.TestCheckResourceAttr(
						"ovh_vps.myvps", "do_not_send_password", "true"),
					resource.TestCheckResourceAttr(
						"ovh_vps.myvps", "image_id", os.Getenv("OVH_VPS_IMAGE_ID")),
				),
			},
		},
	})
}

const testAccVpsUpdateInPlace = `
data "ovh_me" "myaccount" {}

data "ovh_order_cart" "mycart" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
}

data "ovh_order_cart_product_plan" "vps" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "vps"
  plan_code      = "vps-le-2-2-40"
}

resource "ovh_vps" "myvps" {
  display_name = "%s"
  netboot_mode = "%s"

  image_id = "45b2f222-ab10-44ed-863f-720942762b6f"
  public_ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDSD76EaLUzJjf70W8W2uU9FzEyl68di67Bd20qtYfBLJpFJuX/RJC9StI1y1RnXXqC1Lf/Yo+yJzvNx0iqLxCX1G7g0XYex74HkgC6a2QeNhp9M56ANZtA3TKKAbkZ1xobfhOPWpq3lEFp7dgJctcILBPL3l6OjKf6NIxHo5yF67Vy4D0nWl5utumNdWhhlX7MtVQooszLyIwPlNO+DzD3ZnJFCt2Z1jdRkhm/Oobtx17CZ+5SN23tgHXS6pLOgM6w30M11zkI510z95IAIHhRT7MbiXICkvG/0qHuSftz1j/CcHFbttNB27dH86vByumfSEgRKaoRkCqrn64IWrSsFr3Smsf7gZWLBlYLliGPyn8Tsr9bT5pRul6yTvVbfZ31RREBr1I0Lp4q++d+fIpa3LtMGRaMb9huJYy8cwW/Vfzbxsqfz9xzjIOFNcYl7J9l4cvz3hgSlai2Jgngw5ShNVlxcIKUdiynZWm09nQudlYNHgor9ID+JACzCfPkUZ8"

  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
  plan = [
    {
      duration     = "P1M"
      plan_code    = data.ovh_order_cart_product_plan.vps.plan_code
      pricing_mode = "default"

      configuration = [
        {
          label = "vps_datacenter"
          value = "WAW"
        },
        {
          label = "vps_os"
          value = "Debian 10"
        }
      ]
    }
  ]
}
`

const testAccVpsSlaMonitoring = `
data "ovh_me" "myaccount" {}

data "ovh_order_cart" "mycart" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
}

data "ovh_order_cart_product_plan" "vps" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "vps"
  plan_code      = "vps-le-2-2-40"
}

resource "ovh_vps" "myvps" {
  display_name = "%s"
  netboot_mode = "rescue"
  sla_monitoring = true

  image_id = "45b2f222-ab10-44ed-863f-720942762b6f"
  public_ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDSD76EaLUzJjf70W8W2uU9FzEyl68di67Bd20qtYfBLJpFJuX/RJC9StI1y1RnXXqC1Lf/Yo+yJzvNx0iqLxCX1G7g0XYex74HkgC6a2QeNhp9M56ANZtA3TKKAbkZ1xobfhOPWpq3lEFp7dgJctcILBPL3l6OjKf6NIxHo5yF67Vy4D0nWl5utumNdWhhlX7MtVQooszLyIwPlNO+DzD3ZnJFCt2Z1jdRkhm/Oobtx17CZ+5SN23tgHXS6pLOgM6w30M11zkI510z95IAIHhRT7MbiXICkvG/0qHuSftz1j/CcHFbttNB27dH86vByumfSEgRKaoRkCqrn64IWrSsFr3Smsf7gZWLBlYLliGPyn8Tsr9bT5pRul6yTvVbfZ31RREBr1I0Lp4q++d+fIpa3LtMGRaMb9huJYy8cwW/Vfzbxsqfz9xzjIOFNcYl7J9l4cvz3hgSlai2Jgngw5ShNVlxcIKUdiynZWm09nQudlYNHgor9ID+JACzCfPkUZ8"

  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
  plan = [
    {
      duration     = "P1M"
      plan_code    = data.ovh_order_cart_product_plan.vps.plan_code
      pricing_mode = "default"

      configuration = [
        {
          label = "vps_datacenter"
          value = "WAW"
        },
        {
          label = "vps_os"
          value = "Debian 10"
        }
      ]
    }
  ]
}
`

func TestAccResourceVps_updateInPlace(t *testing.T) {
	if os.Getenv("OVH_TESTACC_ORDER_VPS") == "" {
		t.Skip("OVH_TESTACC_ORDER_VPS not set, skipping VPS update-in-place test")
	}

	displayName := acctest.RandomWithPrefix(test_prefix)
	updatedDisplayName := acctest.RandomWithPrefix(test_prefix)

	initialConfig := fmt.Sprintf(testAccVpsUpdateInPlace, displayName, "rescue")
	updatedConfig := fmt.Sprintf(testAccVpsUpdateInPlace, updatedDisplayName, "local")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOrderVPS(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vps.myvps", "display_name", displayName),
					resource.TestCheckResourceAttr("ovh_vps.myvps", "netboot_mode", "rescue"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vps.myvps", "display_name", updatedDisplayName),
					resource.TestCheckResourceAttr("ovh_vps.myvps", "netboot_mode", "local"),
				),
			},
		},
	})
}

func TestAccResourceVps_slaMonitoring(t *testing.T) {
	if os.Getenv("OVH_TESTACC_ORDER_VPS") == "" {
		t.Skip("OVH_TESTACC_ORDER_VPS not set, skipping VPS sla_monitoring test")
	}

	displayName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(testAccVpsSlaMonitoring, displayName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOrderVPS(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vps.myvps", "display_name", displayName),
					resource.TestCheckResourceAttr("ovh_vps.myvps", "sla_monitoring", "true"),
				),
			},
		},
	})
}

func TestAccResourceVps_displayNameTooLong(t *testing.T) {
	if os.Getenv("OVH_TESTACC_ORDER_VPS") == "" {
		t.Skip("OVH_TESTACC_ORDER_VPS not set, skipping VPS display_name validation test")
	}

	// Build a display_name longer than 50 characters to trigger LengthAtMost(50)
	// validation introduced by the W15 fix.
	tooLong := strings.Repeat("a", 51)
	config := fmt.Sprintf(testAccVpsBasic, tooLong)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOrderVPS(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`(?i)(length|50|too long|at most)`),
			},
		},
	})
}

func TestAccResourceVps_reinstallImageOnly(t *testing.T) {
	displayName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccVpsReinstallImageOnly,
		displayName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOrderVPS(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_vps.myvps", "netboot_mode", "rescue"),
					resource.TestCheckResourceAttr(
						"ovh_vps.myvps", "display_name", displayName),
					resource.TestCheckResourceAttr(
						"ovh_vps.myvps", "image_id", "45b2f222-ab10-44ed-863f-720942762b6f"),
					resource.TestCheckResourceAttrSet("ovh_vps.myvps", "id"),
				),
			},
		},
	})
}
