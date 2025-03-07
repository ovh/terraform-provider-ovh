resource "ovh_cloud_project_user" "user" {
  service_name = "XXX"
  description  = "my user"
  role_names   = [
    "objectstore_operator"
  ]
}

resource "ovh_cloud_project_user_s3_credential" "my_s3_credentials" {
  service_name = ovh_cloud_project_user.user.service_name
  user_id      = ovh_cloud_project_user.user.id
}

resource "ovh_cloud_project_user_s3_policy" "policy" {
  service_name = ovh_cloud_project_user.user.service_name
  user_id      = ovh_cloud_project_user.user.id
  policy       = jsonencode({
    "Statement":[{
      "Sid": "RWContainer",
      "Effect": "Allow",
      "Action":["s3:GetObject", "s3:PutObject", "s3:DeleteObject", "s3:ListBucket", "s3:ListMultipartUploadParts", "s3:ListBucketMultipartUploads", "s3:AbortMultipartUpload", "s3:GetBucketLocation"],
      "Resource":["arn:aws:s3:::hp-bucket", "arn:aws:s3:::hp-bucket/*"]
    }]
  })
}
