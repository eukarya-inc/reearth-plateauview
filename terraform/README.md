# Reearth Plateauview for Terraform

Reearth Plateauviewを構築するTerraformです。

## 事前準備

 * コマンドラインツールのインストール
   * gcloud
   * terraform

## gcloud のセットアップ

```
# GCP Project ID
export PROJECT_ID=""
# 使いたいドメイン 例: plateauview.example.com
export DOMAIN=""
# 20文字以内・半角英数ハイフンで自由に決めて良い。例: plateauview-test
export SERVICE_PREFIX=""
# 20文字以内・半角英数ハイフンで自由に決めて良い。SERVICE_PREFIXと同じで問題ない。
export ZONE_NAME=${SERVICE_PREFIX}

gcloud components update
gcloud config configurations create ${SERVICE_PREFIX}
gcloud config set project ${PROJECT_ID}
gcloud auth login

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

gcloud config set compute/region asia-northeast1
gcloud auth application-default login
```

### CloudDNSのセットアップ

```
 gcloud dns managed-zones create ${ZONE_NAME} --dns-name ${DOMAIN} --description "${ZONE_NAME}"
 gcloud dns record-sets list --zone ${ZONE_NAME}
```

```
NAME                           TYPE  TTL    DATA
*********  NS    21600  ns-cloud-a1.googledomains.com.,ns-cloud-a2.googledomains.com.,ns-cloud-a3.googledomains.com.,ns-cloud-a4.googledomains.com.
*********  SOA   21600  ns-cloud-a1.googledomains.com. cloud-dns-hostmaster.google.com. 1 21600 3600 259200 300
```

今回セットアップしたいドメインをcloud dnsでホスティングできるように設定する。
出力されたNSレコードの情報を用いて、ネームサーバーを変更する。手順は各種レジストラの手順を参照。

### MongoDB Atlas のセットアップ

データベースを作成して接続文字列を取得する。DBユーザーの作成やIPアドレスの許可（全IP許可）を忘れずに。

```bash
export REEARTH_DB=""
```

### Auth0 のセットアップ

Auth0テナントを作成した後、[公式のQuick Start](https://github.com/auth0/terraform-provider-auth0/blob/main/docs/guides/quickstart.md)を参考に、M2Mのアプリケーションをセットアップ

```bash
export AUTH0_CLIENT_SECRET=""
```

## Terraform実行手順

### 設定ファイルの準備

今回設定したい環境の設定ファイル(tfvars)を準備する。実行環境費に合わせて必要な情報をセットアップを行う。
[example.tfvars](./env/example.tfvars) をコピーし、必要な設定を追記する。

### backendの作成
リソース全般の作成に必要なSERVICE_PREFIXを決める。
```
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
terraform apply -var-file=env/example.tfvars
```


### terraform実行後のシークレットへの値追加

```bash
echo -n "${REEARTH_DB}" | gcloud secrets versions add reearth-api-REEARTH_DB --data-file=-
echo -n "${REEARTH_DB}" | gcloud secrets versions add reearth-cms-REEARTH_CMS_WORKER_DB --data-file=-
echo -n "${REEARTH_DB}" | gcloud secrets versions add reearth-cms-REEARTH_CMS_DB --data-file=-

# 以下は必要に応じて設定する。設定しなくても次に進むことは可能。
# FMEのトークン
echo -n "${REEARTH_PLATEAUVIEW_FME_TOKEN}" | gcloud secrets versions add reearth-cms-REEARTH_PLATEAUVIEW_FME_TOKEN --data-file=-
# G空間情報センターのAPIトークン
echo -n "${REEARTH_PLATEAUVIEW_CKAN_TOKEN}" | gcloud secrets versions add reearth-cms-REEARTH_PLATEAUVIEW_CKAN_TOKEN --data-file=-
# SendGridのAPIキー
echo -n "${REEARTH_PLATEAUVIEW_SENDGRID_APIKEY}" | gcloud secrets versions add reearth-cms-REEARTH_PLATEAUVIEW_SENDGRID_APIKEY --data-file=-
```

### Cloud Run のデプロイ

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

### DNS・ロードバランサ・証明書のデプロイ完了まで待機する

```
curl https://api.${DOMAIN}/ping
```

を繰り返し試行し `"pong"` が返ってくるまで待つ。

### Auth0 ユーザー作成

先ほど作成したAuth0テナントにてユーザーを作成する。メールアドレスの認証を忘れずに。

### CMS ログイン

CMSにログインする。

https://cms.${DOMAIN}

ログイン後、ワークスペース・Myインテグレーションを作成する。

次に、インテグレーション内に以下の通りWebhookを作成する。作成後、有効化を忘れないこと。

- URL: terraform outputs の plateauview_cms_webhook_url
- シークレット: terraform outputs の plateauview_cms_webhook_secret
- 種類: 全てにチェックを入れる。

作成後、作成したワークスペースに作成したインテグレーションを追加し、オーナー権限に変更する。

### CMS_TOKENの設定

CMSのUIで発行したインテグレーションのトークンを登録する。

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

### 完了

- Re:Earth: https://reearth.${DOMAIN}
- CMS: https://cms.${DOMAIN}

なお、terraform outputs で得られた以下は別のセットアップで使用する。

- plateauview_sdk_token: SDKのトークン。SDKのUIで設定する。
- plateauview_sidebar_token: サイドバーのAPIトークン。エディタ上でサイドバーウィジェットの設定から設定する。
