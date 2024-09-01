output "cloud_function2_webhook_target_uri" {
  description = "Cloud Function2 Webhook Target URI"
  value       = module.cf_tlgrm_msg_upd_forwarder.cloud_function2_uri
}
