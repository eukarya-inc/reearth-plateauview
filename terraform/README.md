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
gcloud services enable secretmanager.googleapis.com
gcloud services enable certificatemanager.googleapis.com
#TODO: その他必要なものもあるが、事前有効化済みだったので後で確認する
```

### CloudDNSのセットアップ

今回セットアップしたいドメインをcloud dnsでホスティングできるように設定する。
手順は各種レジストラの手順を参照。


### Auth0 Management APIの設定

[公式のQuick Start](https://github.com/auth0/terraform-provider-auth0/blob/main/docs/guides/quickstart.md)を参考にセットアップ

## Terraform実行手順

### 設定ファイルの準備

今回設定したい環境の設定ファイル(tfvars)を準備する。実行環境費に合わせて必要な情報をセットアップを行う。
[example.tfvars](./env/example.tfvars) をコピーし、必要な設定を追記する。

### backendの作成

```
export SERVICE_PREFIX=""
gcloud storage buckets create gs://$(SERVICE_PREFIX)-terraform-tfstate
```

[terraform.tf](terraform.tf)の `backend`のbucketを設定する

```diff
  backend "gcs" {
-    bucket = ""
+    bucket = "$(SERVICE_PREFIX)-terraform-tfstate"
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
echo -n "${MONGO_CONNECTION}"| gcloud secrets versions add reearth-api-REEARTH_DB --data-file=-
echo -n "${REEARTH_API_SECRET}"| gcloud secrets versions add reearth-api-REEARTH_MARKETPLACE_SECRET --data-file=-
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
            --image reearth/reearth:0.14.1 \
            --region asia-northeast1 \
            --platform managed \
            --quiet
```