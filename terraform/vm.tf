module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "fnc-1-vpc"
  cidr = "10.0.0.0/16"

  azs             = ["ap-northeast-2a", "ap-northeast-2b"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24"]

  enable_nat_gateway = false
  enable_vpn_gateway = false

  tags = {
  }
}

resource "aws_cloud9_environment_ec2" "api_server" {
  instance_type = "t2.nano"
  name          = "api-server"

  image_id                    = "resolve:ssm:/aws/service/cloud9/amis/amazonlinux-2-x86_64"
  subnet_id                   = module.vpc.public_subnets[0]
  automatic_stop_time_minutes = "60"
}
