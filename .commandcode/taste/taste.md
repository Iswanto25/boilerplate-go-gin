# communication
- Respond in Indonesian, but write code and documentation content in English. Confidence: 0.70

# architecture
- Keep the ExpressJS-inspired feature-based project structure with internal/features/<feature>/{handler,model,repository,service}/ layout. Do not restructure. Confidence: 0.75

# atlas-migrations
- Before running migrate-diff, ensure the target database matches the last applied migration — either by resetting it (dropdb + createdb + reapply migrations) or by verifying no stray columns exist from prior AutoMigrate runs. Confidence: 0.75
- For ADD column: edit model → make run-api (AutoMigrate updates local DB) → make migrate-diff (generates .sql). Do NOT apply the migration to local DB — it's already done by AutoMigrate. The .sql file is saved for other environments (production/staging). Apply to other DBs with: make migrate-up "DB_URL=...". Confidence: 0.80
- For DELETE/EDIT column: AutoMigrate cannot DROP/RENAME/ALTER, so manually alter the database first (ALTER TABLE DROP/ALTER COLUMN), then make migrate-diff, then migrate-up. Confidence: 0.70
