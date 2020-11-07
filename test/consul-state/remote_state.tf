terraform {
  backend "consul" {
    path = "state"
    address = "127.0.0.1:8500"
    scheme = "http"
  }
  required_version = "~> 0.12.13"
}