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
    url    = "https://eoscet74ykdzldt.m.pipedream.net/create"
    body   = jsonencode({ "username" : "test", "email" : "test@example.com" })
  }

  read = {
    method = "POST"
    url    = "https://eoscet74ykdzldt.m.pipedream.net/get/{ID}"
    body   = jsonencode({ "id" : "{ID}", "username" : "test", "email" : "test@example.com" })
  }

  update = {
    method = "POST"
    url    = "https://eoscet74ykdzldt.m.pipedream.net/update"
    body   = jsonencode({ "username" : "test", "email" : "test@example.com" })
  }

  delete = {
    method = "POST"
    url    = "https://eoscet74ykdzldt.m.pipedream.net/delete"
    body   = jsonencode({ "username" : "test", "email" : "test@example.com" })
  }

  lifecycle {
    postcondition {
      condition     = self.statuscode == 200
      error_message = "Received Statuscode should be http/200"
    }
  }
}
