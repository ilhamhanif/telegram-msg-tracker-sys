# Setup Backend Configuration
terraform {
  backend "local" {
    path = "state-telegram-msg-tracker-sys.tfstate"
  }
}

# Google Cloud Project - Config Variable
provider "google" {
  project = var.project_id
  region  = var.region
}

data "google_project" "gcp_project_var" {}

# Pre-Resources
## Google Cloud Project - API
module "gcp_api" {
  source = "./resources/gcp_api"

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

## BigQuery - Dataset OPS
module "bq_dataset_ops" {
  source     = "./resources/bq_dataset_ops"
  depends_on = [module.gcp_api]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

## Google Cloud Storage - Buckets
module "gcp_gcs_bucket" {
  source     = "./resources/gcp_gcs_bucket"
  depends_on = [module.gcp_api]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

## Wait for Pre-Resources to finish
resource "null_resource" "pre_resources_build_complete" {
  depends_on = [
    module.bq_dataset_ops,
    module.gcp_gcs_bucket
  ]
}

# Resources Stack
## Cloud Function2 - Dead Letter Logger
module "cf_ps_dead_letter_logger" {
  source     = "./resources/cf_ps_dead_letter_logger"
  depends_on = [null_resource.pre_resources_build_complete]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

## Cloud Function2 - Utils - File Downloader
module "cf_tlgrm_utils_file_downloader" {
  source     = "./resources/cf_tlgrm_utils_file_downloader"
  depends_on = [null_resource.pre_resources_build_complete]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

## Cloud Function2 - Telegram Message Forwarder
module "cf_tlgrm_msg_upd_forwarder" {
  source     = "./resources/cf_tlgrm_msg_upd_forwarder"
  depends_on = [null_resource.pre_resources_build_complete]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

## Cloud Function2 - Telegram Message Identificator
module "cf_tlgrm_msg_upd_identificator" {
  source     = "./resources/cf_tlgrm_msg_upd_identificator"
  depends_on = [null_resource.pre_resources_build_complete]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

## Cloud Function2 - Telegram Message Logger
module "cf_tlgrm_msg_upd_logger" {
  source     = "./resources/cf_tlgrm_msg_upd_logger"
  depends_on = [null_resource.pre_resources_build_complete]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

## Cloud Function2 - Telegram Action - Send Message
module "cf_tlgrm_act_send_message" {
  source     = "./resources/cf_tlgrm_act_send_message"
  depends_on = [null_resource.pre_resources_build_complete]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}
