# Generates a ZIP compressed file archieve of the source code.
data "archive_file" "source" {
  type        = "zip"
  source_dir  = "../app/${var.cloud_functions2_name}/function"
  output_path = "zip/${var.cloud_functions2_name}.zip"
}

# Create a GCS Bucket to Store Cloud Function ZIP file.
resource "google_storage_bucket" "bucket_cf_zip_source_code" {
  name     = "${var.gcs_zip_project_id}-gcf-zip-source"
  location = var.gcs_zip_bucket_location

  lifecycle_rule {
    condition {
      age = "1"
    }
    action {
      type = "Delete"
    }
  }
}

resource "google_storage_bucket_object" "bucket_cf_zip_source_code_upload" {
  depends_on = [
    data.archive_file.source,
    google_storage_bucket.bucket_cf_zip_source_code
  ]

  name   = "${var.cloud_functions2_name}/${data.archive_file.source.output_md5}.zip"
  bucket = google_storage_bucket.bucket_cf_zip_source_code.name
  source = data.archive_file.source.output_path
}
