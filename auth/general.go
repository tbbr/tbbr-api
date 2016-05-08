package auth

import (
	"strings"

	"github.com/tbbr/tbbr-api/app-error"
	"github.com/tbbr/tbbr-api/database"
	"github.com/tbbr/tbbr-api/models"
)

// GetToken - get's a token with a specific authorization code
func GetToken(authorization string) (models.Token, error) {
	var token models.Token

	if !strings.HasPrefix(authorization, "Bearer") {
		return token, appError.InvalidAuthorization
	}

	accessToken := strings.SplitAfter(authorization, "Bearer ")[1]

	if database.DBCon.Where("access_token = ?", accessToken).Find(&token).RecordNotFound() {
		return token, appError.InvalidAuthorization
	}

	return token, nil
}
