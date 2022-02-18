terraform {
  backend "pg" {
    conn_str = "postgres://postgres:postgres123@localhost:5432/terraform_backend?sslmode=disable"
  }
}