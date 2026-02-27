terraform {
  required_providers {
    atlassian = {
      source  = "surajrajput1024/atlassian"
      version = ">= 0.1"
    }
  }
}

provider "atlassian" {
  domain    = "your-site.atlassian.net"
  email     = "you@example.com"
  api_token = var.atlassian_api_token
}

variable "atlassian_api_token" {
  type      = string
  sensitive = true
}
