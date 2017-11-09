package telefonicaopencloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccFWRuleV2_importBasic(t *testing.T) {
	resourceName := "telefonicaopencloud_fw_rule_v2.rule_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWRuleV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFWRuleV2_basic_2,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
