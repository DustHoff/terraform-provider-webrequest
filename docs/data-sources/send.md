# Data Source `send`

The send data source provides the ability to send web requests like curl.

## Example Usage

### Simple Example
```terraform
data "webrequest_send" "example" {
  url = "https://example.com/"
  method = "GET"
  
}
```
### Komplex Example
````terraform
data "webrequest_send" "example" {
  url = "https://example.com/"
  method = "POST"
  body = jsondecode({body:"value"})
  header = [
    {
      name = "Content-Type"
      value = "application/json"
    }]
  ttl = 3600
}
````

## Attributes Reference

* `url` - (Required) target endpoint of the webservice, incl. scheme host port and resource path
* method - (Optional) request method, default value is GET, Possible Values GET, POST, PUT, DELETE, OPTIONS, HEAD
* body - (Optional) request body, default is a empty body
* header - (Optional) a list of request header. Default is a empty list. See [Header](#Header) below for details
* ttl  - (Optional) time to live of the received response. The Value represents the seconds of validity.

### Header

* name
* value

-> The following attributes are exported.

- `result` contains the response body corresponding to the request
- `expires` a unix timestamp on which invalidate the response 
