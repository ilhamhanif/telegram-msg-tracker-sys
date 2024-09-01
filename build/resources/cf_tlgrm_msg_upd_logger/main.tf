# Setup local variables
locals {
  sa_default_compute_engine           = "${var.project_number}-compute@developer.gserviceaccount.com"
  pubsub_topic_name                   = "telegram_msg_update_logger"
  pubsub_subscriber_name              = local.pubsub_topic_name
  pubsub_subscriber_ack_deadline      = 60
  pubsub_subscriber_expiration_policy = ""
  cf_name                             = "cf-tlgrm-msg-upd-logger"
  cf_entrypoint                       = "TelegramMsgUpdateLogger"
  cf_runtime                          = "go122"
  cf_service_account_id               = local.cf_name
  cf_service_account_name             = "Service Account dedicated for Cloud Function2: ${local.cf_name}"
  cf_service_account_roles = [
    "roles/bigquery.dataEditor"
  ]
  cf_configuration = {
    min_instance_count    = 0
    max_instance_count    = 1
    available_memory      = "128Mi"
    available_cpu         = 1
    timeout_seconds       = 60
    service_account_email = module.cloud_functions2_service_account.service_account_email
  }
  bq_dataset_name = "ops"
  bq_dataset_id   = local.bq_dataset_name
  bq_tables = [
    {
      table_id = "telegram_msg_log_update",
      schema   = <<EOF
      [
        { "name": "update_id", "type": "INTEGER", "mode": "NULLABLE" },
        { "name": "update", "type": "JSON", "mode": "NULLABLE" },
        { "name": "log_date", "type": "DATE", "mode": "NULLABLE" },
        { "name": "log_datetime", "type": "DATETIME", "mode": "NULLABLE" },
        { "name": "log_epoch", "type": "INTEGER", "mode": "NULLABLE" }
      ]
      EOF
      time_partitioning = {
        type                     = "DAY",
        field                    = "log_date",
        require_partition_filter = true,
        expiration_ms            = null,
      },
      range_partitioning = null,
      expiration_time    = null,
      clustering         = ["update_id", "log_epoch"],
      labels             = {}
    }
  ]
}

# Generates a ZIP compressed file archieve of the source code.
module "zip_cf2_gcs" {
  source = "../../modules/zip_cf2_gcs_bucket"

  cloud_functions2_name = local.cf_name
  gcs_zip_project_id    = var.project_id
}

# Create Pub/Sub Topic, and its Subscriber
# with Push method, authenticated with default compute engine service account.
module "pubsub" {
  source  = "terraform-google-modules/pubsub/google"
  version = "~> 6.0"

  topic      = local.pubsub_topic_name
  project_id = var.project_id
  push_subscriptions = [
    {
      name                       = local.pubsub_subscriber_name
      ack_deadline_seconds       = local.pubsub_subscriber_ack_deadline
      push_endpoint              = module.cloud_functions2.function_uri
      oidc_service_account_email = local.sa_default_compute_engine
      expiration_policy          = local.pubsub_subscriber_expiration_policy
    }
  ]
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

# Create BigQuery dataset and table
module "bigquery" {
  source  = "terraform-google-modules/bigquery/google"
  version = "~> 8.1"

  dataset_id          = local.bq_dataset_id
  dataset_name        = local.bq_dataset_name
  project_id          = var.project_id
  location            = var.region
  deletion_protection = false

  tables = local.bq_tables
}