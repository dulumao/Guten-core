package session

import (
	"github.com/fatih/structs"
	"github.com/gookit/rux"
	"github.com/gorilla/sessions"
	"github.com/dulumao/Guten-core/env"
	"net/http"
	"net/url"
)

const defaultMaxMemory = 32 << 20 // 32 MB

type Session struct {
	sess      *sessions.Session
	store     sessions.Store
	r         *http.Request
	w         http.ResponseWriter
	isWritten bool
}

type AlertFlash struct {
	Alert   interface{} `json:"alert"`
	Message interface{} `json:"message"`
}

type RequestFlash struct {
	PostForm url.Values
	Query    url.Values
}

func (self *AlertFlash) IsEmpty() bool {
	if self.Alert != nil || self.Message != nil {
		return false
	}

	return true
}

func (self *AlertFlash) IsError() bool {
	if self.Alert == "error" {
		return true
	}

	return false
}

func (self *AlertFlash) IsSuccess() bool {
	if self.Alert == "success" {
		return true
	}

	return false
}

func (self *AlertFlash) IsWarning() bool {
	if self.Alert == "warning" {
		return true
	}

	return false
}

func (self *AlertFlash) IsInfo() bool {
	if self.Alert == "info" {
		return true
	}

	return false
}

func (self *AlertFlash) Is(v string) bool {
	if self.Alert == v {
		return true
	}

	return false
}

func (s *Session) Get(key interface{}) interface{} {
	return s.getSession().Values[key]
}

func (s *Session) Set(key interface{}, val interface{}, saved ...bool) {
	s.getSession().Values[key] = val
	s.isWritten = true

	if len(saved) > 0 {
		if saved[0] {
			s.Save()
		}
	}
}

func (s *Session) Delete(key interface{}) {
	delete(s.getSession().Values, key)
	s.isWritten = true
}

func (s *Session) Clear() {
	for key := range s.getSession().Values {
		s.Delete(key)
	}
}

func (s *Session) All() map[interface{}]interface{} {
	return s.getSession().Values
}

func (s *Session) AddFlash(value interface{}, vars ...string) {
	s.getSession().AddFlash(value, vars...)
	s.isWritten = true
}

func (s *Session) GetFlashes(vars ...string) []interface{} {
	s.isWritten = true
	flashes := s.getSession().Flashes(vars...)

	return flashes
}

func (s *Session) AddAlertFlash(alert, message string) *Session {
	s.AddFlash(alert, "alert")
	s.AddFlash(message, "message")

	return s
}

func (s *Session) GetAlertFlash() *AlertFlash {
	var alertFlash = new(AlertFlash)

	alerts := s.GetFlashes("alert")
	messages := s.GetFlashes("message")

	if len(alerts) > 0 {
		alertFlash.Alert = alerts[0]
	}

	if len(alerts) > 0 {
		alertFlash.Message = messages[0]
	}

	if !alertFlash.IsEmpty() {
		s.Save()
	}

	return alertFlash
}

func (s *Session) AddRequestFlash() *Session {
	_ = s.r.ParseForm()
	_ = s.r.ParseMultipartForm(defaultMaxMemory)

	s.AddFlash(s.r.PostForm, "_old_post_request")
	s.AddFlash(s.r.URL.Query(), "_old_query_request")

	return s
}

func (s *Session) GetRequestFlash() *RequestFlash {
	var requestFlash = new(RequestFlash)

	postForm := s.GetFlashes("_old_post_request")
	query := s.GetFlashes("_old_query_request")

	if len(postForm) > 0 {
		requestFlash.PostForm = postForm[0].(url.Values)
	}

	if len(query) > 0 {
		requestFlash.Query = query[0].(url.Values)
	}

	s.Save()

	return requestFlash
}

func (s *Session) AddStructFlash(v interface{}, name ...string) *Session {
	if len(name) > 0 {
		s.AddFlash(v, "_struct_"+name[0])
	} else {
		s.AddFlash(v, "_struct")
	}

	return s
}

func (s *Session) GetStructFlash(name ...string) interface{} {
	var v []interface{}

	if len(name) > 0 {
		v = s.GetFlashes("_struct_" + name[0])
	} else {
		v = s.GetFlashes("_struct")
	}

	s.Save()

	if len(v) > 0 {
		return v[0]
	}

	return nil
}

func (s *Session) GetStructMapFlash(name ...string) (vm map[string]interface{}) {
	v := s.GetStructFlash(name...)

	if v != nil {
		vm = structs.Map(v)
	}

	return vm
}

func (s *Session) SetOptions(options *sessions.Options) {
	s.getSession().Options = options
}

func (s *Session) Save() error {
	if s.IsWritten() {
		e := s.getSession().Save(s.r, s.w)

		if e == nil {
			s.isWritten = false
		}

		return e
	}

	return nil
}

func (s *Session) getSession() *sessions.Session {
	if s.sess == nil {
		var err error
		s.sess, err = s.store.Get(s.r, s.Name())

		if err != nil {
			panic(err)
		}
	}

	return s.sess
}

func (s *Session) IsWritten() bool {
	return s.isWritten
}

func (s *Session) Name() string {
	return env.Value.Session.Name
}

func (s *Session) SetRequest(r *http.Request) *Session {
	s.r = r
	return s
}

func (s *Session) SetResponse(w http.ResponseWriter) *Session {
	s.w = w
	return s
}

func NewSession(store sessions.Store) *Session {
	return &Session{
		store: store,
	}
}

// shortcut to get session
func Use(c *rux.Context) *Session {
	return Wrap(c.Resp, c.Req)
}

func Wrap(w http.ResponseWriter, r *http.Request) *Session {
	sess := NewSession(Value)
	sess.SetRequest(r)
	sess.SetResponse(w)
	sess.SetOptions(&sessions.Options{
		Path:     env.Value.Session.Path,
		MaxAge:   env.Value.Session.Lifetime,
		HttpOnly: env.Value.Session.HTTPOnly,
		Secure:   env.Value.Session.Secure,
	})

	return sess
}
