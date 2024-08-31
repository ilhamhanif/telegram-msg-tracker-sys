# Create Service Account
# with roles:
# 1. Cloud Run Function Invoker.
resource "google_service_account" "service_account" {
  account_id   = var.service_account_id
  display_name = var.service_account_display_name
}

data "google_iam_policy" "roles" {
  binding {
    role = "roles/run.invoker"

    members = [
      "serviceAccount:${google_service_account.service_account.email}"
    ]
  }
}
