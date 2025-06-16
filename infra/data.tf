# Find the latest Amazon Linux 2 AMI
data "aws_ami" "amazon_linux2" {
  owners      = ["amazon"]
  most_recent = true

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }
}

data "aws_availability_zones" "azs" {
  state = "available"
}

data "aws_autoscaling_group" "web" {
  name = aws_autoscaling_group.web_asg.name
}

data "aws_instances" "web_asg" {
  instance_tags = {
    "aws:autoscaling:groupName" = aws_autoscaling_group.web_asg.name
  }
}
