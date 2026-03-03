package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq" // PostgreSQLドライバー
)

// 送られてくるデータの形を定義
type ZenRecord struct {
	Task     string `json:"task"`
	Duration int    `json:"duration"`
}

func main() {
	// データベース接続文字列
	connStr := "host=db port=5432 user=sana password=zenpassword dbname=auto_zen_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// テーブル拡張
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS zen_logs (id SERIAL PRIMARY KEY,  task TEXT,  duration INT, timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		log.Fatal(err)
	}

	// 3. データを保存するエンドポイント
	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POSTメゾットのみ受け付けます", http.StatusMethodNotAllowed)
			return
		}

		var record ZenRecord
		// 送られてきたJSONを解析して構造体に入れる
		err := json.NewDecoder(r.Body).Decode(&record)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// データベースに保存
		_, err = db.Exec("INSERT INTO zen_logs (task, duration) VALUES ($1, $2)", record.Task, record.Duration)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, `{"status": "success", "message": "Saved: %s for %d minutes!"}`, record.Task, record.Duration)
	})

	fmt.Println("Server is running on http://localhost:8081 ...")
	http.ListenAndServe(":8081", nil)
}
