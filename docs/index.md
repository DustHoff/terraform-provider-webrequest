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

# disclaimer 

I develop this terraform in my spare time. in any case you find a bug please report it on
github with detailed information. if you already find a solution please feel free to create a
pull request.
