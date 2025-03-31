provider "aws" {
  region     = "ca-central-1"
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key
}

variable "aws_access_key" {
  type = string
}

variable "aws_secret_key" {
  type = string
}

# Data block to fetch existing security group if it exists
data "aws_security_group" "tcp_sg_existing" {
  name   = "allow_tcp_traffic"
  count  = try([data.aws_security_group.tcp_sg_existing[0].id], []) != [] ? 1 : 0 # Check if it exists
}

# Resource block to create security group if it doesn't exist
resource "aws_security_group" "tcp_sg" {
  count       = length(data.aws_security_group.tcp_sg_existing) > 0 ? 0 : 1 # Create only if not found
  name        = "allow_tcp_traffic"
  description = "Allow TCP and SSH inbound traffic"
  ingress {
    from_port   = 10000
    to_port     = 10000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 10001
    to_port     = 10001
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Data block to fetch existing IAM role if it exists
data "aws_iam_role" "ec2_role_existing" {
  name  = "ec2_self_terminate"
  count = try([data.aws_iam_role.ec2_role_existing[0].arn], []) != [] ? 1 : 0 # Check if it exists
}

# Resource block to create IAM role if it doesn't exist
resource "aws_iam_role" "ec2_role" {
  count = length(data.aws_iam_role.ec2_role_existing) > 0 ? 0 : 1 # Create only if not found
  name  = "ec2_self_terminate"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
    }]
  })
}

# IAM role policy (only if creating the role)
resource "aws_iam_role_policy" "terminate_policy" {
  count = length(data.aws_iam_role.ec2_role_existing) > 0 ? 0 : 1 # Attach only if creating role
  role  = aws_iam_role.ec2_role[0].id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = ["ec2:TerminateInstances"]
      Resource = "*"
    }]
  })
}

# Data block to fetch existing instance profile if it exists
data "aws_iam_instance_profile" "ec2_profile_existing" {
  name  = "ec2_terminate_profile"
  count = try([data.aws_iam_instance_profile.ec2_profile_existing[0].arn], []) != [] ? 1 : 0 # Check if it exists
}

# Resource block to create instance profile if it doesn't exist
resource "aws_iam_instance_profile" "ec2_profile" {
  count = length(data.aws_iam_instance_profile.ec2_profile_existing) > 0 ? 0 : 1 # Create only if not found
  name  = "ec2_terminate_profile"
  role  = length(data.aws_iam_role.ec2_role_existing) > 0 ? data.aws_iam_role.ec2_role_existing[0].name : aws_iam_role.ec2_role[0].name
}

resource "aws_instance" "go_tcp_servers" {
  ami                    = "ami-07f7608a8efba8d78"
  instance_type          = "t2.micro"
  key_name               = "my-key-pair-go-tcp-server"
  vpc_security_group_ids = [length(data.aws_security_group.tcp_sg_existing) > 0 ? data.aws_security_group.tcp_sg_existing[0].id : aws_security_group.tcp_sg[0].id]
  iam_instance_profile   = length(data.aws_iam_instance_profile.ec2_profile_existing) > 0 ? data.aws_iam_instance_profile.ec2_profile_existing[0].name : aws_iam_instance_profile.ec2_profile[0].name
  user_data              = <<-EOF
                          #!/bin/bash
                          yum update -y
                          wget https://go.dev/dl/go1.22.1.linux-amd64.tar.gz
                          tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz
                          echo "export PATH=$PATH:/usr/local/go/bin" >> /home/ec2-user/.bashrc
                          source /home/ec2-user/.bashrc
                          yum install git -y
                          yum install awscli -y
                          git clone https://github.com/thisisadityapatel/just-servers.git /home/ec2-user/tcp-server
                          cd /home/ec2-user/tcp-server
                          /usr/local/go/bin/go build -o server server.go
                          nohup ./server &
                          INSTANCE_ID=$(curl -s http://169.254.169.254/latest/meta-data/instance-id)
                          sleep 600
                          aws ec2 terminate-instances --instance-ids $INSTANCE_ID
                          EOF
  tags = {
    Name = "MyGoTCPServer"
  }
}

output "public_ip" {
  value = aws_instance.go_tcp_servers.public_ip
}