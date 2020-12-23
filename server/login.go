package server

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"

	"../redis"

	// red "../redis"
	jr "../response"
)

// Login ...
func (s *Server) Login(w http.ResponseWriter, req *http.Request) {
	// Get The Active User if any Loged In
	user := s.ActiveUser(req)
	if user != nil {
		jr.PrepareResponse(w, http.StatusOK, &jr.RespForm{
			Type:    http.StatusOK,
			Msg:     "You Are Loged In",
			Content: *user,
		})
		return
	}
	username := req.FormValue("Username")
	password := req.FormValue("password")
	if len(username) == 0 || len(password) == 0 || username == "" || password == "" {
		jr.PrepareResponse(w, http.StatusBadRequest, &jr.RespForm{
			Type:    http.StatusBadRequest,
			Msg:     "Please Verify your Credentials",
			Content: "Put Proper Username & Password !!!",
		})
		return
	}

	// Now We need to Get Our SessionID => Cookies
	cookie, userID, err := LoginCookie(&username, &password, s.DB)
	if err != nil {
		jr.PrepareResponse(w, http.StatusBadRequest, &jr.RespForm{
			Type:    http.StatusBadRequest,
			Msg:     "Please Verify your Credentials",
			Content: "Put Proper Username & Password !!!",
		})
		return
	}

	result, err := redis.Redis.Do("SET", cookie.Value, *userID)
	if err != nil {
		jr.PrepareResponse(w, http.StatusInternalServerError, &jr.RespForm{
			Type:    http.StatusInternalServerError,
			Msg:     "Error in Server",
			Content: err.Error(),
		})
		return
	}
	if result == nil {
		jr.PrepareResponse(w, http.StatusUnauthorized, &jr.RespForm{
			Type:    http.StatusUnauthorized,
			Msg:     "You Are Not Authorized !!!",
			Content: "Error Log In !!!",
		})
		return
	}
	http.SetCookie(w, cookie)
	jr.PrepareResponse(w, http.StatusOK, &jr.RespForm{
		Type:    http.StatusOK,
		Msg:     "Success !!!",
		Content: cookie.Value,
	})
	return
}

// IsLogged ...
func (s *Server) IsLogged(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user := s.ActiveUser(req)
		if user == nil {
			jr.PrepareResponse(w, http.StatusUnauthorized, &jr.RespForm{
				Type:    http.StatusUnauthorized,
				Msg:     "You Are Not Authorized !!!",
				Content: "Error Log In !!!",
			})
			return
		}
		next(w, req)
	}
}

// ActiveUser ...
func (s *Server) ActiveUser(req *http.Request) *int {
	cookie, err := req.Cookie("SessionID")

	if err != nil {
		return nil
	}

	result, err := redis.Redis.Do("GET", cookie.Value)
	if err != nil {
		return nil
	}

	if result == nil {
		return nil
	}
	var userID int
	// This Query is Only To Check if User Exists
	query := "select ID from user where ID = ?"
	err = s.DB.QueryRow(query, result).Scan(&userID)
	if err != nil {
		return nil
	}
	return &userID
}

// LoginCookie ...
func LoginCookie(username, password *string, database *sql.DB) (*http.Cookie, *int, error) {
	query := "select ID from user where Username = ? and UserPass = sha1(?);"
	var userID int
	err := database.QueryRow(query, *username, *password).Scan(&userID)
	if err != nil {
		return nil, nil, err
	}

	sessionToken := uuid.New().String()

	cookie := http.Cookie{
		Name:    "SessionID",
		Value:   sessionToken,
		Expires: time.Now().Add(3600 * time.Second),
	}
	return &cookie, &userID, nil
}

// Logout ...
func (s *Server) Logout(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("SessionID")
	if err != nil {
		jr.PrepareResponse(w, http.StatusBadRequest, &jr.RespForm{
			Type:    http.StatusBadRequest,
			Msg:     "Please Verify your Credentials",
			Content: "Put Proper Username & Password !!!",
		})
		return
	}
	result, err := redis.Redis.Do("DEL", cookie.Value)
	if err != nil {
		jr.PrepareResponse(w, http.StatusInternalServerError, &jr.RespForm{
			Type:    http.StatusInternalServerError,
			Msg:     "Error in Server",
			Content: err.Error(),
		})
		return
	}
	if result == nil {
		jr.PrepareResponse(w, http.StatusUnauthorized, &jr.RespForm{
			Type:    http.StatusUnauthorized,
			Msg:     "You Are Not Authorized !!!",
			Content: "Error Log In !!!",
		})
		return
	}
	jr.PrepareResponse(w, http.StatusOK, &jr.RespForm{
		Type:    http.StatusOK,
		Msg:     "Log Out !!!",
		Content: "Logging Out !!!",
	})
	return
}

// IsPermitted ...
func (s *Server) IsPermitted(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user := s.ActiveUser(req)
		if user == nil {
			jr.PrepareResponse(w, http.StatusUnauthorized, &jr.RespForm{
				Type:    http.StatusUnauthorized,
				Msg:     "You Are Not Authorized !!!",
				Content: "Error Log In !!!",
			})
			return
		}
		permission := req.FormValue("Permission")
		query := "select GroupID from user where ID = ?;"
		var groupID int
		err := s.DB.QueryRow(query, *user).Scan(&groupID)
		if err != nil {
			jr.PrepareResponse(w, http.StatusUnauthorized, &jr.RespForm{
				Type:    http.StatusUnauthorized,
				Msg:     "You Are Not Authorized !!!",
				Content: "Error Log In !!!",
			})
			return
		}
		if groupID == 1 || groupID == 2 {
			if groupID == 2 {
				newQuery := "select ID from user_permissions where UserID = ? and PermissionID = ?;"
				exist := -1
				err = s.DB.QueryRow(newQuery, *user, permission).Scan(&exist)
				if err != nil {
					jr.PrepareResponse(w, http.StatusUnauthorized, &jr.RespForm{
						Type:    http.StatusUnauthorized,
						Msg:     "You Are Not Authorized !!!",
						Content: "Error Log In !!!",
					})
					return
				}
				if exist == -1 {
					jr.PrepareResponse(w, http.StatusUnauthorized, &jr.RespForm{
						Type:    http.StatusUnauthorized,
						Msg:     "You Are Not Authorized !!!",
						Content: "Error Log In !!!",
					})
					return
				}
			}
			next(w, req)
		} else {
			jr.PrepareResponse(w, http.StatusUnauthorized, &jr.RespForm{
				Type:    http.StatusUnauthorized,
				Msg:     "You Are Not Authorized !!!",
				Content: "Error Log In !!!",
			})
			return
		}
	}
}
