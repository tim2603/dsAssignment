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
  machine_type = "n2-standard-2"
  zone         = "europe-north1-a"
  tags         = ["ssh"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
    
  }

  service_account {
    email  = "614968029879-compute@developer.gserviceaccount.com"
    scopes = ["https://www.googleapis.com/auth/cloud-platform"]
  }

  metadata_startup_script = "${file("startup-script-worker.sh")}"

  network_interface {
    subnetwork = google_compute_subnetwork.default.id

    access_config {
      # Include this section to give the VM an external IP address
    }
  }  
}

# Dieser Code ist mit Terraform 4.25.0 und Versionen kompatibel, die mit 4.25.0 abwärtskompatibel sind.
# Informationen zum Validieren dieses Terraform-Codes finden Sie unter https://developer.hashicorp.com/terraform/tutorials/gcp-get-started/google-cloud-platform-build#format-and-validate-the-configuration.

resource "google_compute_instance" "worker-e2" {
      count = 8
  boot_disk {
    auto_delete = true
    device_name = "worker-${count.index}"

    initialize_params {
      image = "projects/debian-cloud/global/images/debian-11-bullseye-v20240110"
      size  = 10
      type  = "pd-balanced"
    }

    mode = "READ_WRITE"
  }

  can_ip_forward      = false
  deletion_protection = false
  enable_display      = false

  labels = {
    goog-ec-src = "vm_add-tf"
  }

  machine_type = "e2-medium"
  name         = "worker-${count.index}"

  network_interface {
    access_config {
      network_tier = "PREMIUM"
    }

    queue_count = 0
    stack_type  = "IPV4_ONLY"
    subnetwork  = "projects/distributedsystems-405712/regions/europe-north1/subnetworks/default"
  }

  scheduling {
    automatic_restart   = true
    on_host_maintenance = "MIGRATE"
    preemptible         = false
    provisioning_model  = "STANDARD"
  }

  service_account {
    email  = "614968029879-compute@developer.gserviceaccount.com"
    scopes = ["https://www.googleapis.com/auth/cloud-platform"]
  }

  shielded_instance_config {
    enable_integrity_monitoring = true
    enable_secure_boot          = false
    enable_vtpm                 = true
  }

  zone = "europe-north1-a"
  metadata_startup_script = "sudo apt-get install git -y && git clone https://github.com/tim2603/dsAssignment.git && cd dsAssignment && git checkout tim"
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
