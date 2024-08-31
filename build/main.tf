# Setup Backend Configuration
terraform {
  backend "local" {
    path = "state-telegram-msg-tracker-sys.tfstate"
  }
}

# Setup Project Wide Configuration
provider "google" {
  project = var.project_id
  region  = var.region
}

data "google_project" "gcp_project_var" {}
