# Generates a ZIP compressed file archieve of the source code.
data "archive_file" "source" {
  type        = "zip"
  source_dir  = "../app/${var.cloud_functions2_name}/function"
  output_path = "zip/${var.cloud_functions2_name}.zip"
}

# Create a GCS Bucket to Store Cloud Function ZIP file.
data "google_storage_bucket" "bucket_cf_zip_source_code" {
  name = "${var.gcs_zip_project_id}-gcf-zip-source"
}

resource "google_storage_bucket_object" "bucket_cf_zip_source_code_upload" {
  name   = "${var.cloud_functions2_name}/${data.archive_file.source.output_md5}.zip"
  bucket = data.google_storage_bucket.bucket_cf_zip_source_code.name
  source = data.archive_file.source.output_path
}
