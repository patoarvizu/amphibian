resource "kubernetes_service_account" "amphibian" {
  metadata {
    name = "amphibian"
    namespace = var.create_namespace ? kubernetes_namespace.ns[var.namespace_name].metadata[0].name : data.kubernetes_namespace.ns[var.namespace_name].metadata[0].name
  }
}