package storage

import (
	"testing"
	"time"
)

func TestStorageSet(t *testing.T) {
	store := New()
	store.Set("key", "value", 0)

	item := store.store["key"]

	if item.value != "value" {
		t.Errorf("store.Set(\"test\", \"value\"), get from store = %s, want = value", item.value)
	}
}

func TestStorageGet(t *testing.T) {
	store := New()
	store.store["key"] = item{
		value:      "value",
		expireTime: time.Now(),
	}

	v := store.Get("key")

	if v["key"].value != "value" {
		t.Errorf("store.Get(\"key\") = %s, want = value", v["key"].value)
	}
}

func TestStorageGetOnNotExists(t *testing.T) {
	store := New()
	v := store.Get("key")

	if v["key"].value != "" {
		t.Errorf("store.Get(\"key\") = %s, want = \"\"", v["key"].value)
	}
}

func TestStorageDel(t *testing.T) {
	store := New()
	now := time.Now()
	store.store["key"] = item{
		value:      "value",
		expireTime: now,
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
		value:      "value",
		expireTime: now,
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

	if tmp["key"].value != "value" {
		t.Errorf("store.Dump() with key = key, get = %s, want = value", tmp["key"].value)
	}

	tmp["key"] = item{
		"v",
		now,
	}

	if store.store["key"].value == "v" {
		t.Errorf("store.Dump() change value by key, changed original")
	}
}
