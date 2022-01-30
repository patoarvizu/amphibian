terraform {
  backend "kubernetes" {
    secret_suffix = "state"
  }
}