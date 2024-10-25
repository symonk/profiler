package profiler

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateProfileFile takes the user defined folder (or working dir) if omitted
// and attempts to make the full folder tree. If the folder creation fails, a
// temp folder is created and the file is written to that location.
// File names are currently not customisable and are provided by the caller
// based on the profile mode selected.
func CreateProfileFile(folder string, name string) (*os.File, error) {
	if err := os.MkdirAll(folder, 0777); err != nil {
		// User provided path failed, use a globally unique
		// temp dir
		folder, err = os.MkdirTemp(os.TempDir(), "profiler")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp folder: %w", err)
		}
	}
	joined := filepath.Join(folder, name)
	path, err := os.Create(joined)
	if err != nil {
		return nil, fmt.Errorf("failed to create profile file: %w", err)
	}
	return path, nil
}
