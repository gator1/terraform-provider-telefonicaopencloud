package telefonicaopencloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// PASS
func TestAccBlockStorageV2Volume_importBasic(t *testing.T) {
	resourceName := "telefonicaopencloud_blockstorage_volume_v2.volume_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccBlockStorageV2Volume_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}