# OSS Cert Renewer

This is a **Zero-Code** solution for automatically renewing SSL certificates for static websites hosted on Aliyun OSS (Alibaba Cloud Object Storage Service).

It uses GitHub Actions to deploy a serverless function to Aliyun Function Compute (FC), which then runs periodically to check and renew your Let's Encrypt certificates.

## Features

- **Zero Code Required**: Just fork and configure secrets.
- **Multi-Site Support**: Manage multiple websites/buckets in a single deployment.
- **Automated**: Runs automatically on a schedule (default: every Monday at 03:33 UTC).
- **Secure**: Uses Aliyun RAM Roles with fine-grained permissions. Permanent Access Keys are never stored in the function environment.
- **Low Cost**: Uses Aliyun Function Compute (Serverless), which typically falls within the free tier for low-volume usage.

## How to Use

### 1. Fork this Repository
Click the **Fork** button at the top right of this page to create your own copy of this repository.

### 2. Configure Secrets and Variables
Go to your forked repository's **Settings** > **Secrets and variables** > **Actions**.

#### Secrets
Click **New repository secret** within the **Secrets** tab and add the following:

| Secret Name | Description | Example Value |
|-------------|-------------|---------------|
| `ALIYUN_ACCESS_KEY_ID` | Access Key ID for deployment (requires RAM management permissions). | `LTAI...` |
| `ALIYUN_ACCESS_KEY_SECRET` | Access Key Secret. | `abcde...` |
| `ACME_EMAIL` | Email address for Let's Encrypt registration. | `admin@example.com` |

`ALIYUN_ACCESS_KEY_ID` and `ALIYUN_ACCESS_KEY_SECRET` are the credentials used for deploying the Aliyun Function Compute function.
They are not used by the Aliyun Function Compute function during run time.

#### Variables
Click the **Variables** tab, then **New repository variable** and add the following:

| Variable Name | Description | Example Value |
|---------------|-------------|---------------|
| `ALIYUN_ACCOUNT_ID` | Your Alibaba Cloud Account ID. | `1234567890123456` |
| `OSS_REGION` | The region where your OSS buckets are located. | `oss-cn-hangzhou` |
| `OSS_BUCKETS` | A comma-separated list of your OSS bucket names. | `my-blog,my-shop` |
| `CRON_EXPRESSION` | (Optional) CRON schedule for the function. | `0 33 3 * * 1` (Default: Mo 03:33 UTC) |

### Setup the RAM role for the FC function

The Function Compute (FC) function requires a RAM role with specific permissions to operate correctly. The role name is hardcoded as:

```
acs:ram::${env(ALIYUN_ACCOUNT_ID)}:role/OssCertRenewerFunctionRole
```

**Best Practice:**
Define a custom policy as shown below and associate it with the above RAM role. This ensures the function has the minimum required permissions for OSS and Log Service operations.

```json
{
    "Version": "1",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "oss:ListCname",
                "oss:PutCname",
                "oss:PutObject",
                "oss:DeleteObject"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "log:PostLogStoreLogs",
                "log:CreateLogStore",
                "log:GetLogStore",
                "log:CreateProject",
                "log:GetProject"
            ],
            "Resource": "*"
        }
    ]
}
```

Attach this policy to the RAM role `OssCertRenewerFunctionRole` in your Aliyun account before deploying the function.

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
2. Click the **Variables** tab.
3. Add or Update the repository variable `CRON_EXPRESSION`.
4. Re-run the **Deploy to Aliyun FC** workflow.

The default is `0 33 3 * * 1` (Every Monday at 03:33 UTC).
Common examples:
- Daily at 4AM UTC: `0 0 4 * * *`
- Monthly at 4AM UTC on the 1st: `0 0 4 1 * *`

You can also customise by making code change in your forked repo.
