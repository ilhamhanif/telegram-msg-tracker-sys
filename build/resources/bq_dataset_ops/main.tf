locals {
  bq_dataset_name = "ops"
  bq_dataset_id   = local.bq_dataset_name
  bq_tables = [
    {
      table_id = "gcs_bucket_notif_log",
      schema   = <<EOF
      [
        { "name": "bucket_id", "type": "STRING", "mode": "NULLABLE" },
        { "name": "event_datetime", "type": "DATETIME", "mode": "NULLABLE" },
        { "name": "event_date", "type": "DATE", "mode": "NULLABLE" },
        { "name": "event_epoch", "type": "INTEGER", "mode": "NULLABLE" },
        { "name": "event_type", "type": "STRING", "mode": "NULLABLE" },
        { "name": "notification_config", "type": "STRING", "mode": "NULLABLE" },
        { "name": "object_generation", "type": "STRING", "mode": "NULLABLE" },
        { "name": "payload_format", "type": "STRING", "mode": "NULLABLE" },
        { "name": "object_id", "type": "STRING", "mode": "NULLABLE" },
        { "name": "kind", "type": "STRING", "mode": "NULLABLE" },
        { "name": "id", "type": "STRING", "mode": "NULLABLE" },
        { "name": "self_link", "type": "STRING", "mode": "NULLABLE" },
        { "name": "name", "type": "STRING", "mode": "NULLABLE" },
        { "name": "metageneration", "type": "STRING", "mode": "NULLABLE" },
        { "name": "content_type", "type": "STRING", "mode": "NULLABLE" },
        { "name": "created_datetime", "type": "DATETIME", "mode": "NULLABLE" },
        { "name": "created_date", "type": "DATE", "mode": "NULLABLE" },
        { "name": "created_epoch", "type": "INTEGER", "mode": "NULLABLE" },
        { "name": "updated_datetime", "type": "DATETIME", "mode": "NULLABLE" },
        { "name": "updated_date", "type": "DATE", "mode": "NULLABLE" },
        { "name": "updated_epoch", "type": "INTEGER", "mode": "NULLABLE" },
        { "name": "storage_class", "type": "STRING", "mode": "NULLABLE" },
        { "name": "storage_class_updated_datetime", "type": "DATETIME", "mode": "NULLABLE" },
        { "name": "storage_class_updated_date", "type": "DATE", "mode": "NULLABLE" },
        { "name": "storage_class_updated_epoch", "type": "INTEGER", "mode": "NULLABLE" },
        { "name": "size", "type": "STRING", "mode": "NULLABLE" },
        { "name": "media_link", "type": "STRING", "mode": "NULLABLE" },
        { "name": "log_datetime", "type": "DATETIME", "mode": "NULLABLE" },
        { "name": "log_date", "type": "DATE", "mode": "NULLABLE" },
        { "name": "log_epoch", "type": "INTEGER", "mode": "NULLABLE" },
        { "name": "raw", "type": "JSON", "mode": "NULLABLE" }
      ]
      EOF
      time_partitioning = {
        type                     = "DAY",
        field                    = "event_date",
        require_partition_filter = true,
        expiration_ms            = null,
      },
      range_partitioning = null,
      expiration_time    = null,
      clustering         = ["log_date", "bucket_id", "object_id", "event_type"],
      labels             = {}
    },
    {
      table_id = "pubsub_log_dead_letter",
      schema   = <<EOF
      [
        { "name": "delivery_attempt", "type": "STRING", "mode": "NULLABLE" },
        { "name": "subscription_nm", "type": "STRING", "mode": "NULLABLE" },
        { "name": "subscription_project_id", "type": "STRING", "mode": "NULLABLE" },
        { "name": "publish_datetime", "type": "DATETIME", "mode": "NULLABLE" },
        { "name": "publish_date", "type": "DATE", "mode": "NULLABLE" },
        { "name": "publish_epoch", "type": "INTEGER", "mode": "NULLABLE" },
        { "name": "message_id", "type": "STRING", "mode": "NULLABLE" },
        { "name": "message_data", "type": "STRING", "mode": "NULLABLE" },
        { "name": "message_data_decoded", "type": "JSON", "mode": "NULLABLE" },
        { "name": "dead_letter_datetime", "type": "DATETIME", "mode": "NULLABLE" },
        { "name": "dead_letter_date", "type": "DATE", "mode": "NULLABLE" },
        { "name": "dead_letter_epoch", "type": "INTEGER", "mode": "NULLABLE" },
        { "name": "is_recycled", "type": "BOOL", "mode": "NULLABLE" },
        { "name": "log_datetime", "type": "DATETIME", "mode": "NULLABLE" },
        { "name": "log_date", "type": "DATE", "mode": "NULLABLE" },
        { "name": "log_epoch", "type": "INTEGER", "mode": "NULLABLE" }
      ]
      EOF
      time_partitioning = {
        type                     = "DAY",
        field                    = "dead_letter_date",
        require_partition_filter = true,
        expiration_ms            = null,
      },
      range_partitioning = null,
      expiration_time    = null,
      clustering         = ["log_date", "publish_date", "subscription_nm", "message_id"],
      labels             = {}
    },
    {
      table_id = "telegram_utils_log_file_downloader",
      schema   = <<EOF
      [
        { "name": "update_id", "type": "INTEGER", "mode": "NULLABLE" },
        { "name": "update_epoch", "type": "INTEGER", "mode": "NULLABLE" },
        { "name": "update_date", "type": "DATE", "mode": "NULLABLE" },
        { "name": "update_datetime", "type": "DATETIME", "mode": "NULLABLE" },
        { "name": "file", "type": "JSON", "mode": "NULLABLE" },
        { "name": "file_name", "type": "STRING", "mode": "NULLABLE" },
        { "name": "file_path", "type": "STRING", "mode": "NULLABLE" },
        { "name": "log_date", "type": "DATE", "mode": "NULLABLE" },
        { "name": "log_datetime", "type": "DATETIME", "mode": "NULLABLE" },
        { "name": "log_epoch", "type": "INTEGER", "mode": "NULLABLE" }
      ]
      EOF
      time_partitioning = {
        type                     = "DAY",
        field                    = "update_date",
        require_partition_filter = true,
        expiration_ms            = null,
      },
      range_partitioning = null,
      expiration_time    = null,
      clustering         = ["update_id", "log_date", "log_epoch", "file_name"],
      labels             = {}
    },
    {
      table_id = "telegram_msg_log_update",
      schema   = <<EOF
      [
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
      clustering         = ["log_epoch"],
      labels             = {}
    },
    {
      table_id = "telegram_msg_log_identification",
      schema   = <<EOF
      [
        { "name": "update_id", "type": "INTEGER", "mode": "NULLABLE" },
        { "name": "update_epoch", "type": "INTEGER", "mode": "NULLABLE" },
        { "name": "update_date", "type": "DATE", "mode": "NULLABLE" },
        { "name": "update_datetime", "type": "DATETIME", "mode": "NULLABLE" },
        { "name": "result", "type": "JSON", "mode": "NULLABLE" },
        { "name": "log_date", "type": "DATE", "mode": "NULLABLE" },
        { "name": "log_datetime", "type": "DATETIME", "mode": "NULLABLE" },
        { "name": "log_epoch", "type": "INTEGER", "mode": "NULLABLE" }
      ]
      EOF
      time_partitioning = {
        type                     = "DAY",
        field                    = "update_date",
        require_partition_filter = true,
        expiration_ms            = null,
      },
      range_partitioning = null,
      expiration_time    = null,
      clustering         = ["update_id", "update_epoch", "log_date", "log_epoch"],
      labels             = {}
    },
    {
      table_id = "telegram_act_log_send_message",
      schema   = <<EOF
      [
        { "name": "update_id", "type": "INTEGER", "mode": "NULLABLE" },
        { "name": "update_epoch", "type": "INTEGER", "mode": "NULLABLE" },
        { "name": "update_date", "type": "DATE", "mode": "NULLABLE" },
        { "name": "update_datetime", "type": "DATETIME", "mode": "NULLABLE" },
        { "name": "api_params", "type": "JSON", "mode": "NULLABLE" },
        { "name": "api_result", "type": "JSON", "mode": "NULLABLE" },
        { "name": "log_date", "type": "DATE", "mode": "NULLABLE" },
        { "name": "log_datetime", "type": "DATETIME", "mode": "NULLABLE" },
        { "name": "log_epoch", "type": "INTEGER", "mode": "NULLABLE" }
      ]
      EOF
      time_partitioning = {
        type                     = "DAY",
        field                    = "update_date",
        require_partition_filter = true,
        expiration_ms            = null,
      },
      range_partitioning = null,
      expiration_time    = null,
      clustering         = ["update_id", "update_epoch", "log_date", "log_epoch"],
      labels             = {}
    }
  ]
}

# BigQuery - Dataset Ops
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
