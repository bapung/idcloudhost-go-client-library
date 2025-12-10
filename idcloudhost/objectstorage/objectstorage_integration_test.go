//go:build integration

package objectstorage

import (
	"net/http"
	"os"
	"strconv"
	"testing"
)

func getEnvOrSkip(t *testing.T, key string) string {
	val := os.Getenv(key)
	if val == "" {
		t.Fatalf("Environment variable %s not set", key)
	}
	return val
}

func TestObjectStorageIntegration(t *testing.T) {
	authToken := getEnvOrSkip(t, "IDCLOUDHOST_API_KEY")
	location := getEnvOrSkip(t, "IDCLOUDHOST_LOCATION")
	billingAccountStr := getEnvOrSkip(t, "IDCLOUDHOST_BILLING_ACCOUNT")

	billingAccount, err := strconv.Atoi(billingAccountStr)
	if err != nil {
		t.Fatalf("Invalid billing account: %v", err)
	}

	client := &http.Client{}
	api := ObjectStorageAPI{}
	if err := api.Init(client, authToken, location); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 1. Get S3 user info
	if err := api.GetS3Info(); err != nil {
		t.Fatalf("GetS3Info failed: %v", err)
	}
	t.Logf("S3 User ID: %d", api.S3User.ID)
	t.Logf("Storage Endpoint: %s", api.S3User.StorageEndpoint)

	// 2. Generate S3 access key
	if err := api.GenerateKey(); err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}
	accessKey := api.S3Key.AccessKey
	t.Logf("Generated access key: %s", accessKey)

	// Verify key was created
	if accessKey == "" {
		t.Error("Expected non-empty access key")
	}
	if api.S3Key.SecretKey == "" {
		t.Error("Expected non-empty secret key")
	}

	// 3. List keys and verify existence
	if err := api.GetKeys(); err != nil {
		t.Fatalf("GetKeys failed: %v", err)
	}
	found := false
	for _, key := range api.S3Keys {
		if key.AccessKey == accessKey {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Generated key not found in list")
	}

	// 4. Create bucket
	bucketName := "integration-test-bucket"
	if err := api.CreateBucket(bucketName, "private", billingAccount); err != nil {
		t.Fatalf("CreateBucket failed: %v", err)
	}
	t.Logf("Created bucket: %s", bucketName)

	// Verify bucket was created
	if api.Bucket.Name != bucketName {
		t.Errorf("Expected bucket name %s, got %s", bucketName, api.Bucket.Name)
	}

	// 5. Get bucket details
	if err := api.GetBucket(bucketName); err != nil {
		t.Fatalf("GetBucket failed: %v", err)
	}
	if api.Bucket.Name != bucketName {
		t.Errorf("Expected bucket name %s, got %s", bucketName, api.Bucket.Name)
	}

	// 6. List buckets and verify existence
	if err := api.ListBuckets(); err != nil {
		t.Fatalf("ListBuckets failed: %v", err)
	}
	bucketFound := false
	for _, bucket := range api.Buckets {
		if bucket.Name == bucketName {
			bucketFound = true
			break
		}
	}
	if !bucketFound {
		t.Fatalf("Created bucket not found in list")
	}

	// 7. Modify bucket ACL
	if err := api.ModifyBucket(bucketName, "public-read"); err != nil {
		t.Fatalf("ModifyBucket failed: %v", err)
	}
	if api.Bucket.ACL != "public-read" {
		t.Errorf("Expected ACL public-read, got %s", api.Bucket.ACL)
	}

	// 8. Delete bucket (cleanup)
	if err := api.DeleteBucket(bucketName); err != nil {
		t.Fatalf("DeleteBucket failed: %v", err)
	}
	t.Logf("Deleted bucket: %s", bucketName)

	// 9. Delete access key (cleanup)
	if err := api.DeleteKey(accessKey); err != nil {
		t.Fatalf("DeleteKey failed: %v", err)
	}
	t.Logf("Deleted access key: %s", accessKey)
}
