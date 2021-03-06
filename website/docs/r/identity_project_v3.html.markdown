---
layout: "telefonicaopencloud"
page_title: "TelefonicaOpenCloud: telefonicaopencloud_identity_project_v3"
sidebar_current: "docs-telefonicaopencloud-resource-identity-project-v3"
description: |-
  Manages a V3 Project resource within TelefonicaOpenCloud Keystone.
---

# telefonicaopencloud\_identity\_project_v3

Manages a V3 Project resource within TelefonicaOpenCloud Keystone.

Note: You _must_ have admin privileges in your TelefonicaOpenCloud cloud to use
this resource.

## Example Usage

```hcl
resource "telefonicaopencloud_identity_project_v3" "project_1" {
  name = "project_1"
  description = "A project"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) A description of the project.

* `domain_id` - (Optional) The domain this project belongs to.

* `enabled` - (Optional) Whether the project is enabled or disabled. Valid
  values are `true` and `false`.

* `is_domain` - (Optional) Whether this project is a domain. Valid values
  are `true` and `false`.

* `name` - (Optional) The name of the project.

* `parent_id` - (Optional) The parent of this project.

* `region` - (Optional) The region in which to obtain the V3 Keystone client.
    If omitted, the `region` argument of the provider is used. Changing this
    creates a new User.

## Attributes Reference

The following attributes are exported:

* `domain_id` - See Argument Reference above.
* `parent_id` - See Argument Reference above.

## Import

Projects can be imported using the `id`, e.g.

```
$ terraform import telefonicaopencloud_identity_project_v3.project_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
