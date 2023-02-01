module "reearth-api" {
  source = "./modules/reearth"

  base_domain           = var.base_domain
  gcp_project_name      = var.gcp_project_name
  service_prefix        = var.service_prefix
  dns_managed_zone_name = var.dns_managed_zone_name
}