package authorization

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/dulumao/Guten-core/session"
	"time"
)

const (
	ADMIN    = "admin"
	CUSTOMER = "customer"
)

const authorizationName = "ctx_authorization"

type IAuthorization interface {
	GetId() interface{}
	GetUserName() string
	GetMetadata() map[string]interface{}
}

type Authorization struct {
	session *session.Session
	guard   string
}

type User struct {
	Username  string
	Guard     string
	Metadata  string
	CreatedAt time.Time
	ID        interface{}
}

func Use(sess *session.Session, guard string) *Authorization {
	auth := new(Authorization).setSession(sess)
	auth.guard = guard

	return auth
}

func (self *Authorization) Authenticate(auth IAuthorization) error {
	var user User
	var metadata []byte
	var err error

	metadata, err = json.Marshal(auth.GetMetadata())

	if err != nil {
		return err
	}

	user.Username = auth.GetUserName()
	user.Metadata = string(metadata)
	user.ID = auth.GetId()
	user.Guard = self.guard
	user.CreatedAt = time.Now()

	self.session.Set(fmt.Sprintf("%s_%s", authorizationName, self.guard), &user)

	if err = self.session.Save(); err != nil {
		return err
	}

	return nil
}

func (self Authorization) setSession(sess *session.Session) *Authorization {
	self.session = sess

	return &self
}

func (self *Authorization) Guest() bool {
	if user := self.User(); user != nil {
		if user.Guard == self.guard {
			return false
		}
	}

	return true
}

func (self *Authorization) User() *User {
	authorization := self.session.Get(fmt.Sprintf("%s_%s", authorizationName, self.guard))

	if authorization != nil {
		if user, ok := authorization.(*User); ok {
			return user
		}
	}

	return nil
}

func (self Authorization) Logout() {
	self.session.Delete(fmt.Sprintf("%s_%s", authorizationName, self.guard))
	self.session.Save()
}

func init() {
	gob.Register(new(User))
}
