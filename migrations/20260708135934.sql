-- Create "logs" table
CREATE TABLE "public"."logs" (
  "id" bigserial NOT NULL,
  "date" character varying(20) NULL,
  "name" character varying(100) NULL,
  "role" character varying(50) NULL,
  "host" character varying(255) NULL,
  "status" character varying(10) NULL,
  "data" jsonb NULL,
  "user_id" uuid NULL,
  "ip" character varying(45) NULL,
  "ex" character varying(255) NULL,
  "method" character varying(10) NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_logs_created_at" to table: "logs"
CREATE INDEX "idx_logs_created_at" ON "public"."logs" ("created_at");
-- Create index "idx_logs_user_id" to table: "logs"
CREATE INDEX "idx_logs_user_id" ON "public"."logs" ("user_id");
-- Create "modules" table
CREATE TABLE "public"."modules" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(100) NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_modules_deleted_at" to table: "modules"
CREATE INDEX "idx_modules_deleted_at" ON "public"."modules" ("deleted_at");
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(100) NOT NULL,
  "email" character varying(100) NOT NULL,
  "password" character varying(255) NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");
-- Create "resources" table
CREATE TABLE "public"."resources" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(100) NOT NULL,
  "module_id" uuid NOT NULL,
  "available_action" jsonb NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_resources_module" FOREIGN KEY ("module_id") REFERENCES "public"."modules" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_resources_deleted_at" to table: "resources"
CREATE INDEX "idx_resources_deleted_at" ON "public"."resources" ("deleted_at");
-- Create "roles" table
CREATE TABLE "public"."roles" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(100) NOT NULL,
  "status" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_roles_deleted_at" to table: "roles"
CREATE INDEX "idx_roles_deleted_at" ON "public"."roles" ("deleted_at");
-- Create "rolePermissions" table
CREATE TABLE "public"."rolePermissions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "role_id" uuid NOT NULL,
  "resource_id" uuid NOT NULL,
  "granted_actions" jsonb NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_rolePermissions_resource" FOREIGN KEY ("resource_id") REFERENCES "public"."resources" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_rolePermissions_role" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_rolePermissions_deleted_at" to table: "rolePermissions"
CREATE INDEX "idx_rolePermissions_deleted_at" ON "public"."rolePermissions" ("deleted_at");
