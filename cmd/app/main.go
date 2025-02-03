package main

import (
	"database/sql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"sinartimur-go/config"
	"sinartimur-go/internal/auth"
	"sinartimur-go/internal/employee"
	"sinartimur-go/internal/role"
	"sinartimur-go/internal/user"
	"sinartimur-go/middleware"
	"sinartimur-go/utils"
)

func main() {
	// Start database and Redis connections
	db := config.StartPostgres()
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Failed to close database: %v", err)
		}
	}()

	redisClient := config.NewRedisClient()

	// Register custom validations
	utils.RegisterCustomValidators()

	// Build services
	services := BuildServices(db, redisClient)

	// Initialize routerv1 and middleware
	routerv1 := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	loggedRouter := handlers.CustomLoggingHandler(os.Stdout, routerv1, middleware.Logger)

	// Register routes
	SetupRoutes(routerv1, services)

	// Add CORS middleware
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(loggedRouter)))
}

type Services struct {
	AuthService     *auth.AuthService
	UserService     *user.UserService
	EmployeeService *employee.EmployeeService
	RoleService     *role.RoleService
}

func BuildServices(db *sql.DB, redis *config.RedisClient) *Services {
	authRepo := auth.NewAuthRepository(db)
	authService := auth.NewAuthService(authRepo, redis)

	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)

	roleRepo := role.NewRoleRepository(db)
	roleService := role.NewRoleService(roleRepo)

	employeeRepo := employee.NewEmployeeRepository(db)
	employeeService := employee.NewEmployeeService(employeeRepo)

	return &Services{
		AuthService:     authService,
		UserService:     userService,
		EmployeeService: employeeService, // You'll need to add this service
		RoleService:     roleService,
	}
}
