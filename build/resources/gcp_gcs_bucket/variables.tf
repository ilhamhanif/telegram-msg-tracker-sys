variable "project_id" {
  description = "GCP Project ID Name"
  type        = string
}

variable "project_number" {
  description = "GCP Project ID Number"
  type        = number
}

variable "region" {
  description = "GCP Region Name"
  type        = string
}

variable "pubsub_bucket_notif_topic_name" {
  description = "GCS bucket notif Pubsub topic name"
  type        = string
}
