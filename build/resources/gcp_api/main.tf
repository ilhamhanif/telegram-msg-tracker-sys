locals {
  api_services = [
    "compute.googleapis.com",
    "cloudbuild.googleapis.com",
    "run.googleapis.com",
    "cloudfunctions.googleapis.com",
    "bigquery.googleapis.com"
  ]
}

# API Services
resource "google_project_service" "gcp_api_services" {
  count              = length(local.api_services)
  project            = var.project_id
  service            = local.api_services[count.index]
  disable_on_destroy = false
}

# Wait for All API 
# and Initialization successfully enabled
resource "null_resource" "resource_api_activation_complete" {
  depends_on = [google_project_service.gcp_api_services]
}
