package auth

import (
	"payup/app-error"
	"payup/database"
	"payup/models"
)

// GetToken - get's a token with a specific authorization code
func GetToken(accessToken string) (models.Token, error) {
	var token models.Token
	if database.DBCon.Where("access_token = ?", accessToken).Find(&token).RecordNotFound() {
		
		invalidAuth := appError.InvalidParams
		invalidAuth.Detail = "Invalid Authorization header"
		return token, invalidAuth
	}
	
	return token, nil
}
