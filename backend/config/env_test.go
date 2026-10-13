package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	tempDir := t.TempDir()
	validEnvPath := filepath.Join(tempDir, ".env.valid")
	os.WriteFile(validEnvPath, []byte("MOCK_KEY=MOCK_VALUE\n# COMMENT\nEMPTY=\n"), 0644)

	invalidEnvPath := filepath.Join(tempDir, ".env.invalid")
	// Creating an unreadable directory to force an os error other than ErrNotExist
	os.Mkdir(invalidEnvPath, 0000)

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func()
		verify  func(t *testing.T)
	}{
		{
			name:    "file does not exist (returns nil)",
			args:    args{path: filepath.Join(tempDir, ".env.nonexistent")},
			wantErr: false,
		},
		{
			name:    "valid env file parses and sets vars",
			args:    args{path: validEnvPath},
			wantErr: false,
			setup: func() {
				os.Unsetenv("MOCK_KEY")
			},
			verify: func(t *testing.T) {
				if got := os.Getenv("MOCK_KEY"); got != "MOCK_VALUE" {
					t.Errorf("LoadEnv() failed to set env var; got %v, want %v", got, "MOCK_VALUE")
				}
			},
		},
		{
			name:    "file exists but is unreadable",
			args:    args{path: invalidEnvPath},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if err := LoadEnv(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("LoadEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.verify != nil {
				tt.verify(t)
			}
		})
	}
}
