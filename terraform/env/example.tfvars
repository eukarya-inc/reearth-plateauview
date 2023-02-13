base_domain           = "" #reearthを提供するドメイン名を指定してください。 ex): reearth.io
gcp_project_name      = "" #デプロイするGCPのプロジェクト名を指定してください。
gcp_region            = "" #デプロイ先のregion名を指定してください。ex):asia-northeast1
service_prefix        = "" #作成されるリソース名に付与するprefixを指定してください。　ex): plateauview-dev
dns_managed_zone_name = "" #事前設定で行ったCloudDNSのゾーン名を指定してください。

#reearthでauth0に関する設定
auth0 = {
  domain = ""
}

#Terraformで利用するauth0の設定
auth0_provider = {
  client_id = ""
  domain    = ""
}