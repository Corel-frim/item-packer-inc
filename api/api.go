package api

import (
	"net/http"
	"os"

	"github.com/corel-frim/item-packer-inc/api/handlers"
	"github.com/corel-frim/item-packer-inc/internal/storage"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type API struct {
	orders *handlers.Orders
	packs  *handlers.Packs
}

func NewAPI(storage *storage.PackStorage) *API {
	return &API{
		orders: handlers.NewOrders(storage),
		packs:  handlers.NewPacks(storage),
	}
}

func (api *API) Start() {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessEndpoint:  "/live",
		ReadinessEndpoint: "/ready",
	}))

	swaggerPath := "./docs/swagger.json"
	if envPath := os.Getenv("SWAGGER_PATH"); envPath != "" {
		swaggerPath = envPath
	}
	app.Use(swagger.New(swagger.Config{
		Next: func(_ *fiber.Ctx) bool {
			return os.Getenv("APP_ENV") == "production"
		},
		Path:     "/swagger",
		FilePath: swaggerPath,
		Title:    "Item Packer API",
	}))

	// Register API routes before serving static files
	api.RegisterRoutes(app)

	// Serve static files from the frontend directory
	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.Dir("./frontend"),
		Browse:       false,
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))

	log.Fatal(app.Listen(":8080"))
}

func (api *API) RegisterRoutes(app *fiber.App) {
	api.orders.RegisterRoutes(app)
	api.packs.RegisterRoutes(app)
}
