# Data Source `send`

The send data source provides the ability to send web requests like curl.

## Example Usage

```terraform
data "webrequest_send" "example" {
  url = "https://example.com/"
  method = "GET"
  
}
```
## Attributes Reference

* url
* method
* body
* header a list of request header. See [Header](#Header) below for details

### Header

* name
* value

-> The following attributes are exported.

- `result` contains the response body corresponding to the request
