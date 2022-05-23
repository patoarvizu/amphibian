variable image_version {
  type = string
  default = "latest"
  description = "The label of the image to run."
}

variable create_namespace {
  type = bool
  default = true
  description = "If true, a new namespace will be created with the name set to the value of the namespace_name variable. If false, it will look up an existing namespace with the name of the value of the namespace_name variable."
}

variable namespace_name {
  type = string
  default = "amp"
  description = "The name of the namespace to create or look up."
}

variable terraform_binary_init_container_image {
  type = string
  default = "alpine:3.15.0"
  description = "The image to use for the init container that installs the target `terraform` binary."
}

variable terraform_binary_version {
  type = string
  default = "1.1.2"
  description = "The version of the `terraform` binary. Note that it's not possible to use `latest`, or use 'partial' versions (e.g. `1`, or `1.1`) so you have to specify the full version."
}

variable terraform_binary_operating_system {
  type = string
  default = "linux"
  description = "The operating system for which to download the `terraform` binary."
}

variable terraform_binary_arch {
  type = string
  default = "amd64"
  description = "The architecture for which to download the `terraform` binary."
}

variable watch_namespace {
  type = string
  default = ""
  description = "The value to be set on the `WATCH_NAMESPACE` environment variable."
}

variable enable_prometheus_monitoring {
  type = bool
  default = false
  description = "Create the `Service` and `ServiceMonitor` objects to enable Prometheus monitoring on the operator."
}

variable auth_env_vars {
  type = list(object({
    name = string
    value = string
  }))
  default = []
  description = "Environment variables required for remote state backend authentication."
}

variable auth_env_from_vars {
  type = list(object({
    name = string
    secret_ref_key = string
    secret_ref_name = string
  }))
  default = []
  description = "Environment variables required for remote state backend authentication."
}