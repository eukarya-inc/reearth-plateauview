# Reearth Plateauview for Terraform

Reearth Plateauviewを構築するTerraformです。

## 事前準備

 * コマンドラインツールのインストール
   * gcloud
   * terraform

## gcloudのセットアップ
```
gcloud config configurations create ${SERVICE_PREFIX}
gcloud config set project ${PROJECT_ID}
gcloud auth login
gcloud config set compute/region asia-northeast1
gcloud auth application-default login
```

### GCP APIの有効化

```
gcloud services enable certificatemanager.googleapis.com
gcloud services enable secretmanager.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable cloudresourcemanager.googleapis.com
gcloud services enable cloudtasks.googleapis.com
gcloud services enable compute.googleapis.com
gcloud services enable dns.googleapis.com
gcloud services enable iam.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable sts.googleapis.com
```

### CloudDNSのセットアップ

今回セットアップしたいドメインをcloud dnsでホスティングできるように設定する。
手順は各種レジストラの手順を参照。

```
 gcloud dns managed-zones create ${ZONE_NAME} --dns-name ${DOMAIN} --description "${DESCRIPTION}"
```

### Auth0 Management APIの設定

[公式のQuick Start](https://github.com/auth0/terraform-provider-auth0/blob/main/docs/guides/quickstart.md)を参考にセットアップ

## Terraform実行手順

### 設定ファイルの準備

今回設定したい環境の設定ファイル(tfvars)を準備する。実行環境費に合わせて必要な情報をセットアップを行う。
[example.tfvars](./env/example.tfvars) をコピーし、必要な設定を追記する。

### backendの作成
リソース全般の作成に必要なSERVICE_PREFIXを決める。
```
export SERVICE_PREFIX=""
gcloud storage buckets create gs://${SERVICE_PREFIX}-terraform-tfstate
```

[terraform.tf](terraform.tf)の`backend`のbucketを設定する

```diff
  backend "gcs" {
-    bucket = ""
+    bucket = "${SERVICE_PREFIXで指定した値を入れる}-terraform-tfstate"
  }
```

## Terraform実行後手順

```bash
terraform init
```

```bash
export AUTH0_CLIENT_SECRET_VALUE=""
terraform apply -var-file=env/example.tfvars
```


### terraform実行後のシークレットへの値追加

```bash
echo -n "${REEARTH_DB}"| gcloud secrets versions add reearth-api-REEARTH_DB --data-file=-
echo -n "${REEARTH_API_SECRET}"| gcloud secrets versions add reearth-api-REEARTH_MARKETPLACE_SECRET --data-file=-
echo -n "${REEARTH_CMS_WORKER_DB}" | gcloud secrets versions add reearth-api-reearth-cms-REEARTH_CMS_WORKER_DB --data-file=-
echo -n "${REEARTH_CMS_DB}" | gcloud secrets versions add reearth-api-reearth-cms-REEARTH_CMS_DB --data-file=-
echo -n "${REEARTH_PLATEAUVIEW_CMS_TOKEN}" | gcloud secrets versions add reearth-api-reearth-cms-REEARTH_PLATEAUVIEW_CMS_TOKEN --data-file=-
echo -n "${REEARTH_PLATEAUVIEW_FME_TOKEN}" | gcloud secrets versions add reearth-api-reearth-cms-REEARTH_PLATEAUVIEW_FME_TOKEN --data-file=-
echo -n "${REEARTH_PLATEAUVIEW_CKAN_TOKEN}" | gcloud secrets versions add reearth-api-reearth-cms-REEARTH_PLATEAUVIEW_CKAN_TOKEN --data-file=-
echo -n "${REEARTH_PLATEAUVIEW_SENDGRID_APIKEY}" | gcloud secrets versions add reearth-api-reearth-cms-REEARTH_PLATEAUVIEW_SENDGRID_APIKEY --data-file=-
```

### reearth-app用静的ファイルのアップロード

```bash
export REEARTH_VERSION=0.14.1
curl -o tmp/reearth-web.tar.gz -L https://github.com/reearth/reearth/releases/download/v${REEARTH_VERSION}/reearth-web_v${REEARTH_VERSION}.tar.gz

cd tmp
tar zxvf tmp/reearth-web.tar.gz
gsutil -m -h "Cache-Control:no-store" rsync -x "^reearth_config\\.json$" -dr reearth-web/ gs://${SERVICE_PREFIX}-reearth-app-bucket/
```

### reearthのdeploy
```bash
gcloud run deploy reearth-api \
            --image reearth/reearth:nightly \
            --region asia-northeast1 \
            --platform managed \
            --quiet
```

```
gcloud run deploy reearth-cms-api \
            --image reearth/reearth-cms:nightly \
            --region asia-northeast1 \
            --platform managed \
            --quiet
```

```
gcloud run deploy reearth-cms-worker \
            --image reearth/reearth-cms-worker:nightly \
            --region asia-northeast1 \
            --platform managed \
            --quiet

```

```
gcloud run deploy plateauview-api \
            --image eukarya/plateauview-api:latest \
            --region asia-northeast1 \
            --platform managed \
            --quiet

```


###  CMS_TOKENの設定
CMS UIで発行したTOKENを登録する
```
echo -n "${REEARTH_PLATEAUVIEW_CMS_TOKEN}" | gcloud secrets versions add 	reearth-cms-REEARTH_PLATEAUVIEW_CMS_TOKEN --data-file=-
```

```
gcloud run deploy plateauview-api \
            --image eukarya/plateauview-api:latest \
            --region asia-northeast1 \
            --platform managed \
            --quiet
```