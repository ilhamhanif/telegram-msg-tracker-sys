# Setup local variables
locals {
  sa_default_compute_engine = "${var.project_number}-compute@developer.gserviceaccount.com"
  cf_name                   = "cf-tlgrm-msg-upd-forwarder"
  cf_entrypoint             = "TelegramMsgUpdateForwarder"
  cf_runtime                = "go122"
  cf_service_account_id     = local.cf_name
  cf_service_account_name   = "Service Account dedicated for Cloud Function2: ${local.cf_name}"
  cf_service_account_roles = [
    "roles/pubsub.publisher"
  ]
  cf_configuration = {
    min_instance_count    = 0
    max_instance_count    = 1
    available_memory      = "128Mi"
    available_cpu         = 1
    timeout_seconds       = 60
    service_account_email = module.cloud_functions2_service_account.service_account_email
  }
}

# Generates a ZIP compressed file archieve of the source code.
module "zip_cf2_gcs" {
  source = "../../modules/gcp_cf2_zip_to_gcs"

  cloud_functions2_name = local.cf_name
  gcs_zip_project_id    = var.project_id
}

# Create Cloud Function Gen2
# with custom Service account
module "cloud_functions2_service_account" {
  source = "../../modules/gcp_sa_creator"

  project_id            = var.project_id
  service_account_id    = local.cf_service_account_id
  service_account_name  = local.cf_service_account_name
  service_account_roles = local.cf_service_account_roles
}

module "cloud_functions2" {
  source  = "GoogleCloudPlatform/cloud-functions/google"
  version = "~> 0.6"

  project_id        = var.project_id
  function_name     = local.cf_name
  function_location = var.region
  runtime           = local.cf_runtime
  entrypoint        = local.cf_entrypoint
  storage_source = {
    bucket     = module.zip_cf2_gcs.zip_gcs_bucket_name
    object     = module.zip_cf2_gcs.zip_gcs_object_name
    generation = module.zip_cf2_gcs.zip_gcs_object_generation
  }
  service_config = local.cf_configuration
}

resource "google_cloud_run_service_iam_binding" "cloud_functions2_allUsers_binding" {
  project  = var.project_id
  location = var.region
  service  = module.cloud_functions2.function_name
  role     = "roles/run.invoker"
  members = [
    "allUsers",
  ]
}
