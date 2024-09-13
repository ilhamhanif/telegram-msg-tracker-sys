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
