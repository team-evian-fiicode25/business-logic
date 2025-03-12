package user

import (
	"fmt"

	"github.com/team-evian-fiicode25/business-logic/data"
	"github.com/team-evian-fiicode25/business-logic/database"
)

func Create(authData *data.AuthData) (*data.User, error) {
    if authData == nil {
        return nil, fmt.Errorf("Receieved nil reference to argument (authData)")
    }

    user := &data.User{AuthData: *authData}

    db := database.GetDB()

    err := db.Create(&user).Error

    if err != nil {
        return nil, err
    }

    return user, nil
}
