package storage

import (
	"testing"
	"time"
)

func TestStorageSet(t *testing.T) {
	store := New()
	store.Set("key", "value", 0)

	item := store.store["key"]

	if item.Value != "value" {
		t.Errorf("store.Set(\"test\", \"value\"), get from store = %s, want = value", item.Value)
	}
}

func TestStorageGet(t *testing.T) {
	store := New()
	store.store["key"] = item{
		Value:      "value",
		ExpireTime: time.Now(),
	}

	v := store.Get("key")

	if v["key"].Value != "value" {
		t.Errorf("store.Get(\"key\") = %s, want = value", v["key"].Value)
	}
}

func TestStorageGetOnNotExists(t *testing.T) {
	store := New()
	v := store.Get("key")

	if v["key"].Value != "" {
		t.Errorf("store.Get(\"key\") = %s, want = \"\"", v["key"].Value)
	}
}

func TestStorageDel(t *testing.T) {
	store := New()
	now := time.Now()
	store.store["key"] = item{
		Value:      "value",
		ExpireTime: now,
	}

	store.Del("key")

	if len(store.store) > 0 {
		t.Errorf("store.Del(\"key\"), length of store > 0, want = 0")
	}
}

func TestStorageExists(t *testing.T) {
	store := New()
	now := time.Now()
	store.store["key"] = item{
		Value:      "value",
		ExpireTime: now,
	}

	tmp := store.Exists("key")

	if tmp["key"] != true {
		t.Errorf("store.Exists(\"key\") = false, want = true")
	}
}

func TestStorageDump(t *testing.T) {
	store := New()
	now := time.Now()
	store.store["key"] = item{
		"value",
		now,
	}

	tmp := store.Dump()

	if tmp["key"].Value != "value" {
		t.Errorf("store.Dump() with key = key, get = %s, want = value", tmp["key"].Value)
	}

	tmp["key"] = item{
		"v",
		now,
	}

	if store.store["key"].Value == "v" {
		t.Errorf("store.Dump() change value by key, changed original")
	}
}
