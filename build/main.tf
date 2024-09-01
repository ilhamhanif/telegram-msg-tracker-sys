# Setup Backend Configuration
terraform {
  backend "local" {
    path = "state-telegram-msg-tracker-sys.tfstate"
  }
}

# Setup .locals
locals {
  api_services = [
    "compute.googleapis.com",
    "cloudbuild.googleapis.com",
    "run.googleapis.com",
    "cloudfunctions.googleapis.com",
    "bigquery.googleapis.com"
  ]
}

# Setup Project Wide Configuration
provider "google" {
  project = var.project_id
  region  = var.region
}

data "google_project" "gcp_project_var" {}

# Resources
# API Services
resource "google_project_service" "gcp_api_services" {
  count              = length(local.api_services)
  project            = var.project_id
  service            = local.api_services[count.index]
  disable_on_destroy = false
}

# Google Cloud Storage - Store All Clouc Function2 ZIP files
resource "google_storage_bucket" "bucket_cf_zip_source_code" {
  name     = "${var.project_id}-gcf-zip-source"
  location = var.region

  lifecycle_rule {
    condition {
      age = "1"
    }
    action {
      type = "Delete"
    }
  }
}

# Wait for All API 
# and Initialization successfully enabled
resource "null_resource" "resource_api_activation_complete" {
  depends_on = [
    google_project_service.gcp_api_services,
    google_storage_bucket.bucket_cf_zip_source_code
  ]
}

# Cloud Function2 - Telegram Message Forwarder
module "cf_tlgrm_msg_upd_forwarder" {
  source     = "./resources/cf_tlgrm_msg_upd_forwarder"
  depends_on = [null_resource.resource_api_activation_complete]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

# Cloud Function2 - Telegram Message Orchestrator
module "cf_tlgrm_msg_orchestrator" {
  source     = "./resources/cf_tlgrm_msg_orchestrator"
  depends_on = [null_resource.resource_api_activation_complete]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

# Cloud Function2 - Telegram Message Logger
module "cf_tlgrm_msg_upd_logger" {
  source     = "./resources/cf_tlgrm_msg_upd_logger"
  depends_on = [null_resource.resource_api_activation_complete]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}
