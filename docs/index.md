# amphibian

![Version: 0.0.5](https://img.shields.io/badge/Version-0.0.5-informational?style=flat-square)

Amphibian

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| authEnvVars | string | `nil` | Environment variables required for remote state backend authentication. This is a slice of [`v1.EnvVar`](https://pkg.go.dev/k8s.io/api/core/v1#EnvVar)s. |
| imagePullPolicy | string | `"IfNotPresent"` | The imagePullPolicy to be used on the operator. |
| imageVersion | string | `"latest"` | The image version used for the operator. |
| prometheusMonitoring.enable | bool | `false` | Create the `Service` and `ServiceMonitor` objects to enable Prometheus monitoring on the operator. |
| prometheusMonitoring.serviceMonitor.customLabels | string | `nil` | Custom labels to add to the ServiceMonitor object. |
| rbac.clusterRoleSecretsAccessRules | list | `[{"apiGroups":[""],"resources":["secrets"],"verbs":["create","get","list","patch","update","watch"]}]` | List of `PolicyRule`s for accessing Kubernetes secrets, to be appended to the `amphibian-manager-role` cluster role. |
| resources | object | `nil` | The resources requests/limits to be set on the deployment pod spec template. |
| terraformBinary | object | `{"arch":"amd64","initContainerImage":"alpine:3.15.0","operatingSystem":"linux","version":"1.1.2"}` | Information about the `terraform` binary to inject into the main container. These values will be used to download the binary from `https://releases.hashicorp.com/terraform/<terraformVersion.version>/terraform_<terraformVersion.version>_<terraformVersion.operatingSystem>_<terraformVersion.arch>.zip`. |
| terraformBinary.arch | string | `"amd64"` | The architecture for which to download the `terraform` binary. |
| terraformBinary.initContainerImage | string | `"alpine:3.15.0"` | The image to use for the init container that installs the target `terraform` binary. |
| terraformBinary.operatingSystem | string | `"linux"` | The operating system for which to download the `terraform` binary. |
| terraformBinary.version | string | `"1.1.2"` | The version of the `terraform` binary. Note that it's not possible to use `latest`, or use "partial" versions (e.g. `1`, or `1.1`) so you have to specify the full version. |
| volumeMounts | string | `nil` | List of [`v1.VolumeMount`](https://pkg.go.dev/k8s.io/api/core/v1#VolumeMount) objects to be appended as-is to the amphibian workloads |
| volumes | string | `nil` |  |
| watchNamespace | string | `""` | The value to be set on the `WATCH_NAMESPACE` environment variable. |
