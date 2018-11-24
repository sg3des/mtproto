package mtproto

import (
	"errors"
	"fmt"
	"log"
)

func (m *MTProto) Auth_SendCode(phonenumber string) (string, error) {
	resp, err := m.send(TL_auth_sendCode{
		Flags:          1,
		Current_number: TL_boolFalse{},
		Phone_number:   phonenumber,
		Api_id:         m.appId,
		Api_hash:       m.appHash,
	})
	if err != nil {
		return "", err
	}
	// log.Printf("%+v", resp)

	return resp.(TL_auth_sentCode).Phone_code_hash, nil
}

func (m *MTProto) Auth_SignUp(phonenumber, hash, code, fname, lname string) (TL_auth_authorization, error) {
	resp, err := m.send(TL_auth_signUp{
		Phone_number:    phonenumber,
		Phone_code_hash: hash,
		Phone_code:      code,
		First_name:      fname,
		Last_name:       lname,
	})
	if err != nil {
		return TL_auth_authorization{}, err
	}

	auth, ok := resp.(TL_auth_authorization)
	if !ok {
		return TL_auth_authorization{}, fmt.Errorf("RPC: %#v", resp)
	}

	// userSelf := auth.User.(TL_user)
	// log.Printf("Signed in: id %d name <%s %s>\n", userSelf.Id, userSelf.First_name, userSelf.Last_name)
	return auth, nil
}

func (m *MTProto) Auth_SignIn(phonenumber string, hash, code string) (TL_auth_authorization, error) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_auth_signIn{phonenumber, hash, code},
		resp,
	}
	x := <-resp
	auth, ok := x.(TL_auth_authorization)
	if !ok {
		return TL_auth_authorization{}, fmt.Errorf("RPC: %#v", x)
	}
	// userSelf := auth.User.(TL_user)
	// log.Printf("Signed in: id %d name <%s %s>\n", userSelf.Id, userSelf.First_name, userSelf.Last_name)
	return auth, nil
}

func (m *MTProto) Auth_CheckPhone(phonenumber string) bool {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_auth_checkPhone{
			"989121228718",
		},
		resp,
	}
	x := <-resp
	if v, ok := x.(TL_auth_checkedPhone); ok {
		if toBool(v) {
			return true
		}
	}
	return false
}

func (m *MTProto) Users_GetFullSelf() (*User, error) {
	return m.users_getFullUsers(TL_inputUserSelf{})
}

func (m *MTProto) users_getFullUsers(id TL) (*User, error) {
	resp, err := m.send(TL_users_getFullUser{Id: id})
	if err != nil {
		return nil, err
	}
	user, ok := resp.(TL_userFull)
	if !ok {
		log.Printf("%#v", resp)
		return nil, errors.New("unexpected response")
	}

	return NewUser(user.User), nil
}
