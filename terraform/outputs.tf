# Outputs for easier access to created resources
output "artifact_registry_url" {
  description = "URL of the Artifact Registry repository"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.main.repository_id}"
}

output "gke_service_account_email" {
  description = "Email of the GKE service account"
  value       = google_service_account.gke_nodes.email
}

output "artifact_registry_service_account_email" {
  description = "Email of the Artifact Registry service account"
  value       = google_service_account.artifact_registry.email
}

output "registry_docker_config_command" {
  description = "Command to configure Docker for Artifact Registry"
  value       = "gcloud auth configure-docker ${var.region}-docker.pkg.dev"
}

output "artifact_registry_repo" {
  value = google_artifact_registry_repository.main.repository_id
}

output "project_id" {
  value = var.project_id
}

output "region" {
  value = var.region
}

output "zone" {
  value = var.zone
}