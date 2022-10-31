package main

import (
	"context"
	"day11/connection"
	"day11/middleware"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type MetaData struct {
	Title     string
	IsLogin   bool
	UserName  string
	FlashData string
}

var Data = MetaData{
	Title: "personal Web",
}

type Blog struct {
	Id        int
	Title     string
	Image     string
	Post_date string
	Author    string
	Content   string
	Duration  string
	IsLogin   bool
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

var Blogs = []Blog{
	// {
	// 	Title:     "Dumbways Web App",
	// 	Post_date: "20 October 2022 22:30 WIB",
	// 	Author:    "Rezki Rahman",
	// 	Content:   "Lorem Ipsum",
	// 	Duration:  "2 Bulan",
	// },
	// {
	// 	Title:     "Dumbways Mobile App",
	// 	Post_date: "20 October 2022 22:30 WIB",
	// 	Author:    "Rezki Rahman",
	// 	Content:   "Lorem Ipsum dolor",
	// 	Duration:  "2 Bulan",
	// },
}

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()

	//route path folder untuk public
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	route.HandleFunc("/hello", helloworld).Methods("GET")
	route.HandleFunc("/home", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/blog", blog).Methods("GET")
	route.HandleFunc("/blog-detail/{id}", blogDetail).Methods("GET")

	route.HandleFunc("/form-blog", formAddBlog).Methods("GET")
	route.HandleFunc("/add-blog", middleware.UploadFile(addBlog)).Methods("POST")

	route.HandleFunc("/delete-blog/{id}", deleteBlog).Methods("GET")

	route.HandleFunc("/form-update/{id}", updateForm).Methods("GET")
	route.HandleFunc("/update-project/{id}", updateProject).Methods("POST")

	route.HandleFunc("/register-form", registerForm).Methods("GET")
	route.HandleFunc("/register", register).Methods("POST")

	route.HandleFunc("/login-form", loginForm).Methods("GET")
	route.HandleFunc("/login", login).Methods("POST")

	route.HandleFunc("/logout", logout).Methods("GET")

	fmt.Println("server running on port 8000")
	http.ListenAndServe("localhost:8000", route)
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

//-------------------REGISTER FUNCTION-----------------------------

func registerForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("views/form-register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)

}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var name = r.PostForm.Get("inputName")
	var email = r.PostForm.Get("inputEmail")
	var password = r.PostForm.Get("inputPassword")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	fmt.Println(passwordHash)

	_, err = connection.Conn.Exec(context.Background(),
		"INSERT INTO tb_users(name, email, password) VALUES($1, $2, $3)", name, email, passwordHash)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/login-form", http.StatusMovedPermanently)

}

//----------------------LOGIN FUNCTION--------------------------------

func loginForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/form-login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)

}

func login(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("inputEmail")
	password := r.PostForm.Get("inputPassword")

	user := User{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_users WHERE email=$1", email).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	session.Values["IsLogin"] = true
	session.Values["Name"] = user.Name
	session.Values["ID"] = user.Id
	session.Options.MaxAge = 10800 // 3 hours

	session.AddFlash("Successfully Login!", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

//---------------------LOGOUT FUNCTION------------------------

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logout!")
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

//-----------------------------------------------------------------

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	fmt.Println(session.Values)

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

	rows, _ := connection.Conn.Query(context.Background(),
		"SELECT tb_projects.id, tb_projects.name, description, image, tb_users.name FROM tb_projects LEFT JOIN tb_users ON tb_projects.author_id = tb_users.id ORDER by tb_projects.id")

	var result []Blog //array data

	for rows.Next() {
		var each = Blog{} //memanggil struct

		err := rows.Scan(&each.Id, &each.Title, &each.Content, &each.Image, &each.Author)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.Post_date = "21 October 2022 11:01 WIB"

		// var oneDay = 24*60*60*1000
		// each.Duration = (each.stardate - each.enddate)/oneDay

		if session.Values["IsLogin"] != true {
			each.IsLogin = false
		} else {
			each.IsLogin = session.Values["IsLogin"].(bool)
		}

		result = append(result, each)

	}

	fmt.Println(result)
	fmt.Println(session)

	respData := map[string]interface{}{
		"Data":  Data,
		"Blogs": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("views/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func blog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("views/blog.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func blogDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset-utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/blog-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var BlogDetail = Blog{}

	err = connection.Conn.QueryRow(context.Background(),
		"SELECT id, name, description, image FROM tb_projects WHERE id=$1", id).Scan(
		&BlogDetail.Id, &BlogDetail.Title, &BlogDetail.Content, &BlogDetail.Image,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	data := map[string]interface{}{
		"Blog": BlogDetail,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func formAddBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("views/add-blog.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	fmt.Println(session.Values)

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func addBlog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("inputTitle")
	content := r.PostForm.Get("inputContent")
	startdate := r.PostForm.Get("inputStardate")
	enddate := r.PostForm.Get("inputEnddate")
	startDateForm, _ := time.Parse("2006-01-02", startdate)
	endDateForm, _ := time.Parse("2006-01-02", enddate)

	fmt.Println(startdate, enddate)

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	author := session.Values["ID"].(int)
	author = 3

	dataContex := r.Context().Value("dataFile")
	image := dataContex.(string)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(name,description,image,author_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6)", title, content, image, author, startDateForm, endDateForm)

	//_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(name,description,image,author_id) VALUES ($1, $2, $3, $4)", title, content, image, author)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	//Blogs = append(Blogs, newBlog)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func deleteBlog(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println(id)

	// Blogs = append(Blogs[:index], Blogs[index+1:]...)
	// fmt.Println(Blogs)

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusFound)
}

func updateForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset-utf-8")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/form-update.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var update = Blog{}

	err = connection.Conn.QueryRow(context.Background(),
		"SELECT id, name, description FROM tb_projects WHERE id=$1", id).Scan(
		&update.Id, &update.Title, &update.Content,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	// for i, data := range Blogs {
	// 	if index == i {
	// 		update = Blog{
	// 			Title:     data.Title,
	// 			Content:   data.Content,
	// 			Post_date: data.Post_date,
	// 			Author:    data.Author,
	// 			Duration:  data.Duration,
	// 		}
	// 	}
	// }

	data := map[string]interface{}{
		"Blog": update,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println(id)
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("title : " + r.PostForm.Get("editTitle"))
	fmt.Println("content :" + r.PostForm.Get("editContent"))

	var title = r.PostForm.Get("editTitle")
	var content = r.PostForm.Get("editContent")

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	dataContex := r.Context().Value("dataFile")
	image := dataContex

	_, err = connection.Conn.Exec(context.Background(),
		"UPDATE tb_projects SET name=$1, description=$2, image=$3  WHERE id=$4", title, content, image, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	// var newBlog = Blog{
	// 	Title:     title,
	// 	Content:   content,
	// 	Author:    "Rezki Rahman",
	// 	Post_date: "20 October 2022 22:30 WIB",
	// 	Duration:  "2Bulan",
	// }

	//Blogs[index] = newBlog
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
