locals {
  bq_dataset_name = "ops"
  bq_dataset_id   = local.bq_dataset_name
  bq_tables = [
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
