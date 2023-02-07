//GCS周り
resource "google_storage_bucket" "app" {
  name          = "${var.service_prefix}-reearth-app-bucket"
  location      = "ASIA"
  storage_class = "MULTI_REGIONAL"
  cors {
    max_age_seconds = 60
    method = [
      "GET",
      "HEAD",
      "OPTIONS",
    ]
    origin = [
      local.api_reearth_domain,
      local.reearth_domain,
    ]
    response_header = [
      "Content-Type",
      "Access-Control-Allow-Origin"
    ]
  }

  website {
    main_page_suffix = "index.html"
    not_found_page   = "index.html"
  }
}

resource "google_storage_bucket_iam_binding" "app_public_read" {
  bucket = google_storage_bucket.app.name
  role   = "roles/storage.objectViewer"
  members = [
    "allUsers",
    "serviceAccount:service-${data.google_project.project.number}@compute-system.iam.gserviceaccount.com",
  ]
}


resource "google_storage_bucket" "static" {
  name          = "${var.service_prefix}-reearth-static-bucket"
  location      = "ASIA"
  storage_class = "MULTI_REGIONAL"

  cors {
    max_age_seconds = 60
    method = [
      "GET",
      "HEAD",
      "OPTIONS",
    ]
    origin = [
      "*"
    ]
    response_header = [
      "Content-Type",
      "Access-Control-Allow-Origin"
    ]
  }

  website {
    main_page_suffix = "index.html"
    not_found_page   = "index.html"
  }
}

resource "google_storage_bucket_iam_binding" "static_public_read" {
  bucket = google_storage_bucket.static.name
  role   = "roles/storage.objectViewer"
  members = [
    "allUsers",
    "serviceAccount:service-${data.google_project.project.number}@compute-system.iam.gserviceaccount.com",
  ]
}


resource "google_storage_bucket_object" "reearth_config" {
  bucket        = google_storage_bucket.app.name
  name          = "reearth_config.json"
  cache_control = "no-store"
  content_type  = "application/json"

  content = templatefile("${path.module}/template/reearth_config.template.json",
    {
      "api"                  = "https://${local.api_reearth_domain}/api",
      "plugins"              = "https://${local.static_reearth_domain}/plugins",
      "published"            = "https://{}.${var.base_domain}",
      "auth0ClientId"        = auth0_client.reearth_app.client_id,
      "auth0Domain"          = var.auth0.domain,
      "auth0Audience"        = "https://${local.api_reearth_domain}",
      "cesiumIonAccessToken" = var.cesium_ion_access_token,
  })
}