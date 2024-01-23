provider "google" {
  project     = "distributedsystems-405712"
  region      = "europe-north-1a"
}

resource "google_compute_network" "vpc_network" {
  name                    = "my-custom-mode-network"
  auto_create_subnetworks = false
  mtu                     = 1460
}

resource "google_compute_subnetwork" "default" {
  name          = "my-custom-subnet"
  ip_cidr_range = "10.0.1.0/24"
  region        = "europe-north1"
  network       = google_compute_network.vpc_network.id
}

# Create a single Compute Engine instance
resource "google_compute_instance" "worker" {
    count = 3

  name         = "worker-${count.index}"
  machine_type = "e2-micro"
  zone         = "europe-north1-a"
  tags         = ["ssh"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
    
  }

  metadata_startup_script = "${file("startup-script-worker.sh")}"

  network_interface {
    subnetwork = google_compute_subnetwork.default.id

    access_config {
      # Include this section to give the VM an external IP address
    }
  }
  
}

# Create a single Compute Engine instance
resource "google_compute_instance" "master" {

  name         = "master"
  machine_type = "e2-micro"
  zone         = "europe-north1-a"
  tags         = ["ssh"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
    
  }

  metadata_startup_script = "${file("startup-script-master.sh")}"

  network_interface {
    subnetwork = google_compute_subnetwork.default.id

    access_config {
      # Include this section to give the VM an external IP address
    }
  }
  
}


resource "google_compute_firewall" "ssh" {
  name = "allow-ssh"
  allow {
    ports    = ["22"]
    protocol = "tcp"
  }
  direction     = "INGRESS"
  network       = google_compute_network.vpc_network.id
  priority      = 1000
  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["ssh"]
}

output start_up_script {
  value       = "${file("startup-script-worker.sh")}"
  sensitive   = false
  description = "description"
  depends_on  = []
}
