package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

func TestExecuteCommand(t *testing.T) {
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

			// Execute the command
			ExecuteCommand(tt.command, testStore)

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
		ExecuteCommand(&parser.Command{
			Name: "SET",
			Args: []string{"integration_test", "integration_value"},
		}, testStore)
		
		if !testStore.Exists("integration_test") {
			t.Error("Key should exist after SET")
		}

		// GET
		ExecuteCommand(&parser.Command{
			Name: "GET",
			Args: []string{"integration_test"},
		}, testStore)

		// EXISTS
		ExecuteCommand(&parser.Command{
			Name: "EXISTS",
			Args: []string{"integration_test"},
		}, testStore)

		// EXPIRE
		ExecuteCommand(&parser.Command{
			Name: "EXPIRE",
			Args: []string{"integration_test", "7200"}, // 2 hours
		}, testStore)

		// TTL
		ExecuteCommand(&parser.Command{
			Name: "TTL",
			Args: []string{"integration_test"},
		}, testStore)

		// DEL
		ExecuteCommand(&parser.Command{
			Name: "DEL",
			Args: []string{"integration_test"},
		}, testStore)

		// EXISTS (should return false now)
		ExecuteCommand(&parser.Command{
			Name: "EXISTS",
			Args: []string{"integration_test"},
		}, testStore)

		if testStore.Exists("integration_test") {
			t.Error("Key should not exist after DEL")
		}
	})
}

func TestExecuteCommandEdgeCases(t *testing.T) {
	testStore := store.NewStore()

	t.Run("Empty command", func(t *testing.T) {
		// This should not panic
		ExecuteCommand(&parser.Command{
			Name: "",
			Args: []string{},
		}, testStore)
	})

	t.Run("Unknown command", func(t *testing.T) {
		// This should not panic
		ExecuteCommand(&parser.Command{
			Name: "UNKNOWN",
			Args: []string{"arg1", "arg2"},
		}, testStore)
	})

	t.Run("Commands with insufficient arguments", func(t *testing.T) {
		// These should not panic, but may not work as expected
		ExecuteCommand(&parser.Command{
			Name: "GET",
			Args: []string{}, // No key provided
		}, testStore)

		ExecuteCommand(&parser.Command{
			Name: "SET",
			Args: []string{"key"}, // No value provided
		}, testStore)
	})
}