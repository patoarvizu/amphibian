terraform {
  backend "etcdv3" {
    endpoints = ["localhost:2379"]
    username = "root"
    password = "root123"
  }
}