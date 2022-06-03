resource "aws_security_group" "elb_sg" {
  name   = "${local.project_name}-elb-sg"
  vpc_id = "vpc-01dcfab8b9eca3620"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_lb" "elb" {
  name               = local.project_name
  internal           = false
  load_balancer_type = "network"
  security_groups    = [aws_security_group.elb_sg.id]
  subnets            = ["subnet-0a35954ed7de74c49", "subnet-082af9a032c6f14e1"]

  enable_deletion_protection = false
}

resource "aws_lb_target_group" "api_server" {
  name                 = local.project_name
  port                 = 8080
  protocol             = "TCP"
  vpc_id               = "vpc-01dcfab8b9eca3620"
  deregistration_delay = 60
}

resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.elb.arn
  port              = "443"
  protocol          = "TLS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = aws_acm_certificate.cert.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api_server.arn
  }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.elb.arn
  port              = "80"
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api_server.arn
  }
}

resource "aws_lb_target_group_attachment" "api_server" {
  target_group_arn = aws_lb_target_group.api_server.arn
  target_id        = "i-0bbd969e200b28173"
  port             = 8080
}
