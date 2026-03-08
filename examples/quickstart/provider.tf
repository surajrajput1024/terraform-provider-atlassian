terraform {
  required_providers {
    atlassian = {
      source  = "surajrajput1024/atlassian"
      version = "~> 0.1"
    }
  }
}

provider "atlassian" {
  domain    = var.atlassian_domain
  email     = var.atlassian_email
  api_token = var.atlassian_api_token
}
