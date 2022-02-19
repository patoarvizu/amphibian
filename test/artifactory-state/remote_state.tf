terraform {
  backend "artifactory" {
    username = "admin"
    password = "admin123"
    url = "http://localhost:8082/artifactory"
    repo = "example-repo-local"
    subpath = "/"
  }
}