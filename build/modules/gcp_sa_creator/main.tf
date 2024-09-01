# Create custom SA, attached with roles
resource "google_service_account" "service_account" {
  account_id   = var.service_account_id
  display_name = var.service_account_name
}

resource "google_project_iam_member" "service_account_role_binding" {
  count   = length(var.service_account_roles)
  project = var.project_id
  role    = var.service_account_roles[count.index]
  member  = "serviceAccount:${google_service_account.service_account.email}"
}
