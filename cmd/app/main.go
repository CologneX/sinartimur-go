package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sinartimur-go/config"
	"sinartimur-go/internal/auth"
	"sinartimur-go/internal/category"
	"sinartimur-go/internal/customer"
	"sinartimur-go/internal/employee"
	"sinartimur-go/internal/finance"
	"sinartimur-go/internal/inventory"
	"sinartimur-go/internal/product"
	"sinartimur-go/internal/purchase"
	purchase_order "sinartimur-go/internal/purchase/purchase-order"
	"sinartimur-go/internal/sales"
	"sinartimur-go/internal/unit"
	"sinartimur-go/internal/user"
	"sinartimur-go/internal/wage"
	"sinartimur-go/middleware"
	"sinartimur-go/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

	// Initialize v1 and middleware
	v1 := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	loggedRouter := handlers.CustomLoggingHandler(os.Stdout, v1, middleware.Logger)

	// Register routes
	SetupRoutes(v1, services)

	log.Fatal(http.ListenAndServe(":8080", middleware.CORSMiddleware()(loggedRouter)))
}

type Services struct {
	AuthService     *auth.AuthService
	UserService     *user.UserService
	EmployeeService *employee.EmployeeService
	//RoleService     *role.RoleService
	WageService          *wage.WageService
	ProductService       *product.ProductService
	CategoryService      *category.CategoryService
	UnitService          *unit.UnitService
	SupplierService      *purchase.SupplierService
	PurchaseOrderService *purchase_order.PurchaseOrderService
	InventoryService     *inventory.StorageService
	CustomerService      *customer.CustomerService
	FinanceService       *finance.FinanceService
	SalesService         *sales.SalesService
}

func BuildServices(db *sql.DB, redis *config.RedisClient) *Services {
	authRepo := auth.NewAuthRepository(db)
	authService := auth.NewAuthService(authRepo, redis)

	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)

	//roleRepo := role.NewRoleRepository(db)
	//roleService := role.NewRoleService(roleRepo)

	employeeRepo := employee.NewEmployeeRepository(db)
	employeeService := employee.NewEmployeeService(employeeRepo)

	wageRepo := wage.NewWageRepository(db)
	wageService := wage.NewWageService(wageRepo)

	productRepo := product.NewProductRepository(db)
	productService := product.NewProductService(productRepo)

	categoryRepo := category.NewCategoryRepository(db)
	categoryService := category.NewCategoryService(categoryRepo)

	unitRepo := unit.NewUnitRepository(db)
	unitService := unit.NewUnitService(unitRepo)

	supplierRepo := purchase.NewSupplierRepository(db)
	supplierService := purchase.NewSupplierService(supplierRepo)

	purchaseOrderRepo := purchase_order.NewPurhaseOrderRepository(db)
	purchaseOrderService := purchase_order.NewPurchaseOrderService(purchaseOrderRepo, db)

	inventoryRepo := inventory.NewStorageRepository(db)
	inventoryService := inventory.NewStorageService(inventoryRepo)

	customerRepo := customer.NewCustomerRepository(db)
	customerService := customer.NewCustomerService(customerRepo)

	financeRepo := finance.NewFinanceTransactionRepository(db)
	financeService := finance.NewFinanceTransactionService(financeRepo)

	salesRepo := sales.NewSalesRepository(db)
	salesService := sales.NewSalesService(salesRepo)

	return &Services{
		AuthService:     authService,
		UserService:     userService,
		EmployeeService: employeeService,
		//RoleService:     roleService,
		WageService:          wageService,
		ProductService:       productService,
		CategoryService:      categoryService,
		UnitService:          unitService,
		SupplierService:      supplierService,
		PurchaseOrderService: purchaseOrderService,
		InventoryService:     inventoryService,
		CustomerService:      customerService,
		FinanceService:       financeService,
		SalesService:         salesService,
	}
}
