package main

import (
	"context"
	"day12/connection"
	"day12/middleware"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type SessionData struct {
	IsLogin   bool
	UserName  string
	FlashData string
	userID	int
}

var Data = SessionData{}

type Project struct {
	Id            int
	Author        string
	NameProject   string
	StartDate     time.Time
	EndDate       time.Time
	StartDateForm string
	EndDateForm   string
	Duration      string
	Description   string
	Technologies  []string
	Reactjs       string
	Javascript    string
	Golang        string
	Nodejs        string
	Image         string
	IsLogin bool
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

func main() {

	route := mux.NewRouter()

	connection.DatabaseConnect()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/home", home).Methods("GET")
	route.HandleFunc("/addProject", addProject).Methods("GET")
	route.HandleFunc("/addProject", middleware.UploadFile(addProjectPost)).Methods("POST")
	route.HandleFunc("/contact", contactMe).Methods("GET")
	route.HandleFunc("/projectDetail/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/editProject/{id}", editProject).Methods("GET")
	route.HandleFunc("/update-project/{id}", middleware.UploadFile(submitEdit)).Methods("POST")
	route.HandleFunc("/deleteProject/{id}", deleteProject).Methods("GET")
	route.HandleFunc("/register", register).Methods("GET")
	route.HandleFunc("/submit-register", submitRegister).Methods("POST")
	route.HandleFunc("/login", login).Methods("GET")
	route.HandleFunc("/submit-login", submitLogin).Methods("POST")
	route.HandleFunc("/logout", logout).Methods("GET")

	fmt.Println("server running on port 5000")
	http.ListenAndServe("localhost:5000", route)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("views/home.html")
	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
		Data.userID = session.Values["ID"].(int)
		
	}

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, f1 := range fm {

			flashes = append(flashes, f1.(string))
		}
	}
	Data.FlashData = strings.Join(flashes, " ")
	println(flashes) 
	
	var result []Project
	if session.Values["Islogin"] != true {
		data, _ := connection.Conn.Query(context.Background(), `SELECT tb_project.id, tb_project.name_project, start_date, end_date, duration, description, technologies, image, tb_user.name FROM tb_project LEFT JOIN tb_user ON tb_project.author_id = tb_user.id WHERE "author_id"=$1`, Data.userID)//  
		for data.Next() {
			var each = Project{}
			err := data.Scan(&each.Id, &each.NameProject, &each.StartDate, &each.EndDate, &each.Duration, &each.Description, &each.Technologies, &each.Image, &each.Author)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			each.IsLogin = Data.IsLogin
			result = append(result, each)
		}
	} 
	resData := map[string]interface{}{
		"DataSession": Data,
		"Projects":    result,
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resData)
}
func contactMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("views/add-project.html")

	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func addProjectPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var nameProject = r.PostForm.Get("input-nameProject")
	var description = r.PostForm.Get("description")

	var startDate = r.PostForm.Get("input-startDate")
	var endDate = r.PostForm.Get("input-endDate")
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	technologies := []string{r.PostForm.Get("react"), r.PostForm.Get("javascript"), r.PostForm.Get("golang"), r.PostForm.Get("nodejs")}

	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	author := session.Values["ID"].(int)

	
	//Get duration
	hours := end.Sub(start).Hours()
	days := hours / 24
	weeks := math.Floor(days / 7)
	months := math.Floor(days / 30)
	years := math.Floor(days / 365)

	var duration string

	if years > 0 {
		duration = strconv.FormatFloat(years, 'f', 0, 64) + " Years"
	} else if months > 0 {
		duration = strconv.FormatFloat(months, 'f', 0, 64) + " Months"
	} else if weeks > 0 {
		duration = strconv.FormatFloat(weeks, 'f', 0, 64) + " Weeks"
	} else if days > 0 {
		duration = strconv.FormatFloat(days, 'f', 0, 64) + " Days"
	} else if hours > 0 {
		duration = strconv.FormatFloat(hours, 'f', 0, 64) + " Hours"
	} else {
		duration = "0 Days"
	}
	println(hours)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO public.tb_project (author_id, name_project, start_date, end_date, duration, description, technologies, image) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		author, nameProject, start, end, duration, description, technologies, image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("views/project-detail.html")

	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	var ProjectDetail = Project{}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name_project, start_date, end_date, duration, description, technologies, image FROM tb_project WHERE tb_project.id = $1", id).Scan(
		&ProjectDetail.Id, &ProjectDetail.NameProject, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Duration, &ProjectDetail.Description, &ProjectDetail.Technologies, &ProjectDetail.Image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	ProjectDetail.StartDateForm = ProjectDetail.StartDate.Format("2006-01-02")
	ProjectDetail.EndDateForm = ProjectDetail.EndDate.Format("2006-01-02")

	data := map[string]interface{}{
		"ProjectDetail": ProjectDetail,
	}

	tmpl.Execute(w, data)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/update-project.html")

	if err != nil {
		w.Write([]byte("message :" + err.Error()))
		return
	}

	var editProject = Project{}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name_project, start_date, end_date, duration, description, technologies, image FROM public.tb_project WHERE id = $1", id).Scan(
		&editProject.Id, &editProject.NameProject, &editProject.StartDate, &editProject.EndDate, &editProject.Duration, &editProject.Description, &editProject.Technologies, &editProject.Image)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	editProject.StartDateForm = editProject.StartDate.Format("2006-01-02")
	editProject.EndDateForm = editProject.EndDate.Format("2006-01-02")

	data := map[string]interface{}{
		"editProject": editProject,
	}
	tmpl.Execute(w, data)
}

func submitEdit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	nameProject := r.PostForm.Get("input-nameProject")
	description := r.PostForm.Get("description")

	//var startDate = r.PostForm.Get("input-startDate")
	//var endDate = r.PostForm.Get("input-endDate")

	var reactjs = r.PostForm.Get("react")
	var javascript = r.PostForm.Get("javascript")
	var golang = r.PostForm.Get("golang")
	var nodejs = r.PostForm.Get("nodejs")
	var technologies = []string{reactjs, javascript, golang, nodejs}
	var startDate = r.PostForm.Get("input-startDate")
	var endDate = r.PostForm.Get("input-endDate")
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	//GET DURATION
	hours := end.Sub(start).Hours()
	days := hours / 24
	weeks := math.Floor(days / 7)
	months := math.Floor(days / 30)
	years := math.Floor(days / 365)

	var duration string

	if years > 0 {
		duration = strconv.FormatFloat(years, 'f', 0, 64) + " Years"
	} else if months > 0 {
		duration = strconv.FormatFloat(months, 'f', 0, 64) + " Months"
	} else if weeks > 0 {
		duration = strconv.FormatFloat(weeks, 'f', 0, 64) + " Weeks"
	} else if days > 0 {
		duration = strconv.FormatFloat(days, 'f', 0, 64) + " Days"
	} else if hours > 0 {
		duration = strconv.FormatFloat(hours, 'f', 0, 64) + " Hours"
	} else {
		duration = "0 Days"
	}
	println(hours)

	_, err = connection.Conn.Exec(context.Background(), "UPDATE public.tb_project SET name_project = $1, start_date = $2, end_date = $3, duration = $4, description = $5, technologies = $6, image = $7 WHERE id = $8",
		nameProject, start, end, duration, description, technologies, image, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM public.tb_project WHERE id = $1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("views/register.html")

	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func submitRegister(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var name = r.PostForm.Get("inputName")
	var email = r.PostForm.Get("inputEmail")
	var password = r.PostForm.Get("inputPassword")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("views/login.html")

	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, f1 := range fm {

			flashes = append(flashes, f1.(string))
		}
	}
	Data.FlashData = strings.Join(flashes, " ")

	resData := map[string]interface{}{
		"DataLogin": Data,
	}

	tmpl.Execute(w, resData)
}

func submitLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	var email = r.PostForm.Get("inputEmail")
	var password = r.PostForm.Get("inputPassword")

	user := User{}

	// mengambil data email, dan melakukan pengecekan email
	err = connection.Conn.QueryRow(context.Background(),
		"SELECT * FROM tb_user WHERE email=$1", email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)

	if err != nil {
		fmt.Println("Email belum terdaftar")
		var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")

		session.AddFlash("Email belum terdaftar", "message")
		session.Save(r, w)

		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	// melakukan pengecekan password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("Password salah")
		var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")

		session.AddFlash("Password anda salah", "message")
		session.Save(r, w)

		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	//berfungsi untuk menyimpan data kedalam sessions browser
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	session.Values["Name"] = user.Name
	session.Values["Email"] = user.Email
	session.Values["ID"] = user.Id
	session.Values["IsLogin"] = true
	session.Options.MaxAge = 10800 // 3 JAM

	session.AddFlash("succesfull login", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logout")
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
