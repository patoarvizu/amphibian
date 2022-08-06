<!-- BEGIN_TF_DOCS -->

## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 0.14.9 |
| <a name="requirement_kubernetes"></a> [kubernetes](#requirement\_kubernetes) | ~> 2.7.1 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_kubernetes"></a> [kubernetes](#provider\_kubernetes) | ~> 2.7.1 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [kubernetes_cluster_role.amphibian_manager_role](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/cluster_role) | resource |
| [kubernetes_cluster_role_binding.amphibian_manager_rolebinding](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/cluster_role_binding) | resource |
| [kubernetes_deployment.amphibian](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/deployment) | resource |
| [kubernetes_manifest.customresourcedefinition_terraformstates_terraform_patoarvizu_dev](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/manifest) | resource |
| [kubernetes_manifest.servicemonitor_amphibian_metrics](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/manifest) | resource |
| [kubernetes_namespace.ns](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/namespace) | resource |
| [kubernetes_role.leader_election_role](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/role) | resource |
| [kubernetes_role_binding.leader_election_rolebinding](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/role_binding) | resource |
| [kubernetes_service.amphibian_metrics](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/service) | resource |
| [kubernetes_service_account.amphibian](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/service_account) | resource |
| [kubernetes_namespace.ns](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/data-sources/namespace) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_auth_env_from_vars"></a> [auth\_env\_from\_vars](#input\_auth\_env\_from\_vars) | Environment variables required for remote state backend authentication. | <pre>list(object({<br>    name = string<br>    secret_ref_key = string<br>    secret_ref_name = string<br>  }))</pre> | `[]` | no |
| <a name="input_auth_env_vars"></a> [auth\_env\_vars](#input\_auth\_env\_vars) | Environment variables required for remote state backend authentication. | <pre>list(object({<br>    name = string<br>    value = string<br>  }))</pre> | `[]` | no |
| <a name="input_create_namespace"></a> [create\_namespace](#input\_create\_namespace) | If true, a new namespace will be created with the name set to the value of the namespace\_name variable. If false, it will look up an existing namespace with the name of the value of the namespace\_name variable. | `bool` | `true` | no |
| <a name="input_enable_prometheus_monitoring"></a> [enable\_prometheus\_monitoring](#input\_enable\_prometheus\_monitoring) | Create the `Service` and `ServiceMonitor` objects to enable Prometheus monitoring on the operator. | `bool` | `false` | no |
| <a name="input_image_version"></a> [image\_version](#input\_image\_version) | The label of the image to run. | `string` | `"latest"` | no |
| <a name="input_namespace_name"></a> [namespace\_name](#input\_namespace\_name) | The name of the namespace to create or look up. | `string` | `"amp"` | no |
| <a name="input_service_monitor_custom_labels"></a> [service\_monitor\_custom\_labels](#input\_service\_monitor\_custom\_labels) | Custom labels to add to the `ServiceMonitor` object. | `map` | `{}` | no |
| <a name="input_terraform_binary_arch"></a> [terraform\_binary\_arch](#input\_terraform\_binary\_arch) | The architecture for which to download the `terraform` binary. | `string` | `"amd64"` | no |
| <a name="input_terraform_binary_init_container_image"></a> [terraform\_binary\_init\_container\_image](#input\_terraform\_binary\_init\_container\_image) | The image to use for the init container that installs the target `terraform` binary. | `string` | `"alpine:3.15.0"` | no |
| <a name="input_terraform_binary_operating_system"></a> [terraform\_binary\_operating\_system](#input\_terraform\_binary\_operating\_system) | The operating system for which to download the `terraform` binary. | `string` | `"linux"` | no |
| <a name="input_terraform_binary_version"></a> [terraform\_binary\_version](#input\_terraform\_binary\_version) | The version of the `terraform` binary. Note that it's not possible to use `latest`, or use 'partial' versions (e.g. `1`, or `1.1`) so you have to specify the full version. | `string` | `"1.1.2"` | no |
| <a name="input_watch_namespace"></a> [watch\_namespace](#input\_watch\_namespace) | The value to be set on the `WATCH_NAMESPACE` environment variable. | `string` | `""` | no |

## Outputs

No outputs.
<!-- END_TF_DOCS -->