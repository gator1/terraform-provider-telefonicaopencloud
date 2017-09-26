package huaweicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/firewall_groups"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccFWFirewallGroupV2_basic(t *testing.T) {
	var epolicyID *string
	var ipolicyID *string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFWFirewallGroupV2_basic_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2("huaweicloud_fw_firewall_group_v2.fw_1", "", "", ipolicyID, epolicyID),
				),
			},
			resource.TestStep{
				Config: testAccFWFirewallGroupV2_basic_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2(
						"huaweicloud_fw_firewall_group_v2.fw_1", "fw_1", "terraform acceptance test", ipolicyID, epolicyID),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_port(t *testing.T) {
	var firewall_group FirewallGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFWFirewallV2_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("huaweicloud_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallRouterCount(&firewall_group, 1),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_no_ports(t *testing.T) {
	var firewall_group FirewallGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFWFirewallV2_no_ports,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("huaweicloud_fw_firewall_group_v2.fw_1", &firewall_group),
					resource.TestCheckResourceAttr("huaweicloud_fw_firewall_group_v2.fw_1", "description", "firewall router test"),
					testAccCheckFWFirewallRouterCount(&firewall_group, 0),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_port_update(t *testing.T) {
	var firewall_group FirewallGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFWFirewallV2_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("huaweicloud_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallRouterCount(&firewall_group, 1),
				),
			},
			resource.TestStep{
				Config: testAccFWFirewallV2_port_add,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("huaweicloud_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallRouterCount(&firewall_group, 2),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_port_remove(t *testing.T) {
	var firewall_group FirewallGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFWFirewallV2_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("huaweicloud_fw_firewall_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallRouterCount(&firewall_group, 1),
				),
			},
			resource.TestStep{
				Config: testAccFWFirewallV2_port_remove,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("huaweicloud_fw_firewall_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallRouterCount(&firewall_group, 0),
				),
			},
		},
	})
}

func testAccCheckFWFirewallGroupV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_firewall_group" {
			continue
		}

		_, err = firewall_groups.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Firewall group (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckFWFirewallGroupV2Exists(n string, firewall_group *FirewallGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Exists) Error creating HuaweiCloud networking client: %s", err)
		}

		var found FirewallGroup
		err = firewall_groups.Get(networkingClient, rs.Primary.ID).ExtractInto(&found)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Firewall group not found")
		}

		*firewall_group = found

		return nil
	}
}

func testAccCheckFWFirewallRouterCount(firewall_group *FirewallGroup, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(firewall_group.PortIDs) != expected {
			return fmt.Errorf("Expected %d Ports, got %d", expected, len(firewall_group.PortIDs))
		}

		return nil
	}
}

func testAccCheckFWFirewallGroupV2(n, expectedName, expectedDescription string, ipolicyID *string, epolicyID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Exists) Error creating HuaweiCloud networking client: %s", err)
		}

		var found *firewall_groups.FirewallGroup
		for i := 0; i < 5; i++ {
			// Firewall creation is asynchronous. Retry some times
			// if we get a 404 error. Fail on any other error.
			found, err = firewall_groups.Get(networkingClient, rs.Primary.ID).Extract()
			if err != nil {
				if _, ok := err.(gophercloud.ErrDefault404); ok {
					time.Sleep(time.Second)
					continue
				}
				return err
			}
			break
		}

		switch {
		case found.Name != expectedName:
			err = fmt.Errorf("Expected Name to be <%s> but found <%s>", expectedName, found.Name)
		case found.Description != expectedDescription:
			err = fmt.Errorf("Expected Description to be <%s> but found <%s>",
				expectedDescription, found.Description)
		case found.IngressPolicyID == "":
			err = fmt.Errorf("Ingress Policy should not be empty")
		case found.EgressPolicyID == "":
			err = fmt.Errorf("Egress Policy should not be empty")
		case ipolicyID != nil && found.IngressPolicyID == *ipolicyID:
			err = fmt.Errorf("Ingress Policy had not been correctly updated. Went from <%s> to <%s>",
				expectedName, found.Name)
		case epolicyID != nil && found.EgressPolicyID == *epolicyID:
			err = fmt.Errorf("Egress Policy had not been correctly updated. Went from <%s> to <%s>",
				expectedName, found.Name)
		}

		if err != nil {
			return err
		}

		ipolicyID = &found.IngressPolicyID
		epolicyID = &found.EgressPolicyID

		return nil
	}
}

const testAccFWFirewallGroupV2_basic_1 = `
resource "huaweicloud_fw_firewall_group_v2" "fw_1" {
  ingress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"
  egress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "huaweicloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
}
`

const testAccFWFirewallGroupV2_basic_2 = `
resource "huaweicloud_fw_firewall_group_v2" "fw_1" {
  name = "fw_1"
  description = "terraform acceptance test"
  ingress_policy_id = "${huaweicloud_fw_policy_v2.policy_2.id}"
  egress_policy_id = "${huaweicloud_fw_policy_v2.policy_2.id}"
  admin_state_up = true

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "huaweicloud_fw_policy_v2" "policy_2" {
  name = "policy_2"
}
`

const testAccFWFirewallV2_port = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}

resource "huaweicloud_networking_port_v2" "port_2" {
  name = "port_2"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.24"
  }
}

resource "huaweicloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "huaweicloud_fw_firewall_group_v2" "fw_1" {
  name = "firewall_1"
  description = "firewall router test"
  ingress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"
  egress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"
  ports = [
	"${huaweicloud_networking_port_v2.port_1.id}"
  ]
}
`

const testAccFWFirewallV2_port_add = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}

resource "huaweicloud_networking_port_v2" "port_2" {
  name = "port_2"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.24"
  }
}

resource "huaweicloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "huaweicloud_fw_firewall_group_v2" "fw_1" {
  name = "firewall_1"
  description = "firewall router test"
  ingress_policy_id = "${huaweicloud_fw_policy_v1.policy_1.id}"
  egress_policy_id = "${huaweicloud_fw_policy_v1.policy_1.id}"
  ports = [
	"${huaweicloud_networking_port_v2.port_1.id}",
	"${huaweicloud_networking_port_v2.port_2.id}"
  ]
}
`

const testAccFWFirewallV2_port_remove = `
resource "huaweicloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "huaweicloud_fw_firewall_group_v2" "fw_1" {
  name = "firewall_1"
  description = "firewall router test"
  ingress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"
  egress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"
}
`

const testAccFWFirewallV2_no_ports = `
resource "huaweicloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "huaweicloud_fw_firewall_group_v2" "fw_1" {
  name = "firewall_1"
  description = "firewall router test"
  ingress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"
  egress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"
}
`