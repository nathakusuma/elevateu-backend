package midtrans

import (
	"github.com/midtrans/midtrans-go"
	"github.com/nathakusuma/elevateu-backend/internal/infra/env"
)

func SetupMidtrans() {
	midtrans.ServerKey = env.GetEnv().MidtransServerKey
	midtrans.Environment = env.GetEnv().MidtransEnvironment
}
