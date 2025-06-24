package main

import (
	"github.com/corel-frim/item-packer-inc/api"
	"github.com/corel-frim/item-packer-inc/internal/storage"
)

// @title Item Packer API
// @version 1.0
// nolint:errcheck
func main() {
	// Create a new storage instance
	packStorage := storage.NewPackStorage()

	// Add some default packs
	packStorage.AddPack(250)
	packStorage.AddPack(500)
	packStorage.AddPack(1000)
	packStorage.AddPack(2000)
	packStorage.AddPack(5000)

	newAPI := api.NewAPI(packStorage)
	newAPI.Start()
}
