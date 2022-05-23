resource "kubernetes_cluster_role" "amphibian_manager_role" {
  metadata {
    name = "amphibian-manager-role"
  }

  rule {
    verbs      = ["create", "get", "list", "patch", "update", "watch"]
    api_groups = [""]
    resources  = ["configmaps"]
  }

  rule {
    verbs      = ["create", "delete", "get", "list", "patch", "update", "watch"]
    api_groups = ["terraform.patoarvizu.dev"]
    resources  = ["terraformstates"]
  }

  rule {
    verbs      = ["get", "patch", "update"]
    api_groups = ["terraform.patoarvizu.dev"]
    resources  = ["terraformstates/status"]
  }

  rule {
    verbs      = ["create", "get", "list", "patch", "update", "watch"]
    api_groups = [""]
    resources  = ["secrets"]
  }
}

resource "kubernetes_cluster_role_binding" "amphibian_manager_rolebinding" {
  metadata {
    name = "amphibian-manager-rolebinding"
  }

  subject {
    kind      = "ServiceAccount"
    name      = "amphibian"
    namespace = var.create_namespace ? kubernetes_namespace.ns[var.namespace_name].metadata[0].name : data.kubernetes_namespace.ns[var.namespace_name].metadata[0].name
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = "amphibian-manager-role"
  }
}

resource "kubernetes_role" "leader_election_role" {
  metadata {
    name = "leader-election-role"
    namespace = var.create_namespace ? kubernetes_namespace.ns[var.namespace_name].metadata[0].name : data.kubernetes_namespace.ns[var.namespace_name].metadata[0].name
  }

  rule {
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
    api_groups = [""]
    resources  = ["configmaps"]
  }

  rule {
    verbs      = ["get", "update", "patch"]
    api_groups = [""]
    resources  = ["configmaps/status"]
  }

  rule {
    verbs      = ["create", "patch"]
    api_groups = [""]
    resources  = ["events"]
  }
}

resource "kubernetes_role_binding" "leader_election_rolebinding" {
  metadata {
    name = "leader-election-rolebinding"
  }

  subject {
    kind      = "ServiceAccount"
    name      = "amphibian"
    namespace = var.create_namespace ? kubernetes_namespace.ns[var.namespace_name].metadata[0].name : data.kubernetes_namespace.ns[var.namespace_name].metadata[0].name
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "Role"
    name      = "leader-election-role"
  }
}