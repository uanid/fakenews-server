resource "aws_dynamodb_table" "fnc-database" {
  name = "fnc1-db"

  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5

  hash_key = "Key"

  attribute {
    name = "Key"
    type = "S"
  }
}

resource "aws_sqs_queue" "fnc-queue" {
  name = "fnc1-queue.fifo"

  fifo_queue = true

  receive_wait_time_seconds  = 10    #Consumer 수신 대기 시간
  visibility_timeout_seconds = 30    #Consumer 처리 제한 시간
  delay_seconds              = 0     #Queue 메시지 전송 지연??
  message_retention_seconds  = 86400 #24시간
  max_message_size           = 2048  #2KB
}

resource "aws_ecr_repository" "fnc-ecr" {
  name = "${local.project_name}-server"

  image_tag_mutability = "IMMUTABLE"
}
