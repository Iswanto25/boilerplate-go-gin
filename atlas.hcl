variable "db_url" {
  type    = string
  default = getenv("DB_URL", "postgres://postgres:postgres@localhost:5432/boilerplate?sslmode=disable")
}

env "local" {
  url   = var.db_url
  src   = "file://migrations"
  dev   = getenv("DEV_DB_URL", "docker://postgres/15/dev?search_path=public")

  migration {
    dir = "file://migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

lint {
  latest {
    destructive {
      error = false
    }
  }
}
