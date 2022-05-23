resource "kubernetes_service" "amphibian_metrics" {
  for_each = var.enable_prometheus_monitoring ? {"monitor": true} : {}
  metadata {
    name = "amphibian-metrics"
    namespace = var.create_namespace ? kubernetes_namespace.ns[var.namespace_name].metadata[0].name : data.kubernetes_namespace.ns[var.namespace_name].metadata[0].name

    labels = {
      app = "amphibian"
    }
  }

  spec {
    port {
      name        = "http-metrics"
      protocol    = "TCP"
      port        = 8080
      target_port = "http-metrics"
    }

    selector = {
      app = "amphibian"
    }

    type = "ClusterIP"
  }
}

resource "kubernetes_manifest" "servicemonitor_amphibian_metrics" {
  for_each = var.enable_prometheus_monitoring ? {"monitor": true} : {}
  manifest = {
    apiVersion = "monitoring.coreos.com/v1"
    kind = "ServiceMonitor"
    metadata = {
      name = "amphibian-metrics"
      namespace = var.create_namespace ? kubernetes_namespace.ns[var.namespace_name].metadata[0].name : data.kubernetes_namespace.ns[var.namespace_name].metadata[0].name
    }
    spec = {
      endpoints = [
        {
          path = "/metrics"
          port = "http-metrics"
        },
      ]
      selector = {
        matchLabels = {
          app = "amphibian"
        }
      }
    }
  }
}