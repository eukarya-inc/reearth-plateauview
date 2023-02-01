locals {
  reearth_api_secret = [
    "REEARTH_DB",
    "REEARTH_AUTH0_CLIENTID",
    "REEARTH_AUTH0_CLIENTSECRET",
    "REEARTH_SIGNUPSECRET",
    "REEARTH_MARKETPLACE_SECRET",
  ]
}

resource "google_cloud_run_service" "reearth_api" {
  name                       = "reearth-api"
  location                   = var.gcp_region
  autogenerate_revision_name = true
  metadata {
    annotations = {
      "run.googleapis.com/launch-stage"   = "BETA"
      "run.googleapis.com/ingress"        = "all"
      "run.googleapis.com/ingress-status" = "all"
    }
  }

  template {
    spec {
      service_account_name = google_service_account.reearth_api.email
      containers {
        # 最初はGCRにコンテナが存在しないケースも有るため、dummyのコンテナを立ち上げる。。
        # コンテナイメージはCD(GithubAction)から更新する
        # ignore_changesにimageを指定しているため、構築以降Terraformで更新してもこのコンテナに戻ることはない
        image = "gcr.io/cloudrun/hello"
        resources {
          limits = {
            cpu    = "1000m"
            memory = "512Mi"
          }

        }
        dynamic "env" {
          for_each = { for i in local.reearth_api_secret : i => i }
          content {
            name = env.value
            value_from {
              secret_key_ref {
                name = google_secret_manager_secret.reearth_api[env.value].secret_id
                key  = "latest"
              }
            }
          }
        }
        env {
          name  = "REEARTH_AUTH0_DOMAIN"
          value = "https://reearth-dev.auth0.com/"
        }
        env {
          name  = "REEARTH_GCS_BUCKETNAME"
          value = "static.${var.base_domain}"
        }
        env {
          name  = "REEARTH_ASSETBASEURL"
          value = "https://static.${var.base_domain}"
        }
        env {
          name  = "REEARTH_TRACERSAMPLE"
          value = ".0"
        }
        env {
          name  = "REEARTH_GCS_PUBLICATIONCACHECONTROL"
          value = "no-store"
        }
        env {
          name  = "REEARTH_ORIGINS"
          value = "http://localhost:3000,https://app.${var.base_domain},https://*.netlify.app,https://marketplace.${var.base_domain}"
        }
        env {
          name  = "REEARTH_AUTHSRV_DISABLED"
          value = "true"
        }
        env {
          name  = "REEARTH_HOST"
          value = "https://api.${var.base_domain}"
        }
        env {
          name  = "REEARTH_HOST_WEB"
          value = "https://app.${var.base_domain}"
        }
        env {
          name  = "REEARTH_AUTH_AUD"
          value = "https://api.${var.base_domain}"
        }
        env {
          name  = "REEARTH_MARKETPLACE_ENDPOINT"
          value = "https://api.marketplace.${var.base_domain}"
        }
        env {
          name  = "GOOGLE_CLOUD_PROJECT"
          value = var.gcp_project_name
        }
      }
    }
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"         = "100"
        "run.googleapis.com/execution-environment" = "gen2"
      }
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
  lifecycle {
    ignore_changes = [
      metadata[0].annotations,
      template[0].spec[0].containers[0].image,
      template[0].metadata[0].annotations["run.googleapis.com/client-name"],
      template[0].metadata[0].annotations["client.knative.dev/user-image"]
    ]
  }
  depends_on = [
    google_secret_manager_secret_version.reearth_api_dummy
  ]
}

resource "google_secret_manager_secret" "reearth_api" {
  for_each  = toset(local.reearth_api_secret)
  secret_id = "reearth-api-${each.value}"
  labels = {
    label = "reearth-api"
  }
  replication {
    user_managed {
      replicas {
        location = "asia-northeast2"
      }
    }
  }
}

//MEMO: secret_managerに値が入っていないとcloudrunが起動エラーになるので、
//      あとから手動で値を入れるものに関しては先にdummyの値を入れておく
resource "google_secret_manager_secret_version" "reearth_api_dummy" {
  for_each = toset([
    "REEARTH_DB",
    "REEARTH_AUTH0_CLIENTID",
    "REEARTH_AUTH0_CLIENTSECRET",
    "REEARTH_SIGNUPSECRET",
    "REEARTH_MARKETPLACE_SECRET",
  ])
  secret = google_secret_manager_secret.reearth_api[each.value].id

  secret_data = "dummy"
  lifecycle {
    ignore_changes = [
      secret_data
    ]
  }
}