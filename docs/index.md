# amphibian

![Version: 0.0.1](https://img.shields.io/badge/Version-0.0.1-informational?style=flat-square)

Amphibian

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| authEnvVars | string | `nil` | Environment variables required for remote state backend authentication. This is a slice of [`v1.EnvVar`](https://pkg.go.dev/k8s.io/api/core/v1#EnvVar)s. |
| imagePullPolicy | string | `"IfNotPresent"` | The imagePullPolicy to be used on the operator. |
| imageVersion | string | `"latest"` | The image version used for the operator. |
| prometheusMonitoring.enable | bool | `false` | Create the `Service` and `ServiceMonitor` objects to enable Prometheus monitoring on the operator. |
| rbac.clusterRoleSecretsAccessRules | list | `[{"apiGroups":[""],"resources":["secrets"],"verbs":["create","get","list","patch","update","watch"]}]` | List of `PolicyRule`s for accessing Kubernetes secrets, to be appended to the `amphibian-manager-role` cluster role. |
| resources | object | `nil` | The resources requests/limits to be set on the deployment pod spec template. |
| watchNamespace | string | `""` | The value to be set on the `WATCH_NAMESPACE` environment variable. |
