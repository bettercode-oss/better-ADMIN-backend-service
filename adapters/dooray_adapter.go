package adapters

import (
	"better-admin-backend-service/config"
	"better-admin-backend-service/domain"
	"better-admin-backend-service/dtos"
	"fmt"
	"github.com/bettercode-oss/rest"
	"github.com/go-ldap/ldap/v3"
	"github.com/pkg/errors"
)

type DoorayAdapter struct {
}

func (DoorayAdapter) Authenticate(doorayDomain, token, signId, password string) (dtos.DoorayMember, error) {
	ldapConn, err := ldap.DialURL(config.Config.Dooray.LdapDialUrl)
	if err != nil {
		return dtos.DoorayMember{}, errors.Wrap(err, "ldap conn error")
	}

	defer ldapConn.Close()

	if err := ldapConn.Bind(fmt.Sprint(fmt.Sprintf("%s\\", doorayDomain), signId), password); err != nil {
		if ldap.IsErrorWithCode(err, ldap.ErrorNetwork) || err.Error() == "ldap: connection timed out" {
			return dtos.DoorayMember{}, errors.Wrap(err, "ldap network error")
		}

		return dtos.DoorayMember{}, domain.ErrAuthentication
	}

	result := map[string]interface{}{}

	client := rest.Client{}
	err = client.
		Request().
		SetHeader("Authorization", fmt.Sprintf("dooray-api %s", token)).
		SetResult(&result).
		Get(fmt.Sprintf("https://api.dooray.com/common/v1/members?userCode=%s", signId))

	if err != nil {
		return dtos.DoorayMember{}, errors.Wrap(err, "find dooray member error")
	}

	resultHeader := result["header"].(map[string]interface{})
	if resultHeader["resultCode"].(float64) == 0 && resultHeader["isSuccessful"].(bool) == true && result["totalCount"].(float64) == 1 {
		//// TODO 프로필 이미지 가져오기
		//// 화면 로그인 후에 해당 세션을 가지고 이미지를 조회한다.
		//// https://bettercode.dooray.com/profile-image/1879346658407346013 / 두레이 ID
		//// 그리고 해당 이미지를 DB Blob 로 저장 한다.

		user := result["result"].([]interface{})[0].(map[string]interface{})
		return dtos.DoorayMember{
			Id:                   user["id"].(string),
			UserCode:             signId,
			Name:                 user["name"].(string),
			ExternalEmailAddress: user["externalEmailAddress"].(string),
		}, nil
	}

	return dtos.DoorayMember{}, domain.ErrAuthentication
}
