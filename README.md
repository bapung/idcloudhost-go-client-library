# idcloudhost-go-client-library

Golang client library for [IdCloudHost](https://idcloudhost.com/)

## Example usages

### Create VM

```go
var vmApi = idcloudhost.VirtualMachineAPI{}

vmApi.Init(&http.Client{}, "definitely-legit-auth-token", "jkt01")
var newVM = idcloudhost.NewVM{
  Backup:          false,
  Name:            "testvm",
  OSName:          "ubuntu",
  OSVersion:       "16.04",
  Disks:           20,
  VCPU:            1,
  Memory:          1024,
  Username:        "example",
  InitialPassword: "Password123",
  BillingAccount:  9999, //make sure it is linked to correct authToken
}
err := vmApi.Create(NewVM)
if err != nil {
  //handle error
}

```

### Attach Disk to VM

```go
var (
  diskApi = idcloudhost.DiskAPI{}
  sizeGb = 200
)

diskApi.Bind("existing-vm-in-interest")
err := diskApi.Create(sizeGb)
if err != nil {
  //handle error
}
```

### Create and attach Floating IP

```go
var (
  ipApi = idcloudhost.FloatingIPAPI{}
  billingAcc = 1111
)

ipApi.Init(&http.Client{}, "definitely-legit-auth-token", "jkt01")
err := ipApi.Create("desired name", billingAcc)
if err != nil {
  //handle error
}
myIP := ipApi.FloatingIP.Address
err = ipApi.Assign(myIP, "valid-existing-vm-uuid")
if err != nil {
  //handle error
}
```

## Testing

This library includes both unit tests and integration tests.

### Unit Tests

Unit tests use mocked HTTP clients and don't require API credentials:

```bash
go test ./...
```

### Integration Tests

Integration tests create real resources in your IDCloudHost account. Set up environment variables first:

```bash
export IDCLOUDHOST_API_KEY="your-api-key"
export IDCLOUDHOST_LOCATION="jkt01"
export IDCLOUDHOST_BILLING_ACCOUNT="1234567890"
```

Run integration tests:

```bash
# Run all integration tests
go test -tags=integration ./...

# Or use the convenience script
./run_integration_tests.sh

# Run specific resource tests
./run_integration_tests.sh firewall
./run_integration_tests.sh floatingip
./run_integration_tests.sh network
./run_integration_tests.sh loadbalancer
./run_integration_tests.sh objectstorage
```
