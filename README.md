# AWS S3 Multi-File Upload Service

A lightweight, reusable Go package for uploading and managing multiple files on AWS S3. This service provides a clean API for handling file uploads with support for concurrent operations, progress tracking, and error handling.

## Features

- 📤 Upload single or multiple files concurrently
- 🔒 Presigned URL generation for secure downloads
- 📊 File metadata retrieval
- 🛡️ Comprehensive error handling
- 🔧 Easy integration into larger projects
- 📝 Support for custom content types
- 🎯 Context-aware operations

## Prerequisites

- Go 1.24 or higher
- AWS Account with S3 access
- AWS credentials configured (Access Key ID and Secret Access Key)

## Installation

```bash
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/s3
go get github.com/aws/aws-sdk-go-v2/feature/s3/manager
```

## Configuration

Set up your AWS credentials using one of these methods:

### Environment Variables
```bash
export AWS_ACCESS_KEY=your_access_key
export AWS_SECRET_KEY=your_secret_key
export AWS_REGION=us-east-1
export AWS_BUCKET_NAME=your-bucket-name
export AWS_ARN=your-aws-arn
```

## Error Handling

The service returns descriptive errors for common scenarios:
- Invalid bucket configuration
- Network failures
- Permission issues
- File not found

Always check returned errors in production code.


## Security Best Practices

- Never commit AWS credentials to version control
- Use IAM roles when deploying to AWS services
- Implement bucket policies for access control
- Enable S3 bucket versioning
- Use presigned URLs for temporary access
- Validate file types and sizes before upload

**Note**: This is a foundational service designed for easy integration. Customize according to your specific requirements.