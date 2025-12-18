package db

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var guestCartFile = filepath.Join(".", "db", "guest_carts.json")
var gcMutex sync.Mutex

type guestCartStore map[string]map[int]int // sessionID -> productID -> quantity

func loadGuestStore() (guestCartStore, error) {
	gcMutex.Lock()
	defer gcMutex.Unlock()

	store := make(guestCartStore)
	if _, err := os.Stat(guestCartFile); os.IsNotExist(err) {
		log.Printf("GuestCart: file not found, returning empty store: %s", guestCartFile)
		return store, nil
	}
	b, err := ioutil.ReadFile(guestCartFile)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		log.Printf("GuestCart: file empty, returning empty store: %s", guestCartFile)
		return store, nil
	}
	if err := json.Unmarshal(b, &store); err != nil {
		log.Printf("GuestCart: failed to unmarshal file %s: %v", guestCartFile, err)
		return nil, err
	}
	log.Printf("GuestCart: loaded store with %d sessions", len(store))
	return store, nil
}

func saveGuestStore(store guestCartStore) error {
	gcMutex.Lock()
	defer gcMutex.Unlock()
	b, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	// ensure dir
	if err := os.MkdirAll(filepath.Dir(guestCartFile), 0o755); err != nil {
		log.Printf("GuestCart: mkdir error for %s: %v", guestCartFile, err)
		return err
	}
	if err := ioutil.WriteFile(guestCartFile, b, 0o644); err != nil {
		log.Printf("GuestCart: write file error %s: %v", guestCartFile, err)
		return err
	}
	log.Printf("GuestCart: saved store (%d sessions) to %s", len(store), guestCartFile)
	return nil
}

// AddItem increases quantity for product in session cart
func AddGuestItem(sessionID string, productID, qty int) error {
	store, err := loadGuestStore()
	if err != nil {
		return err
	}
	if _, ok := store[sessionID]; !ok {
		store[sessionID] = make(map[int]int)
	}
	store[sessionID][productID] += qty
	err = saveGuestStore(store)
	if err == nil {
		log.Printf("GuestCart: added product %d (qty %d) to session %s", productID, qty, sessionID)
	}
	return err
}

// UpdateGuestItem sets quantity (or removes if qty <=0)
func UpdateGuestItem(sessionID string, productID, qty int) error {
	store, err := loadGuestStore()
	if err != nil {
		return err
	}
	if _, ok := store[sessionID]; !ok {
		store[sessionID] = make(map[int]int)
	}
	if qty <= 0 {
		delete(store[sessionID], productID)
	} else {
		store[sessionID][productID] = qty
	}
	err = saveGuestStore(store)
	if err == nil {
		log.Printf("GuestCart: updated product %d qty to %d for session %s", productID, qty, sessionID)
	}
	return err
}

// RemoveGuestItem removes a product from session cart
func RemoveGuestItem(sessionID string, productID int) error {
	return UpdateGuestItem(sessionID, productID, 0)
}

// ListGuestItems returns the map productID->quantity for session
func ListGuestItems(sessionID string) (map[int]int, error) {
	store, err := loadGuestStore()
	if err != nil {
		return nil, err
	}
	m, ok := store[sessionID]
	if !ok {
		log.Printf("GuestCart: no items for session %s", sessionID)
		return map[int]int{}, nil
	}
	log.Printf("GuestCart: listing %d items for session %s", len(m), sessionID)
	return m, nil
}

// ClearGuestCart removes all items for a session
func ClearGuestCart(sessionID string) error {
	store, err := loadGuestStore()
	if err != nil {
		return err
	}
	delete(store, sessionID)
	err = saveGuestStore(store)
	if err == nil {
		log.Printf("GuestCart: cleared cart for session %s", sessionID)
	}
	return err
}
