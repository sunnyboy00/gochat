package protos

// 包括握手,心跳,登录

import (
	"encoding/binary"
	"errors"
	"fmt"
)

import (
	"misc/zpack"
	"share"
	"types"
)

func handle_shake(user *types.User, msg []byte) (ack []byte, err error) {
	// msg[0]: msgType
	subType := msg[1]

	if subType == 0 {
		ack = zpack.Pack('>', []interface{}{byte(0), byte(1), user.Coder.CryptKey})
	} else if subType == 2 {
		user.Coder.Shaked = true
		fmt.Println("Shaked:")
	} else {
		err = errors.New("handle_shake: unknown subType")
	}

	return
}

func handle_nop(user *types.User, msg []byte) (ack []byte, err error) {
	// TODO: 发送到hub服务器,维持用户在线信息
	return
}

func handle_login(user *types.User, msg []byte) (ack []byte, err error) {
	if !user.Coder.Shaked {
		err = errors.New("handle_login: not shaked")
		return
	}

	// msg[0]: msgType

	uid := binary.BigEndian.Uint32(msg[1:5])
	password := string(msg[5:])

	if true { // TODO: 向hub服务器发送用户名密码请求登录(http接口?)
		fmt.Println("Login:", uid, password)
		user.UID = uid
		user.Password = password
		user.Online = true
		user.Logined = true
		share.Clients.Set(user.UID, user)
		// TODO: 通知hub服务器该用户登录
		ack = zpack.Pack('>', []interface{}{byte(2), byte(0)})
	} else {
		//err = errors.New("login failed")
		ack = zpack.Pack('>', []interface{}{byte(2), byte(1)})
	}

	return
}
