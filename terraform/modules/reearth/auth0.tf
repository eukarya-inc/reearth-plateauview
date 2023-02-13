resource "auth0_client" "reearth_app" {
  name = "plateau-reearth"

  app_type                   = "spa"
  initiate_login_uri         = "https://${local.reearth_domain}"
  token_endpoint_auth_method = "none"
  allowed_logout_urls = [
    "https://${local.reearth_domain}"
  ]
  allowed_origins = [
    "https://${local.reearth_domain}"
  ]
  callbacks = [
    "https://${local.reearth_domain}"
  ]
  web_origins = [
    "https://${local.reearth_domain}"
  ]
}

resource "auth0_client" "reearth_api" {
  name = "plateau-reearth-m2m"

  app_type           = "non_interactive"
  initiate_login_uri = "https://${local.reearth_domain}"

  allowed_logout_urls = [
    "https://${local.reearth_domain}"
  ]
  allowed_origins = [
    "https://${local.reearth_domain}"
  ]
  callbacks = [
    "https://${local.reearth_domain}"
  ]
  web_origins = [
    "https://${local.reearth_domain}"
  ]
}

resource "auth0_resource_server" "reearth_api" {
  name       = "plateau-reearth"
  identifier = "https://${local.api_reearth_domain}"
}

resource "auth0_action" "reearth_app" {
  name    = "signup-reearth"
  runtime = "node16"
  deploy  = true
  code    = file("${path.module}/source/signup-reearth.js")

  supported_triggers {
    id      = "post-user-registration"
    version = "v2"
  }

  dependencies {
    name    = "axios"
    version = "1.3.2"
  }

  secrets {
    name  = "secret"
    value = random_string.reeart_app_action_secret.result
  }

  secrets {
    name  = "api"
    value = "https://${local.api_reearth_domain}/api/signup"
  }
}

resource "random_string" "reeart_app_action_secret" {
  length  = 32
  special = false
}

# resource "auth0_trigger_binding" "reearth_login" {
#   trigger = "post-user-registration"

#   actions {
#     id           = auth0_action.reearth_app.id
#     display_name = auth0_action.action_foo.name
#   }
# }