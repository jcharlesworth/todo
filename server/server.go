// Implements a simple HTTP server providing a REST API to a task handler.
//
//		GET		/task/			Retrieves all the tasks.
//		POST	/task/			Creates a new task given a title.
//		GET		/task/{taskID}	Retrieves the task with the given id.
//		PUT		/task/{taskID}	Updates the task with the given id.
//

package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jcharlesworth/todo/task"
	"github.com/gorilla/mux"	
)

var tasks = task.NewTaskManager()

const PathPrefix = "/task/"

func RegisterHandlers() {
	r := mux.NewRouter()
	r.HandleFunc(PathPrefix, errorHandler(ListTasks)).Methods("GET")
	r.HandleFunc(PathPrefix, errorHandler(NewTask)).Methods("POST")
	r.HandleFunc(PathPrefix+"{id}", errorHandler(GetTask)).Methods("GET")
	r.HandleFunc(PathPrefix+"{id}", errorHandler(UpdateTask)).Methods("PUT")
	http.Handle(PathPrefix, r)
}

// badRequest is handled by setting the status code in the reply to StatusBadRequest.
type badRequest struct{ error }

// notFound is handled by setting the status code in the reply to StatusNotFound
type notFound struct{ error }

// errorHandler wraps a function returning an error by handling the error and
// returning the http.Handler.
// If the error is a type defined above, it is handled as descrived for every type.
// If the error is another type, it is considered an internal error and its message is logged.
func errorHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		switch err.(type) {
		case badRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case notFound:
			http.Error(w, "task not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(w, "oops", http.StatusInternalServerError)
		}
	}
}

// Returns an object with a Tasks field containing a list of tasks.
// Example:
// 		req: GET /task/
//		res: 200 {"Tasks": [
//				{"ID": 1, "Title": "Learn Go", "Done": false},
//				{"ID": 2, "Title": "Buy bread", "Done": true}
//			]}
func ListTasks(w http.ResponseWriter, r *http.Request) error {
	res := struct{ Tasks []*task.Task}{tasks.All()}
	return json.NewEncoder(w).Encode(res)
}

func NewTask(w http.ResponseWriter, r *http.Request) error {
	req := struct{ Title string }{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequest{err}
	}
	t, err := task.NewTask(req.Title)
	if err != nil {
		return badRequest{err}
	}
	return tasks.Save(t)
}

// Obtain the id variable from the given request url,
// parses the text and returns the result.
func parseID(r *http.Request) (int64, error) {
	txt, ok := mux.Vars(r)["id"]
	if !ok {
		return 0, fmt.Errorf("task id not found")
	}
	return strconv.ParseInt(txt, 10, 0)
}

func GetTask(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	log.Println("task is ", id)
	if err != nil {
		return badRequest{err}
	}
	t, ok := tasks.Find(id)
	log.Println("Found", ok)

	if !ok {
		return notFound{}
	}
	return json.NewEncoder(w).Encode(t)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return badRequest{err}
	}
	var t task.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return badRequest{err}
	}
	if t.ID != id {
		return badRequest{fmt.Errorf("inconsistent task IDs")}
	}
	if _, ok := tasks.Find(id); !ok {
		return notFound{}
	}
	return tasks.Save(&t)
}