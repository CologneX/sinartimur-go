package main

import (
	"github.com/gorilla/mux"
	v1 "sinartimur-go/api/v1"
	"sinartimur-go/internal/auth"
	"sinartimur-go/internal/category"
	"sinartimur-go/internal/employee"
	"sinartimur-go/internal/product"
	"sinartimur-go/internal/purchase"
	"sinartimur-go/internal/unit"
	"sinartimur-go/internal/user"
	"sinartimur-go/internal/wage"
	"sinartimur-go/middleware"
)

func RegisterAuthRoutes(router *mux.Router, userService *auth.AuthService) {
	router.HandleFunc("/login", v1.LoginHandler(userService)).Methods("GET")
	router.HandleFunc("/refresh", v1.RefreshTokenHandler(userService)).Methods("GET")
}

func RegisterUserRoutes(router *mux.Router, userService *user.UserService) {
	router.HandleFunc("/user", v1.CreateUserHandler(userService)).Methods("POST")
	router.HandleFunc("/users", v1.GetAllUsersHandler(userService)).Methods("GET")
	router.HandleFunc("/user/{id}", v1.UpdateUserHandler(userService)).Methods("PUT")
	router.HandleFunc("/user-credential/{id}", v1.UpdateUserCredentialHandler(userService)).Methods("PUT")
}

//func RegisterRoleRoutes(router *mux.Router, roleService *role.RoleService) {
//	router.HandleFunc("/role", v1.CreateRoleHandler(roleService)).Methods("POST")
//	router.HandleFunc("/role/{id}", v1.UpdateRoleHandler(roleService)).Methods("PUT")
//	router.HandleFunc("/roles", v1.GetAllRolesHandler(roleService)).Methods("GET")
//	router.HandleFunc("/role/assign", v1.AssignRoleToUserHandler(roleService)).Methods("POST")
//	router.HandleFunc("/role/unassign", v1.UnassignRoleFromUserHandler(roleService)).Methods("POST")
//}

func RegisterPurchaseOrderRoutes(router *mux.Router, purchaseOrderService *purchase.PurchaseOrderService, purchaseDetailService *purchase.PurchaseOrderDetailService) {
	router.HandleFunc("/order", v1.CreatePurchaseOrderHandler(purchaseOrderService)).Methods("POST")
	router.HandleFunc("/orders", v1.GetAllPurchaseOrdersHandler(purchaseOrderService)).Methods("GET")
	//router.HandleFunc("/order/{id}", v1.GetPurchaseOrderByIDHandler(purchaseOrderService)).Methods("GET")
	router.HandleFunc("/order/{id}", v1.UpdatePurchaseOrderHandler(purchaseOrderService)).Methods("PUT")
	router.HandleFunc("/order/{id}/cancel", v1.CancelPurchaseOrderHandler(purchaseOrderService)).Methods("POST")
	router.HandleFunc("/order/{id}/receive", v1.ReceivePurchaseOrderHandler(purchaseOrderService)).Methods("POST")

	// Purchase order detail routes
	router.HandleFunc("/order/{id}/detail", v1.GetPurchaseOrderDetailHandler(purchaseOrderService)).Methods("GET")
	router.HandleFunc("/order/{order_id}/item", v1.AddPurchaseOrderDetailHandler(purchaseDetailService)).Methods("POST")
	router.HandleFunc("/order/item/{id}", v1.UpdatePurchaseOrderDetailHandler(purchaseDetailService)).Methods("PUT")
	router.HandleFunc("/order/item/{id}", v1.DeletePurchaseOrderDetailHandler(purchaseDetailService)).Methods("DELETE")
}

func RegisterSupplierRoutes(router *mux.Router, supplierService *purchase.SupplierService) {
	router.HandleFunc("/supplier", v1.CreateSupplierHandler(supplierService)).Methods("POST")
	router.HandleFunc("/supplier/{id}", v1.UpdateSupplierHandler(supplierService)).Methods("PUT")
	router.HandleFunc("/supplier/{id}", v1.DeleteSupplierHandler(supplierService)).Methods("DELETE")
	//router.HandleFunc("/supplier/{id}", v1.GetSupplierByIDHandler(supplierService)).Methods("GET")
	router.HandleFunc("/suppliers", v1.GetAllSuppliersHandler(supplierService)).Methods("GET")
}

func RegisterEmployeeRoutes(router *mux.Router, employeeService *employee.EmployeeService) {
	router.HandleFunc("/employee", v1.CreateEmployeeHandler(employeeService)).Methods("POST")
	router.HandleFunc("/employee/{id}", v1.UpdateEmployeeHandler(employeeService)).Methods("PUT")
	router.HandleFunc("/employee/{id}", v1.DeleteEmployeeHandler(employeeService)).Methods("DELETE")
	router.HandleFunc("/employees", v1.GetAllEmployeesHandler(employeeService)).Methods("GET")
}

func RegisterWageRoutes(router *mux.Router, wageService *wage.WageService) {
	router.HandleFunc("/wage", v1.CreateWageHandler(wageService)).Methods("POST")
	router.HandleFunc("/wage/{id}", v1.UpdateWageHandler(wageService)).Methods("PUT")
	router.HandleFunc("/wage/{id}", v1.DeleteWageHandler(wageService)).Methods("DELETE")
	router.HandleFunc("/wage/{id}", v1.GetWageDetailHandler(wageService)).Methods("GET")
	router.HandleFunc("/wages", v1.GetAllWagesHandler(wageService)).Methods("GET")
}

func RegisterProductRoutes(router *mux.Router, productService *product.ProductService) {
	router.HandleFunc("/product", v1.CreateProductHandler(productService)).Methods("POST")
	router.HandleFunc("/product/{id}", v1.UpdateProductHandler(productService)).Methods("PUT")
	router.HandleFunc("/product/{id}", v1.DeleteProductHandler(productService)).Methods("DELETE")
	router.HandleFunc("/products", v1.GetAllProductHandler(productService)).Methods("GET")
	router.HandleFunc("/product/batch/{id}", v1.GetProductBatchHandler(productService)).Methods("GET")
}

func RegisterUnitRoutes(router *mux.Router, unitService *unit.UnitService) {
	router.HandleFunc("/unit", v1.CreateUnitHandler(unitService)).Methods("POST")
	router.HandleFunc("/unit/{id}", v1.UpdateUnitHandler(unitService)).Methods("PUT")
	router.HandleFunc("/unit/{id}", v1.DeleteUnitHandler(unitService)).Methods("DELETE")
	router.HandleFunc("/units", v1.GetAllUnitHandler(unitService)).Methods("GET")
}

func RegisterCategoryRoutes(router *mux.Router, categoryService *category.CategoryService) {
	router.HandleFunc("/category", v1.CreateCategoryHandler(categoryService)).Methods("POST")
	router.HandleFunc("/category/{id}", v1.UpdateCategoryHandler(categoryService)).Methods("PUT")
	router.HandleFunc("/category/{id}", v1.DeleteCategoryHandler(categoryService)).Methods("DELETE")
	router.HandleFunc("/categories", v1.GetAllCategoryHandler(categoryService)).Methods("GET")
}

// SetupRoutes registers all API routes
func SetupRoutes(router *mux.Router, services *Services) {
	// Auth Routes
	authRouter := router.PathPrefix("/auth").Subrouter()
	RegisterAuthRoutes(authRouter, services.AuthService)

	// Global auth middleware
	protectedRoutes := router.PathPrefix("").Subrouter()
	protectedRoutes.Use(middleware.AuthMiddleware)

	/// Role-based auth middleware
	// HR middleware setup
	HRRoutes := router.PathPrefix("/hr").Subrouter()
	HRRoutes.Use(middleware.RoleMiddleware("hr"))
	RegisterEmployeeRoutes(HRRoutes, services.EmployeeService)
	RegisterWageRoutes(HRRoutes, services.WageService)

	// Admin middleware setup
	AdminRoutes := router.PathPrefix("/admin").Subrouter()
	AdminRoutes.Use(middleware.RoleMiddleware())
	RegisterUserRoutes(AdminRoutes, services.UserService)
	//RegisterRoleRoutes(AdminRoutes, services.RoleService)

	// Inventory middleware setup
	InventoryRoutes := router.PathPrefix("/inventory").Subrouter()
	InventoryRoutes.Use(middleware.RoleMiddleware("inventory"))
	RegisterProductRoutes(InventoryRoutes, services.ProductService)
	RegisterCategoryRoutes(InventoryRoutes, services.CategoryService)
	RegisterUnitRoutes(InventoryRoutes, services.UnitService)

	// Finance middleware setup
	FinanceRoutes := router.PathPrefix("/finance").Subrouter()
	FinanceRoutes.Use(middleware.RoleMiddleware("finance"))

	// Sales middleware setup
	SalesRoutes := router.PathPrefix("/sales").Subrouter()
	SalesRoutes.Use(middleware.RoleMiddleware("sales"))

	// Purchase middleware setup
	PurchaseRoutes := router.PathPrefix("/purchase").Subrouter()
	PurchaseRoutes.Use(middleware.RoleMiddleware("purchase"))
	RegisterSupplierRoutes(PurchaseRoutes, services.SupplierService)
	RegisterPurchaseOrderRoutes(PurchaseRoutes, services.PurchaseOrderService, services.PurchaseOrderDetailService)

	//// Purchase middleware setup
	//PurchaseRoutes := router.PathPrefix("/purchase").Subrouter()
	//PurchaseRoutes.Use(middleware.RoleMiddleware("purchase"))

}
