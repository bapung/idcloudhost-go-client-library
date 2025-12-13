//go:build integration

package objectstorage

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
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

	// 1. Generate S3 access key
	if err := api.GenerateKey(); err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}
	accessKey := api.S3Key.AccessKey
	secretKey := api.S3Key.SecretKey
	t.Logf("Generated access key: %s", accessKey)
	t.Logf("Generated secret key: %s", secretKey)

	// Verify key was created
	if accessKey == "" {
		t.Error("Expected non-empty access key")
	}
	if secretKey == "" {
		t.Error("Expected non-empty secret key")
	}

	// 2. List keys and verify existence
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

	// 3. Create bucket with unique name to avoid 409 conflicts
	bucketName := fmt.Sprintf("integration-test-%d", time.Now().Unix())
	if err := api.CreateBucket(bucketName, billingAccount); err != nil {
		t.Fatalf("CreateBucket failed: %v", err)
	}
	t.Logf("Created bucket: %s", bucketName)

	// Verify bucket was created
	if api.Bucket.Name != bucketName {
		t.Errorf("Expected bucket name %s, got %s", bucketName, api.Bucket.Name)
	}

	// 4. List buckets and verify existence
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

	// 5. Delete bucket (cleanup)
	if err := api.DeleteBucket(bucketName); err != nil {
		t.Fatalf("DeleteBucket failed: %v", err)
	}
	t.Logf("Deleted bucket: %s", bucketName)

	// 6. Delete access key (cleanup)
	if err := api.DeleteKey(accessKey); err != nil {
		t.Fatalf("DeleteKey failed: %v", err)
	}
	t.Logf("Deleted access key: %s", accessKey)
}
