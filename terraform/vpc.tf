resource "google_compute_network" "vpc" {
    name = "main"
    routing_mode = "REGIONAL"
    auto_create_subnetworks = false 
    delete_default_routes_on_create = true

    depends_on = [google_project_service.api]
}

resource "google_compute_route" "default_route" {
    name = "default-route"
    dest_range = "0.0.0.0/0"
    network = google_compute_network.vpc.name
    next_hop_gateway = "default-internet-gateway"
}

resource "google_compute_subnetwork" "private" {
    name = "private"
    # I inputed the public subnetwork's IP range to an online CIDR calculator to find 
    # available ranges that don't overlap
    ip_cidr_range = "10.0.32.0/19"
    region = var.region
    network = google_compute_network.vpc.name
    private_ip_google_access = true
    stack_type = "IPV4_ONLY"

    # This may help if our pods need to access any VM outside kubernetes
    secondary_ip_range {
      range_name = "k8s-pods"
      ip_cidr_range = "172.16.0.0/14"
    }

    # When exposing services to other applications we use the ClusterIP service
    secondary_ip_range {
      range_name = "k8s-services"
      ip_cidr_range = "172.20.0.0/18"
    }
}

resource "google_compute_subnetwork" "proxy_only" {
  name          = "proxy-only-subnet"
  ip_cidr_range = "10.129.0.0/23"  # No range overlap
  region        = var.region
  purpose       = "REGIONAL_MANAGED_PROXY"
  role          = "ACTIVE"
  network       = google_compute_network.vpc.name
}

resource "google_compute_address" "nat" {
    name = "nat"
    address_type = "EXTERNAL"
    network_tier = "PREMIUM"
    region = var.region 
    depends_on = [google_project_service.api]
}

# Router is used to advertise this NAT gateway to the prior subnet
resource "google_compute_router" "router" {
    name = "router"
    region = var.region
    network = google_compute_network.vpc.id
}

# NAT gateway translates private IP addresses to shared public IP addresses, allowing them to reach the internet
resource "google_compute_router_nat" "nat" {
    name = "nat"
    region = var.region
    router = google_compute_router.router.name

    nat_ip_allocate_option = "MANUAL_ONLY"
    source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"
    nat_ips = [google_compute_address.nat.self_link]

    subnetwork {
      name = google_compute_subnetwork.private.self_link
      source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
    }
}


resource "google_compute_firewall" "allow_internal" {
  name    = "allow-internal"
  network = google_compute_network.vpc.name

  allow {
    protocol = "icmp"
  }
  allow {
    protocol = "tcp"
    ports    = ["1-65535"]
  }
  allow {
    protocol = "udp"
    ports    = ["1-65535"]
  }

  source_ranges = ["10.0.32.0/19"] # Private subnet range
}

resource "google_compute_firewall" "allow_gcp_health_checks" {
  name    = "allow-gcp-health-checks"
  network = google_compute_network.vpc.name

  allow {
    protocol = "tcp"
    ports    = ["3000", "8080"]
  }

  # Google Cloud health check IP ranges
  source_ranges = [
    "35.191.0.0/16",
    "130.211.0.0/22",
  ]

  target_tags = [] 
}