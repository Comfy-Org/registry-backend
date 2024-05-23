package mapper

import (
	"registry-backend/drip"
	"registry-backend/ent"
	"strings"
)

func DbToApiPersonalAccessToken(dbToken *ent.PersonalAccessToken) *drip.PersonalAccessToken {
	maskedToken := maskToken(dbToken.Token)
	return &drip.PersonalAccessToken{
		Id:          &dbToken.ID,
		Name:        &dbToken.Name,
		CreatedAt:   &dbToken.CreateTime,
		Description: &dbToken.Description,
		Token:       &maskedToken,
	}
}

func maskToken(token string) string {
	tokenLength := len(token)
	if tokenLength <= 8 {
		return strings.Repeat("*", tokenLength)
	}
	return token[:4] + strings.Repeat("*", 3)
}
