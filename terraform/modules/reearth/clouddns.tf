data "google_dns_managed_zone" "reearth" {
  name = var.dns_managed_zone_name
}

resource "google_dns_record_set" "api" {
  name = "api.reearth.${data.google_dns_managed_zone.reearth.dns_name}"
  type = "A"
  ttl  = 60

  managed_zone = data.google_dns_managed_zone.reearth.name
  rrdatas      = [google_compute_global_address.reearth_lb.address]
}

resource "google_dns_record_set" "static" {
  name = "static.reearth.${data.google_dns_managed_zone.reearth.dns_name}"
  type = "A"
  ttl  = 60

  managed_zone = data.google_dns_managed_zone.reearth.name
  rrdatas      = [google_compute_global_address.reearth_lb.address]
}

resource "google_dns_record_set" "app" {
  name = "reearth.${data.google_dns_managed_zone.reearth.dns_name}"
  type = "A"
  ttl  = 60

  managed_zone = data.google_dns_managed_zone.reearth.name
  rrdatas      = [google_compute_global_address.reearth_lb.address]
}