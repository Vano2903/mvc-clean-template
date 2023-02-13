package model

// This struct is the pure data structure for the application,
// You are free to add any tag you want or leave it emtpy
// if you leave it emtpy your http handlers can have a copy of the model with the tags needed
// for example you could omit certain fields to an user that is not an admin while showing everything to an admin
// or you could decide to define them here and use them in the packages.
// You choose.
type (
	User struct {
		ID                        int    //`json:"id"`
		FirstName                 string //`json:"first_name"`
		LastName                  string //`json:"last_name"`
		Pfp                       string //`json:"pfp"`
		Email                     string //`json:"email"`
		Password                  string //`json:"password"`
		Role                      string //`json:"role"`
		IsBanned                  bool   //`json:"is_banned"`
		SpecialMagicalSecretField string //`json:"special_magical_secret_field"`
	}
)

const (
	RoleAdmin       = "admin"
	RoleUser        = "user"
	RoleUnupdatable = "unupdatable"
)
