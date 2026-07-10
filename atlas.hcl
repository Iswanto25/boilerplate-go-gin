variable "db_url" {
  type    = string
  default = getenv("DB_URL")
}

variable "dev_url" {
  type    = string
  default = getenv("DEV_DB_URL")
}

env "local" {
  url   = var.db_url
  src   = "file://migrations"
  dev   = var.dev_url

  migration {
    dir = "file://migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }

  schemas = ["public"]

  diff {
    skip {
      add_schema  = true
      drop_schema = true
    }
  }
}

lint {
  destructive {
    error = false
  }
}
