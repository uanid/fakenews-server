resource "aws_acm_certificate" "cert" {
  domain_name       = "*.fnc-1.link"
  validation_method = "DNS"
}

resource "aws_route53_record" "example" {
  for_each = {
  for dvo in aws_acm_certificate.cert.domain_validation_options : dvo.domain_name => {
    name   = dvo.resource_record_name
    record = dvo.resource_record_value
    type   = dvo.resource_record_type
  }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = "Z01518507GPOW8LJO1VS"
}

resource "aws_route53_record" "arecord" {
  name    = "api.fnc-1.link"
  type    = "A"
  zone_id = "Z01518507GPOW8LJO1VS"

  alias {
    evaluate_target_health = true
    name                   = aws_lb.elb.dns_name
    zone_id                = aws_lb.elb.zone_id
  }
}
