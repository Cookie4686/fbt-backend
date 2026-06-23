package user_test

// func TestUser(t *testing.T) {
// 	ctx, _, conn := test.NewTestConnection(t, 1234)

// 	session := test.SetupUser(t, ctx, conn)

// 	client := pb.NewUserClient(conn)

// 	t.Run("GetByUsername", func(t *testing.T) {
// 		ctx := metadata.AppendToOutgoingContext(ctx, "session_id", session.Id)

// 		res, err := client.GetByUsername(ctx, &pb.GetByUsernameRequest{})
// 		require.NoError(t, err)

// 		assert.Equal(t, session.UserID, res.User.Id)
// 	})
// }
