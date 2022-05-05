terraform {
  backend "s3" {
    bucket  = "fnc-tfstates"
    key     = "fakenews-server/users.tfstate"
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
