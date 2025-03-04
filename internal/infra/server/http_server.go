package server

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	authhnd "github.com/nathakusuma/elevateu-backend/internal/app/auth/handler"
	authrepo "github.com/nathakusuma/elevateu-backend/internal/app/auth/repository"
	authsvc "github.com/nathakusuma/elevateu-backend/internal/app/auth/service"
	categoryhnd "github.com/nathakusuma/elevateu-backend/internal/app/category/handler"
	categoryrepo "github.com/nathakusuma/elevateu-backend/internal/app/category/repository"
	categorysvc "github.com/nathakusuma/elevateu-backend/internal/app/category/service"
	coursehnd "github.com/nathakusuma/elevateu-backend/internal/app/course/handler"
	courserepo "github.com/nathakusuma/elevateu-backend/internal/app/course/repository"
	coursesvc "github.com/nathakusuma/elevateu-backend/internal/app/course/service"
	paymenthnd "github.com/nathakusuma/elevateu-backend/internal/app/payment/handler"
	paymentrepo "github.com/nathakusuma/elevateu-backend/internal/app/payment/repository"
	paymentsvc "github.com/nathakusuma/elevateu-backend/internal/app/payment/service"
	userhnd "github.com/nathakusuma/elevateu-backend/internal/app/user/handler"
	userrepo "github.com/nathakusuma/elevateu-backend/internal/app/user/repository"
	usersvc "github.com/nathakusuma/elevateu-backend/internal/app/user/service"
	"github.com/nathakusuma/elevateu-backend/internal/infra/cache"
	"github.com/nathakusuma/elevateu-backend/internal/infra/env"
	"github.com/nathakusuma/elevateu-backend/internal/infra/gcp"
	"github.com/nathakusuma/elevateu-backend/internal/middleware"
	"github.com/nathakusuma/elevateu-backend/pkg/bcrypt"
	"github.com/nathakusuma/elevateu-backend/pkg/fileutil"
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
	MountRoutes(db *sqlx.DB, cache cache.ICache)
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

func (s *httpServer) MountRoutes(db *sqlx.DB, cache cache.ICache) {
	gcpClient := gcp.NewStorageClient()
	bcryptInstance := bcrypt.GetBcrypt()
	fileUtil := fileutil.NewFileUtil(gcpClient)
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

	userRepository := userrepo.NewUserRepository(db)
	authRepository := authrepo.NewAuthRepository(db)
	paymentRepository := paymentrepo.NewPaymentRepository(db, cache)
	categoryRepository := categoryrepo.NewCategoryRepository(db)
	courseRepository := courserepo.NewCourseRepository(db)
	courseContentRepository := courserepo.NewCourseContentRepository(db)

	userService := usersvc.NewUserService(userRepository, bcryptInstance, fileUtil, uuidInstance)
	authService := authsvc.NewAuthService(authRepository, userService, bcryptInstance, cache, fileUtil, jwtAccess, mailer,
		randomGenerator, uuidInstance)
	midtransService := paymentsvc.NewMidtransService()
	paymentService := paymentsvc.NewPaymentService(paymentRepository, midtransService, uuidInstance)
	categoryService := categorysvc.NewCategoryService(categoryRepository, uuidInstance)
	courseService := coursesvc.NewCourseService(courseRepository, fileUtil, uuidInstance)
	courseContentService := coursesvc.NewCourseContentService(courseContentRepository, fileUtil, uuidInstance)

	userhnd.InitUserHandler(v1, middlewareInstance, validatorInstance, userService)
	authhnd.InitAuthHandler(v1, middlewareInstance, validatorInstance, authService)
	paymenthnd.InitPaymentHandler(v1, paymentService, midtransService)
	categoryhnd.InitCategoryHandler(v1, categoryService, middlewareInstance, validatorInstance)
	coursehnd.InitCourseHandler(v1, middlewareInstance, validatorInstance, courseService)
	coursehnd.InitCourseContentHandler(v1, middlewareInstance, validatorInstance, courseContentService)
}
