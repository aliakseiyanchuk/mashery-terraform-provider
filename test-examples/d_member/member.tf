data "mashery_member" "root_member" {
  search = {
    "username" = var.username
  }
}

output "member_id" {
  value = data.mashery_member.root_member.id
}