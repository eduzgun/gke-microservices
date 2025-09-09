resource "google_container_cluster" "primary" {
  name     = "my-gke-cluster"
  location = var.zone  # Use variable reference

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count       = 1
  
  network = google_compute_network.vpc.name
  subnetwork = google_compute_subnetwork.private.name # GKE nodes go in private subnet

  addons_config {
    http_load_balancing {
      disabled = false
    }

    horizontal_pod_autoscaling {
      disabled = true
    }

  }
  
  gateway_api_config {
    channel = "CHANNEL_STANDARD"
  }

  deletion_protection = false

  release_channel {
    channel = "REGULAR"
  }

  # Create permissions to individual pods
  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }

  #Gateway API needs VPC-native cluster with IP allocation
  ip_allocation_policy {
    cluster_secondary_range_name = "k8s-pods"
    services_secondary_range_name = "k8s-services"
  }

  private_cluster_config {
    enable_private_nodes = true
    # Since we don't have a bastion or VPN we use public endpoint
    enable_private_endpoint = false
    # Allocate a small range of just 16 IP addresses for the control plane
    # This range can later be used to grant access to external CI/CD tools
    # For example with jenkins we will need to open the firewall and also whitelist the subnet
    master_ipv4_cidr_block = "192.168.0.0/28"
  }


}

# Node pool
resource "google_container_node_pool" "small_nodes" {
  name       = "small-pool"
  cluster    = google_container_cluster.primary.name
  location   = google_container_cluster.primary.location

  node_config {
    machine_type = "e2-small"
    disk_size_gb = 15
    preemptible  = true # For cost savings
    
    # Use the custom service account
    service_account = google_service_account.gke_nodes.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]

    # Create custom labels which will be useful to migrate applications to another node group
    # with the same label later on
    labels = {
        role = "general"
    }
  }

  autoscaling {
    min_node_count = 1
    max_node_count = 2  
  }

  management {
    auto_repair = true
    auto_upgrade = true
  }
}