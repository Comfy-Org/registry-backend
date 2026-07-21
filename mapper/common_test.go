package mapper

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseUUIDParam(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	parsedUUID := uuid.MustParse(validUUID)

	tests := []struct {
		name    string
		input   *string
		wantID  uuid.UUID
		wantErr bool
	}{
		{
			name:    "nil pointer returns uuid.Nil without error",
			input:   nil,
			wantID:  uuid.Nil,
			wantErr: false,
		},
		{
			name:    "valid UUID string is parsed correctly",
			input:   &validUUID,
			wantID:  parsedUUID,
			wantErr: false,
		},
		{
			name:    "garbage string returns error",
			input:   strPtr("not-a-uuid"),
			wantID:  uuid.Nil,
			wantErr: true,
		},
		{
			name:    "empty string returns error",
			input:   strPtr(""),
			wantID:  uuid.Nil,
			wantErr: true,
		},
		{
			name:    "truncated UUID returns error",
			input:   strPtr("550e8400-e29b-41d4-a716"),
			wantID:  uuid.Nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUUIDParam(tt.input)
			if tt.wantErr {
				require.Error(t, err, "expected an error but got none")
				assert.Equal(t, uuid.Nil, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantID, got)
			}
		})
	}
}

func strPtr(s string) *string { return &s }
