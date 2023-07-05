provider "aws" {
  region = local.region
}

locals {
  name        = "eks-demo"
  vpc_name    = "eks-vpc"
  region      = "us-east-1"
  vpc_cidr = "10.0.0.0/16"
  azs      = slice(data.aws_availability_zones.available.names, 0, 3)
}

data "aws_availability_zones" "available" {}
