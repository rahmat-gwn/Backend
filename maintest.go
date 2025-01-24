package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Struct: Bentuk
type Bentuk interface {
	Luas() float64
}

// Struct: Persegi
type Persegi struct {
	Sisi float64
}

func (p Persegi) Luas() float64 {
	return p.Sisi * p.Sisi
}

// Struct: Lingkaran
type Lingkaran struct {
	JariJari float64
}

func (l Lingkaran) Luas() float64 {
	return math.Pi * l.JariJari * l.JariJari
}

// Fungsi untuk validasi input
func validasiInput(input string) (float64, error) {
	value, err := strconv.ParseFloat(strings.TrimSpace(input), 64)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("input tidak valid, masukkan angka positif")
	}
	return value, nil
}

// Fungsi untuk menyimpan hasil ke file
func simpanKeFile(namaFile string, data string) error {
	file, err := os.OpenFile(namaFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data + "\n")
	return err
}

// Fungsi utama
func main() {
	fmt.Println("=== Aplikasi Perhitungan Luas ===")
	fmt.Println("1. Hitung luas persegi")
	fmt.Println("2. Hitung luas lingkaran")
	fmt.Print("Pilih opsi (1/2): ")

	// Membaca input dari pengguna
	reader := bufio.NewReader(os.Stdin)
	opsiInput, _ := reader.ReadString('\n')
	opsi := strings.TrimSpace(opsiInput)

	var bentuk Bentuk
	switch opsi {
	case "1":
		fmt.Print("Masukkan panjang sisi persegi: ")
		sisiInput, _ := reader.ReadString('\n')
		sisi, err := validasiInput(sisiInput)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		bentuk = Persegi{Sisi: sisi}

	case "2":
		fmt.Print("Masukkan jari-jari lingkaran: ")
		jariInput, _ := reader.ReadString('\n')
		jariJari, err := validasiInput(jariInput)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		bentuk = Lingkaran{JariJari: jariJari}

	default:
		fmt.Println("Opsi tidak valid. Program keluar.")
		return
	}

	// Menghitung luas dengan goroutine
	resultChannel := make(chan string)
	go func(b Bentuk) {
		luas := b.Luas()
		resultChannel <- fmt.Sprintf("Luas: %.2f (Dihitung pada %s)", luas, time.Now().Format(time.RFC1123))
	}(bentuk)

	// Menampilkan hasil
	result := <-resultChannel
	fmt.Println(result)

	// Menyimpan hasil ke file
	namaFile := "hasil_luas.txt"
	id := uuid.New()
	data := fmt.Sprintf("%s | UUID: %s", result, id)
	if err := simpanKeFile(namaFile, data); err != nil {
		fmt.Println("Gagal menyimpan hasil ke file:", err)
	} else {
		fmt.Println("Hasil berhasil disimpan ke file:", namaFile)
	}
}
