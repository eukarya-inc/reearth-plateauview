//GCS周り
resource "google_storage_bucket" "assets" {
  name          = "${var.service_prefix}-cms-assets-bucket"
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

resource "google_storage_bucket_iam_binding" "assets_public_read" {
  bucket = google_storage_bucket.assets.name
  role   = "roles/storage.objectViewer"
  members = [
    "allUsers",
    "serviceAccount:service-${data.google_project.project.number}@compute-system.iam.gserviceaccount.com",
  ]
}