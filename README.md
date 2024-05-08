
# Storage Service

Handles blob storage to the Olympsis platform


This service contains two http endpoints:

## POST /v1/storage/{fileBucket}
- Uploads an object to storage
- Requires X-Filename header to name file appropriately

## DELETE /v1/storage/{fileBucket}
- Deletes an object from storage
- Requires X-Filename header to find file appropriately

