io_mode = "async"

service "http" {
  listen_addr = "127.0.0.1:8080"

  process "main" {
    command = ["/usr/local/bin/awesome-app", "server"]
  }

  process "mgmt" {
    command = ["/usr/local/bin/awesome-app", "mgmt"]
  }
}

/*
output like:
  {
    "io_mode": "async",
    "service": {
      "http": {
        "web_proxy": {
          "listen_addr": "127.0.0.1:8080",
          "process": {
            "main": {
              "command": ["/usr/local/bin/awesome-app", "server"]
            },
            "mgmt": {
              "command": ["/usr/local/bin/awesome-app", "mgmt"]
            },
          }
        }
      }
    }
  }
*/