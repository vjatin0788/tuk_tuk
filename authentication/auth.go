package authentication

import (
	"context"
	"log"

	"github.com/TukTuk/errs"
	"github.com/TukTuk/model"
)

func (auth *TukTukAuth) Authentication(ctx context.Context, user, driver bool, authToken string) (AuthUser, error) {
	var (
		err         error
		cust, duser User
	)

	defaultAuth := AuthUser{}

	if authToken == "" {
		log.Println("[Authentication][Error] Empty Auth Token ", authToken)
		return defaultAuth, errs.Err("TT_AU_401")
	}

	if user {
		cust, err = auth.getCustomerData(ctx, authToken)
		if err != nil {
			return defaultAuth, err
		}

		if cust.Token != authToken {
			log.Println("[Authentication][Error] Token mismatch ", cust.Token)
			return defaultAuth, errs.Err("TT_AU_400")
		}
	}

	if driver {
		duser, err = auth.getDriverData(ctx, authToken)
		if err != nil {
			return defaultAuth, err
		}

		if duser.Token != authToken {
			log.Println("[Authentication][Error]driver Token mismatch ", cust.Token)
			return defaultAuth, errs.Err("TT_AU_400")
		}
	}

	defaultAuth.Customer = cust
	defaultAuth.Driver = duser

	return defaultAuth, err
}

func (auth *TukTukAuth) getCustomerData(ctx context.Context, authToken string) (User, error) {
	var (
		user User
		err  error
	)

	userData, err := model.TukTuk.GetCustomerByToken(ctx, authToken)
	if err != nil {
		log.Println("[doCustomerAuth][Error] DB error", err)
		return user, err
	}

	user = User{
		Id:    userData.CustomerId,
		Token: userData.Token,
	}

	return user, err
}

func (auth *TukTukAuth) getDriverData(ctx context.Context, authToken string) (User, error) {
	var (
		user User
		err  error
	)

	userData, err := model.TukTuk.GetDriverByToken(ctx, authToken)
	if err != nil {
		log.Println("[getDriverData][Error] DB error", err)
		return user, err
	}

	user = User{
		Id:    userData.Userid,
		Token: userData.Token,
	}
	return user, err
}
