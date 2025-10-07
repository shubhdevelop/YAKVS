package store

import "time"

type KvObjectDict map[string]kvObj

type ExpiryDict map[string]int64  // Separate expires dictionary: key -> unix_timestamp

type Store struct {
	Dict   *KvObjectDict
	Expiry *ExpiryDict
}

type StoreInterface interface {
	GetValue(key string) interface{}
	SetValue(key string, value interface{})
	DeleteValue(key string) bool
	Exists(key string) bool
	GetTTL(key string) int
	SetTTL(key string, ttl int) bool
	RemoveExpiry(key string) bool
}

func NewStore() *Store {
	dict := make(KvObjectDict, 0)
	expiry := make(ExpiryDict, 0)
	return &Store{
		Dict:   &dict, 
		Expiry: &expiry,
	}
}

// for the given key get the kvObject and return the value
func (s *Store) GetValue(key string) interface{} {

	// if it exists in the expiry dictionary, check if it has expired
	if _, exists := (*s.Expiry)[key]; exists {
		if time.Now().Unix() > (*s.Expiry)[key] {
			delete(*s.Expiry, key)
			delete(*s.Dict, key)
			return nil // Key has expired
		}
	}
	// only return if the ref count if greater than 0
	if obj, exists := (*s.Dict)[key]; exists && obj.refcount > 0 {
		// Handle different encodings based on the object's encoding
		switch obj.getEncoding() {
		case OBJ_ENCODING_INT:
			// For integer encoding, the ptr points to an int
			return *(*int)(obj.ptr)
		case OBJ_ENCODING_RAW:
			// For raw string encoding, the ptr points to a string
			return *(*string)(obj.ptr)
		default:
			// For other encodings, return the pointer as-is
			return obj.ptr
		}
	}
	return nil
}

func (s *Store) SetValue(key string, value interface{}) {
	if strVal, ok := value.(string); ok {
		kvObj := createStringObj(strVal)
		(*s.Dict)[key] = *kvObj
	} else if intVal, ok := value.(int); ok {
		kvObj := createIntObj(intVal)
		(*s.Dict)[key] = *kvObj
	}	
}


func (s *Store) DeleteValue(key string) bool {
	if obj, exists := (*s.Dict)[key]; exists {
		obj.refcount = 0 // we can remove the key from the dictionary
		delete(*s.Dict, key)
		delete(*s.Expiry, key)
		return true
	}
	return false
}

func (s *Store) Exists(key string) bool {
	obj, exists := (*s.Dict)[key]
	return exists && obj.refcount > 0
}

func (s *Store) GetTTL(key string) int {
	if _, exists := (*s.Dict)[key]; !exists {
		return -2 // Key doesn't exist at all
	}
	
	// Check if key has expiry set
	ttl, hasExpiry := (*s.Expiry)[key]
	if !hasExpiry {
		return -1 // Key exists but has no expiry
	}
	
	// calculate the time difference between the current time and the expiry time
	timeDiff := time.Until(time.Unix(ttl, 0))

	if timeDiff.Seconds() < 0 {
		delete(*s.Expiry, key)
		delete(*s.Dict, key)
		return -2 // Key has expired
	}

	return int(timeDiff.Seconds())
}

// SetTTL sets the time-to-live for a key in seconds
func (s *Store) SetTTL(key string, ttl int64) bool {
	// Check if the key exists in the main dictionary
	if _, exists := (*s.Dict)[key]; !exists {
		return false // Key doesn't exist
	}
	
	// Set the expiry
	(*s.Expiry)[key] = ttl
	return true
}

// RemoveExpiry removes the TTL from a key, making it persistent
func (s *Store) RemoveExpiry(key string) bool {
	// Check if the key exists in the main dictionary
	if _, exists := (*s.Dict)[key]; !exists {
		return false // Key doesn't exist
	}
	
	// Remove from expiry dictionary
	delete(*s.Expiry, key)
	return true
}

