resource "webrequest_restcall" "call" {
  ignorestatuscode = true
  filter           = "//data"
  key              = "//json/username"

  header = {
    Content-Type = "application/json"
    Accept       = "application/json"
  }

  create = {
    method = "POST"
    url    = "https://httpbin.org/post"
    body   = jsonencode({ "username" : "test", "email" : "test@example.com" })
  }

  read = {
    method = "GET"
    url    = "https://httpbin.org/get/{ID}"
    body   = jsonencode({ "id" : "{ID}", "username" : "test", "email" : "test@example.com" })
  }

  update = {
    method = "PUT"
    url    = "https://httpbin.org/put"
    body   = jsonencode({ "username" : "test", "email" : "test@example.com" })
  }

  delete = {
    method = "DELETE"
    url    = "https://httpbin.org/delete"
    body   = jsonencode({ "username" : "test", "email" : "test@example.com" })
  }

  lifecycle {
    postcondition {
      condition     = self.statuscode == 200
      error_message = "Received Statuscode should be http/200"
    }
  }
}
