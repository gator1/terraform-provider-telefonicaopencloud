package telefonicaopencloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	//"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
	//"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	//"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

// PASS
func TestAccNetworkingV2Network_basic(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Network_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("telefonicaopencloud_networking_network_v2.network_1", &network),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Network_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"telefonicaopencloud_networking_network_v2.network_1", "name", "network_2"),
				),
			},
		},
	})
}

// PASS
func TestAccNetworkingV2Network_netstack(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Network_netstack,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("telefonicaopencloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("telefonicaopencloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("telefonicaopencloud_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2RouterInterfaceExists(
						"telefonicaopencloud_networking_router_interface_v2.ri_1"),
				),
			},
		},
	})
}

// PASS
func TestAccNetworkingV2Network_timeout(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Network_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("telefonicaopencloud_networking_network_v2.network_1", &network),
				),
			},
		},
	})
}

// SKIP needs admin user
func TestAccNetworkingV2Network_multipleSegmentMappings(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Network_multipleSegmentMappings,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("telefonicaopencloud_networking_network_v2.network_1", &network),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2NetworkDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating TelefonicaOpenCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "telefonicaopencloud_networking_network_v2" {
			continue
		}

		_, id := ExtractValFromNid(rs.Primary.ID)
		_, err := networks.Get(networkingClient, id).Extract()
		if err == nil {
			return fmt.Errorf("Network still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2NetworkExists(n string, network *networks.Network) resource.TestCheckFunc {
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
			return fmt.Errorf("Error creating TelefonicaOpenCloud networking client: %s", err)
		}

		_, id := ExtractValFromNid(rs.Primary.ID)
		found, err := networks.Get(networkingClient, id).Extract()
		if err != nil {
			return err
		}

		if found.ID != id {
			return fmt.Errorf("Network not found")
		}

		*network = *found

		return nil
	}
}

const testAccNetworkingV2Network_basic = `
resource "telefonicaopencloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}
`

const testAccNetworkingV2Network_update = `
resource "telefonicaopencloud_networking_network_v2" "network_1" {
  name = "network_2"
  # Can't do this to a network on OTC
  #admin_state_up = "false"
}
`

const testAccNetworkingV2Network_netstack = `
resource "telefonicaopencloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "telefonicaopencloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.10.0/24"
  ip_version = 4
  network_id = "${telefonicaopencloud_networking_network_v2.network_1.id}"
}

resource "telefonicaopencloud_networking_router_v2" "router_1" {
  name = "router_1"
}

resource "telefonicaopencloud_networking_router_interface_v2" "ri_1" {
  router_id = "${telefonicaopencloud_networking_router_v2.router_1.id}"
  subnet_id = "${telefonicaopencloud_networking_subnet_v2.subnet_1.id}"
}
`

const testAccNetworkingV2Network_fullstack = `
resource "telefonicaopencloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "telefonicaopencloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${telefonicaopencloud_networking_network_v2.network_1.id}"
}

resource "telefonicaopencloud_compute_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "a security group"
  rule {
    from_port = 22
    to_port = 22
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "telefonicaopencloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  security_group_ids = ["${telefonicaopencloud_compute_secgroup_v2.secgroup_1.id}"]
  network_id = "${telefonicaopencloud_networking_network_v2.network_1.id}"

  fixed_ip {
    "subnet_id" =  "${telefonicaopencloud_networking_subnet_v2.subnet_1.id}"
    "ip_address" =  "192.168.199.23"
  }
}

resource "telefonicaopencloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["${telefonicaopencloud_compute_secgroup_v2.secgroup_1.name}"]

  network {
    port = "${telefonicaopencloud_networking_port_v2.port_1.id}"
  }
}
`

const testAccNetworkingV2Network_timeout = `
resource "telefonicaopencloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

const testAccNetworkingV2Network_multipleSegmentMappings = `
resource "telefonicaopencloud_networking_network_v2" "network_1" {
  name = "network_1"
  segments =[
    {
      segmentation_id = 2,
      network_type = "vxlan"
    }
  ],
  admin_state_up = "true"
}
`
