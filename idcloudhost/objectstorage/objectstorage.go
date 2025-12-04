package objectstorage

import (
	"encoding/json"
	"errors"
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
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	UserID    int    `json:"user_id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
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
	ID             int    `json:"id,omitempty"`
	Name           string `json:"name"`
	BillingAccount int    `json:"billing_account_id,omitempty"`
	UserID         int    `json:"user_id,omitempty"`
	SizeBytes      int64  `json:"size_bytes,omitempty"`
	Region         string `json:"region,omitempty"`
	ACL            string `json:"acl"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
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
	s.ApiEndpoint = "https://api.idcloudhost.com/v1/object-storage"

	r, err := http.Get(s.ApiEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode == http.StatusNotFound {
		return fmt.Errorf("object-storage API endpoint not available")
	}
	return nil
}

// GetS3Info returns S3 API information for the user
func (s *ObjectStorageAPI) GetS3Info() error {
	req, err := http.NewRequest("GET", s.ApiEndpoint, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", s.AuthToken)
	r, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v", r.StatusCode))
	}
	return json.NewDecoder(r.Body).Decode(&s.S3User)
}

// CreateBucket creates a new S3 bucket
func (s *ObjectStorageAPI) CreateBucket(name string, acl string, billingAccountID int) error {
	data := url.Values{}
	data.Set("name", name)
	data.Set("acl", acl)
	data.Set("billing_account_id", strconv.Itoa(billingAccountID))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/buckets", s.ApiEndpoint),
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apiKey", s.AuthToken)

	r, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v", r.StatusCode))
	}

	return json.NewDecoder(r.Body).Decode(&s.Bucket)
}

// ModifyBucket updates a bucket's ACL
func (s *ObjectStorageAPI) ModifyBucket(name string, acl string) error {
	data := url.Values{}
	data.Set("acl", acl)

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/buckets/%s", s.ApiEndpoint, name),
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apiKey", s.AuthToken)

	r, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v", r.StatusCode))
	}

	return json.NewDecoder(r.Body).Decode(&s.Bucket)
}

// DeleteBucket deletes a bucket
func (s *ObjectStorageAPI) DeleteBucket(name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/buckets/%s", s.ApiEndpoint, name), nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", s.AuthToken)

	r, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v", r.StatusCode))
	}

	return nil
}

// GetBucket gets information about a specific bucket
func (s *ObjectStorageAPI) GetBucket(name string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/buckets/%s", s.ApiEndpoint, name), nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", s.AuthToken)

	r, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v", r.StatusCode))
	}

	return json.NewDecoder(r.Body).Decode(&s.Bucket)
}

// ListBuckets lists all buckets owned by the user
func (s *ObjectStorageAPI) ListBuckets() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/buckets", s.ApiEndpoint), nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", s.AuthToken)

	r, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v", r.StatusCode))
	}

	return json.NewDecoder(r.Body).Decode(&s.Buckets)
}

// GetKeys retrieves the user's S3 access keys
func (s *ObjectStorageAPI) GetKeys() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/keys", s.ApiEndpoint), nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", s.AuthToken)

	r, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v", r.StatusCode))
	}

	return json.NewDecoder(r.Body).Decode(&s.S3Keys)
}

// GenerateKey generates a new S3 access key
func (s *ObjectStorageAPI) GenerateKey() error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/keys", s.ApiEndpoint), nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", s.AuthToken)

	r, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v", r.StatusCode))
	}

	return json.NewDecoder(r.Body).Decode(&s.S3Key)
}

// DeleteKey deletes an S3 access key
func (s *ObjectStorageAPI) DeleteKey(accessKey string) error {
	data := url.Values{}
	data.Set("access_key", accessKey)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/keys", s.ApiEndpoint),
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apiKey", s.AuthToken)

	r, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v", r.StatusCode))
	}

	return nil
}
