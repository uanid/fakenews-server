terraform {
  backend "s3" {
    bucket  = "fnc-tfstates"
    key     = "fakenews-server/generic.tfstate"
    region  = "ap-northeast-2"
    profile = "fnc"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

locals {
  project_name = "fnc-1"
}

provider "aws" {
  region  = "ap-northeast-2"
  profile = "fnc"

  default_tags {
    tags = {
      Project = "fnc-1"
      Owner   = "musong"
    }
  }
}

data "aws_caller_identity" "current" {
}
