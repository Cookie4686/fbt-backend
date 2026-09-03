package test

import (
	"context"
	"testing"

	"fbt/backend/internal/domain/auth/features/credentials"
	"fbt/backend/internal/domain/bookkeeping/features/account"
	"fbt/backend/internal/domain/bookkeeping/model"
	bookkeepingService "fbt/backend/internal/domain/bookkeeping/service"
	"fbt/backend/internal/interceptor"
	"fbt/backend/internal/util"

	authv1 "fbt/backend/gen/proto/go/auth/v1"
	bookkeepingv1 "fbt/backend/gen/proto/go/bookkeeping/v1"

	authService "fbt/backend/internal/domain/auth/service"

	"github.com/stretchr/testify/require"
)

type TestUtil struct {
	AuthService           authService.Service
	CredentialsController credentials.Controller

	BookkeepingService bookkeepingService.Service
	AccountController  account.Controller

	Init Init
}

type Init struct {
	User    UserInit
	Account AccountInit
}

type UserInit struct {
	Username string
	Email    string
	Password string
}

type AccountInit struct {
	Cash model.Account
	Loan model.Account
}

func NewTestUtil(d *util.Dependency) *TestUtil {
	authService := authService.NewService(d)
	bookkeepingService := bookkeepingService.NewService(d)

	return &TestUtil{
		AuthService:           authService,
		CredentialsController: *credentials.NewController(authService, credentials.NewRepo(d.DB)),

		AccountController: *account.NewController(bookkeepingService, account.NewRepo(d.DB)),

		Init: Init{
			User: UserInit{
				Username: "test",
				Email:    "test@email.com",
				Password: "12345678",
			},

			Account: AccountInit{
				Cash: model.Account{
					Code:    "101",
					Name:    "Cash",
					IsDebit: true,
				},
				Loan: model.Account{
					Code:    "201",
					Name:    "Loan",
					IsDebit: false,
				},
			},
		},
	}
}

func (s *TestUtil) NewAuthContext(ctx context.Context, sessionID string) context.Context {
	auth, err := s.AuthService.Validate(ctx, sessionID)
	if err != nil {
		return ctx
	}

	return interceptor.NewAuthContext(ctx, auth)
}

func (s *TestUtil) SetupUser(t *testing.T) *authv1.Session {
	res, err := s.CredentialsController.Register(t.Context(), &authv1.CredentialServiceRegisterRequest{
		Username: s.Init.User.Username,
		Email:    s.Init.User.Email,
		Password: s.Init.User.Password,
	})
	require.NoError(t, err)

	return res.Session
}

func (s *TestUtil) SetupAccount(t *testing.T, sessionID string) AccountInit {
	accs := []*model.Account{
		&s.Init.Account.Cash,
		&s.Init.Account.Loan,
	}

	for _, a := range accs {
		ctx := s.NewAuthContext(t.Context(), sessionID)

		res, err := s.AccountController.Create(ctx, &bookkeepingv1.AccountServiceCreateRequest{
			Code:    a.Code,
			Name:    a.Name,
			IsDebit: a.IsDebit,
		})
		require.NoError(t, err)

		a.ID = res.Account.Id
	}

	return s.Init.Account
}
