-- Modify "logs" table
ALTER TABLE "public"."logs" ADD COLUMN "req_id" character varying(15) NULL;
-- Create index "idx_logs_req_id" to table: "logs"
CREATE INDEX "idx_logs_req_id" ON "public"."logs" ("req_id");
