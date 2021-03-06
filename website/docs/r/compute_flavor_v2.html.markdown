---
layout: "telefonicaopencloud"
page_title: "TelefonicaOpenCloud: telefonicaopencloud_compute_flavor_v2"
sidebar_current: "docs-telefonicaopencloud-resource-compute-flavor-v2"
description: |-
  Manages a V2 flavor resource within TelefonicaOpenCloud.
---

# telefonicaopencloud\_compute\_flavor_v2

Manages a V2 flavor resource within TelefonicaOpenCloud.

## Example Usage

```hcl
resource "telefonicaopencloud_compute_flavor_v2" "test-flavor" {
  name  = "my-flavor"
  ram   = "8"
  vcpus = "2"
  disk  = "20"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    Flavors are associated with accounts, but a Compute client is needed to
    create one. If omitted, the `region` argument of the provider is used.
    Changing this creates a new flavor.

* `name` - (Required) A unique name for the flavor. Changing this creates a new
    flavor.

* `ram` - (Required) The amount of RAM to use, in megabytes. Changing this
    creates a new flavor.

* `vcpus` - (Required) The number of virtual CPUs to use. Changing this creates
    a new flavor.

* `disk` - (Required) The amount of disk space in gigabytes to use for the root
    (/) partition. Changing this creates a new flavor.

* `swap` - (Optional) The amount of disk space in megabytes to use. If
    unspecified, the default is 0. Changing this creates a new flavor.

* `rx_tx_factor` - (Optional) RX/TX bandwith factor. The default is 1. Changing
    this creates a new flavor.

* `is_public` - (Optional) Whether the flavor is public. Changing this creates
    a new flavor.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `ram` - See Argument Reference above.
* `vcpus` - See Argument Reference above.
* `disk` - See Argument Reference above.
* `swap` - See Argument Reference above.
* `rx_tx_factor` - See Argument Reference above.
* `is_public` - See Argument Reference above.

## Import

Flavors can be imported using the `ID`, e.g.

```
$ terraform import telefonicaopencloud_compute_flavor_v2.my-flavor 4142e64b-1b35-44a0-9b1e-5affc7af1106
```
