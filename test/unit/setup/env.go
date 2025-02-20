package setup

import (
	"github.com/nathakusuma/elevateu-backend/internal/infra/env"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
)

func init() {
	// Set up test environment
	env.SetEnv(&env.Env{
		AppEnv: "test",
	})

	log.NewLogger()
}
