provider "google" {
  project = var.gcp_project_name
  region  = var.gcp_region
}

provider "google-beta" {
  project = var.gcp_project_name
  region  = var.gcp_region
}
