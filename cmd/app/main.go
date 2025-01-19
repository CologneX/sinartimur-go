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
	// register services
	services := registerServices(db, redisClient)

	// define routes
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	// Add logging middleware from gorilla/mux
	loggedRouter := handlers.CustomLoggingHandler(os.Stdout, router, middleware.Logger)
	v1.RegisterAuthRoutes(router, services.AuthService)
	// Add auth middleware
	protectedRoutes := router.PathPrefix("").Subrouter()
	protectedRoutes.Use(middleware.AuthMiddleware)
	v1.RegisterUserRoutes(protectedRoutes, services.UserService)
	v1.RegisterEmployeeRoutes(protectedRoutes, services.EmployeeService)
	v1.RegisterRoleRoutes(protectedRoutes, services.RoleService)

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
