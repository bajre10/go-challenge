package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"shortcut-challenge/api"
	"shortcut-challenge/api/appErrors"
	"shortcut-challenge/database"
	"shortcut-challenge/models"
	"strconv"
	"time"
)

const RESTOCK_TIMEFRAME = -24 * time.Hour
const RESTOCK_MAX = 3
const LOW_STOCK_THRESHOLD = 20

func restockQuotaReached(itemId uint, r *http.Request) (bool, error) {
	db := database.GetDBInstance(r)

	var count int64
	threshold := time.Now().Add(RESTOCK_TIMEFRAME)

	result := db.Model(&models.Restock{}).
		Where("item_id = ? AND created_at > ?", itemId, threshold).
		Count(&count)

	if result.Error != nil {
		return false, errors.New("database error")
	}

	if count >= RESTOCK_MAX {
		return true, nil
	}

	return false, nil
}

type CreateItemDTO struct {
	Name        string `json:"name"`
	Quantity    uint   `json:"quantity"`
	Description string `json:"description"`
}

// CreateItem godoc
// @Security BearerAuth
// @Summary Create a new inventory item
// @Description Creates a new item in the inventory with the provided details, validates input, and handles errors such as duplicate items.
// @Tags Inventory
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param item body CreateItemDTO true "Inventory item details"
// @Success 201 {object} models.InventoryItem "Item successfully created"
// @Failure 400 {object} api.RequestError "Invalid input or item already exists"
// @Failure 500 {object} api.RequestError "Internal server error"
// @Router /inventory [post]
func CreateItem(w http.ResponseWriter, r *http.Request) {
	db := database.GetDBInstance(r)

	var itemDto CreateItemDTO

	if err := json.NewDecoder(r.Body).Decode(&itemDto); err != nil {
		api.ThrowRequestError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	item := models.InventoryItem{
		Name:        itemDto.Name,
		Quantity:    itemDto.Quantity,
		Description: itemDto.Description,
	}

	if err := db.Create(&item).Error; err != nil {
		var dbError *appErrors.DBError
		if errors.As(err, &dbError) {
			api.ThrowRequestError(w, dbError.Message, http.StatusBadRequest)
			return
		} else if errors.Is(err, gorm.ErrDuplicatedKey) {
			api.ThrowRequestError(w, "Item already exists", http.StatusConflict)
		} else {
			api.ThrowRequestError(w, "Error creating item", http.StatusInternalServerError)
			return
		}

		return
	}

	api.SendJson(w, http.StatusCreated, item)
}

type RestockItemDTO struct {
	Quantity uint `json:"quantity"`
}

// RestockItem godoc
// @Security BearerAuth
// @Summary Restock an inventory item
// @Description Restocks an inventory item with the provided quantity, checks if restock quota is reached, and handles errors such as invalid itemID, quota exceeded, or database issues.
// @Tags Inventory
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param itemID path int true "Inventory Item ID"
// @Param restock body RestockItemDTO true "Restock details"
// @Success 201 {object} models.Restock "Item successfully restocked"
// @Failure 400 {object} api.RequestError "Invalid itemID or invalid input"
// @Failure 404 {object} api.RequestError "Item not found"
// @Failure 409 {object} api.RequestError "Item already exists"
// @Failure 429 {object} api.RequestError "Item quota reached"
// @Failure 500 {object} api.RequestError "Internal server error"
// @Router /inventory/{itemID}/restock [post]
func RestockItem(w http.ResponseWriter, r *http.Request) {
	db := database.GetDBInstance(r)

	itemIDStr := chi.URLParam(r, "itemID")

	if itemIDStr == "" {
		api.ThrowRequestError(w, "Invalid itemID", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		api.ThrowRequestError(w, "Invalid itemID", http.StatusBadRequest)
		return
	}

	var restockItemDTO RestockItemDTO

	if err := json.NewDecoder(r.Body).Decode(&restockItemDTO); err != nil {
		api.ThrowRequestError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Get corresponding item
	var item models.InventoryItem

	if err := db.Where("id = ?", uint(itemID)).First(&item).Error; err != nil {
		api.ThrowRequestError(w, "Invalid itemID", http.StatusBadRequest)
		return
	}

	// Create restock object
	restock := models.Restock{
		ItemID:   item.ID,
		Quantity: restockItemDTO.Quantity,
	}

	// Check if restock quota is reached
	qr, err := restockQuotaReached(restock.ItemID, r)

	if err != nil {
		api.ThrowInternalError(w)
		return
	}

	if qr == true {
		api.ThrowRequestError(w, "Item quota reached", http.StatusTooManyRequests)
		return
	}

	// Add restock quantity to item
	item.Quantity += restock.Quantity

	// Start transaction
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&restock).Error; err != nil {
			// Rollback after error
			return err
		}

		// Set lastRestock
		item.LastRestock = restock.CreatedAt

		if err := tx.Save(&item).Error; err != nil {
			return err
		}

		// Return nil; commit transaction
		return nil
	}); err != nil {
		var dbError *appErrors.DBError
		if errors.As(err, &dbError) {
			api.ThrowRequestError(w, dbError.Message, http.StatusBadRequest)
			return
		} else if errors.Is(err, gorm.ErrDuplicatedKey) {
			api.ThrowRequestError(w, "Item already exists", http.StatusConflict)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			api.ThrowRequestError(w, "Item not found", http.StatusNotFound)
			return
		} else {
			api.ThrowRequestError(w, "Error restocking item", http.StatusInternalServerError)
			return
		}
	}

	api.SendJson(w, http.StatusCreated, restock)
}

// GetInventoryItems godoc
// @Security BearerAuth
// @Summary Get a list of inventory items
// @Description Retrieves a list of inventory items from the database. Optionally filters the results to show items with low stock, based on the query parameter `lowStock`.
// @Tags Inventory
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param lowStock query string false "Filter by low stock (quantity <= LOW_STOCK_THRESHOLD)"
// @Success 200 {array} models.InventoryItem "List of inventory items"
// @Failure 500 {object} api.RequestError "Internal server error"
// @Router /inventory [get]
func GetInventoryItems(w http.ResponseWriter, r *http.Request) {
	db := database.GetDBInstance(r)

	var items []models.InventoryItem
	query := db

	lowStock := r.URL.Query().Get("lowStock")

	if lowStock != "" {
		query = query.Where("quantity <= ?", LOW_STOCK_THRESHOLD)
	}

	if err := query.Find(&items).Error; err != nil {
		http.Error(w, "Could not fetch inventory items", http.StatusInternalServerError)
		return
	}

	api.SendJson(w, http.StatusOK, items)
}

type RestockHistoryDTO struct {
	ItemID        uint      `json:"itemID"`
	Name          string    `json:"name"`
	RestockAmount uint      `json:"amount"`
	Time          time.Time `json:"time"`
}

// GetRestockHistory godoc
// @Security BearerAuth
// @Summary Get restock history for inventory items
// @Description Retrieves the restock history for inventory items. Optionally filters the history based on the `itemId` query parameter.
// @Tags Inventory, Restock
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param itemId query string false "Filter restock history by item ID"
// @Success 200 {array} RestockHistoryDTO "List of restock history records"
// @Failure 400 {object} api.RequestError "Invalid input"
// @Failure 404 {object} api.RequestError "Item not found"
// @Failure 500 {object} api.RequestError "Internal server error"
// @Router /inventory/restock [get]
func GetRestockHistory(w http.ResponseWriter, r *http.Request) {
	db := database.GetDBInstance(r)

	var restocks []models.Restock
	query := db

	itemId := r.URL.Query().Get("itemId")

	if itemId != "" {
		query = query.Where("item_id <= ?", itemId)
	}

	if err := query.Joins("InventoryItem").Find(&restocks).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.ThrowRequestError(w, "Item not found", http.StatusNotFound)
			return
		} else {
			api.ThrowRequestError(w, "Could not get inventory items", http.StatusInternalServerError)
			return
		}
	}

	var res []RestockHistoryDTO

	for _, restock := range restocks {
		res = append(res, RestockHistoryDTO{
			Name:          restock.InventoryItem.Name,
			ItemID:        restock.ItemID,
			RestockAmount: restock.Quantity,
			Time:          restock.CreatedAt,
		})
	}

	api.SendJson(w, http.StatusOK, res)
}
