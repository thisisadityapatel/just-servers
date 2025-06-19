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

resource "aws_security_group" "tcp_sg" {
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
    from_port   = 10002
    to_port     = 10002
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 10003
    to_port     = 10003
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

  lifecycle {
    ignore_changes = [name] # Avoid conflicts with existing resource
  }
}

resource "aws_iam_role" "ec2_role" {
  name = "ec2_self_terminate"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
    }]
  })

  lifecycle {
    ignore_changes = [name] # Avoid conflicts with existing role
  }
}

resource "aws_iam_role_policy" "terminate_policy" {
  role = aws_iam_role.ec2_role.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = ["ec2:TerminateInstances"]
      Resource = "*"
    }]
  })
}

resource "aws_iam_instance_profile" "ec2_profile" {
  name = "ec2_terminate_profile"
  role = aws_iam_role.ec2_role.name

  lifecycle {
    ignore_changes = [name] # Avoid conflicts with existing profile
  }
}

resource "aws_instance" "go_tcp_servers" {
  ami                    = "ami-07f7608a8efba8d78"
  instance_type          = "t2.micro"
  key_name               = "my-key-pair-go-tcp-server"
  vpc_security_group_ids = [aws_security_group.tcp_sg.id]
  iam_instance_profile   = aws_iam_instance_profile.ec2_profile.name
  user_data              = <<-EOF
                          #!/bin/bash
                          yum update -y
                          wget https://go.dev/dl/go1.22.1.linux-amd64.tar.gz
                          tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz
                          echo "export PATH=$PATH:/usr/local/go/bin" >> /home/ec2-user/.bashrc
                          echo "export HOME=/home/ec2-user" >> /home/ec2-user/.bashrc
                          source /home/ec2-user/.bashrc
                          yum install git -y
                          yum install awscli -y
                          git clone https://github.com/thisisadityapatel/just-servers.git /home/ec2-user/tcp-server
                          cd /home/ec2-user/tcp-server
                          nohup /usr/local/go/bin/go run main.go >> /home/ec2-user/server.log 2>&1 &
                          TOKEN=$(curl -s -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")
                          INSTANCE_ID=$(curl -s -H "X-aws-ec2-metadata-token: $TOKEN" http://169.254.169.254/latest/meta-data/instance-id 2>> /home/ec2-user/terminate.log)
                          sleep 600
                          aws ec2 terminate-instances --instance-ids $INSTANCE_ID &> /home/ec2-user/server.log
                          EOF
  tags = {
    Name = "MyGoTCPServer"
  }
}

output "public_ip" {
  value = aws_instance.go_tcp_servers.public_ip
}