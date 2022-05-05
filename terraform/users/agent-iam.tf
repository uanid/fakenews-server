resource "aws_iam_user" "agent" {
  name = "fnc-agent"
}

data "aws_iam_policy_document" "agent_ddb" {
  statement {
    effect  = "Allow"
    actions = [
      "dynamodb:BatchGetItem",
      "dynamodb:BatchWriteItem",
      "dynamodb:ConditionCheckItem",
      "dynamodb:PutItem",
      "dynamodb:DescribeTable",
      "dynamodb:DeleteItem",
      "dynamodb:GetItem",
      "dynamodb:Scan",
      "dynamodb:Query",
      "dynamodb:UpdateItem",
    ]
    resources = ["*"]
  }
}

data "aws_iam_policy_document" "agent_sqs" {
  statement {
    effect  = "Allow"
    actions = [
      "sqs:SendMessage",
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage",
    ]
    resources = ["*"]
  }
}

resource "aws_iam_user_policy" "agent_ddb" {
  user        = aws_iam_user.agent.id
  name_prefix = "fnc-agent-ddb"

  policy = data.aws_iam_policy_document.agent_ddb.json
}

resource "aws_iam_user_policy" "agent_sqs" {
  user        = aws_iam_user.agent.id
  name_prefix = "fnc-agent-sqs"

  policy = data.aws_iam_policy_document.agent_sqs.json
}

resource "aws_iam_access_key" "agent" {
  user = aws_iam_user.agent.id
}

output "agent_access_key" {
  value = aws_iam_access_key.agent.id
}

output "agent_secret_key" {
  value     = aws_iam_access_key.agent.secret
  sensitive = true
}
