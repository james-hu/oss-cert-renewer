# OSS Cert Renewer

This is a **Zero-Code** solution for automatically renewing SSL certificates for static websites hosted on Aliyun OSS (Alibaba Cloud Object Storage Service).

It uses GitHub Actions to deploy a serverless function to Aliyun Function Compute (FC), which then runs daily to check and renew your Let's Encrypt certificates.

## Features

- **Zero Code Required**: Just fork and configure secrets.
- **Multi-Site Support**: Manage multiple websites/buckets in a single deployment.
- **Automated**: Runs automatically on a schedule (default: daily).
- **Secure**: Sensitive credentials are stored in GitHub Secrets and passed securely to the runtime.
- **Low Cost**: Uses Aliyun Function Compute (Serverless), which typically falls within the free tier for low-volume usage.

## How to Use

### 1. Fork this Repository
Click the **Fork** button at the top right of this page to create your own copy of this repository.

### 2. Configure Secrets
Go to your forked repository's **Settings** > **Secrets and variables** > **Actions**.
Click **New repository secret** and add the following secrets:

| Secret Name | Description | Example Value |
|-------------|-------------|---------------|
| `ALIYUN_ACCOUNT_ID` | Your Alibaba Cloud Account ID. | `1234567890123456` |
| `ALIYUN_ACCESS_KEY_ID` | Access Key ID with permissions for OSS and FC. | `LTAI...` |
| `ALIYUN_ACCESS_KEY_SECRET` | Access Key Secret. | `abcde...` |
| `OSS_REGION` | The region where your OSS buckets are located. | `oss-cn-hangzhou` |
| `OSS_BUCKETS` | A comma-separated list of your OSS bucket names. | `my-blog,my-shop,landing-page` |
| `ACME_EMAIL` | Email address for Let's Encrypt registration. | `admin@example.com` |
| `CRON_EXPRESSION` | (Optional) CRON schedule for the function. | `0 0 4 * * 0` (Default: Weekly) |

> **Note**: Your Aliyun RAM user needs permissions to `AliyunOSSFullAccess` and `AliyunFCFullAccess` (or simpler custom policies if you prefer).

### 3. Deploy
1. Go to the **Actions** tab in your repository.
2. Select the **Deploy to Aliyun FC** workflow on the left.
3. Click the **Run workflow** button.

That's it! GitHub Actions will build the tool and deploy it to your Aliyun Function Compute account.

## How it Works
- **Deployment**: The GitHub Action builds the Go binary and manages the Aliyun Function Compute resource using Serverless Devs.
- **Execution**: The deployed function runs on a schedule (defined in `s.yaml`) automatically on Aliyun.
- **Renewal Process**:
    1. It iterates through each bucket listed in `OSS_BUCKETS`.
    2. It checks the SSL certificate expiration for the custom domain (CNAME) attached to the bucket.
    3. If a certificate is expiring (less than 30 days left) or missing, it triggers the renewal process using the ACME protocol (Let's Encrypt).
    4. It solves the HTTP-01 challenge by uploading a temporary file to your OSS bucket.
    5. It uploads the new certificate back to OSS updates the domain configuration.

## Customization (Optional)
To change the renewal schedule:
1. Go to your repository **Settings** > **Secrets and variables** > **Actions**.
2. Add or Update the repository secret `CRON_EXPRESSION`.
3. Re-run the **Deploy to Aliyun FC** workflow.

The default is `0 0 4 * * 0` (Every Sunday at 04:00 UTC).
Common examples:
- Daily: `0 0 4 * * *`
- Monthly: `0 0 4 1 * *`
