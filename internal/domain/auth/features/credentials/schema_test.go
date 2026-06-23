package credentials_test

/*
var validate = validator.New(validator.WithRequiredStructEnabled())

 func TestRegisterSchema(t *testing.T) {
	t.Run("Missing fields", func(t *testing.T) {
		for _, invalidPayload := range []credentials.RegisterPayload{
			{Username: mock.ValidUser.Username},
			{Password: mock.ValidUser.Password},
			{Email: mock.ValidUser.Email},
		} {
			err := validate.Struct(invalidPayload)
			assert.Error(t, err)
		}
	})

	t.Run("Username Constraint", func(t *testing.T) {
		for _, invalidUsername := range []string{
			"g", "qe", "f--", "qr041$",
		} {
			err := validate.Struct(credentials.RegisterPayload{
				Username: invalidUsername,
				Password: mock.ValidUser.Password,
				Email:    mock.ValidUser.Email,
			})
			assert.Error(t, err)
		}
	})

	t.Run("Password Constraint", func(t *testing.T) {
		for _, invalidPassword := range []string{
			"g", "qe", "f--", "qr041$",
		} {
			err := validate.Struct(credentials.RegisterPayload{
				Username: mock.ValidUser.Username,
				Password: invalidPassword,
				Email:    mock.ValidUser.Email,
			})
			assert.Error(t, err)
		}
	})
}
*/
