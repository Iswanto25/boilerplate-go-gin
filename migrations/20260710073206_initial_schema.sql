-- Create "logs" table
CREATE TABLE "public"."logs" (
  "id" bigserial NOT NULL,
  "date" character varying(20) NULL,
  "name" character varying(100) NULL,
  "role" character varying(50) NULL,
  "host" character varying(255) NULL,
  "status" character varying(10) NULL,
  "data" jsonb NULL,
  "userID" uuid NULL,
  "ip" character varying(45) NULL,
  "method" character varying(10) NULL,
  "reqID" character varying(50) NULL,
  "createdAt" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_logs_created_at" to table: "logs"
CREATE INDEX "idx_logs_created_at" ON "public"."logs" ("createdAt");
-- Create index "idx_logs_date" to table: "logs"
CREATE INDEX "idx_logs_date" ON "public"."logs" ("date");
-- Create index "idx_logs_method" to table: "logs"
CREATE INDEX "idx_logs_method" ON "public"."logs" ("method");
-- Create index "idx_logs_req_id" to table: "logs"
CREATE INDEX "idx_logs_req_id" ON "public"."logs" ("reqID");
-- Create index "idx_logs_status" to table: "logs"
CREATE INDEX "idx_logs_status" ON "public"."logs" ("status");
-- Create index "idx_logs_user_id" to table: "logs"
CREATE INDEX "idx_logs_user_id" ON "public"."logs" ("userID");
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(100) NOT NULL,
  "email" character varying(100) NOT NULL,
  "password" character varying(255) NOT NULL,
  "role" character varying(20) NOT NULL DEFAULT 'user',
  "createdAt" timestamptz NULL,
  "updatedAt" timestamptz NULL,
  "deletedAt" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deletedAt");
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");
-- Create "profiles" table
CREATE TABLE "public"."profiles" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "userID" uuid NOT NULL,
  "phone" character varying(15) NULL,
  "address" character varying(255) NULL,
  "photo" character varying(255) NULL,
  "createdAt" timestamptz NULL,
  "updatedAt" timestamptz NULL,
  "deletedAt" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_profile" FOREIGN KEY ("userID") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_profiles_deleted_at" to table: "profiles"
CREATE INDEX "idx_profiles_deleted_at" ON "public"."profiles" ("deletedAt");
-- Create index "idx_profiles_user_id" to table: "profiles"
CREATE INDEX "idx_profiles_user_id" ON "public"."profiles" ("userID");
-- Create "modules" table
CREATE TABLE "public"."modules" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(100) NOT NULL,
  "createdAt" timestamptz NULL,
  "updatedAt" timestamptz NULL,
  "deletedAt" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_modules_deleted_at" to table: "modules"
CREATE INDEX "idx_modules_deleted_at" ON "public"."modules" ("deletedAt");
-- Create "resources" table
CREATE TABLE "public"."resources" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(100) NOT NULL,
  "moduleId" uuid NOT NULL,
  "availableAction" jsonb NULL,
  "createdAt" timestamptz NULL,
  "updatedAt" timestamptz NULL,
  "deletedAt" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_resources_module" FOREIGN KEY ("moduleId") REFERENCES "public"."modules" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_resources_deleted_at" to table: "resources"
CREATE INDEX "idx_resources_deleted_at" ON "public"."resources" ("deletedAt");
-- Create index "idx_resources_module_id" to table: "resources"
CREATE INDEX "idx_resources_module_id" ON "public"."resources" ("moduleId");
-- Create "roles" table
CREATE TABLE "public"."roles" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(100) NOT NULL,
  "status" boolean NULL DEFAULT true,
  "createdAt" timestamptz NULL,
  "updatedAt" timestamptz NULL,
  "deletedAt" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_roles_deleted_at" to table: "roles"
CREATE INDEX "idx_roles_deleted_at" ON "public"."roles" ("deletedAt");
-- Create "rolePermissions" table
CREATE TABLE "public"."rolePermissions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "roleId" uuid NOT NULL,
  "resourceId" uuid NOT NULL,
  "grantedActions" jsonb NULL,
  "createdAt" timestamptz NULL,
  "updatedAt" timestamptz NULL,
  "deletedAt" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_rolePermissions_resource" FOREIGN KEY ("resourceId") REFERENCES "public"."resources" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_rolePermissions_role" FOREIGN KEY ("roleId") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_rolePermissions_deleted_at" to table: "rolePermissions"
CREATE INDEX "idx_rolePermissions_deleted_at" ON "public"."rolePermissions" ("deletedAt");
-- Create index "idx_rolePermissions_resource_id" to table: "rolePermissions"
CREATE INDEX "idx_rolePermissions_resource_id" ON "public"."rolePermissions" ("resourceId");
-- Create index "idx_rolePermissions_role_id" to table: "rolePermissions"
CREATE INDEX "idx_rolePermissions_role_id" ON "public"."rolePermissions" ("roleId");
