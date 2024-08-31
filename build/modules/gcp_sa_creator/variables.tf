variable "project_id" {
  description = "Project ID"
  type        = string
}

variable "service_account_id" {
  description = "Service Account ID"
  type        = string
}

variable "service_account_name" {
  description = "Service Account Name (Display)"
  type        = string
}

variable "service_account_roles" {
  description = "Service Account Roles"
  type        = list(string)
}
