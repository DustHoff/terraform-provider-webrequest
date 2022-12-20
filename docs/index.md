# webrequest Provider

The webrequest provider is used to interact with external services. 
This provider support any service which is accessible via http or https.

* use providers to fetch data from a web service by sending GET, POST, PUT or Delete requests.

## Example Usage

```terraform
provider "webrequest" {
  timeout = 30
}
```
