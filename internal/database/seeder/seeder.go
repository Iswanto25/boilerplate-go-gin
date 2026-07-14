package seeder

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Seeder interface {
	Name() string
	Enabled() bool
	Seed(db *gorm.DB) (int, error)
}

type Runner struct {
	db      *gorm.DB
	seeders []Seeder
}

func NewRunner(db *gorm.DB) *Runner {
	return &Runner{db: db}
}

func (r *Runner) Register(s ...Seeder) {
	r.seeders = append(r.seeders, s...)
}

func (r *Runner) RunAll() {
	total := len(r.seeders)
	if total == 0 {
		fmt.Println("🌱 [SEEDER] No seeders registered.")
		return
	}

	fmt.Println()
	fmt.Println("🌱 [SEEDER] Starting seeders...")
	fmt.Println("────────────────────────────────────────")

	start := time.Now()
	success, skipped, failed := 0, 0, 0

	for i, s := range r.seeders {
		prefix := fmt.Sprintf("[%d/%d]", i+1, total)

		if !s.Enabled() {
			fmt.Printf("⏭️  %s %s | SKIPPED (disabled)\n", prefix, s.Name())
			skipped++
			continue
		}

		fmt.Printf("🌱 %s %s | STARTED\n", prefix, s.Name())
		seedStart := time.Now()

		rows, err := s.Seed(r.db)
		elapsed := time.Since(seedStart).Milliseconds()

		if err != nil {
			fmt.Printf("❌ %s %s | ERROR: %v | %dms\n", prefix, s.Name(), err, elapsed)
			failed++
		} else {
			fmt.Printf("✅ %s %s | SUCCESS — %d rows upserted | %dms\n", prefix, s.Name(), rows, elapsed)
			success++
		}
	}

	totalMs := time.Since(start).Milliseconds()
	fmt.Println("────────────────────────────────────────")
	fmt.Printf("🏁 [SEEDER] Completed: %d success, %d skipped, %d failed | %dms\n\n",
		success, skipped, failed, totalMs)
}
