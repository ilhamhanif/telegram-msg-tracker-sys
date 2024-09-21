locals {
  gcs_bucket_name = [
    "gcf-zip-source",
    "telegram-object-sent"
  ]
}

# GCS Buckets
module "gcs_buckets" {
  source     = "terraform-google-modules/cloud-storage/google"
  version    = "~> 6.1"
  project_id = var.project_id
  prefix     = var.project_id
  names      = local.gcs_bucket_name
  location   = var.region
  bucket_lifecycle_rules = {

    gcf-zip-source = [
      {
        condition = {
          age = "1"
        }
        action = {
          type = "Delete"
        }
      }
    ]

  }
}

# Add default GCS service account a PubSub publisher role.
data "google_storage_project_service_account" "gcs_service_account" {}

resource "google_pubsub_topic_iam_member" "gcs_service_account_role_binding" {
  project = var.project_id
  topic   = var.pubsub_bucket_notif_topic_name
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:${data.google_storage_project_service_account.gcs_service_account.email_address}"
}

# Add Bucket Notification for all bucket.
resource "google_storage_notification" "gcs_notif_bucket_notification" {
  count          = length(local.gcs_bucket_name)
  bucket         = "${var.project_id}-${local.gcs_bucket_name[count.index]}"
  payload_format = "JSON_API_V1"
  topic          = var.pubsub_bucket_notif_topic_name
  event_types    = ["OBJECT_FINALIZE", "OBJECT_METADATA_UPDATE", "OBJECT_ARCHIVE", "OBJECT_DELETE"]
}
