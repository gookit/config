io_mode = "async"
pkg_name = "config"

service "http" {
  listen_addr = "127.0.0.1:9999"

  process "main" {
    command = [
      "/usr/local/bin/awesome-app",
      "server"]
  }
}
