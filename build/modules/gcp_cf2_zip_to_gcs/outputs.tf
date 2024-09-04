output "zip_gcs_bucket_name" {
  description = "GCP GCS ZIP Bucket Name"
  value       = google_storage_bucket_object.bucket_cf_zip_source_code_upload.bucket
}

output "zip_gcs_object_name" {
  description = "GCP GCS ZIP Object Name"
  value       = google_storage_bucket_object.bucket_cf_zip_source_code_upload.name
}

output "zip_gcs_object_generation" {
  description = "GCP GCS ZIP Object Generation"
  value       = google_storage_bucket_object.bucket_cf_zip_source_code_upload.generation
}
