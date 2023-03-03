package db

import "walltrack/model"

func SelectUserByEmail(email string) (model.User, bool, error) {
	ctx, cancel := getDbContext()
	defer cancel()

	var user model.User
	err := db.Collection(userCollection).FindOne(ctx, model.User{
		Email: email,
	}).Decode(&user)

	if err != nil {
		if IsNotFoundError(err) {
			return user, false, nil
		}
		return user, false, err
	}

	return user, true, err
}
