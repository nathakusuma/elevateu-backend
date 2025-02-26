package server

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/nathakusuma/elevateu-backend/internal/infra/gcp"
	"github.com/nathakusuma/elevateu-backend/pkg/fileutil"

	"github.com/redis/go-redis/v9"

	authhnd "github.com/nathakusuma/elevateu-backend/internal/app/auth/handler"
	authrepo "github.com/nathakusuma/elevateu-backend/internal/app/auth/repository"
	authsvc "github.com/nathakusuma/elevateu-backend/internal/app/auth/service"
	paymenthnd "github.com/nathakusuma/elevateu-backend/internal/app/payment/handler"
	paymentrepo "github.com/nathakusuma/elevateu-backend/internal/app/payment/repository"
	paymentsvc "github.com/nathakusuma/elevateu-backend/internal/app/payment/service"
	storagerepo "github.com/nathakusuma/elevateu-backend/internal/app/storage/repository"
	userhnd "github.com/nathakusuma/elevateu-backend/internal/app/user/handler"
	userrepo "github.com/nathakusuma/elevateu-backend/internal/app/user/repository"
	usersvc "github.com/nathakusuma/elevateu-backend/internal/app/user/service"
	"github.com/nathakusuma/elevateu-backend/internal/infra/env"
	"github.com/nathakusuma/elevateu-backend/internal/middleware"
	"github.com/nathakusuma/elevateu-backend/pkg/bcrypt"
	"github.com/nathakusuma/elevateu-backend/pkg/jwt"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
	"github.com/nathakusuma/elevateu-backend/pkg/mail"
	"github.com/nathakusuma/elevateu-backend/pkg/randgen"
	"github.com/nathakusuma/elevateu-backend/pkg/uuidpkg"
	"github.com/nathakusuma/elevateu-backend/pkg/validator"
)

type HTTPServer interface {
	Start(part string)
	MountMiddlewares()
	MountRoutes(db *sqlx.DB, rds *redis.Client)
	GetApp() *fiber.App
}

type httpServer struct {
	app *fiber.App
}

func NewHTTPServer() HTTPServer {
	config := fiber.Config{
		AppName:      "ElevateU",
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
		ErrorHandler: errorHandler,
	}

	app := fiber.New(config)

	return &httpServer{
		app: app,
	}
}

func (s *httpServer) GetApp() *fiber.App {
	return s.app
}

func (s *httpServer) Start(port string) {
	if port[0] != ':' {
		port = ":" + port
	}

	err := s.app.Listen(port)

	if err != nil {
		log.Fatal(map[string]interface{}{
			"error": err.Error(),
		}, "[SERVER][Start] failed to start server")
	}
}

func (s *httpServer) MountMiddlewares() {
	s.app.Use(middleware.LoggerConfig())
	s.app.Use(middleware.Helmet())
	s.app.Use(middleware.Compress())
	s.app.Use(middleware.Cors())
	s.app.Use(middleware.RecoverConfig())
}

func (s *httpServer) MountRoutes(db *sqlx.DB, rds *redis.Client) {
	gcpClient := gcp.NewStorageClient()
	bcryptInstance := bcrypt.GetBcrypt()
	fileUtil := fileutil.NewFileUtil()
	jwtAccess := jwt.NewJwt(env.GetEnv().JwtAccessExpireDuration, env.GetEnv().JwtAccessSecretKey)
	mailer := mail.NewMailDialer()
	randomGenerator := randgen.GetRandGen()
	uuidInstance := uuidpkg.GetUUID()
	validatorInstance := validator.NewValidator()
	middlewareInstance := middleware.NewMiddleware(jwtAccess)

	s.app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString("ElevateU Healthy")
	})

	api := s.app.Group("/api")
	v1 := api.Group("/v1")

	storageRepository := storagerepo.NewStorageRepository(gcpClient)
	userRepository := userrepo.NewUserRepository(db)
	authRepository := authrepo.NewAuthRepository(db, rds)

	userService := usersvc.NewUserService(userRepository, storageRepository, bcryptInstance, fileUtil, uuidInstance)
	authService := authsvc.NewAuthService(authRepository, userService, bcryptInstance, jwtAccess, mailer,
		randomGenerator, uuidInstance)

	paymentRepository := paymentrepo.NewPaymentRepository(db, rds)
	midtransService := paymentsvc.NewMidtransService()
	paymentService := paymentsvc.NewPaymentService(paymentRepository, midtransService, uuidInstance)

	userhnd.InitUserHandler(v1, middlewareInstance, validatorInstance, userService)
	authhnd.InitAuthHandler(v1, middlewareInstance, validatorInstance, authService)
	paymenthnd.InitPaymentHandler(v1, paymentService, midtransService)
}
