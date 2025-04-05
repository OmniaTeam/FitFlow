
locals {
  db_host          = getenv("PG_HOST")
  db_port          = getenv("PG_PORT")
  db_user          = getenv("PG_USER")
  db_password      = getenv("PG_PASSWORD")
  db_name          = getenv("PG_DATABASE")
  db_name_for_test = "atlas"
}

env "server" {
  url = "postgres://${local.db_user}:${local.db_password}@${local.db_host}:${local.db_port}/${local.db_name}?search_path=public&sslmode=disable"
  dev = "postgres://${local.db_user}:${local.db_password}@${local.db_host}/atlas?search_path=public&sslmode=disable"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
