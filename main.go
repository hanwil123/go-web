package go_web

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// Pengguna adalah struktur data pengguna
type User struct {
	ID       int    `db:"id" json:"id"`
	Nama     string `db:"nama" json:"nama"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type Database struct {
	*sqlx.DB
}

func main() {
	db := initDB()

	router := mux.NewRouter()
	router.HandleFunc("/api/register", registerHandler(db)).Methods("POST")
	router.HandleFunc("/api/login", loginHandler(db)).Methods("POST")
	router.HandleFunc("/api/logout", logoutHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func initDB() *Database {
	db, err := sqlx.Connect("mysql", "user:password@tcp(localhost:3306)/database")
	if err != nil {
		log.Fatal(err)
	}

	return &Database{db}
}

func registerHandler(db *Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pengguna User
		err := json.NewDecoder(r.Body).Decode(&pengguna)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validasi apakah email sudah terdaftar
		existingUser := User{}
		err = db.Get(&existingUser, "SELECT * FROM user WHERE email = ?", pengguna.Email)
		if err == nil {
			http.Error(w, "Email sudah terdaftar", http.StatusBadRequest)
			return
		} else if err != sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Enkripsi kata sandi
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pengguna.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pengguna.Password = string(hashedPassword)

		// Simpan pengguna ke database
		_, err = db.NamedExec("INSERT INTO user (nama, email, password) VALUES (:nama, :email, :password)", &pengguna)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pengguna)
	}
}

func loginHandler(db *Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pengguna User
		err := json.NewDecoder(r.Body).Decode(&pengguna)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Cari pengguna berdasarkan email
		existingUser := User{}
		err = db.Get(&existingUser, "SELECT * FROM user WHERE email = ?", pengguna.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Email atau kata sandi salah", http.StatusBadRequest)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Periksa kecocokan kata sandi
		err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(pengguna.Password))
		if err != nil {
			http.Error(w, "Email atau kata sandi salah", http.StatusBadRequest)
			return
		}

		// Generate token otentikasi dan kirim sebagai respons
		// ...

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingUser)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Hapus token otentikasi

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
