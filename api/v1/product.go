package v1

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/product"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"
)

// GetAllProductHandler fetches all products
func GetAllProductHandler(productService *product.ProductService) http.HandlerFunc {
	return utils.NewPaginatedHandler(func(w http.ResponseWriter, r *http.Request, page, pageSize int, sortBy, sortOrder string) {
		var req product.GetProductRequest
		req.Name, req.Category, req.Unit = r.URL.Query().Get("name"), r.URL.Query().Get("category"), r.URL.Query().Get("unit")
		req.Page = page
		req.PageSize = pageSize
		req.SortBy = sortBy
		req.SortOrder = sortOrder
		// Validate req
		validationErrors := utils.ValidateStruct(req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		products, totalItems, err := productService.GetAllProducts(req)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		utils.WritePaginationJSON(w, http.StatusOK, req.Page, totalItems, req.PageSize, products)
	})

	//return func(w http.ResponseWriter, r *http.Request) {
	//	var req product.GetProductRequest
	//	req.Name, req.Category, req.Unit = r.URL.Query().Get("name"), r.URL.Query().Get("category"), r.URL.Query().Get("unit")
	//	req.Page = page
	//	req.PageSize = pageSize
	//	req.SortBy = sortBy
	//	req.SortOrder = sortOrder
	//	// Validate req
	//	validationErrors := utils.ValidateStruct(req)
	//	if validationErrors != nil {
	//		utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
	//		return
	//	}
	//
	//	products, totalItems, err := productService.GetAllProducts(req)
	//	if err != nil {
	//		utils.ErrorJSON(w, err)
	//		return
	//	}
	//
	//	utils.WritePaginationJSON(w, http.StatusOK, products, totalItems)
	//}
}

// CreateProductHandler creates a new product
func CreateProductHandler(productService *product.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req product.CreateProductRequest

		validationErrors := utils.DecodeAndValidate(r, &req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		createProduct, errSer := productService.CreateProduct(req)
		if errSer != nil {
			utils.ErrorJSON(w, errSer)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"message": fmt.Sprintf("Produk %s berhasil didaftarkan", createProduct.Name),
		})
	}
}

// UpdateProductHandler updates a product
func UpdateProductHandler(productService *product.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from parameter
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"general": "ID tidak valid",
			}))
			return
		}
		var req product.UpdateProductRequest
		req.ID = id

		validationErrors := utils.DecodeAndValidate(r, &req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		updateProduct, errService := productService.UpdateProduct(req)
		if errService != nil {
			utils.ErrorJSON(w, errService)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Produk berhasil diupdate ke %s", updateProduct.Name)})
	}
}

// DeleteProductHandler soft deletes a product
func DeleteProductHandler(productService *product.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from parameter
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"general": "ID tidak valid",
			}))
			return
		}
		var req product.DeleteProductRequest
		req.ID = id

		validationErrors := utils.ValidateStruct(&req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		errService := productService.DeleteProduct(req)
		if errService != nil {
			utils.ErrorJSON(w, errService)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Produk berhasil dihapus"})
	}
}
