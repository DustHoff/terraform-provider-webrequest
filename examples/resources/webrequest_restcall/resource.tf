resource "webrequest_restcall" "call" {
  url              = "https://eoscet74ykdzldt.m.pipedream.net"
  ignorestatuscode = true
  body             = jsonencode({ "username" : "test", "email" : "test@example.com" })
  header = {
    Content-Type = "application/json"
    Accept       = "application/json"
  }

  create = {
    method = "POST"
    url    = "https://eoscet74ykdzldt.m.pipedream.net/create"
  }

  read = {
    method = "POST"
    url    = "https://eoscet74ykdzldt.m.pipedream.net/get"
  }

  update = {
    method = "POST"
    url    = "https://eoscet74ykdzldt.m.pipedream.net/update"
  }

  delete = {
    method = "POST"
    url    = "https://eoscet74ykdzldt.m.pipedream.net/delete"
  }

  lifecycle {
    postcondition {
      condition     = self.statuscode == 200
      error_message = "Received Statuscode should be http/200"
    }
  }
}
