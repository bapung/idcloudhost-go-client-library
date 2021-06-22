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

### Disclaimer
For now, I am developing this client library as prerequisites to my other project: [Terraform Provider for idCloudHost](https://github.com/bapung/terraform-provider-idcloudhost). The reason is I want to codify my project in IdCloudHost (afaik the cheapest Cloud VPS in Indonesia right now, which fit my needs). 
This is a fun project and not guaranted for stability, if you found bugs or need to add feature feel free to create issue or PR.
