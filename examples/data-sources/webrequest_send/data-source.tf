data "webrequest_send" "call" {
  url    = "https://eoscet74ykdzldt.m.pipedream.net"
  method = "POST"
  ttl    = 3600
  body   = jsonencode({ "username" : "test", "email" : "test@example.com" })
  header = {
    Content-Type = "application/json"
    Accept       = "application/json"
  }
}
