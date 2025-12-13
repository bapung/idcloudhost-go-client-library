package objectstorage

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type S3Key struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	UserID    string `json:"userId,omitempty"`
}

type S3User struct {
	ID              int     `json:"id,omitempty"`
	UserID          int     `json:"user_id,omitempty"`
	BillingAccount  int     `json:"billing_account_id"`
	TotalStorageGB  float64 `json:"total_storage_gb,omitempty"`
	UsedStorageGB   float64 `json:"used_storage_gb,omitempty"`
	StorageEndpoint string  `json:"storage_endpoint,omitempty"`
	ServiceEndpoint string  `json:"service_endpoint,omitempty"`
	CreatedAt       string  `json:"created_at,omitempty"`
	UpdatedAt       string  `json:"updated_at,omitempty"`
}

type Bucket struct {
	Name           string `json:"name"`
	BillingAccount int    `json:"billing_account_id,omitempty"`
	UserID         int    `json:"user_id,omitempty"`
	SizeBytes      int64  `json:"size_bytes,omitempty"`
	NumObjects     int    `json:"num_objects,omitempty"`
	Owner          string `json:"owner,omitempty"`
	IsSuspended    bool   `json:"is_suspended,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	ModifiedAt     string `json:"modified_at,omitempty"`
	ACL            string `json:"acl,omitempty"`
}

type ObjectStorageAPI struct {
	c           HTTPClient
	AuthToken   string
	ApiEndpoint string
	S3User      *S3User
	S3Key       *S3Key
	S3Keys      []S3Key
	Bucket      *Bucket
	Buckets     []Bucket
}

func (s *ObjectStorageAPI) Init(c HTTPClient, authToken string, location string) error {
	s.c = c
	s.AuthToken = authToken
	s.ApiEndpoint = "https://api.idcloudhost.com/v1/storage/api/s3"

	req, err := http.NewRequest("GET", s.ApiEndpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	r, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to verify endpoint: %v", err)
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode == http.StatusNotFound {
		return fmt.Errorf("object-storage API endpoint not available")
	}
	return nil
}

// CreateBucket creates a new S3 bucket
func (s *ObjectStorageAPI) CreateBucket(name string, billingAccountID int) error {
	data := url.Values{}
	data.Set("name", name)
	data.Set("billing_account_id", strconv.Itoa(billingAccountID))

	req, err := http.NewRequest("PUT", "https://api.idcloudhost.com/v1/storage/bucket",
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apikey", s.AuthToken)

	r, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if r.StatusCode == http.StatusConflict {
		return fmt.Errorf("bucket name '%s' already exists (bucket names are globally unique)", name)
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %v", r.StatusCode)
	}

	return json.NewDecoder(r.Body).Decode(&s.Bucket)
}

// DeleteBucket deletes a bucket
func (s *ObjectStorageAPI) DeleteBucket(name string) error {
	data := url.Values{}
	data.Set("name", name)

	req, err := http.NewRequest("DELETE", "https://api.idcloudhost.com/v1/storage/bucket",
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apikey", s.AuthToken)

	r, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if r.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%v", r.StatusCode)
	}

	return nil
}

// ListBuckets lists all buckets owned by the user
func (s *ObjectStorageAPI) ListBuckets() error {
	req, err := http.NewRequest("GET", "https://api.idcloudhost.com/v1/storage/bucket/list", nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apikey", s.AuthToken)

	r, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}

	return json.NewDecoder(r.Body).Decode(&s.Buckets)
}

// GetKeys retrieves the user's S3 access keys
func (s *ObjectStorageAPI) GetKeys() error {
	req, err := http.NewRequest("GET", "https://api.idcloudhost.com/v1/storage/user/keys", nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apikey", s.AuthToken)

	r, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}

	return json.NewDecoder(r.Body).Decode(&s.S3Keys)
}

// GenerateKey generates a new S3 access key
func (s *ObjectStorageAPI) GenerateKey() error {
	req, err := http.NewRequest("POST", "https://api.idcloudhost.com/v1/storage/user/keys", nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apikey", s.AuthToken)

	r, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}

	var keys []S3Key
	if err := json.NewDecoder(r.Body).Decode(&keys); err != nil {
		return err
	}

	if len(keys) > 0 {
		s.S3Key = &keys[len(keys)-1] // Get the last (newest) key
	}

	return nil
}

// DeleteKey deletes an S3 access key
func (s *ObjectStorageAPI) DeleteKey(accessKey string) error {
	data := url.Values{}
	data.Set("access_key", accessKey)

	req, err := http.NewRequest("DELETE", "https://api.idcloudhost.com/v1/storage/user/keys",
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apikey", s.AuthToken)

	r, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}

	return nil
}
