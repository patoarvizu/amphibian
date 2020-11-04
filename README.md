# Amphibian

<!-- TOC -->

- [Amphibian](#amphibian)
  - [Intro](#intro)
  - [Configuration](#configuration)
    - [Backends](#backends)
      - [[Remote (Terraform Cloud)](https://www.terraform.io/docs/backends/types/remote.html) (`type: remote`, `remoteConfig`)](#remote-terraform-cloudhttpswwwterraformiodocsbackendstypesremotehtml-type-remote-remoteconfig)
      - [[S3](https://www.terraform.io/docs/backends/types/s3.html) (`type: s3`, `s3Config`)](#s3httpswwwterraformiodocsbackendstypess3html-type-s3-s3config)
      - [[Consul](https://www.terraform.io/docs/backends/types/consul.html) (`type: consul`, `consulConfig`)](#consulhttpswwwterraformiodocsbackendstypesconsulhtml-type-consul-consulconfig)
    - [Target](#target)

<!-- /TOC -->

## Intro

Similar to how [amphibians](https://en.wikipedia.org/wiki/Amphibian) can inhabit both land and water, this project aims to close the interface gap between Terraform outputs and Kubernetes configuration discovery. The existing [terraform-helm](https://github.com/hashicorp/terraform-helm) and [aws-controllers-k8s](https://github.com/aws/aws-controllers-k8s) projects don't yet have the full functionality and flexibility that Amphibian provides.

A [Custom resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) of kind `TerraformState` deployed on Kubernetes clusters will create a new `ConfigMap` and populate it with the [outputs](https://www.terraform.io/docs/configuration/outputs.html) of the corresponding remote Terraform state.

## Configuration

The configuration required to discover the state will depend on the type of backend (defined by the `type` field), and will match the options available to each backend kind. Any exceptions or additional configuration requirements are noted in the corresponding section below. Keep in mind that any options that can be supplied via environment variables will also be honored. In addition to setting the `type`, a configuration block corresponding to the given type needs to be provided.

### Backends

#### [Remote (Terraform Cloud)](https://www.terraform.io/docs/backends/types/remote.html) (`type: remote`, `remoteConfig`)

The only additional option required for this backend type is the Terraform Cloud token. This needs to be injected as an environment variable called `TERRAFORM_CLOUD_TOKEN`.

#### [S3](https://www.terraform.io/docs/backends/types/s3.html) (`type: s3`, `s3Config`)

The following fields can be alternatively be set as environment variables (as documented in the link above):

- `region` (`AWS_DEFAULT_REGION`/`AWS_REGION`)
- `access_key` (`AWS_ACCESS_KEY_ID`)
- `secret_key` (`AWS_SECRET_ACCESS_KEY`)
- `iam_endpoint` (`AWS_IAM_ENDPOINT`)
- `profile` (`AWS_PROFILE`)
- `sts_endpoint` (`AWS_STS_ENDPOINT`)
- `sse_customer_key` (`AWS_SSE_CUSTOMER_KEY`)

Additionally, in the case of `access_key` and `secret_key`, they can also be set via other mechanisms like shared credentials files, EC2 instance profiles, etc., as officially [documented](https://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/credentials.html).

Lastly, the following options are not available since they're irrelevant for looking up remote states:

- `acl`
- `encrypt`
- `dynamodb_endpoint`
- `dynamodb_table`

#### [Consul](https://www.terraform.io/docs/backends/types/consul.html) (`type: consul`, `consulConfig`)

The following fields can be alternatively be set as environment variables (as documented in the link above):

- `access_token` (`CONSUL_HTTP_TOKEN`)
- `address` (`CONSUL_HTTP_ADDR`)
- `scheme` (`CONSUL_HTTP_SSL`)
- `http_auth` (`CONSUL_HTTP_AUTH`)
- `ca_file` (`CONSUL_CACERT`)
- `cert_file` (`CONSUL_CLIENT_CERT`)
- `key_file` (`CONSUL_CLIENT_KEY`)

Additionally, the following options are not available since they're irrelevant for looking up remote states:

- `gzip`
- `lock`

### Target

The `target` field represents the location where the outputs from the upstream state will be projected.

- `configMapName`: The name of the `ConfigMap` that will hold the `outputs` map.
- `namespace`: The name where the `ConfigMap` above will be placed.