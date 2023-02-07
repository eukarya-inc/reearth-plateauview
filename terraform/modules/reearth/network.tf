
resource "google_compute_target_http_proxy" "reearth" {
  name       = "reearth-common-http-targetproxy"
  proxy_bind = "false"
  url_map    = google_compute_url_map.reearth.id
}


resource "google_compute_target_https_proxy" "reearth" {
  name             = "reearth-common-https-targetproxy"
  url_map          = google_compute_url_map.reearth.id
  ssl_certificates = [google_compute_managed_ssl_certificate.common.id]
  # lifecycle {
  #   ignore_changes = [
  #     ssl_certificates,
  #   ] #証明書を動的に追加するので2回目以降は無視させる
  # }
}

resource "google_compute_global_address" "reearth_lb" {
  name = "reearth-common-lb"
}

resource "google_compute_managed_ssl_certificate" "common" {
  name = "reearth-common-cert"

  managed {
    domains = [
      local.api_reearth_domain,
      local.reearth_domain,
      local.static_reearth_domain,
    ] #TODO: ワイルドカード対応
  }
  # lifecycle {
  #   create_before_destroy = true
  # }
}

resource "google_compute_global_forwarding_rule" "reearth_https" {
  name       = "reearth-common-https"
  target     = google_compute_target_https_proxy.reearth.self_link
  port_range = "443"
  ip_address = google_compute_global_address.reearth_lb.address

  depends_on = [google_compute_url_map.reearth]
}

resource "google_compute_url_map" "reearth_redirect" {
  name = "reearth-https-redirect"
  default_url_redirect {
    https_redirect         = "true"
    redirect_response_code = "MOVED_PERMANENTLY_DEFAULT"
    strip_query            = "false"
  }

  description = "HTTP to HTTPS redirect forwarding rule"
}

resource "google_compute_url_map" "reearth" {
  name        = "reearth-common-urlmap"
  description = "reearth common urlmap"

  default_service = google_compute_backend_bucket.app_backend.self_link

  host_rule {
    hosts = [
      local.reearth_domain,
    ]
    path_matcher = "path-matcher-1"
  }

  path_matcher {
    default_service = google_compute_backend_bucket.app_backend.self_link
    name            = "path-matcher-1"
  }

  host_rule {
    hosts = [
      local.static_reearth_domain,
    ]
    path_matcher = "path-matcher-2"
  }

  path_matcher {
    default_service = google_compute_backend_bucket.static_backend.self_link
    name            = "path-matcher-2"
  }

  host_rule {
    hosts = [
      local.api_reearth_domain,
    ]
    path_matcher = "path-matcher-3"
  }

  path_matcher {
    default_service = google_compute_backend_service.reearth_api.self_link
    name            = "path-matcher-3"
  }
}

resource "google_compute_backend_bucket" "app_backend" {
  name        = "reearth-app-backend"
  bucket_name = google_storage_bucket.app.name
}

resource "google_compute_backend_bucket" "static_backend" {
  name        = "reearth-static-backend"
  bucket_name = google_storage_bucket.static.name
}


resource "google_compute_region_network_endpoint_group" "reearth_api" {
  name                  = "reearth-api-neg"
  network_endpoint_type = "SERVERLESS"
  region                = "asia-northeast1"
  cloud_run {
    service = google_cloud_run_service.reearth_api.name
  }
}

resource "google_compute_backend_service" "reearth_api" {
  affinity_cookie_ttl_sec = "0"

  backend {
    balancing_mode               = "UTILIZATION"
    capacity_scaler              = "0"
    group                        = google_compute_region_network_endpoint_group.reearth_api.id
    max_connections              = "0"
    max_connections_per_endpoint = "0"
    max_connections_per_instance = "0"
    max_rate                     = "0"
    max_rate_per_endpoint        = "0"
    max_rate_per_instance        = "0"
    max_utilization              = "0"
  }

  connection_draining_timeout_sec = "0"
  description                     = "reearth-api-neg"
  enable_cdn                      = "false"
  load_balancing_scheme           = "EXTERNAL"

  log_config {
    enable      = "true"
    sample_rate = "1"
  }

  name             = "reearth-api-backend"
  port_name        = "http"
  protocol         = "HTTPS"
  session_affinity = "NONE"
  timeout_sec      = "30"
}