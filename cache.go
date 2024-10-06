package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type Cache struct {
	ProcessedPermalinks map[string]bool `json:"processed_permalinks"`
	LastCacheUpdate  time.Time       	`json:"last_cache_update"`
	LastPersisted       time.Time       `json:"last_persisted"`
	mu                  sync.Mutex
}

func NewCache() *Cache {
	return &Cache{
		ProcessedPermalinks: make(map[string]bool),
	}
}

func (c *Cache) Load(filename string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("No cache file found, starting with an empty cache")
			return nil // No cache file exists, start with an empty cache
		}
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(c); err != nil {
		return err
	}

	fileInfo, err := os.Stat(filename)
	if err == nil {
		c.LastPersisted = fileInfo.ModTime()
	}

	log.Printf("Successfully loaded %d entries from cache", len(c.ProcessedPermalinks))
	return nil
}

func (c *Cache) Save(filename string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if the cache has been updated since the last save
	fileInfo, err := os.Stat(filename)
	if err == nil && fileInfo.ModTime().After(c.LastCacheUpdate) {
		log.Println("Cache has not changed since the last save, skipping persistence")
		return nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(c); err != nil {
		return err
	}

	c.LastPersisted = time.Now()
	log.Println("Cache successfully persisted to disk")
	return nil
}

func (c *Cache) AddProcessedPermalink(permalink string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.ProcessedPermalinks) >= 100 {
		// Remove the oldest entry
		for k := range c.ProcessedPermalinks {
			delete(c.ProcessedPermalinks, k)
			break
		}
	}

	c.ProcessedPermalinks[permalink] = true
	c.LastCacheUpdate = time.Now()
}

func (c *Cache) IsProcessed(permalink string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.ProcessedPermalinks[permalink]
}


func archiveCorruptedCache(cacheFile string) {
    archiveFile := cacheFile + ".archive.bak"
    if err := os.Rename(cacheFile, archiveFile); err != nil {
        log.Printf("Error archiving corrupted cache: %v", err)
    } else {
        log.Printf("Archived corrupted cache to: %s", archiveFile)
    }
}


func (c *Cache) CreateNew(filename string) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    // Initialize a new cache
    c.ProcessedPermalinks = make(map[string]bool)
    c.LastCacheUpdate = time.Now()
    c.LastPersisted = time.Now()

    // Save the new cache to the specified file
    return c.Save(filename)
}


func (c *Cache) EnsureUsable(filename string) error {
    if err := c.Load(filename); err != nil {
        log.Printf("Error loading cache: %v", err)
        archiveCorruptedCache(filename)
        if err := c.CreateNew(filename); err != nil {
            return err
        }
    }
    return nil
}