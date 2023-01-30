terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.50"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 4.50"
    }
    random = {
      source = "hashicorp/random"
    }
  }
  required_version = ">= v1.3.7"
}
