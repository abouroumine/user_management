package server

import (
	"database/sql"
	"net/http"
	"strconv"

	db "../database"
	red "../redis"
	jr "../response"
)

// Server ...
type Server struct {
	Server *http.ServeMux
	DB     *sql.DB
}

// Run ...
func (s *Server) Run(add string) {
	http.ListenAndServe(add, s.Server)
}

// InitializeDB ...
func (s *Server) InitializeDB() {
	var err error
	s.DB, err = db.ConnectMySQL()
	if err != nil {
		panic(err)
	}
}

// Initialize ...
func (s *Server) Initialize() {
	s.InitializeDB()

	s.Server = http.NewServeMux()

	red.InitRedis()

	s.Server.HandleFunc("/hello", s.Hello)
	s.Server.HandleFunc("/logedhello", s.IsLogged(s.LogedHello))
	s.Server.HandleFunc("/login", s.Login)
	s.Server.HandleFunc("/logout", s.IsLogged(s.Logout))
	s.Server.HandleFunc("/create", s.IsPermitted(s.IsLogged(s.CreateUser)))
	s.Server.HandleFunc("/viewall", s.IsPermitted(s.IsLogged(s.ViewAll)))
	s.Server.HandleFunc("/viewone", s.IsPermitted(s.IsLogged(s.ViewOne)))
	s.Server.HandleFunc("/edit", s.IsPermitted(s.IsLogged(s.Edit)))
	s.Server.HandleFunc("/delete", s.IsPermitted(s.IsLogged(s.Delete)))

	s.Run(":8080")
}

// Hello ...
func (s *Server) Hello(w http.ResponseWriter, r *http.Request) {
	jr.PrepareResponse(w, http.StatusOK, &jr.RespForm{
		Type:    http.StatusOK,
		Msg:     "Hello World !!!",
		Content: "New Content !!!",
	})
	return
}

// LogedHello ...
func (s *Server) LogedHello(w http.ResponseWriter, r *http.Request) {
	jr.PrepareResponse(w, http.StatusOK, &jr.RespForm{
		Type:    http.StatusOK,
		Msg:     "Hello World You have loged in successefully!!!",
		Content: "New Content 2!!!",
	})
	return
}

// CreateUser ...
func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	fullname := r.FormValue("fullname")
	age := r.FormValue("age")
	address := r.FormValue("address")
	groupID := r.FormValue("groupID")
	statement, err := s.DB.Prepare("insert into user(Username, UserPass, FullName, Age, Address, GroupID) values (?, sha1(?), ?, ?, ?, ?);")
	if err != nil {
		jr.PrepareResponse(w, http.StatusBadRequest, &jr.RespForm{
			Type:    http.StatusBadRequest,
			Msg:     "Error in the Query Definition",
			Content: err.Error(),
		})
		return
	}
	result, err := statement.Exec(username, password, fullname, age, address, groupID)
	if err != nil {
		jr.PrepareResponse(w, http.StatusBadRequest, &jr.RespForm{
			Type:    http.StatusBadRequest,
			Msg:     "Error in the Query Execution",
			Content: err.Error(),
		})
		return
	}
	jr.PrepareResponse(w, http.StatusOK, &jr.RespForm{
		Type:    http.StatusOK,
		Msg:     "The New User ID is: ",
		Content: result,
	})
	return
}

// ViewAll ...
func (s *Server) ViewAll(w http.ResponseWriter, r *http.Request) {
	users, err := GetAll(s.DB)
	if err != nil {
		jr.PrepareResponse(w, http.StatusBadRequest, &jr.RespForm{
			Type:    http.StatusBadRequest,
			Msg:     "Error in the Query Execution",
			Content: err.Error(),
		})
		return
	}
	jr.PrepareResponse(w, http.StatusOK, &jr.RespForm{
		Type:    http.StatusOK,
		Msg:     "The User List is: ",
		Content: *users,
	})
	return
}

// ViewOne ...
func (s *Server) ViewOne(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("UserID")
	user, err := GetOne(s.DB, &username)
	if err != nil {
		jr.PrepareResponse(w, http.StatusBadRequest, &jr.RespForm{
			Type:    http.StatusBadRequest,
			Msg:     "Error in the Query Execution",
			Content: err.Error(),
		})
		return
	}
	jr.PrepareResponse(w, http.StatusOK, &jr.RespForm{
		Type:    http.StatusOK,
		Msg:     "The User is: ",
		Content: *user,
	})
	return
}

// Edit ...
func (s *Server) Edit(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("UserID")
	FullName := r.FormValue("FullName")
	id, err := strconv.Atoi(userID)
	if err != nil {
		jr.PrepareResponse(w, http.StatusBadRequest, &jr.RespForm{
			Type:    http.StatusBadRequest,
			Msg:     "Error in the Information Send",
			Content: err.Error(),
		})
		return
	}
	user, err := EditUser(s.DB, &id, &FullName)
	if err != nil {
		jr.PrepareResponse(w, http.StatusBadRequest, &jr.RespForm{
			Type:    http.StatusBadRequest,
			Msg:     "Error in the Query Execution",
			Content: err.Error(),
		})
		return
	}
	jr.PrepareResponse(w, http.StatusOK, &jr.RespForm{
		Type:    http.StatusOK,
		Msg:     "The User is: ",
		Content: *user,
	})
	return
}

// Delete ...
func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("UserID")
	id, err := strconv.Atoi(userID)
	if err != nil {
		jr.PrepareResponse(w, http.StatusBadRequest, &jr.RespForm{
			Type:    http.StatusBadRequest,
			Msg:     "Error in the Information Send",
			Content: err.Error(),
		})
		return
	}
	userIDDeleted, err := DeleteUser(s.DB, &id)
	if err != nil {
		jr.PrepareResponse(w, http.StatusBadRequest, &jr.RespForm{
			Type:    http.StatusBadRequest,
			Msg:     "Error in the Query Execution",
			Content: err.Error(),
		})
		return
	}
	jr.PrepareResponse(w, http.StatusOK, &jr.RespForm{
		Type:    http.StatusOK,
		Msg:     "The User ID is: ",
		Content: *userIDDeleted,
	})
	return
}

// GetAll ...
func GetAll(database *sql.DB) (*[]db.UserAPI, error) {
	var users = make([]db.UserAPI, 0)
	rows, err := database.Query("SELECT Username, UserPass, FullName, Age, Address, Name FROM user, user_group where user.GroupID = user_group.ID;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user db.UserAPI
		err = rows.Scan(&user.Username, &user.Password, &user.FullName, &user.Age, &user.Address, &user.GroupName)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return &users, nil
}

// GetOne ...
func GetOne(database *sql.DB, id *string) (*db.UserAPI, error) {
	var user db.UserAPI
	err := database.QueryRow("SELECT Username, UserPass, FullName, Age, Address, Name FROM userManagement.user, userManagement.user_group where user.GroupID = user_group.ID and Username = ?;", *id).
		Scan(&user.Username, &user.Password, &user.FullName, &user.Age, &user.Address, &user.GroupName)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// EditUser ...
func EditUser(database *sql.DB, userID *int, fullName *string) (*db.UserAPI, error) {
	statement, err := database.Prepare("update user set FullName = ? where ID = ?;")
	if err != nil {
		return nil, err
	}
	_, err = statement.Exec(*fullName, *userID)
	if err != nil {
		return nil, err
	}
	var user db.UserAPI
	err = database.QueryRow("SELECT Username, UserPass, FullName, Age, Address, Name FROM userManagement.user, userManagement.user_group where user.GroupID = user_group.ID and user.ID = ?;", *userID).
		Scan(&user.Username, &user.Password, &user.FullName, &user.Age, &user.Address, &user.GroupName)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUser ...
func DeleteUser(database *sql.DB, userID *int) (*int, error) {
	statement, err := database.Prepare("delete from user where ID = ?;")
	if err != nil {
		return nil, err
	}
	_, err = statement.Exec(*userID)
	if err != nil {
		return nil, err
	}
	return userID, nil
}
