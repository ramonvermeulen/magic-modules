package identityplatform_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccIdentityPlatformTenant_identityPlatformTenantUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIdentityPlatformTenantDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityPlatformTenant_identityPlatformTenantBasic(context),
			},
			{
				ResourceName:      "google_identity_platform_tenant.tenant",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIdentityPlatformTenant_identityPlatformTenantUpdate(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_identity_platform_tenant.tenant", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_identity_platform_tenant.tenant", "display_name", "my-tenant"),
					resource.TestCheckResourceAttr("google_identity_platform_tenant.tenant", "allow_password_signup", "false"),
					resource.TestCheckResourceAttr("google_identity_platform_tenant.tenant", "enable_email_link_signin", "true"),
					resource.TestCheckResourceAttr("google_identity_platform_tenant.tenant", "disable_auth", "true"),
					resource.TestCheckResourceAttr("google_identity_platform_tenant.tenant", "client.0.permissions.0.disabled_user_signup", "true"),
					resource.TestCheckResourceAttr("google_identity_platform_tenant.tenant", "client.0.permissions.0.disabled_user_deletion", "true"),
				),
			},
			{
				ResourceName:      "google_identity_platform_tenant.tenant",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIdentityPlatformTenant_identityPlatformTenantBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_identity_platform_tenant" "tenant" {
  display_name          = "tenant"
  allow_password_signup = true
}
`, context)
}

func testAccIdentityPlatformTenant_identityPlatformTenantUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_identity_platform_tenant" "tenant" {
  display_name             = "my-tenant"
  allow_password_signup    = false
  enable_email_link_signin = true
  disable_auth             = true
  client {
    permissions {
      disabled_user_signup   = true
      disabled_user_deletion = true
    }
  }
}
`, context)
}
