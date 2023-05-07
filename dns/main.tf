terraform {
  required_providers {
    godaddy = {
      source  = "terraform.local/local/godaddy"
      version = "1.8.7"
    }
  }
}

provider "godaddy" {
}


resource "godaddy_domain_record" "gd-fancy-domain" {
  domain   = "myplayk8s.live"

  // required if provider key does not belong to customer

  // specify zero or more record blocks
  // a record block allows you to configure A, or NS records with a custom time-to-live value
  // a record block also allow you to configure AAAA, CNAME, TXT, or MX records
  dynamic record {
    for_each = toset(var.records)
    content {
      name = record.value
      type = "CNAME"
      data = "@"
      ttl = 600
    }
  }
  record {
    name = "@"
    data = "192.168.1.2"
    port = 0
    priority = 0
    ttl = 600
    weight = 0
    type = "A"
  }
  record {
      data     = "_domainconnect.gd.domaincontrol.com"
      name     = "_domainconnect"
      port     = 0
      priority = 0
      ttl      = 3600
      type     = "CNAME"
      weight   = 0
  }  
  
}

variable "records" {
  type = list 
  default = [
    "radarr",
    "sonarr",
    "plex",
    "jackett",
    "gorunner",
    "argocd",
    "minio",
    "traefik",
    "transmission",
    "version-discover",
    "game",
    "flaresolver",
    "iptvparser",
    "emby",
    "kodi",
    "qbittorrent",
  ]
}
