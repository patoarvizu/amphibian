# Amphibian

<!-- TOC -->

- [Amphibian](#amphibian)
  - [Intro](#intro)
  - [Design](#design)
  - [Configuration](#configuration)
    - [Backends](#backends)
      - [Remote (Terraform Cloud)](#remote-terraform-cloud)
      - [S3](#s3)
      - [Consul](#consul)
    - [Target](#target)
      - [Values](#values)
  - [For security nerds](#for-security-nerds)
    - [Docker images are signed and published to Docker Hub's Notary server](#docker-images-are-signed-and-published-to-docker-hubs-notary-server)
    - [Docker images are labeled with Git and GPG metadata](#docker-images-are-labeled-with-git-and-gpg-metadata)
  - [Multi-architecture images](#multi-architecture-images)

<!-- /TOC -->

## Intro

The adoption of Terraform in many organizations predates the adoption of Kubernetes, or in some cases they're separate efforts owned by different teams. In addition to that, the integration between both systems consists of manual copy/pasting of values, since there's no clearly defined discovery mechanism between the two.

Just like [amphibians](https://en.wikipedia.org/wiki/Amphibian) can inhabit both land and water, this project aims to close the interface gap between Terraform outputs and Kubernetes configuration discovery. The existing [terraform-helm](https://github.com/hashicorp/terraform-helm) and [aws-controllers-k8s](https://github.com/aws/aws-controllers-k8s) projects don't yet have the full functionality and flexibility that Amphibian provides.

A [Custom resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) of kind `TerraformState` deployed on Kubernetes clusters will create a new `ConfigMap` or `Secret` and populate it with the [outputs](https://www.terraform.io/docs/configuration/outputs.html) of the corresponding remote Terraform state.

## Design

Even though Terraform has a `struct` for capturing a module's [output values](https://github.com/hashicorp/terraform/blob/v0.13.5/states/output_value.go) programmatically, that API can't be considered public and guaranteed.

Since the only guaranteed interface is the command line, the way this controller gets the outputs from the remote state is by creating a `data.tf` and an `outputs.tf` file, running `terraform apply`, followed by `terraform output -json`, and then unmarshaling that output back into a Go `struct`. The controller then uses those outputs to create a new `ConfigMap` or `Secret` in the location defined by `target`, that has the exact contents returned by Terraform.

**Note:** Keep in mind that Terraform only returns the [root-level outputs](https://registry.terraform.io/providers/hashicorp/terraform/latest/docs/data-sources/remote_state#root-outputs-only). If you need to consume the outputs of a submodule, you'll have to expose it all the way to the root level so they can be discovered in Kubernetes.

Since a single controller can handle multiple `TerraformState` objects, each different set of `data.tf`/`outputs.tf` set of files is created in a subdirectory corresponding to its namespace and name, to avoid collisions and overwrites between resources. However, this should generally be transparent to the end user, and the choice of name or namespace shouldn't have any effect.

## Configuration

The configuration required to discover the state will depend on the type of backend (defined by the `type` field), and will match the options available to each backend kind. Any exceptions or additional configuration requirements are noted in the corresponding section below. Keep in mind that any options that can be supplied via environment variables will also be honored. In addition to setting the `type`, a configuration block corresponding to the given type needs to be provided.

### Backends

#### Remote (Terraform Cloud)

- [Documentation](https://www.terraform.io/docs/backends/types/remote.html)
- `type: remote`
- Configuration block name: `remoteConfig`

The only additional option required for this backend type is the Terraform Cloud token. This needs to be injected as an environment variable called `TERRAFORM_CLOUD_TOKEN`.

#### S3

- [Documentation](https://www.terraform.io/docs/backends/types/s3.html)
- `type: s3`
- Configuration block name: `s3Config`

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

#### Consul

- [Documentation](https://www.terraform.io/docs/backends/types/consul.html)
- `type: consul`
- Configuration block name: `consulConfig`

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

The `target` field represents the location and type of object where the outputs from the upstream state will be projected.

- `type`: The type of object where the outputs will be projected. It supports either `configmap` or `secret` (all lowercase in both cases).
- `name`: The name of either the `ConfigMap` or the `Secret` that will hold the `outputs` map.

#### Values

The values of each field in the projected ConfigMap will depend on the output type. If it's a `string`, it'll be set as-is on the ConfigMap. If it's a `map` or `list`, it'll be set to it's JSON-ified string, as generated by `json.Marshal()`.

## For security nerds

### Docker images are signed and published to Docker Hub's Notary server

The [Notary](https://github.com/theupdateframework/notary) project is a CNCF incubating project that aims to provide trust and security to software distribution. Docker Hub runs a Notary server at https://notary.docker.io for the repositories it hosts.

[Docker Content Trust](https://docs.docker.com/engine/security/trust/content_trust/) is the mechanism used to verify digital signatures and enforce security by adding a validating layer.

You can inspect the signed tags for this project by doing `docker trust inspect --pretty docker.io/patoarvizu/amphibian`, or (if you already have `notary` installed) `notary -d ~/.docker/trust/ -s https://notary.docker.io list docker.io/patoarvizu/amphibian`.

If you run `docker pull` with `DOCKER_CONTENT_TRUST=1`, the Docker client will only pull images that come from registries that have a Notary server attached (like Docker Hub).

### Docker images are labeled with Git and GPG metadata

In addition to the digital validation done by Docker on the image itself, you can do your own human validation by making sure the image's content matches the Git commit information (including tags if there are any) and that the GPG signature on the commit matches the key on the commit on github.com.

For example, if you run `docker pull patoarvizu/amphibian:054e78a77b4923dd8fbd1ace79714152024ee8c4` to pull the image tagged with that commit id, then run `docker inspect patoarvizu/amphibian:054e78a77b4923dd8fbd1ace79714152024ee8c4 | jq -r '.[0].Config.Labels'` (assuming you have [jq](https://stedolan.github.io/jq/) installed) you should see that the `GIT_COMMIT` label matches the tag on the image. Furthermore, if you go to https://github.com/patoarvizu/amphibian/commit/054e78a77b4923dd8fbd1ace79714152024ee8c4 (notice the matching commit id), and click on the **Verified** button, you should be able to confirm that the GPG key ID used to match this commit matches the value of the `SIGNATURE_KEY` label, and that the key belongs to the `AUTHOR_EMAIL` label. When an image belongs to a commit that was tagged, it'll also include a `GIT_TAG` label, to further validate that the image matches the content.

Keep in mind that this isn't tamper-proof. A malicious actor with access to publish images can create one with malicious content but with values for the labels matching those of a valid commit id. However, when combined with Docker Content Trust, the certainty of using a legitimate image is increased because the chances of a bad actor having both the credentials for publishing images, as well as Notary signing credentials are significantly lower and even in that scenario, compromised signing keys can be revoked or rotated.

Here's the list of included Docker labels:

- `AUTHOR_EMAIL`
- `COMMIT_TIMESTAMP`
- `GIT_COMMIT`
- `GIT_TAG`
- `SIGNATURE_KEY`

## Multi-architecture images

Manifests published with the semver tag (e.g. `patoarvizu/amphibian:v0.0.0`), as well as `latest` are multi-architecture manifest lists. In addition to those, there are architecture-specific tags that correspond to an image manifest directly, tagged with the corresponding architecture as a suffix, e.g. `v0.0.0-amd64`. Both types (image manifests or manifest lists) are signed with Notary as described above.

Here's the list of architectures the images are being built for, and their corresponding suffixes for images:

- `linux/amd64`, `-amd64`
- `linux/arm64`, `-arm64`
- `linux/arm/v7`, `-arm7`