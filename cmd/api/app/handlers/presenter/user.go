package presenter

type jsonUser struct {
	Email string `json:"email"`
}

func User(email string) *jsonUser {
	toReturn := &jsonUser{
		Email: email,
	}

	return toReturn
}
