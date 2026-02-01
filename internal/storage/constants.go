package storage

import "errors"

// Page constants
const (
	PageSize       = 4096
	PageHeaderSize = 100
	MaxRecordSize  = PageSize - PageHeaderSize
)

// Buffer pool constants
const (
	DefaultBufferPoolSize = 1000
)

// Custom errors for MaverickDB storage layer
var (
	ErrPageNotFound = errors.New("storage: page not found")

	ErrBufferPoolFull = errors.New("storage: buffer pool is full, no available frames")

	ErrRecordTooLarge = errors.New("storage: record size exceeds maximum allowed size")

	ErrInvalidPageID = errors.New("storage: invalid page ID")

	ErrPageCorrupted = errors.New("storage: page data is corrupted")

	ErrDiskFull = errors.New("storage: no space left on disk")

	ErrFileNotOpen = errors.New("storage: database file is not open")
)

// StorageError wraps errors with additional context
type StorageError struct {
	Op     string // Operation that failed (e.g., "read", "write", "flush")
	PageID int64  // Page ID involved, -1 if not applicable
	Err    error  // Underlying error
}

func (e *StorageError) Error() string {
	if e.PageID >= 0 {
		return "storage: " + e.Op + " page " + formatPageID(e.PageID) + ": " + e.Err.Error()
	}
	return "storage: " + e.Op + ": " + e.Err.Error()
}

func (e *StorageError) Unwrap() error {
	return e.Err
}

// NewStorageError creates a new StorageError
func NewStorageError(op string, pageID int64, err error) *StorageError {
	return &StorageError{
		Op:     op,
		PageID: pageID,
		Err:    err,
	}
}

// formatPageID converts page ID to string without importing strconv
func formatPageID(id int64) string {
	if id == 0 {
		return "0"
	}

	var result []byte
	negative := id < 0
	if negative {
		id = -id
	}

	for id > 0 {
		result = append([]byte{byte('0' + id%10)}, result...)
		id /= 10
	}

	if negative {
		result = append([]byte{'-'}, result...)
	}
	return string(result)
}
