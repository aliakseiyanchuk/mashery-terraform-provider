### Case 1

resource "mashery_service" "minimal_config" {
  name_prefix = "testService"
  description = "some description"
  version = "0.0.1a"
}


# Case 2
resource "mashery_service" "with-cache" {
  name_prefix = "testService"
  description = "some description"
  version = "0.0.1a"
  cache_ttl: 30
}

# Case 3
resource "mashery_service" "with-oauth" {
  name_prefix = "testService"
  description = "some description"
  version = "0.0.1a"
  oauth {
    access_token_ttl_enabled =  true,
    access_token_ttl =          "1h",
  }
}

# Case 4
data "mashery_role" "role" {
  search = {
    name = "role",
  }
}

resource "mashery_service" "with-role" {
  name_prefix = "testService"
  description = "some description"
  version = "0.0.1a"
  iodocs_accessed_by = [
    data.mashery_role.role.read_permission
  ]
}