package executor

import (
	"fmt"
	"testing"
	"time"

	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

func TestExecuteCommandAysnc(t *testing.T) {
	tests := []struct {
		name     string
		command  *parser.Command
		setup    func(*store.Store) // Optional setup function
		verify   func(*store.Store) // Verification function
	}{
		{
			name: "SET command",
			command: &parser.Command{
				Name: "SET",
				Args: []string{"testkey", "testvalue"},
			},
			verify: func(s *store.Store) {
				value := s.GetValue("testkey")
				if value == nil {
					t.Error("Expected value to be set, got nil")
				}
				if value != "testvalue" {
					t.Errorf("Expected 'testvalue', got %v", value)
				}
			},
		},
		{
			name: "GET command - existing key",
			command: &parser.Command{
				Name: "GET",
				Args: []string{"testkey"},
			},
			setup: func(s *store.Store) {
				s.SetValue("testkey", "testvalue")
			},
			verify: func(s *store.Store) {
				// The GET command prints to stdout, so we can't easily test the output
				// But we can verify the key exists
				if !s.Exists("testkey") {
					t.Error("Expected key to exist")
				}
			},
		},
		{
			name: "GET command - non-existing key",
			command: &parser.Command{
				Name: "GET",
				Args: []string{"nonexistent"},
			},
			verify: func(s *store.Store) {
				// Verify the key doesn't exist
				if s.Exists("nonexistent") {
					t.Error("Expected key to not exist")
				}
			},
		},
		{
			name: "DEL command - existing key",
			command: &parser.Command{
				Name: "DEL",
				Args: []string{"testkey"},
			},
			setup: func(s *store.Store) {
				s.SetValue("testkey", "testvalue")
			},
			verify: func(s *store.Store) {
				if s.Exists("testkey") {
					t.Error("Expected key to be deleted")
				}
			},
		},
		{
			name: "DEL command - non-existing key",
			command: &parser.Command{
				Name: "DEL",
				Args: []string{"nonexistent"},
			},
			verify: func(s *store.Store) {
				// Should not cause any errors
				if s.Exists("nonexistent") {
					t.Error("Expected key to not exist")
				}
			},
		},
		{
			name: "EXISTS command - existing key",
			command: &parser.Command{
				Name: "EXISTS",
				Args: []string{"testkey"},
			},
			setup: func(s *store.Store) {
				s.SetValue("testkey", "testvalue")
			},
			verify: func(s *store.Store) {
				if !s.Exists("testkey") {
					t.Error("Expected key to exist")
				}
			},
		},
		{
			name: "EXISTS command - non-existing key",
			command: &parser.Command{
				Name: "EXISTS",
				Args: []string{"nonexistent"},
			},
			verify: func(s *store.Store) {
				if s.Exists("nonexistent") {
					t.Error("Expected key to not exist")
				}
			},
		},
		{
			name: "TTL command - key with no expiry",
			command: &parser.Command{
				Name: "TTL",
				Args: []string{"testkey"},
			},
			setup: func(s *store.Store) {
				s.SetValue("testkey", "testvalue")
			},
			verify: func(s *store.Store) {
				ttl := s.GetTTL("testkey")
				if ttl != -1 {
					t.Errorf("Expected TTL to be -1 (no expiry), got %d", ttl)
				}
			},
		},
		{
			name: "TTL command - non-existing key",
			command: &parser.Command{
				Name: "TTL",
				Args: []string{"nonexistent"},
			},
			verify: func(s *store.Store) {
				ttl := s.GetTTL("nonexistent")
				if ttl != -2 {
					t.Errorf("Expected TTL to be -2 (key doesn't exist), got %d", ttl)
				}
			},
		},
		{
			name: "EXPIRE command - valid TTL",
			command: &parser.Command{
				Name: "EXPIRE",
				Args: []string{"testkey", "3600"}, // 1 hour
			},
			setup: func(s *store.Store) {
				s.SetValue("testkey", "testvalue")
			},
			verify: func(s *store.Store) {
				ttl := s.GetTTL("testkey")
				// TTL should be approximately 3600 seconds remaining
				if ttl < 3599 || ttl > 3601 {
					t.Errorf("Expected TTL to be around 3600 seconds, got %d", ttl)
				}
			},
		},
		{
			name: "EXPIRE command - non-existing key",
			command: &parser.Command{
				Name: "EXPIRE",
				Args: []string{"nonexistent", "3600"},
			},
			verify: func(s *store.Store) {
				// Should not set TTL for non-existing key
				ttl := s.GetTTL("nonexistent")
				if ttl != -2 {
					t.Errorf("Expected TTL to be -2 (key doesn't exist), got %d", ttl)
				}
			},
		},
		{
			name: "EXPIRE command - invalid TTL",
			command: &parser.Command{
				Name: "EXPIRE",
				Args: []string{"testkey", "invalid"},
			},
			setup: func(s *store.Store) {
				s.SetValue("testkey", "testvalue")
			},
			verify: func(s *store.Store) {
				// Should not change the TTL due to parsing error
				ttl := s.GetTTL("testkey")
				if ttl != -1 {
					t.Errorf("Expected TTL to be -1 (no expiry), got %d", ttl)
				}
			},
		},
		{
			name: "EXPIREAT command - valid timestamp",
			command: &parser.Command{
				Name: "EXPIREAT",
				Args: []string{"testkey", fmt.Sprintf("%d", time.Now().Unix()+3600)}, // 1 hour from now
			},
			setup: func(s *store.Store) {
				s.SetValue("testkey", "testvalue")
			},
			verify: func(s *store.Store) {
				ttl := s.GetTTL("testkey")
				// Should return remaining seconds until the timestamp (around 3600)
				if ttl < 3599 || ttl > 3601 {
					t.Errorf("Expected TTL to be around 3600 seconds, got %d", ttl)
				}
			},
		},
		{
			name: "EXPIREAT command - non-existing key",
			command: &parser.Command{
				Name: "EXPIREAT",
				Args: []string{"nonexistent", fmt.Sprintf("%d", time.Now().Unix()+3600)},
			},
			verify: func(s *store.Store) {
				// Should not set TTL for non-existing key
				ttl := s.GetTTL("nonexistent")
				if ttl != -2 {
					t.Errorf("Expected TTL to be -2 (key doesn't exist), got %d", ttl)
				}
			},
		},
		{
			name: "EXPIREAT command - invalid timestamp",
			command: &parser.Command{
				Name: "EXPIREAT",
				Args: []string{"testkey", "invalid"},
			},
			setup: func(s *store.Store) {
				s.SetValue("testkey", "testvalue")
			},
			verify: func(s *store.Store) {
				// Should not change the TTL due to parsing error
				ttl := s.GetTTL("testkey")
				if ttl != -1 {
					t.Errorf("Expected TTL to be -1 (no expiry), got %d", ttl)
				}
			},
		},
		{
			name: "TTL command - expired key (automatic cleanup)",
			command: &parser.Command{
				Name: "TTL",
				Args: []string{"expiredkey"},
			},
			setup: func(s *store.Store) {
				s.SetValue("expiredkey", "testvalue")
				// Set expiry to a past timestamp to simulate expired key
				s.SetTTL("expiredkey", time.Now().Unix()-3600) // 1 hour ago
			},
			verify: func(s *store.Store) {
				// Key should be automatically deleted when expired
				ttl := s.GetTTL("expiredkey")
				if ttl != -2 {
					t.Errorf("Expected TTL to be -2 (key doesn't exist - expired), got %d", ttl)
				}
				// Key should not exist in the store
				if s.Exists("expiredkey") {
					t.Error("Expected expired key to be automatically deleted")
				}
			},
		},
		{
			name: "PERSIST command - key with TTL",
			command: &parser.Command{
				Name: "PERSIST",
				Args: []string{"testkey"},
			},
			setup: func(s *store.Store) {
				s.SetValue("testkey", "testvalue")
				s.SetTTL("testkey", time.Now().Unix()+3600) // 1 hour from now
			},
			verify: func(s *store.Store) {
				// TTL should be -1 (no expiry) after PERSIST
				ttl := s.GetTTL("testkey")
				if ttl != -1 {
					t.Errorf("Expected TTL to be -1 (no expiry), got %d", ttl)
				}
				// Key should still exist
				if !s.Exists("testkey") {
					t.Error("Expected key to still exist after PERSIST")
				}
			},
		},
		{
			name: "PERSIST command - key without TTL",
			command: &parser.Command{
				Name: "PERSIST",
				Args: []string{"testkey"},
			},
			setup: func(s *store.Store) {
				s.SetValue("testkey", "testvalue")
				// No TTL set
			},
			verify: func(s *store.Store) {
				// TTL should still be -1 (no expiry)
				ttl := s.GetTTL("testkey")
				if ttl != -1 {
					t.Errorf("Expected TTL to be -1 (no expiry), got %d", ttl)
				}
				// Key should still exist
				if !s.Exists("testkey") {
					t.Error("Expected key to still exist after PERSIST")
				}
			},
		},
		{
			name: "PERSIST command - non-existing key",
			command: &parser.Command{
				Name: "PERSIST",
				Args: []string{"nonexistent"},
			},
			verify: func(s *store.Store) {
				// Key should not exist
				if s.Exists("nonexistent") {
					t.Error("Expected key to not exist")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize a fresh store for each test
			testStore := store.NewStore()
			
			// Run setup if provided
			if tt.setup != nil {
				tt.setup(testStore)
			}

			// Execute the command with a channel to wait for completion
			resultChan := make(chan ResultWithError, 1)
			ExecuteCommandAysnc(tt.command, testStore, resultChan)
			
			// Wait for the command to complete
			<-resultChan

			// Run verification if provided
			if tt.verify != nil {
				tt.verify(testStore)
			}
		})
	}
}

func TestExecuteCommandIntegration(t *testing.T) {
	// Integration test that tests multiple commands in sequence
	testStore := store.NewStore()

	// Test sequence: SET -> GET -> EXISTS -> EXPIRE -> TTL -> DEL -> EXISTS
	t.Run("Complete workflow", func(t *testing.T) {
		// SET
		resultChan := make(chan ResultWithError, 1)
		ExecuteCommandAysnc(&parser.Command{
			Name: "SET",
			Args: []string{"integration_test", "integration_value"},
		}, testStore, resultChan)
		<-resultChan
		
		if !testStore.Exists("integration_test") {
			t.Error("Key should exist after SET")
		}

		// GET
		resultChan = make(chan ResultWithError, 1)
		ExecuteCommandAysnc(&parser.Command{
			Name: "GET",
			Args: []string{"integration_test"},
		}, testStore, resultChan)
		<-resultChan

		// EXISTS
		resultChan = make(chan ResultWithError, 1)
		ExecuteCommandAysnc(&parser.Command{
			Name: "EXISTS",
			Args: []string{"integration_test"},
		}, testStore, resultChan)
		<-resultChan

		// EXPIRE
		resultChan = make(chan ResultWithError, 1)
		ExecuteCommandAysnc(&parser.Command{
			Name: "EXPIRE",
			Args: []string{"integration_test", "7200"}, // 2 hours
		}, testStore, resultChan)
		<-resultChan

		// TTL
		resultChan = make(chan ResultWithError, 1)
		ExecuteCommandAysnc(&parser.Command{
			Name: "TTL",
			Args: []string{"integration_test"},
		}, testStore, resultChan)
		<-resultChan

		// DEL
		resultChan = make(chan ResultWithError, 1)
		ExecuteCommandAysnc(&parser.Command{
			Name: "DEL",
			Args: []string{"integration_test"},
		}, testStore, resultChan)
		<-resultChan

		// EXISTS (should return false now)
		resultChan = make(chan ResultWithError, 1)
		ExecuteCommandAysnc(&parser.Command{
			Name: "EXISTS",
			Args: []string{"integration_test"},
		}, testStore, resultChan)
		<-resultChan

		if testStore.Exists("integration_test") {
			t.Error("Key should not exist after DEL")
		}
	})
}

func TestExecuteCommandEdgeCases(t *testing.T) {
	testStore := store.NewStore()

	t.Run("Empty command", func(t *testing.T) {
		// This should not panic
		resultChan := make(chan ResultWithError, 1)
		ExecuteCommandAysnc(&parser.Command{
			Name: "",
			Args: []string{},
		}, testStore, resultChan)
		<-resultChan
	})

	t.Run("Unknown command", func(t *testing.T) {
		// This should not panic
		resultChan := make(chan ResultWithError, 1)
		ExecuteCommandAysnc(&parser.Command{
			Name: "UNKNOWN",
			Args: []string{"arg1", "arg2"},
		}, testStore, resultChan)
		<-resultChan
	})

	t.Run("Commands with insufficient arguments", func(t *testing.T) {
		// These should not panic, but may not work as expected
		resultChan := make(chan ResultWithError, 1)
		ExecuteCommandAysnc(&parser.Command{
			Name: "GET",
			Args: []string{}, // No key provided
		}, testStore, resultChan)
		<-resultChan

		resultChan = make(chan ResultWithError, 1)
		ExecuteCommandAysnc(&parser.Command{
			Name: "SET",
			Args: []string{"key"}, // No value provided
		}, testStore, resultChan)
		<-resultChan
	})
}

func TestNumericValueEncoding(t *testing.T) {
	testStore := store.NewStore()

	t.Run("SET numeric value as string should be stored as integer", func(t *testing.T) {
		// Set a numeric value as string
		testStore.SetValue("numkey", "123")
		
		// Verify it's stored as integer
		value := testStore.GetValue("numkey")
		if value == nil {
			t.Error("Expected value to be set, got nil")
		}
		
		// Check if it's an integer
		if intVal, ok := value.(int); !ok {
			t.Errorf("Expected integer value, got %T: %v", value, value)
		} else if intVal != 123 {
			t.Errorf("Expected 123, got %d", intVal)
		}
	})

	t.Run("SET non-numeric string should be stored as string", func(t *testing.T) {
		// Set a non-numeric value
		testStore.SetValue("strkey", "hello")
		
		// Verify it's stored as string
		value := testStore.GetValue("strkey")
		if value == nil {
			t.Error("Expected value to be set, got nil")
		}
		
		// Check if it's a string
		if strVal, ok := value.(string); !ok {
			t.Errorf("Expected string value, got %T: %v", value, value)
		} else if strVal != "hello" {
			t.Errorf("Expected 'hello', got %s", strVal)
		}
	})

	t.Run("SET negative number should be stored as integer", func(t *testing.T) {
		// Set a negative numeric value as string
		testStore.SetValue("negkey", "-456")
		
		// Verify it's stored as integer
		value := testStore.GetValue("negkey")
		if value == nil {
			t.Error("Expected value to be set, got nil")
		}
		
		// Check if it's an integer
		if intVal, ok := value.(int); !ok {
			t.Errorf("Expected integer value, got %T: %v", value, value)
		} else if intVal != -456 {
			t.Errorf("Expected -456, got %d", intVal)
		}
	})

	t.Run("SET zero should be stored as integer", func(t *testing.T) {
		// Set zero as string
		testStore.SetValue("zerokey", "0")
		
		// Verify it's stored as integer
		value := testStore.GetValue("zerokey")
		if value == nil {
			t.Error("Expected value to be set, got nil")
		}
		
		// Check if it's an integer
		if intVal, ok := value.(int); !ok {
			t.Errorf("Expected integer value, got %T: %v", value, value)
		} else if intVal != 0 {
			t.Errorf("Expected 0, got %d", intVal)
		}
	})
}

func TestIncrByWithNumericEncoding(t *testing.T) {
	testStore := store.NewStore()

	t.Run("INCRBY should work with numeric values stored as integers", func(t *testing.T) {
		// First set a numeric value (which should be stored as integer)
		testStore.SetValue("numkey", "10")
		
		// Verify it's stored as integer
		value := testStore.GetValue("numkey")
		if intVal, ok := value.(int); !ok || intVal != 10 {
			t.Errorf("Expected integer 10, got %T: %v", value, value)
		}
		
		// Now test INCRBY
		result, err := testStore.IncreBy("numkey", 5)
		if err != nil {
			t.Errorf("INCRBY failed: %v", err)
		}
		if result != 15 {
			t.Errorf("Expected 15, got %d", result)
		}
		
		// Verify the value was updated
		updatedValue := testStore.GetValue("numkey")
		if intVal, ok := updatedValue.(int); !ok || intVal != 15 {
			t.Errorf("Expected integer 15, got %T: %v", updatedValue, updatedValue)
		}
	})
}

func TestIncrCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  *parser.Command
		setup    func(*store.Store)
		verify   func(*store.Store)
	}{
		{
			name: "INCR command - new key",
			command: &parser.Command{
				Name: "INCR",
				Args: []string{"newkey"},
			},
			verify: func(s *store.Store) {
				value := s.GetValue("newkey")
				if value == nil {
					t.Error("Expected value to be set, got nil")
				}
				if intVal, ok := value.(int); !ok || intVal != 1 {
					t.Errorf("Expected integer 1, got %T: %v", value, value)
				}
			},
		},
		{
			name: "INCR command - existing numeric key",
			command: &parser.Command{
				Name: "INCR",
				Args: []string{"existingkey"},
			},
			setup: func(s *store.Store) {
				s.SetValue("existingkey", "5")
			},
			verify: func(s *store.Store) {
				value := s.GetValue("existingkey")
				if value == nil {
					t.Error("Expected value to be set, got nil")
				}
				if intVal, ok := value.(int); !ok || intVal != 6 {
					t.Errorf("Expected integer 6, got %T: %v", value, value)
				}
			},
		},
		{
			name: "INCR command - insufficient arguments",
			command: &parser.Command{
				Name: "INCR",
				Args: []string{},
			},
			verify: func(s *store.Store) {
				// Should not cause any errors, just print error message
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStore := store.NewStore()
			
			if tt.setup != nil {
				tt.setup(testStore)
			}

			resultChan := make(chan ResultWithError, 1)
			ExecuteCommandAysnc(tt.command, testStore, resultChan)
			<-resultChan

			if tt.verify != nil {
				tt.verify(testStore)
			}
		})
	}
}

func TestDecrCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  *parser.Command
		setup    func(*store.Store)
		verify   func(*store.Store)
	}{
		{
			name: "DECR command - new key",
			command: &parser.Command{
				Name: "DECR",
				Args: []string{"newkey"},
			},
			verify: func(s *store.Store) {
				value := s.GetValue("newkey")
				if value == nil {
					t.Error("Expected value to be set, got nil")
				}
				if intVal, ok := value.(int); !ok || intVal != -1 {
					t.Errorf("Expected integer -1, got %T: %v", value, value)
				}
			},
		},
		{
			name: "DECR command - existing numeric key",
			command: &parser.Command{
				Name: "DECR",
				Args: []string{"existingkey"},
			},
			setup: func(s *store.Store) {
				s.SetValue("existingkey", "10")
			},
			verify: func(s *store.Store) {
				value := s.GetValue("existingkey")
				if value == nil {
					t.Error("Expected value to be set, got nil")
				}
				if intVal, ok := value.(int); !ok || intVal != 9 {
					t.Errorf("Expected integer 9, got %T: %v", value, value)
				}
			},
		},
		{
			name: "DECR command - insufficient arguments",
			command: &parser.Command{
				Name: "DECR",
				Args: []string{},
			},
			verify: func(s *store.Store) {
				// Should not cause any errors, just print error message
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStore := store.NewStore()
			
			if tt.setup != nil {
				tt.setup(testStore)
			}

			resultChan := make(chan ResultWithError, 1)
			ExecuteCommandAysnc(tt.command, testStore, resultChan)
			<-resultChan

			if tt.verify != nil {
				tt.verify(testStore)
			}
		})
	}
}

func TestIncrByCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  *parser.Command
		setup    func(*store.Store)
		verify   func(*store.Store)
	}{
		{
			name: "INCRBY command - new key with positive increment",
			command: &parser.Command{
				Name: "INCRBY",
				Args: []string{"newkey", "5"},
			},
			verify: func(s *store.Store) {
				value := s.GetValue("newkey")
				if value == nil {
					t.Error("Expected value to be set, got nil")
				}
				if intVal, ok := value.(int); !ok || intVal != 5 {
					t.Errorf("Expected integer 5, got %T: %v", value, value)
				}
			},
		},
		{
			name: "INCRBY command - existing key with positive increment",
			command: &parser.Command{
				Name: "INCRBY",
				Args: []string{"existingkey", "3"},
			},
			setup: func(s *store.Store) {
				s.SetValue("existingkey", "7")
			},
			verify: func(s *store.Store) {
				value := s.GetValue("existingkey")
				if value == nil {
					t.Error("Expected value to be set, got nil")
				}
				if intVal, ok := value.(int); !ok || intVal != 10 {
					t.Errorf("Expected integer 10, got %T: %v", value, value)
				}
			},
		},
		{
			name: "INCRBY command - existing key with negative increment",
			command: &parser.Command{
				Name: "INCRBY",
				Args: []string{"existingkey", "-2"},
			},
			setup: func(s *store.Store) {
				s.SetValue("existingkey", "8")
			},
			verify: func(s *store.Store) {
				value := s.GetValue("existingkey")
				if value == nil {
					t.Error("Expected value to be set, got nil")
				}
				if intVal, ok := value.(int); !ok || intVal != 6 {
					t.Errorf("Expected integer 6, got %T: %v", value, value)
				}
			},
		},
		{
			name: "INCRBY command - zero increment",
			command: &parser.Command{
				Name: "INCRBY",
				Args: []string{"existingkey", "0"},
			},
			setup: func(s *store.Store) {
				s.SetValue("existingkey", "5")
			},
			verify: func(s *store.Store) {
				value := s.GetValue("existingkey")
				if value == nil {
					t.Error("Expected value to be set, got nil")
				}
				if intVal, ok := value.(int); !ok || intVal != 5 {
					t.Errorf("Expected integer 5, got %T: %v", value, value)
				}
			},
		},
		{
			name: "INCRBY command - insufficient arguments",
			command: &parser.Command{
				Name: "INCRBY",
				Args: []string{"key"},
			},
			verify: func(s *store.Store) {
				// Should not cause any errors, just print error message
			},
		},
		{
			name: "INCRBY command - invalid increment value",
			command: &parser.Command{
				Name: "INCRBY",
				Args: []string{"key", "invalid"},
			},
			verify: func(s *store.Store) {
				// Should not cause any errors, just print error message
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStore := store.NewStore()
			
			if tt.setup != nil {
				tt.setup(testStore)
			}

			resultChan := make(chan ResultWithError, 1)
			ExecuteCommandAysnc(tt.command, testStore, resultChan)
			<-resultChan

			if tt.verify != nil {
				tt.verify(testStore)
			}
		})
	}
}

func TestDecrByCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  *parser.Command
		setup    func(*store.Store)
		verify   func(*store.Store)
	}{
		{
			name: "DECRBY command - new key with positive decrement",
			command: &parser.Command{
				Name: "DECRBY",
				Args: []string{"newkey", "3"},
			},
			verify: func(s *store.Store) {
				value := s.GetValue("newkey")
				if value == nil {
					t.Error("Expected value to be set, got nil")
				}
				if intVal, ok := value.(int); !ok || intVal != -3 {
					t.Errorf("Expected integer -3, got %T: %v", value, value)
				}
			},
		},
		{
			name: "DECRBY command - existing key with positive decrement",
			command: &parser.Command{
				Name: "DECRBY",
				Args: []string{"existingkey", "2"},
			},
			setup: func(s *store.Store) {
				s.SetValue("existingkey", "10")
			},
			verify: func(s *store.Store) {
				value := s.GetValue("existingkey")
				if value == nil {
					t.Error("Expected value to be set, got nil")
				}
				if intVal, ok := value.(int); !ok || intVal != 8 {
					t.Errorf("Expected integer 8, got %T: %v", value, value)
				}
			},
		},
		{
			name: "DECRBY command - existing key with negative decrement (increment)",
			command: &parser.Command{
				Name: "DECRBY",
				Args: []string{"existingkey", "-4"},
			},
			setup: func(s *store.Store) {
				s.SetValue("existingkey", "5")
			},
			verify: func(s *store.Store) {
				value := s.GetValue("existingkey")
				if value == nil {
					t.Error("Expected value to be set, got nil")
				}
				if intVal, ok := value.(int); !ok || intVal != 9 {
					t.Errorf("Expected integer 9, got %T: %v", value, value)
				}
			},
		},
		{
			name: "DECRBY command - zero decrement",
			command: &parser.Command{
				Name: "DECRBY",
				Args: []string{"existingkey", "0"},
			},
			setup: func(s *store.Store) {
				s.SetValue("existingkey", "7")
			},
			verify: func(s *store.Store) {
				value := s.GetValue("existingkey")
				if value == nil {
					t.Error("Expected value to be set, got nil")
				}
				if intVal, ok := value.(int); !ok || intVal != 7 {
					t.Errorf("Expected integer 7, got %T: %v", value, value)
				}
			},
		},
		{
			name: "DECRBY command - insufficient arguments",
			command: &parser.Command{
				Name: "DECRBY",
				Args: []string{"key"},
			},
			verify: func(s *store.Store) {
				// Should not cause any errors, just print error message
			},
		},
		{
			name: "DECRBY command - invalid decrement value",
			command: &parser.Command{
				Name: "DECRBY",
				Args: []string{"key", "invalid"},
			},
			verify: func(s *store.Store) {
				// Should not cause any errors, just print error message
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStore := store.NewStore()
			
			if tt.setup != nil {
				tt.setup(testStore)
			}

			resultChan := make(chan ResultWithError, 1)
			ExecuteCommandAysnc(tt.command, testStore, resultChan)
			<-resultChan

			if tt.verify != nil {
				tt.verify(testStore)
			}
		})
	}
}