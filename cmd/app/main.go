package main

import (
	"database/sql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	v1 "sinartimur-go/api/v1"
	"sinartimur-go/config"
	"sinartimur-go/internal/auth"
	"sinartimur-go/internal/employee"
	"sinartimur-go/internal/role"
	"sinartimur-go/internal/user"
	"sinartimur-go/middleware"
	"sinartimur-go/utils"
)

type Services struct {
	AuthService     *auth.AuthService
	UserService     *user.UserService
	EmployeeService *employee.EmployeeService
	RoleService     *role.RoleService
}

func main() {
	// start a connection to Postgres
	db := config.StartPostgres()
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Failed to close database: %v", err)
		}
	}()
	// start a connection to Redis
	redisClient := config.NewRedisClient()

	// Register Custom Validations
	utils.RegisterCustomValidators()
	// register services
	services := registerServices(db, redisClient)

	// define routes
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	// Add logging middleware from gorilla/mux
	loggedRouter := handlers.CustomLoggingHandler(os.Stdout, router, middleware.Logger)
	v1.RegisterAuthRoutes(router, services.AuthService)

	// Global auth middleware
	protectedRoutes := router.PathPrefix("").Subrouter()
	protectedRoutes.Use(middleware.AuthMiddleware)

	/// Role-based auth middleware
	// HR middleware setup
	HRRoutes := router.PathPrefix("/hr").Subrouter()
	HRRoutes.Use(middleware.RoleMiddleware("hr"))
	v1.RegisterEmployeeRoutes(HRRoutes, services.EmployeeService)

	// Admin middleware setup
	AdminRoutes := router.PathPrefix("/admin").Subrouter()
	AdminRoutes.Use(middleware.RoleMiddleware())
	v1.RegisterUserRoutes(AdminRoutes, services.UserService)
	v1.RegisterRoleRoutes(AdminRoutes, services.RoleService)

	// Inventory middleware setup
	InventoryRoutes := router.PathPrefix("/inventory").Subrouter()
	InventoryRoutes.Use(middleware.RoleMiddleware("inventory"))

	// Finance middleware setup
	FinanceRoutes := router.PathPrefix("/finance").Subrouter()
	FinanceRoutes.Use(middleware.RoleMiddleware("finance"))

	// Sales middleware setup
	SalesRoutes := router.PathPrefix("/sales").Subrouter()
	SalesRoutes.Use(middleware.RoleMiddleware("sales"))

	// Purchase middleware setup
	PurchaseRoutes := router.PathPrefix("/purchase").Subrouter()
	PurchaseRoutes.Use(middleware.RoleMiddleware("purchase"))

	// serve the router on port 8080
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(loggedRouter)))
}

func registerServices(
	db *sql.DB,
	redis *config.RedisClient,
) *Services {
	authRepo := auth.NewAuthRepository(db)
	authService := auth.NewAuthService(authRepo, redis)
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	employeeRepo := employee.NewEmployeeRepository(db)
	employeeService := employee.NewEmployeeService(employeeRepo)
	roleRepo := role.NewRoleRepository(db)
	roleService := role.NewRoleService(roleRepo)

	return &Services{
		AuthService:     authService,
		UserService:     userService,
		EmployeeService: employeeService,
		RoleService:     roleService,
	}
}
