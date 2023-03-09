package main

import (
	"encoding/json"
	"fmt"
	"go-book-webapp/models"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var message string = ""

func welcomeHandler(w http.ResponseWriter, 
	r *http.Request) {
	templ, error := template.ParseFiles("template/welcome.html")
	if error != nil {
		log.Fatal()
	}

	error = templ.Execute(w, nil)
}

func interactHandler(w http.ResponseWriter, 
	r *http.Request) {
	type context struct {
		Books []models.Book
		Message string
		Msg string
	}
	
	templ, error := template.ParseFiles("template/interact.html")
	if error != nil {
		log.Fatal()
	}
	var ctx context
	if len(message) > 0 {
		ctx.Msg = message
		message = ""
	}
	msg, ok := r.URL.Query()["msg"]
	if ok && len(msg[0])>0 {
		ctx.Message = msg[0]
	}
	ctx.Books = models.GetAllBooks()
	error = templ.Execute(w, ctx)
}

func newHandler(w http.ResponseWriter, 
	r *http.Request) {
	type context struct {
		Book models.Book
	}
	templ, error := template.ParseFiles("template/new.html")
	if error != nil {
		log.Fatal()
	}

	error = templ.Execute(w, nil)
}

func createHandler(w http.ResponseWriter, 
	r *http.Request) {

	title := r.FormValue("Title")
	author := r.FormValue("Author")
	publication := r.FormValue("Publication")

	newBook := &models.Book{}
	newBook.Title = title
	newBook.Author = author
	newBook.Publication = publication

	newBook.CreateBook()
	message = "New book " + title + " added"
	http.Redirect(w, r, "/interact?msg=New book added", http.StatusFound)
}

func deleteHandler(w http.ResponseWriter, 
	r *http.Request) {
	
	key, ok := r.URL.Query()["bookId"]

	if ok && len(key[0]) > 0 {
		id, err := strconv.ParseInt(key[0], 10, 64)
		if err == nil {
			models.DeleteBook(id)
		}
		
	}

	http.Redirect(w, r, "/interact?msg=Book deleted", http.StatusFound)
}

func editHandler(w http.ResponseWriter, 
	r *http.Request) {
	type context struct {
		Book models.Book
	}
	var ctx context
	key, ok := r.URL.Query()["bookId"]

	if ok && len(key[0]) > 0 {
		id, err := strconv.ParseInt(key[0], 10, 64)
		if err == nil {
			ctx.Book,_ = models.GetBook(id)
		}		
	}
	templ, _ := template.ParseFiles("template/update.html")
	templ.Execute(w, ctx)
}

func updateHandler(w http.ResponseWriter, 
	r *http.Request) {
	key, ok := r.URL.Query()["bookId"]

	if ok && len(key[0]) > 0 {
		id, err := strconv.ParseInt(key[0], 10, 64)
		if err == nil {
			title := r.FormValue("Title")
			author := r.FormValue("Author")
			publication := r.FormValue("Publication")
			newBook,count := models.GetBook(id)
			if count>0 {
				newBook.Title = title
				newBook.Author = author
				newBook.Publication = publication
				newBook.UpdateBook()
			}

		}
		
	}

	http.Redirect(w, r, "/interact?msg=Book updated", http.StatusFound)
}

func apiBookGetAll(w http.ResponseWriter, 
	r *http.Request) {
	newBooks := models.GetAllBooks()
	res, _ := json.Marshal(newBooks)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func apiBookDelete(w http.ResponseWriter, 
	r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 0, 0)
	if err == nil {
		book := models.DeleteBook(ID)
		res, _ := json.Marshal(book)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write(nil)
}
func ParseBody(r *http.Request, x interface{}) {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}

func apiBookCreate(w http.ResponseWriter, 
	r *http.Request) {
	newBook := &models.Book{}
	ParseBody(r, newBook)
	crBook := newBook.CreateBook()

	res, _ := json.Marshal(crBook)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func apiBookUpdate(w http.ResponseWriter, 
	r *http.Request) {
	newBook := &models.Book{}
	ParseBody(r, newBook)
	
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 10, 64)
	if err == nil {
		book, rowEffected := models.GetBook(ID)
		if rowEffected > 0{
			book.Title = newBook.Title
			book.Author = newBook.Author
			book.Publication = newBook.Publication
			book.UpdateBook()
			res, _ := json.Marshal(book)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(res)
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write(nil)
}

func main() {
	router := mux.NewRouter()
	http.Handle("/", router)
	fileServer:= http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static").Handler(http.StripPrefix("/static/resources", fileServer))

	router.HandleFunc("/", welcomeHandler)

	router.HandleFunc("/interact", interactHandler)
	router.HandleFunc("/new", newHandler)
	router.HandleFunc("/create", createHandler)
	router.HandleFunc("/delete", deleteHandler)
	router.HandleFunc("/edit", editHandler)
	router.HandleFunc("/doUpdate", updateHandler)

	router.HandleFunc("/api/book/getall", apiBookGetAll).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/book/delete/{bookId}", apiBookDelete).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/book/update/{bookId}", apiBookUpdate).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/book/add", apiBookCreate).Methods("POST", "OPTIONS")
	
	
	fmt.Println("Server started. Listening on port 8085")
	err := http.ListenAndServe("localhost:8085", nil)
	log.Fatal(err)
}