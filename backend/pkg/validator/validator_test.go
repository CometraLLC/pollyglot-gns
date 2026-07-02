package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sampleDTO struct {
	Name  string `validate:"required,min=2,max=10"`
	Email string `validate:"required,email"`
	Kind  string `validate:"omitempty,oneof=basic pro"`
}

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		name    string
		input   sampleDTO
		wantErr []string
	}{
		{
			name:    "valid struct passes",
			input:   sampleDTO{Name: "Marc", Email: "marc@example.com", Kind: "pro"},
			wantErr: nil,
		},
		{
			name:    "missing required fields",
			input:   sampleDTO{},
			wantErr: []string{"name is required", "email is required"},
		},
		{
			name:    "invalid email",
			input:   sampleDTO{Name: "Marc", Email: "not-an-email"},
			wantErr: []string{"email must be a valid email"},
		},
		{
			name:    "name below minimum length",
			input:   sampleDTO{Name: "M", Email: "marc@example.com"},
			wantErr: []string{"name must be at least 2 characters"},
		},
		{
			name:    "name above maximum length",
			input:   sampleDTO{Name: "MarcMarcMarc", Email: "marc@example.com"},
			wantErr: []string{"name must be at most 10 characters"},
		},
		{
			name:    "value outside oneof set",
			input:   sampleDTO{Name: "Marc", Email: "marc@example.com", Kind: "enterprise"},
			wantErr: []string{"kind must be one of: basic pro"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.input)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				return
			}

			require.Error(t, err)
			var vErr *ValidationError
			require.ErrorAs(t, err, &vErr)
			assert.Equal(t, tt.wantErr, vErr.Errors)
		})
	}
}
