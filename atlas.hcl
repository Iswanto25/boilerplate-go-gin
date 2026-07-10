data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./tools/loader",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/15/dev"
  migration {
    dir = "file://migrations"
  }
}
