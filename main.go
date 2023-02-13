package main

import (
	"context"
	"time"

	"fmt"

	"strconv"

	"log"

	"net/http"

	"github.com/gorilla/mux"

	"html/template"

	"Personal-Web/connection"
)

var Data = map[string]interface{}{
	"Title": "Personal Web",
}

type Project struct {
	Id           int
	Title        string
	Start_date   time.Time
	End_date     time.Time
	Description  string
	Technologies string
	NodeJs       string
	Java         string
	Php          string
	Laravel      string
	Image        string
	Format_start string
	Format_end   string
}

//type Projects []Project

//func NewProject() *Project {
//	return &Project{
//		Start_date: time.Date(2022, 5, 12, 21, 0, 0, 0, time.Local),
//		End_date:   time.Date(2022, 5, 12, 21, 0, 0, 0, time.Local),
//	}
//}

//var Projects = []Project{
//	{
//		Title:       "Pembelajaran Online",
//		Duration:    "Duration : 3 Weeks",
//		Author:      " | Bagas",
//		Description: "Sangat sulit sekali hehehehe",
//	},
//}

// function routing
func main() {
	router := mux.NewRouter()

	connection.DatabaseConnect()

	// Create Folder
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/Project", project).Methods("GET")
	router.HandleFunc("/addProject", addProject).Methods("POST")
	router.HandleFunc("/contactMe", contactMe).Methods("GET")
	router.HandleFunc("/projectDetail/{id}", projectDetail).Methods("GET")
	router.HandleFunc("/delete-project/{id}", deleteProject).Methods("DELETE")

	fmt.Println("Server Running Successfully")
	http.ListenAndServe("localhost:5000", router)
}

// function handling index.html
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	//Parsing template html file
	var tmpl, err = template.ParseFiles("index.html")
	//Error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	rows, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM public.tb_project;")

	var result []Project
	for rows.Next() {
		var each = Project{}

		var err = rows.Scan(&each.Id, &each.Title, &each.Start_date, &each.End_date, &each.Description, &each.Technologies, &each.Image)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.Format_start = each.Start_date.Format("12 May 2001")
		each.Format_end = each.End_date.Format("12 May 2001")

		result = append(result, each)
	}

	resp := map[string]interface{}{
		"Title":    Data,
		"Projects": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

// function handling myproject.html
func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	//Parsing template html file
	var tmpl, err = template.ParseFiles("myproject.html")
	//Error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

// function handling contactMe.html
func contactMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	//Parsing template html file
	var tmpl, err = template.ParseFiles("contactMe.html")
	//Error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

// function handling myproiect-detail.html
func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	//Parsing template html file
	var tmpl, err = template.ParseFiles("myproject-detail.html")
	//Error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	ProjectDetail := Project{}

	connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM public.tb_project WHERE id=$1", id).
		Scan(&ProjectDetail.Id, &ProjectDetail.Title, &ProjectDetail.Start_date, &ProjectDetail.End_date, &ProjectDetail.Description, &ProjectDetail.Technologies, &ProjectDetail.Image)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	resp := map[string]interface{}{
		"Data":    Data,
		"Project": ProjectDetail,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("title")
	description := r.PostForm.Get("description")
	startDate := r.PostForm.Get("start-date")
	endDate := r.PostForm.Get("end-date")
	nodeJs := r.PostForm.Get("NodeJs")
	java := r.PostForm.Get("Java")
	php := r.PostForm.Get("Php")
	laravel := r.PostForm.Get("Laravel")

	connection.Conn.Exec(context.Background(), "INSERT INTO tb_project(name, start_date, end_date, description, technologies, image) VALUES($1, $2, $3, $4, $5, $6)", title, description, startDate, endDate, nodeJs, java, php, laravel)

	//	var newProject = Project{
	//		Title:       title,
	//		Start_date:  startDate,
	//		End_date:    endDate,
	//		Description: description,
	//		NodeJs:      nodeJs,
	//		Java:        java,
	//		Php:         php,
	//		Laravel:     laravel,
	//	}

	//Projects = append(Projects, newProject)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	//	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	//	Projects = append(Projects[:id], Projects[id+1:]...)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
