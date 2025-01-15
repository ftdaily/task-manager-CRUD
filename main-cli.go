package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Task struct {
	ID          int
	Title       string
	Description string
	CreatedAt   string
}

var db *sql.DB

func init() {
	// Inisialisasi koneksi ke database
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/task_manager"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}
}

func tampilkanMenu() {
	fmt.Println("\nManajer Tugas CLI")
	fmt.Println("===================")
	fmt.Println("1. Lihat semua tugas")
	fmt.Println("2. Tambahkan tugas baru")
	fmt.Println("3. Perbarui tugas")
	fmt.Println("4. Hapus tugas")
	fmt.Println("5. Keluar")
	fmt.Print("Pilih opsi: ")
}

func bacaInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func lihatTugas() {
	rows, err := db.Query("SELECT id, title, description, created_at FROM tasks")
	if err != nil {
		log.Fatalf("Gagal mendapatkan data tugas: %v", err)
	}
	defer rows.Close()

	tugas := []Task{}
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt); err != nil {
			log.Fatalf("Gagal membaca data tugas: %v", err)
		}
		tugas = append(tugas, task)
	}

	if len(tugas) == 0 {
		fmt.Println("Tidak ada tugas yang ditemukan.")
		return
	}

	fmt.Println("\nDaftar Tugas:")
	for _, task := range tugas {
		fmt.Printf("ID: %d | Judul: %s | Deskripsi: %s | Dibuat Pada: %s\n",
			task.ID, task.Title, task.Description, task.CreatedAt)
	}
}

func tambahTugas() {
	judul := bacaInput("Masukkan judul tugas: ")
	if judul == "" {
		fmt.Println("Judul tidak boleh kosong!")
		return
	}

	deskripsi := bacaInput("Masukkan deskripsi tugas: ")
	if deskripsi == "" {
		fmt.Println("Deskripsi tidak boleh kosong!")
		return
	}

	_, err := db.Exec("INSERT INTO tasks (title, description) VALUES (?, ?)", judul, deskripsi)
	if err != nil {
		log.Fatalf("Gagal menambahkan tugas: %v", err)
	}

	fmt.Println("Tugas berhasil ditambahkan!")
}

func perbaruiTugas() {
	id := bacaInput("Masukkan ID tugas yang akan diperbarui: ")
	judul := bacaInput("Masukkan judul baru: ")
	if judul == "" {
		fmt.Println("Judul tidak boleh kosong!")
		return
	}

	deskripsi := bacaInput("Masukkan deskripsi baru: ")
	if deskripsi == "" {
		fmt.Println("Deskripsi tidak boleh kosong!")
		return
	}

	_, err := db.Exec("UPDATE tasks SET title = ?, description = ? WHERE id = ?", judul, deskripsi, id)
	if err != nil {
		log.Fatalf("Gagal memperbarui tugas: %v", err)
	}

	fmt.Println("Tugas berhasil diperbarui!")
}

func hapusTugas() {
	id := bacaInput("Masukkan ID tugas yang akan dihapus: ")
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		log.Fatalf("Gagal menghapus tugas: %v", err)
	}

	fmt.Println("Tugas berhasil dihapus!")
}

func main() {
	for {
		tampilkanMenu()

		var pilihan int
		fmt.Scan(&pilihan)

		switch pilihan {
		case 1:
			lihatTugas()
		case 2:
			tambahTugas()
		case 3:
			perbaruiTugas()
		case 4:
			hapusTugas()
		case 5:
			fmt.Println("Sampai jumpa!")
			os.Exit(0)
		default:
			fmt.Println("Pilihan tidak valid. Silakan coba lagi.")
		}
	}
}
