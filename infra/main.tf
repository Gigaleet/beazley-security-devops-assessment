# Create a new VPC
resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "main-vpc"
  }
}

# Create one subnet per AZ
resource "aws_subnet" "public" {
  count                   = 3
  vpc_id                  = aws_vpc.main.id
  cidr_block              = cidrsubnet(var.vpc_cidr, 8, count.index)
  availability_zone       = data.aws_availability_zones.azs.names[count.index]
  map_public_ip_on_launch = true

  tags = {
    Name = "public-subnet-${count.index}"
  }
}

# Public Internet access
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.main.id
  tags = { Name = "main-igw" }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = { Name = "public-rt" }
}

resource "aws_route_table_association" "public" {
  count          = length(aws_subnet.public)
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

# Launch Template
resource "aws_launch_template" "web" {
  name_prefix   = "web-lt-"
  image_id      = data.aws_ami.amazon_linux2.id
  instance_type = var.instance_type

  tag_specifications {
    resource_type = "instance"
    tags = {
      Role = "web-server"
    }
  }
}

# Auto Scaling Group for exactly 3 instances (1 per AZ)
resource "aws_autoscaling_group" "web_asg" {
  name               = "web-asg"
  max_size           = 3
  min_size           = 3
  desired_capacity   = 3
  health_check_type  = "EC2"
  health_check_grace_period = 300

  vpc_zone_identifier = aws_subnet.public[*].id

  launch_template {
    id      = aws_launch_template.web.id
    version = "$Latest"
  }

  tag {
    key                 = "Name"
    value               = "web-server"
    propagate_at_launch = true
  }
}
