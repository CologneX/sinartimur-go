package main

import (
	v1 "sinartimur-go/api/v1"
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

	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(router *mux.Router, userService *auth.AuthService) {
	router.HandleFunc("/login", v1.LoginHandler(userService)).Methods("POST")
	router.HandleFunc("/refresh", v1.RefreshTokenHandler(userService)).Methods("GET")
}

func RegisterUserRoutes(router *mux.Router, userService *user.UserService) {
	router.HandleFunc("/user", v1.CreateUserHandler(userService)).Methods("POST")
	router.HandleFunc("/users", v1.GetAllUsersHandler(userService)).Methods("GET")
	router.HandleFunc("/user/{id}", v1.UpdateUserHandler(userService)).Methods("PUT")
	router.HandleFunc("/user-credential/{id}", v1.UpdateUserCredentialHandler(userService)).Methods("PUT")
}

func RegisterPurchaseOrderRoutes(router *mux.Router, purchaseOrderService *purchase_order.PurchaseOrderService, productService *product.ProductService, storageService *inventory.StorageService) {
	// Purchase Orders
	router.HandleFunc("/orders", v1.GetAllPurchaseOrderHandler(purchaseOrderService)).Methods("GET")
	router.HandleFunc("/order", v1.CreatePurchaseOrderHandler(purchaseOrderService)).Methods("POST")
	router.HandleFunc("/order/{id}", v1.UpdatePurchaseOrderHandler(purchaseOrderService)).Methods("PUT")
	router.HandleFunc("/order/{id}", v1.GetPurchaseOrderDetailHandler(purchaseOrderService)).Methods("GET")
	// router.HandleFunc("/order/{id}/receive", v1.ReceivePurchaseOrderHandler(purchaseOrderService)).Methods("GET")
	router.HandleFunc("/order/{id}/cancel", v1.CancelPurchaseOrderHandler(purchaseOrderService)).Methods("PUT")
	router.HandleFunc("/order/{id}/check", v1.CheckPurchaseOrderHandler(purchaseOrderService)).Methods("PUT")
	router.HandleFunc("/order/returns", v1.GetAllPurchaseOrderReturnHandler(purchaseOrderService)).Methods("GET")
	router.HandleFunc("/order/return", v1.CreatePurchaseOrderReturnHandler(purchaseOrderService)).Methods("POST")
	router.HandleFunc("/order/return/{id}/cancel", v1.CancelPurchaseOrderReturnHandler(purchaseOrderService)).Methods("PUT")
	// Add route for completing full purchase order
	// router.HandleFunc("/api/v1/purchase-orders/{id}/complete-full", middleware.AuthHandler(middleware.RoleHandler("purchase", v1.CompleteFullPurchaseOrderHandler(purchaseOrderService)))).Methods("POST")
	router.HandleFunc("/order/{id}/complete", v1.CompleteFullPurchaseOrderHandler(purchaseOrderService)).Methods("POST")

	// Purchase Order Items
	router.HandleFunc("/order/items/{id}", v1.DeletePurchaseOrderItemHandler(purchaseOrderService)).Methods("DELETE")
	router.HandleFunc("/order/items/{id}", v1.UpdatePurchaseOrderItemHandler(purchaseOrderService)).Methods("PUT")
	router.HandleFunc("/order/{id}/item", v1.CreatePurchaseOrderItemHandler(purchaseOrderService)).Methods("POST")

	// Product
	router.HandleFunc("/products", v1.GetAllProductHandler(productService)).Methods("GET")
	router.HandleFunc("/storages", v1.GetAllStoragesHandler(storageService)).Methods("GET")
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

func RegisterInventoryRoutes(router *mux.Router, storageService *inventory.StorageService) {
	router.HandleFunc("/storage", v1.CreateStorageHandler(storageService)).Methods("POST")
	router.HandleFunc("/storage/{id}", v1.UpdateStorageHandler(storageService)).Methods("PUT")
	router.HandleFunc("/storage/{id}", v1.DeleteStorageHandler(storageService)).Methods("DELETE")
	router.HandleFunc("/storages", v1.GetAllStoragesHandler(storageService)).Methods("GET")
	router.HandleFunc("/batches", v1.GetAllBatchHandler(storageService)).Methods("GET")
	router.HandleFunc("/move-batch", v1.MoveBatchHandler(storageService)).Methods("POST")
	router.HandleFunc("/logs", v1.GetAllInventoryLogHandler(storageService)).Methods("GET")
	router.HandleFunc("/logs/refresh", v1.RefreshInventoryLogViewHandler(storageService)).Methods("POST")
}

func RegisterCustomerRoutes(router *mux.Router, customerService *customer.CustomerService) {
	router.HandleFunc("/customer", v1.CreateCustomerHandler(customerService)).Methods("POST")
	router.HandleFunc("/customer/{id}", v1.UpdateCustomerHandler(customerService)).Methods("PUT")
	router.HandleFunc("/customer/{id}", v1.DeleteCustomerHandler(customerService)).Methods("DELETE")
	router.HandleFunc("/customers", v1.GetAllCustomersHandler(customerService)).Methods("GET")
}
func RegisterFinanceTransactionRoutes(router *mux.Router, service *finance.FinanceService) {
	router.HandleFunc("/transaction", v1.CreateFinanceTransactionHandler(service)).Methods("POST")
	router.HandleFunc("/transactions", v1.GetAllFinanceTransactionsHandler(service)).Methods("GET")
	router.HandleFunc("/transaction/cancel/{id}", v1.CancelFinanceTransactionHandler(service)).Methods("POST")
	router.HandleFunc("/transactions/summary", v1.GetFinanceTransactionSummaryHandler(service)).Methods("GET")
	router.HandleFunc("/transactions/refresh", v1.RefreshFinanceTransactionViewHandler(service)).Methods("POST")
}

func RegisterSalesRoutes(router *mux.Router, salesService *sales.SalesService) {
	// Sales Order endpoints
	router.HandleFunc("/orders", v1.GetSalesOrdersHandler(salesService)).Methods("GET")
	router.HandleFunc("/order", v1.CreateSalesOrderHandler(salesService)).Methods("POST")
	router.HandleFunc("/order/{id}", v1.UpdateSalesOrderHandler(salesService)).Methods("PUT")
	router.HandleFunc("/order/{id}", v1.GetSalesOrderDetailsHandler(salesService)).Methods("GET")
	router.HandleFunc("/order/{id}/cancel", v1.CancelSalesOrderHandler(salesService)).Methods("POST")

	// Sales Order Item endpoints
	router.HandleFunc("/order/item/{id}", v1.AddItemToSalesOrderHandler(salesService)).Methods("POST")
	router.HandleFunc("/order/item/{id}", v1.UpdateSalesOrderItemHandler(salesService)).Methods("PUT")
	router.HandleFunc("/order/{order_id}/item/{detail_id}", v1.DeleteSalesOrderItemHandler(salesService)).Methods("DELETE")

	// Sales Invoice endpoints
	router.HandleFunc("/invoices", v1.GetSalesInvoicesHandler(salesService)).Methods("GET")
	router.HandleFunc("/invoice", v1.CreateSalesInvoiceHandler(salesService)).Methods("POST")
	router.HandleFunc("/invoice/cancel", v1.CancelSalesInvoiceHandler(salesService)).Methods("POST")

	// Sales Invoice Return endpoints
	router.HandleFunc("/invoice/return", v1.ReturnInvoiceItemsHandler(salesService)).Methods("POST")
	router.HandleFunc("/invoice/return/{return_id}/cancel", v1.CancelInvoiceReturnHandler(salesService)).Methods("POST")

	// Delivery Note endpoints
	router.HandleFunc("/delivery-note", v1.CreateDeliveryNoteHandler(salesService)).Methods("POST")
	router.HandleFunc("/delivery-note/{delivery_note_id}/cancel", v1.CancelDeliveryNoteHandler(salesService)).Methods("POST")

	// Get products and batches
	router.HandleFunc("/batches", v1.GetSalesOrderBatchesHandler(salesService)).Methods("GET")
	//router.HandleFunc("/product/{id}/batches", v1.GetProductBatchHandler(salesService)).Methods("GET")
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
	RegisterFinanceTransactionRoutes(AdminRoutes, services.FinanceService)

	// Inventory middleware setup
	InventoryRoutes := router.PathPrefix("/inventory").Subrouter()
	InventoryRoutes.Use(middleware.RoleMiddleware("inventory"))
	RegisterProductRoutes(InventoryRoutes, services.ProductService)
	RegisterCategoryRoutes(InventoryRoutes, services.CategoryService)
	RegisterUnitRoutes(InventoryRoutes, services.UnitService)
	RegisterInventoryRoutes(InventoryRoutes, services.InventoryService)

	//// Finance middleware setup
	//FinanceRoutes := router.PathPrefix("/finance").Subrouter()
	//FinanceRoutes.Use(middleware.RoleMiddleware("finance"))

	// Sales middleware setup
	SalesRoutes := router.PathPrefix("/sales").Subrouter()
	SalesRoutes.Use(middleware.RoleMiddleware("sales"))
	RegisterProductRoutes(SalesRoutes, services.ProductService)
	RegisterCustomerRoutes(SalesRoutes, services.CustomerService)
	RegisterSalesRoutes(SalesRoutes, services.SalesService)

	// Purchase middleware setup
	PurchaseRoutes := router.PathPrefix("/purchase").Subrouter()
	PurchaseRoutes.Use(middleware.RoleMiddleware("purchase"))
	RegisterSupplierRoutes(PurchaseRoutes, services.SupplierService)
	RegisterPurchaseOrderRoutes(PurchaseRoutes, services.PurchaseOrderService, services.ProductService, services.InventoryService)

	//// Purchase middleware setup
	//PurchaseRoutes := router.PathPrefix("/purchase").Subrouter()
	//PurchaseRoutes.Use(middleware.RoleMiddleware("purchase"))

}
