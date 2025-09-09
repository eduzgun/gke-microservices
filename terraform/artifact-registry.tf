resource "google_artifact_registry_repository" "main" {
  location      = var.region
  repository_id = "main-repo"
  description   = "Main Docker repository for applications"
  format        = "DOCKER"

  cleanup_policies {
    id     = "delete-old-images"
    action = "DELETE"
    
    condition {
      tag_state  = "TAGGED"
      older_than = "2500000s"
    }
  }

  cleanup_policies {
    id     = "keep-recent"
    action = "KEEP"
    
    most_recent_versions {
      keep_count = 10
    }
  }
}

resource "google_project_service" "artifact_registry" {
  service = "artifactregistry.googleapis.com"
  
  disable_dependent_services = true
}