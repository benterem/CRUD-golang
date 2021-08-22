package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"regexp"
	"strconv"
	//    "github.com/joho/godotenv"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "square8_db"
)

type Container struct {
	ID          uint64 `json:"id"`
	ContainerId string `json:"containerId"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Size        string `json:"size"`
	Deleted     bool   `json:"-"`
}

type response struct {
	Message string `json:"message,omitempty"`
}

// create connection with db
func connectDb() *sql.DB {
	stmt := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// conect to DB
	db, err := sql.Open("postgres", stmt)

	if err != nil {
		fmt.Printf(err.Error())
	}

	// check connection
	err = db.Ping()
	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Printf("connected to postgress db %s on port %d", dbname, port)
	return db
}

//get all containers
func GetContainers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := connectDb()
	defer db.Close()
	rows, err := db.Query(`SELECT * FROM "containers"`)
    defer rows.Close()
	
    if err != nil {
        errorResponse("Couldn't access DB", http.StatusInternalServerError, w)
		return
	}

	var containers []Container
	for rows.Next() {
		var c Container
		err = rows.Scan(&c.ID, &c.ContainerId, &c.Type, &c.Status, &c.Size, &c.Deleted)
		if err != nil {
			errorResponse(err.Error(), http.StatusInternalServerError, w)
			return
		}
		if !c.Deleted {
			containers = append(containers, c)
		}
	}
	//response
	json.NewEncoder(w).Encode(containers)
}

// get specific container by id
func GetContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	idStr, err := strconv.Atoi(params["id"])
	if err != nil {
		errorResponse("Invalid id", http.StatusBadRequest, w)
		return
	}

	id := uint64(idStr)

	db := connectDb()
	defer db.Close()

	row := db.QueryRow(`SELECT * FROM "containers" WHERE id=$1`, id)

	var container Container
	err = row.Scan(&container.ID, &container.ContainerId, &container.Type, &container.Status, &container.Size, &container.Deleted)
	if err != nil {
		errorResponse(err.Error(), http.StatusInternalServerError, w)
		return
	}

	if container.Deleted {
		errorResponse("container was previously deleted", http.StatusBadRequest, w)
		return
	}

	json.NewEncoder(w).Encode(container)

}

// create a new container
func CreateContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

    // decode request body into Container struct
	var container Container
	_ = json.NewDecoder(r.Body).Decode(&container)

	// add to db
	db := connectDb()
	defer db.Close()

	// validate that containerId is composed of 9 latin letters and digits
	matched, err := regexp.MatchString(`^[A-Z0-9]{9}$`, container.ContainerId)
	if !matched || err != nil {
		errorResponse("A valid container Id must be provided: valid IDs are composed of A-Z and 0-9", http.StatusBadRequest, w)
		return
	}

	// validate that type is only "Full" or "Empty"
	switch container.Type {
	case "Full":
	case "Empty":
	default:
		errorResponse("Please specify if container is 'Full' or 'Empty'", http.StatusBadRequest, w)
		return
	}

	var id uint64
	stmt := `insert into "containers"("container_id", "type", "status", "size") values($1, $2, $3, $4) returning id`
	err = db.QueryRow(stmt, container.ContainerId, container.Type, container.Status, container.Size).Scan(&id)
	if err != nil {
		errorResponse(err.Error(), http.StatusInternalServerError, w)
		return
	}

	//response
	w.WriteHeader(http.StatusCreated)
	container.ID = id
	json.NewEncoder(w).Encode(container)
}

// delete a container
func DeleteContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	idStr, err := strconv.Atoi(params["id"])
	if err != nil {
		errorResponse("Invalid id", http.StatusBadRequest, w)
		return
	}

	id := uint64(idStr)

	db := connectDb()
	defer db.Close()

	stmt := `UPDATE "containers" SET deleted = TRUE WHERE id = $1`

	_, err = db.Exec(stmt, id)
	if err != nil {
		errorResponse("Container does not exit in DB", http.StatusBadRequest, w)
	}

	w.WriteHeader(http.StatusNoContent)
}

// update a container
func UpdateContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	idStr, err := strconv.Atoi(params["id"])
	if err != nil {
		errorResponse("Invalid id", http.StatusBadRequest, w)
		return
	}

	id := uint64(idStr)

	db := connectDb()
	defer db.Close()

	var container Container
	_ = json.NewDecoder(r.Body).Decode(&container)
	stmt := `UPDATE "containers" SET container_id=$2, type=$3, status=$4, size=$5 WHERE id=$1`
	_, err = db.Exec(stmt, id, container.ContainerId, container.Type, container.Status, container.Size)

	if err != nil {
		errorResponse("Container does not exist in DB", http.StatusBadRequest, w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func errorResponse(message string, statusCode int, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	if message != "" {
		json.NewEncoder(w).Encode(response{Message: message})
	}
}