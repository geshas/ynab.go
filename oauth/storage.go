package oauth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// TokenStorage defines the interface for token persistence
type TokenStorage interface {
	// SaveToken persists a token
	SaveToken(token *Token) error

	// LoadToken retrieves a token
	LoadToken() (*Token, error)

	// ClearToken removes the stored token
	ClearToken() error

	// HasToken checks if a token is stored
	HasToken() bool
}

// MemoryStorage implements in-memory token storage (not persistent)
type MemoryStorage struct {
	mu    sync.RWMutex
	token *Token
}

// NewMemoryStorage creates a new in-memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

// SaveToken saves the token in memory
func (s *MemoryStorage) SaveToken(token *Token) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = token
	return nil
}

// LoadToken loads the token from memory
func (s *MemoryStorage) LoadToken() (*Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.token == nil {
		return nil, fmt.Errorf("no token stored")
	}
	return s.token, nil
}

// ClearToken clears the token from memory
func (s *MemoryStorage) ClearToken() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = nil
	return nil
}

// HasToken checks if a token is stored in memory
func (s *MemoryStorage) HasToken() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.token != nil
}

// FileStorage implements file-based token storage
type FileStorage struct {
	filePath string
	fileMode os.FileMode
}

// NewFileStorage creates a new file-based storage
func NewFileStorage(filePath string) *FileStorage {
	return &FileStorage{
		filePath: filePath,
		fileMode: 0600, // Read/write for owner only
	}
}

// WithFileMode sets the file permissions for the token file
func (s *FileStorage) WithFileMode(mode os.FileMode) *FileStorage {
	s.fileMode = mode
	return s
}

// SaveToken saves the token to a file
func (s *FileStorage) SaveToken(token *Token) error {
	if token == nil {
		return fmt.Errorf("token cannot be nil")
	}

	// Ensure directory exists
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Serialize token
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// Write to file with secure permissions
	if err := os.WriteFile(s.filePath, data, s.fileMode); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// LoadToken loads the token from a file
func (s *FileStorage) LoadToken() (*Token, error) {
	// Check if file exists
	if !s.HasToken() {
		return nil, fmt.Errorf("no token file found")
	}

	// Read file
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	// Deserialize token
	var token Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &token, nil
}

// ClearToken removes the token file
func (s *FileStorage) ClearToken() error {
	if !s.HasToken() {
		return nil // Already cleared
	}

	if err := os.Remove(s.filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove token file: %w", err)
	}

	return nil
}

// HasToken checks if the token file exists
func (s *FileStorage) HasToken() bool {
	_, err := os.Stat(s.filePath)
	return err == nil
}

// GetFilePath returns the file path
func (s *FileStorage) GetFilePath() string {
	return s.filePath
}

// DefaultTokenPath returns the default token file path
func DefaultTokenPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".ynab_token.json"
	}

	return filepath.Join(homeDir, ".config", "ynab", "token.json")
}

// EncryptedFileStorage implements encrypted file-based token storage
type EncryptedFileStorage struct {
	*FileStorage
	key []byte
}

// NewEncryptedFileStorage creates a new encrypted file-based storage.
// key must be 16, 24, or 32 bytes (for AES-128, AES-192, or AES-256).
func NewEncryptedFileStorage(filePath string, key []byte) (*EncryptedFileStorage, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("encryption key must be 16, 24, or 32 bytes (got %d)", len(key))
	}
	return &EncryptedFileStorage{
		FileStorage: NewFileStorage(filePath),
		key:         key,
	}, nil
}

// SaveToken saves the encrypted token to a file
func (s *EncryptedFileStorage) SaveToken(token *Token) error {
	if token == nil {
		return fmt.Errorf("token cannot be nil")
	}

	// Serialize token
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// Encrypt data using AES-GCM
	encrypted, err := s.encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt token: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write encrypted data to file
	if err := os.WriteFile(s.filePath, encrypted, s.fileMode); err != nil {
		return fmt.Errorf("failed to write encrypted token file: %w", err)
	}

	return nil
}

// LoadToken loads and decrypts the token from a file
func (s *EncryptedFileStorage) LoadToken() (*Token, error) {
	// Check if file exists
	if !s.HasToken() {
		return nil, fmt.Errorf("no encrypted token file found")
	}

	// Read encrypted file
	encrypted, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read encrypted token file: %w", err)
	}

	// Decrypt data
	data, err := s.decrypt(encrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token file: %w", err)
	}

	// Deserialize token
	var token Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal decrypted token: %w", err)
	}

	return &token, nil
}

// encrypt encrypts data using AES-GCM authenticated encryption.
// The output format is: nonce (12 bytes) || ciphertext+tag.
func (s *EncryptedFileStorage) encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// decrypt decrypts data encrypted by encrypt.
func (s *EncryptedFileStorage) decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token: %w", err)
	}

	return plaintext, nil
}

// ChainedStorage implements a chain of storage backends with fallback
type ChainedStorage struct {
	storages []TokenStorage
}

// NewChainedStorage creates a new chained storage
func NewChainedStorage(storages ...TokenStorage) *ChainedStorage {
	return &ChainedStorage{
		storages: storages,
	}
}

// SaveToken saves the token to all storages in the chain
func (s *ChainedStorage) SaveToken(token *Token) error {
	var firstError error

	for _, storage := range s.storages {
		if err := storage.SaveToken(token); err != nil && firstError == nil {
			firstError = err
		}
	}

	return firstError
}

// LoadToken loads the token from the first available storage
func (s *ChainedStorage) LoadToken() (*Token, error) {
	for _, storage := range s.storages {
		if storage.HasToken() {
			token, err := storage.LoadToken()
			if err == nil {
				return token, nil
			}
		}
	}

	return nil, fmt.Errorf("no token found in any storage")
}

// ClearToken clears the token from all storages
func (s *ChainedStorage) ClearToken() error {
	var firstError error

	for _, storage := range s.storages {
		if err := storage.ClearToken(); err != nil && firstError == nil {
			firstError = err
		}
	}

	return firstError
}

// HasToken checks if any storage has a token
func (s *ChainedStorage) HasToken() bool {
	for _, storage := range s.storages {
		if storage.HasToken() {
			return true
		}
	}
	return false
}

// StorageOptions provides configuration for creating storage instances
type StorageOptions struct {
	Type       string // "memory", "file", "encrypted"
	FilePath   string
	FileMode   os.FileMode
	EncryptKey []byte
}

// NewStorage creates a new storage instance based on options
func NewStorage(opts StorageOptions) (TokenStorage, error) {
	switch opts.Type {
	case "memory":
		return NewMemoryStorage(), nil

	case "file":
		path := opts.FilePath
		if path == "" {
			path = DefaultTokenPath()
		}

		storage := NewFileStorage(path)
		if opts.FileMode != 0 {
			storage.WithFileMode(opts.FileMode)
		}
		return storage, nil

	case "encrypted":
		path := opts.FilePath
		if path == "" {
			path = DefaultTokenPath()
		}

		if len(opts.EncryptKey) == 0 {
			return nil, fmt.Errorf("encryption key is required for encrypted storage")
		}

		return NewEncryptedFileStorage(path, opts.EncryptKey)

	default:
		return nil, fmt.Errorf("unknown storage type: %s", opts.Type)
	}
}
