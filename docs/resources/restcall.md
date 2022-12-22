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

To use a rest api please only provide the base url to the service. any object specific url will be generated based on the 
configuration. Default this resource assumes that the received JSON object contains a id field. This field is used to generated
object specific urls.

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
  key = "user_id"
}
````

## Attributes Reference

* `url` - (Required) target endpoint of the webservice, incl. scheme host port and resource path
* `body` - (Required) request body
* `header` - (Optional) a list of request header. Default is a empty list. See [Header](#Header) below for details
* `key` - (Optional) primary key object field of the received data

### Header

* name
* value

-> The following attributes are exported.

- `result` contains the response body corresponding to the request
- `objectid` the primary key value of the received object 
