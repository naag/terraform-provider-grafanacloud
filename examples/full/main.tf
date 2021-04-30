provider "grafanacloud" {
  organisation = var.organisation
}

resource "grafanacloud_portal_api_key" "prometheus_remote_write" {
  name = "prometheus-remote-write"
  role = "MetricsPublisher"
}

resource "grafanacloud_grafana_api_key" "api_client" {
  name  = "api_client"
  role  = "Editor"
  stack = var.stack
}

data "grafanacloud_stack" "demo" {
  name = var.stack
}

data "grafanacloud_stacks" "all" {
}
