# [1] Setup Backend Configuration
terraform {
  backend "local" {
    path = "state-telegram-msg-tracker-sys.tfstate"
  }
}

# [2] Build Pre-Resources
# [2.1] Google Cloud Project - Config Variable
provider "google" {
  project = var.project_id
  region  = var.region
}

data "google_project" "gcp_project_var" {}

# [2.2] Google Cloud Project - API
module "gcp_api" {
  source = "./resources/gcp_api"

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

# [2.3] BigQuery - Dataset OPS
module "bq_dataset_ops" {
  source     = "./resources/bq_dataset_ops"
  depends_on = [module.gcp_api]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

# [2.4] Cloud Function2 - Bucket Notif Logger
module "cf_gcs_bucket_notif_logger" {
  source     = "./resources/cf_gcs_bucket_notif_logger"
  depends_on = [module.gcp_api]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

# [2.5] Google Cloud Storage - Buckets
# All bucket has bucket notif configuration.
module "gcp_gcs_bucket" {
  source = "./resources/gcp_gcs_bucket"

  project_id                     = var.project_id
  project_number                 = data.google_project.gcp_project_var.number
  region                         = var.region
  pubsub_bucket_notif_topic_name = module.cf_gcs_bucket_notif_logger.pubsub_topic_name
}

# [2.6] Cloud Function2 - Dead Letter Logger
module "cf_ps_dead_letter_logger" {
  source     = "./resources/cf_ps_dead_letter_logger"
  depends_on = [module.gcp_api]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

# [2.7] Wait for Pre-Resources to finish
resource "null_resource" "pre_resources_build_complete" {
  depends_on = [
    module.bq_dataset_ops,
    module.gcp_gcs_bucket,
    module.cf_ps_dead_letter_logger
  ]
}

# [3] Resources Stack
# [3.1] Cloud Function2 - Utils - File Downloader
module "cf_tlgrm_utils_file_downloader" {
  source = "./resources/cf_tlgrm_utils_file_downloader"
  depends_on = [
    null_resource.pre_resources_build_complete,
    module.cf_ps_dead_letter_logger
  ]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

# [3.2] Cloud Function2 - Telegram Message Forwarder
module "cf_tlgrm_msg_upd_forwarder" {
  source = "./resources/cf_tlgrm_msg_upd_forwarder"
  depends_on = [
    null_resource.pre_resources_build_complete,
    module.cf_ps_dead_letter_logger
  ]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

# [3.3] Cloud Function2 - Telegram Message Identificator
module "cf_tlgrm_msg_upd_identificator" {
  source = "./resources/cf_tlgrm_msg_upd_identificator"
  depends_on = [
    null_resource.pre_resources_build_complete,
    module.cf_ps_dead_letter_logger
  ]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

# [3.4] Cloud Function2 - Telegram Message Logger
module "cf_tlgrm_msg_upd_logger" {
  source = "./resources/cf_tlgrm_msg_upd_logger"
  depends_on = [
    null_resource.pre_resources_build_complete,
    module.cf_ps_dead_letter_logger
  ]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}

# [3.5] Cloud Function2 - Telegram Action - Send Message
module "cf_tlgrm_act_send_message" {
  source = "./resources/cf_tlgrm_act_send_message"
  depends_on = [
    null_resource.pre_resources_build_complete,
    module.cf_ps_dead_letter_logger
  ]

  project_id     = var.project_id
  project_number = data.google_project.gcp_project_var.number
  region         = var.region
}
