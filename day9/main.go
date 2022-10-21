package main

import (
	"context"
	"day9/connection"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var Data = map[string]interface{}{
	"title": "personal Web",
}

type Blog struct {
	id 		int
	Title     string
	Post_date string
	Author    string
	Content   string
	Duration  string
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
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/hello", helloworld).Methods("GET")
	route.HandleFunc("/home", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/blog", blog).Methods("GET")
	route.HandleFunc("/blog-detail/{index}", blogDetail).Methods("GET")
	route.HandleFunc("/form-blog", formAddBlog).Methods("GET")
	route.HandleFunc("/add-blog", addBlog).Methods("POST")
	route.HandleFunc("/delete-blog/{index}", deleteBlog).Methods("GET")
	route.HandleFunc("/form-update/{index}", updateForm).Methods("GET")
	route.HandleFunc("/update-project/{index}", updateProject).Methods("POST")

	fmt.Println("server running on port 8000")
	http.ListenAndServe("localhost:8000", route)
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	//DATABASE
	rows, _ :=connection.Conn.Query(context.Background(), "SELECT id, name, description FROM tb_projects")
	
	var result []Blog //array data

	for rows.Next(){
		var each = Blog{} //memanggil struct

		err := rows.Scan(&each.id, &each.Title, &each.Content)

		if err != nil{
			fmt.Println(err.Error())
			return
		}

		each.Author = "Rezki Rahman"
		each.Post_date = "21 October 2022 11:01 WIB"
		result = append(result, each)
	}

	respData := map[string]interface{}{
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

	var tmpl, err = template.ParseFiles("views/blog-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var BlogDetail = Blog{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range Blogs {
		if index == i {
			BlogDetail = Blog{
				Title:     data.Title,
				Content:   data.Content,
				Post_date: data.Post_date,
				Author:    data.Author,
				Duration:  data.Duration,
			}
		}
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

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func addBlog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("title : " + r.PostForm.Get("inputTitle"))
	fmt.Println("content :" + r.PostForm.Get("inputContent"))

	var title = r.PostForm.Get("inputTitle")
	var content = r.PostForm.Get("inputContent")

	var newBlog = Blog{
		Title:     title,
		Content:   content,
		Author:    "Rezki Rahman",
		Post_date: "20 October 2022 22:30 WIB",
		Duration:  "2Bulan",
	}

	Blogs = append(Blogs, newBlog)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func deleteBlog(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	fmt.Println(index)

	Blogs = append(Blogs[:index], Blogs[index+1:]...)
	fmt.Println(Blogs)

	http.Redirect(w, r, "/home", http.StatusFound)
}

func updateForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("views/form-update.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var update = Blog{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range Blogs {
		if index == i {
			update = Blog{
				Title:     data.Title,
				Content:   data.Content,
				Post_date: data.Post_date,
				Author:    data.Author,
				Duration:  data.Duration,
			}
		}
	}

	data := map[string]interface{}{
		"Blog": update,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	fmt.Println(index)
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("title : " + r.PostForm.Get("editTitle"))
	fmt.Println("content :" + r.PostForm.Get("editContent"))

	var title = r.PostForm.Get("editTitle")
	var content = r.PostForm.Get("editContent")

	var newBlog = Blog{
		Title:     title,
		Content:   content,
		Author:    "Rezki Rahman",
		Post_date: "20 October 2022 22:30 WIB",
		Duration:  "2Bulan",
	}

	Blogs[index] = newBlog

	Blogs = append(Blogs)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
