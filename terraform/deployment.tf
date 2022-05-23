resource "kubernetes_deployment" "amphibian" {
  metadata {
    name = "amphibian"
    namespace = var.create_namespace ? kubernetes_namespace.ns[var.namespace_name].metadata[0].name : data.kubernetes_namespace.ns[var.namespace_name].metadata[0].name

    labels = {
      app = "amphibian"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "amphibian"
      }
    }

    template {
      metadata {
        labels = {
          app = "amphibian"
        }
      }

      spec {
        volume {
          name = "terraform"

          empty_dir {
            medium = "Memory"
          }
        }

        volume {
          name = "terraform-bin"

          empty_dir {
            medium = "Memory"
          }
        }

        init_container {
          name    = "install-terraform"
          image   = var.terraform_binary_init_container_image
          command = ["sh", "-c", "wget https://releases.hashicorp.com/terraform/${var.terraform_binary_version}/terraform_${var.terraform_binary_version}_${var.terraform_binary_operating_system}_${var.terraform_binary_arch}.zip && unzip terraform_${var.terraform_binary_version}_${var.terraform_binary_operating_system}_${var.terraform_binary_arch}.zip && cp terraform /terraform-bin/"]

          volume_mount {
            name       = "terraform-bin"
            mount_path = "/terraform-bin"
          }
        }

        container {
          name    = "manager"
          image   = "patoarvizu/amphibian:${var.image_version}"
          command = ["/manager"]
          args    = ["--enable-leader-election"]

          port {
            name           = "http-metrics"
            container_port = 8080
          }

          env {
            name  = "TF_CLI_CONFIG_FILE"
            value = "/terraform/.terraformrc"
          }

          env {
            name = "WATCH_NAMESPACE"
            value = var.watch_namespace
          }

          dynamic "env" {
            for_each = var.auth_env_vars
            content {
              name = env.value["name"]
              value = env.value["value"]
            }
          }

          dynamic "env" {
            for_each = var.auth_env_from_vars
            content {
              name = env.value["name"]
              value_from {
                secret_key_ref {
                  key = env.value["secret_ref_key"]
                  name = env.value["secret_ref_name"]
                }
              }
            }
          }

          volume_mount {
            name       = "terraform"
            mount_path = "/terraform"
          }

          volume_mount {
            name       = "terraform-bin"
            mount_path = "/terraform-bin"
          }

          image_pull_policy = "IfNotPresent"
        }

        service_account_name = kubernetes_service_account.amphibian.metadata[0].name
      }
    }
  }
}

