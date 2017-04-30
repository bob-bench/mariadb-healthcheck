/**
 * I am a very very basic health check to be used bu the health probe
 * of the Microsoft Azure LoadBalancer
 */
package main

import "database/sql"
import "fmt"
import _ "github.com/go-sql-driver/mysql"
import "net/http"
import "time"
import "log"
import "context"

// Connect to MariaDB Galera and check if this specific node is
// considered ready and accepting writes
func checkMysql(ctx context.Context, user string, password string, host string, port int) (int, string) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, password, host, port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Print(err, "\n")
		return 500, "Failed to open MySQL"
	}

	// When the node returns a value of ON it can accept write-sets from the cluster
	rows, err := db.QueryContext(ctx, "SHOW GLOBAL STATUS LIKE 'wsrep_ready'")
	if err != nil {
		db.Close()
		fmt.Print(err, "\n")
		return 500, "Failed to query wsrep_ready"
	}
	if !rows.Next() {
		rows.Close()
		db.Close()
		return 500, "No query result"
	}
	var key string
	var val string
	rows.Scan(&key, &val)
	rows.Close()
	db.Close()

	if val == "ON" {
		return 200, "OK"
	} else {
		return 500, val
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/selfcheck", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// This drives the health check..
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ctx, cancel := context.WithTimeout(ctx, 1 * time.Second)
		code, reason := checkMysql(ctx, "user", "pass", "localhost", 3306)
		cancel()
		w.WriteHeader(code)
		w.Write([]byte(reason))
	})

	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
