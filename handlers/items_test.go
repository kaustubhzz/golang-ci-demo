package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hs-heilbronn-devsecops/acetlisto/stores"
	"github.com/stretchr/testify/assert"
)

func TestListItems(t *testing.T) {
	store := stores.NewMemoryItemStore()
	item, _ := store.CreateItem(stores.CreateItemRequest{Name: "Test Item", Description: "This is a test item"})
	router := New(store)

	req, _ := http.NewRequest("GET", "/items/", nil)
	response := executeRequest(router, req)

	assert.Equal(t, http.StatusOK, response.Code)

	var items []stores.Item
	err := json.Unmarshal(response.Body.Bytes(), &items)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, *item, items[0])
}

func TestCreateItem(t *testing.T) {
	store := stores.NewMemoryItemStore()
	router := New(store)

	payload := `{"name":"New Item","description":"This is a new item"}`
	req, _ := http.NewRequest("POST", "/items/", http.NoBody)
	req.Body = io.NopCloser(strings.NewReader(payload))
	response := executeRequest(router, req)

	assert.Equal(t, http.StatusCreated, response.Code)

	var item stores.Item
	err := json.Unmarshal(response.Body.Bytes(), &item)
	assert.NoError(t, err)
	assert.Equal(t, "New Item", item.Name)
	assert.Equal(t, "This is a new item", item.Description)
	assert.NotEmpty(t, item.ID)
}

func TestGetItem(t *testing.T) {
	store := stores.NewMemoryItemStore()
	item, _ := store.CreateItem(stores.CreateItemRequest{Name: "Test Item", Description: "This is a test item"})
	router := New(store)

	req, _ := http.NewRequest("GET", "/items/"+item.ID, nil)
	response := executeRequest(router, req)

	assert.Equal(t, http.StatusOK, response.Code)

	var fetchedItem stores.Item
	err := json.Unmarshal(response.Body.Bytes(), &fetchedItem)
	assert.NoError(t, err)
	assert.Equal(t, *item, fetchedItem)
}

func TestUpdateItem(t *testing.T) {
	store := stores.NewMemoryItemStore()
	item, _ := store.CreateItem(stores.CreateItemRequest{Name: "Test Item", Description: "This is a test item"})
	router := New(store)

	payload := `{"name":"Updated Item","description":"This item has been updated"}`
	req, _ := http.NewRequest("PUT", "/items/"+item.ID, http.NoBody)
	req.Body = io.NopCloser(strings.NewReader(payload))
	response := executeRequest(router, req)

	assert.Equal(t, http.StatusOK, response.Code)

	var updatedItem stores.Item
	err := json.Unmarshal(response.Body.Bytes(), &updatedItem)
	assert.NoError(t, err)
	assert.Equal(t, item.ID, updatedItem.ID)
	assert.Equal(t, "Updated Item", updatedItem.Name)
	assert.Equal(t, "This item has been updated", updatedItem.Description)
}

func TestDeleteItem(t *testing.T) {
	store := stores.NewMemoryItemStore()
	item, _ := store.CreateItem(stores.CreateItemRequest{Name: "Test Item", Description: "This is a test item"})
	router := New(store)

	req, _ := http.NewRequest("DELETE", "/items/"+item.ID, nil)
	response := executeRequest(router, req)

	assert.Equal(t, http.StatusOK, response.Code)

	// Verify item is deleted
	_, err := store.GetItem(item.ID)
	assert.Error(t, err)
}

func TestGetNonExistentItem(t *testing.T) {
	store := stores.NewMemoryItemStore()
	router := New(store)

	req, _ := http.NewRequest("GET", "/items/nonexistent", nil)
	response := executeRequest(router, req)

	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestUpdateNonExistentItem(t *testing.T) {
	store := stores.NewMemoryItemStore()
	router := New(store)

	payload := `{"name":"Updated Item","description":"This item has been updated"}`
	req, _ := http.NewRequest("PUT", "/items/nonexistent", http.NoBody)
	req.Body = io.NopCloser(strings.NewReader(payload))
	response := executeRequest(router, req)

	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestDeleteNonExistentItem(t *testing.T) {
	store := stores.NewMemoryItemStore()
	router := New(store)

	req, _ := http.NewRequest("DELETE", "/items/nonexistent", nil)
	response := executeRequest(router, req)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

func TestCreateItemInvalidPayload(t *testing.T) {
	store := stores.NewMemoryItemStore()
	router := New(store)

	payload := `{"name":123,"description":"This is a new item"}`
	req, _ := http.NewRequest("POST", "/items/", http.NoBody)
	req.Body = io.NopCloser(strings.NewReader(payload))
	response := executeRequest(router, req)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestUpdateItemInvalidPayload(t *testing.T) {
	store := stores.NewMemoryItemStore()
	item, _ := store.CreateItem(stores.CreateItemRequest{Name: "Test Item", Description: "This is a test item"})
	router := New(store)

	payload := `{"name":123,"description":"This item has been updated"}`
	req, _ := http.NewRequest("PUT", "/items/"+item.ID, http.NoBody)
	req.Body = io.NopCloser(strings.NewReader(payload))
	response := executeRequest(router, req)

	assert.Equal(t, http.StatusBadRequest, response.Code)

}

func TestListItemsStoreError(t *testing.T) {
	// Create a mock store that returns an error on GetAllItems
	store := &mockItemStore{
		getAllItemsFunc: func() ([]stores.Item, error) {
			return nil, assert.AnError
		},
	}
	router := New(store)

	req, _ := http.NewRequest("GET", "/items/", nil)
	response := executeRequest(router, req)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

func TestCreateItemStoreError(t *testing.T) {
	// Create a mock store that returns an error on CreateItem
	store := &mockItemStore{
		createItemFunc: func(req stores.CreateItemRequest) (*stores.Item, error) {
			return nil, assert.AnError
		},
	}
	router := New(store)

	payload := `{"name":"New Item","description":"This is a new item"}`
	req, _ := http.NewRequest("POST", "/items/", http.NoBody)
	req.Body = io.NopCloser(strings.NewReader(payload))
	response := executeRequest(router, req)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

// mockItemStore implements stores.ItemStore for testing error cases
type mockItemStore struct {
	getAllItemsFunc func() ([]stores.Item, error)
	createItemFunc  func(stores.CreateItemRequest) (*stores.Item, error)
	getItemFunc     func(string) (*stores.Item, error)
	updateItemFunc  func(string, stores.CreateItemRequest) (*stores.Item, error)
	deleteItemFunc  func(string) error
}

func (m *mockItemStore) GetAllItems() ([]stores.Item, error) {
	if m.getAllItemsFunc != nil {
		return m.getAllItemsFunc()
	}
	return nil, nil
}
func (m *mockItemStore) CreateItem(req stores.CreateItemRequest) (*stores.Item, error) {
	if m.createItemFunc != nil {
		return m.createItemFunc(req)
	}
	return nil, nil
}
func (m *mockItemStore) GetItem(id string) (*stores.Item, error) {
	if m.getItemFunc != nil {
		return m.getItemFunc(id)
	}
	return nil, nil
}
func (m *mockItemStore) UpdateItem(id string, req stores.CreateItemRequest) (*stores.Item, error) {
	if m.updateItemFunc != nil {
		return m.updateItemFunc(id, req)
	}
	return nil, nil
}
func (m *mockItemStore) DeleteItem(id string) error {
	if m.deleteItemFunc != nil {
		return m.deleteItemFunc(id)
	}
	return nil
}
