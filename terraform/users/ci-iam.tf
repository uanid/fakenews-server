resource "aws_iam_user" "ci" {
  name = "fnc-ci"
}

data "aws_iam_policy" "ecr_power_user" {
  arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser"
}

resource "aws_iam_user_policy_attachment" "ci_ecr" {
  user       = aws_iam_user.ci.id
  policy_arn = data.aws_iam_policy.ecr_power_user.arn
}

resource "aws_iam_access_key" "ci" {
  user = aws_iam_user.ci.id
}

output "ci_access_key" {
  value = aws_iam_access_key.ci.id
}

output "ci_secret_key" {
  value     = aws_iam_access_key.ci.secret
  sensitive = true
}
