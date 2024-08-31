output "zip_gcs_bucket_name" {
  description = "GCP GCS Bucket Name"
  value       = google_storage_bucket_object.upload_to_bucket_cf_zip_source_code.bucket
}

output "zip_gcs_object_name" {
  description = "GCP GCS Object Name"
  value       = google_storage_bucket_object.upload_to_bucket_cf_zip_source_code.name
}
