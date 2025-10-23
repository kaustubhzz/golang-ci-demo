package stores

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialState(t *testing.T) {
	m := NewMemoryItemStore()

	items, _ := m.GetAllItems()
	assert.Equal(t, 0, len(items), "Store should initially have no items")

}

func TestCreate(t *testing.T) {
	m := NewMemoryItemStore()

	item, err := m.CreateItem(CreateItemRequest{
		Name:        "Oil",
		Description: "Some beautiful oil.",
	})
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Greater(t, len(item.ID), 6, "ID must be longer than 6 characters.")
}

func TestCreateAndLoad(t *testing.T) {
	m := NewMemoryItemStore()

	item, _ := m.CreateItem(CreateItemRequest{
		Name:        "Oil",
		Description: "Some beautiful oil.",
	})

	loadedItem, err := m.GetItem(item.ID)
	assert.NoError(t, err)
	assert.NotNil(t, loadedItem)
	assert.Equal(t, loadedItem, item)
}

func TestGetItemForNonExistingID(t *testing.T) {
	m := NewMemoryItemStore()

	loadedItem, err := m.GetItem("dummy")
	assert.Error(t, err)
	assert.Nil(t, loadedItem)
}

func TestUpdateItem(t *testing.T) {
	m := NewMemoryItemStore()

	item, _ := m.CreateItem(CreateItemRequest{
		Name:        "Oil",
		Description: "Some beautiful oil.",
	})

	updatedItem, err := m.UpdateItem(item.ID, CreateItemRequest{
		Name:        "Another Oil",
		Description: "Must be saucy!",
	})
	assert.NoError(t, err)
	assert.NotNil(t, updatedItem)
	assert.Equal(t, &Item{
		ID:          item.ID,
		Name:        "Another Oil",
		Description: "Must be saucy!",
	},
		updatedItem)
}

func TestUpdateItemForNonExistingID(t *testing.T) {
	m := NewMemoryItemStore()

	updatedItem, err := m.UpdateItem("dummy", CreateItemRequest{
		Name:        "Another Oil",
		Description: "Must be saucy!",
	})
	assert.Error(t, err)
	assert.Nil(t, updatedItem)
}

func TestDeleteItem(t *testing.T) {
	m := NewMemoryItemStore()

	item, _ := m.CreateItem(CreateItemRequest{
		Name:        "Oil",
		Description: "Some beautiful oil.",
	})

	items, _ := m.GetAllItems()
	assert.Equal(t, 1, len(items))

	err := m.DeleteItem(item.ID)
	assert.NoError(t, err)

	items, _ = m.GetAllItems()
	assert.Equal(t, 0, len(items))
}

func TestDeleteItemForNonExistingID(t *testing.T) {
	m := NewMemoryItemStore()

	err := m.DeleteItem("dummy")
	assert.Error(t, err)
}
