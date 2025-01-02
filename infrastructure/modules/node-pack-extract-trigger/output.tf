output "topic_id" {
  value = google_pubsub_topic.topic.id
}

output "bucket_notification_id" {
  value = google_storage_notification.notification.id
}

output "trigger_id" {
  value = google_cloudbuild_trigger.trigger.id
}

output "backfill_scheduler_id" {
  value = google_cloud_scheduler_job.backfill.id
}
