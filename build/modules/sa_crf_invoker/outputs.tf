output "service_account_id" {
  description = "SA ID"
  value       = google_service_account.service_account.id
}

output "service_account_email" {
  description = "SA Email"
  value       = google_service_account.service_account.email
}
