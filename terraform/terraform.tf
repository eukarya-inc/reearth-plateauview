terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.8"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 5.8"
    }
    random = {
      source = "hashicorp/random"
    }
    auth0 = {
      source  = "auth0/auth0"
      version = "1.1.1"
    }
  }
  required_version = ">= 1.6.5"

  backend "gcs" {
    bucket = "plateau-test2-terraform-tfstate"
  }
}
