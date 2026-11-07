package storage

import "testing"

func TestStorageSet(t *testing.T) {
	store := New()
	store.Set("key", "value")
	v := store.store["key"]

	if v != "value" {
		t.Errorf("store.Set(\"test\", \"value\"), get from store = %s, want = value", v)
	}
}

func TestStorageGet(t *testing.T) {
	store := New()
	store.store["key"] = "value"
	v := store.Get("key")

	if v["key"] != "value" {
		t.Errorf("store.Get(\"key\") = %s, want = value", v)
	}
}

func TestStorageGetOnNotExists(t *testing.T) {
	store := New()
	v := store.Get("key")

	if v["key"] != "" {
		t.Errorf("store.Get(\"key\") = %s, want = \"\"", v)
	}
}

func TestStorageDel(t *testing.T) {
	store := New()
	store.store["key"] = "value"
	store.Del("key")

	if len(store.store) > 0 {
		t.Errorf("store.Del(\"key\"), length of store > 0, want = 0")
	}
}

func TestStorageExists(t *testing.T) {
	store := New()
	store.store["key"] = "value"

	tmp := store.Exists("key")

	if tmp["key"] != true {
		t.Errorf("store.Exists(\"key\") = false, want = true")
	}
}
