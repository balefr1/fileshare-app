Fileshare app to be deployed on ECS

Simple REST API in Go (Gin+GORM) with Reactjs Frontend.

The Api allows to register/manage users, each user can upload files choosing whether to save them to S3 or on local filesystem.

The Go application needs the following env vars:
"AWS_REGION"   --> the aws region
"S3_BUCKET"    --> Name of the bucket where to save user uploads *
"USER_FILE_PATH"  --> path on local fs where to save user uploads
"DB_HOST" --> MySQL db host to connect our Go app to **




---
*make sure the AWS role in use is granted s3:PutObject,s3:GetObject on all bucket objects,
this application is meant to run on Docker containers with access to AWS (via ECS task role)

**Db access config is located at Config/Database.go (ovverride db credentials here)
