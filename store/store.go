package store

type KvObjectDict map[string]kvObj

type ExpiryDict map[string]int

type Store struct {
	Dict   *KvObjectDict
	Expiry *ExpiryDict
}

type StoreInterface interface {
	GetValue(key string) interface{}
	SetValue(key string, value interface{})
	DeleteValue(key string) bool
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
		/* 
		Doing this will remove the key from the dictionary, ideally we should decrement 
		the refcount and later batch remove the keys with refcount 0 
		*/
		obj.refcount = 0
		(*s.Dict)[key] = obj
		return true
	}
	return false
}

func (s *Store) Exists(key string) bool {
	obj, exists := (*s.Dict)[key]
	return exists && obj.refcount > 0
}

