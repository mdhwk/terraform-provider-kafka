provider "kafka" {
  bootstrap_servers = [
    "broker-1:9094",
    "broker-2:9094",
    "broker-3:9094"
  ]
}