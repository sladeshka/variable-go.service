package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Data struct {
	FileName     string `json:"file_name"`
	VariableName string `json:"variable_name"`
	Data         string `json:"data"`
}

type ApiError struct {
	Code    string `json:"code"`
	Data    string `json:"data"`
	Query1  string `json:"query1"`
	Query2  string `json:"query2"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func main() {
	variableDB, err := sql.Open("sqlite3", "./variable.db")
	if err != nil {
		log.Fatal(err)
	}
	defer func(variableDB *sql.DB) {
		err := variableDB.Close()
		if err != nil {

		}
	}(variableDB)

	_, err = variableDB.Exec(`CREATE TABLE IF NOT EXISTS variable (id INTEGER PRIMARY KEY, file_name TEXT, variable_name TEXT, data TEXT)`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/v1/variable", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			rData, err := io.ReadAll(req.Body)
			if err != nil {
				http.Error(res, "Error reading request body", http.StatusBadRequest)
				return
			}

			var data Data
			if err := json.Unmarshal(rData, &data); err != nil {
				http.Error(res, "Error unmarshalling JSON", http.StatusBadRequest)
				return
			}
			queryString := "INSERT OR REPLACE INTO variable (id, file_name, variable_name, data) VALUES ((SELECT id FROM variable WHERE file_name = ? AND variable_name = ?),?,?,?)"
			stmt, err := variableDB.Prepare(queryString)
			if err != nil {
				http.Error(res, "Error preparing SQL statement", http.StatusInternalServerError)
				return
			}

			_, err = stmt.Exec(strings.ToLower(data.FileName), data.VariableName, strings.ToLower(data.FileName), data.VariableName, data.Data)
			if err != nil {
				http.Error(res, "Error inserting variable into database", http.StatusInternalServerError)
				return
			}

			res.WriteHeader(http.StatusCreated)
		}

		if req.Method == "GET" {
			query := req.URL.Query()

			row := variableDB.QueryRow("SELECT file_name, variable_name, data FROM variable WHERE file_name = ? AND variable_name LIKE ?", strings.ToLower(query.Get("file_name")), query.Get("variable_name"))

			var data Data
			err := row.Scan(&data.FileName, &data.VariableName, &data.Data)
			if err != nil {
				res.Header().Set("Content-Type", "application/json")
				var apiError ApiError
				apiError.Code = "404"
				apiError.Data = err.Error()
				apiError.Query1 = query.Get("file_name")
				apiError.Query2 = query.Get("variable_name")
				apiError.Error = err.Error()
				apiError.Message = "Not found"
				err := json.NewEncoder(res).Encode(apiError)
				if err != nil {
					return
				}
				return
			}

			res.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(res).Encode(data)
			if err != nil {
				return
			}
			if err != nil {
				http.Error(res, "Error sending variable to client", http.StatusInternalServerError)
				return
			}
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
