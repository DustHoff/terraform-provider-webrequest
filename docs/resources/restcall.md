# Resource `restcall`

The resource function provides the ability to interact with a rest API.
It integrates all parts of a rest API.

While resource creation, a rest post API call is executed. After creation, any manipulation of this resource results in a
PUT API call. In any case, the resource gets deleted a delete API call is executed.

The resource requires to fulfill the rest API convention.


POST   - scheme://host:port/users

GET    - scheme://host:port/users/1

PUT    - scheme://host:port/users/1

DELETE - scheme://host:port/users/1

## Example Usage

### Simple Example
```terraform
data "webrequest_restcall" "example" {
  url = "https://example.com/"
  body = jsondecode({username:"example",email:"simple@example.com"})
  
}
```
### Complex Example
````terraform
data "webrequest_restcall" "example" {
  url = "https://example.com/"
  body = jsondecode({username:"example",email:"simple@example.com"})
  header = [
    {
      name = "Content-Type"
      value = "application/json"
    }]
}
````

## Attributes Reference

* `url` - (Required) target endpoint of the webservice, incl. scheme host port and resource path
* `body` - (Required) request body
* `header` - (Optional) a list of request header. Default is a empty list. See [Header](#Header) below for details
* `key` - (Optional) object field ob the received data

### Header

* name
* value

-> The following attributes are exported.

- `result` contains the response body corresponding to the request
- `objectid` a unix timestamp on which invalidate the response 
