package main

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

// Struct: Definisi tipe data Orang
type Orang struct {
	Nama string
	Umur int
}

// Fungsi dengan Struct: Menampilkan informasi Orang
func (o Orang) Perkenalan() string {
	return fmt.Sprintf("Halo, nama saya %s dan saya berumur %d tahun.", o.Nama, o.Umur)
}

// Interface: Bentuk dengan metode Luas()
type Bentuk interface {
	Luas() float64
}

// Struct: Persegi dan Lingkaran
type Persegi struct {
	Sisi float64
}

func (p Persegi) Luas() float64 {
	return p.Sisi * p.Sisi
}

type Lingkaran struct {
	JariJari float64
}

func (l Lingkaran) Luas() float64 {
	return math.Pi * l.JariJari * l.JariJari
}

// Fungsi Utama
func main() {
	// Hello, World!
	fmt.Println("Hello, World!")

	// Variabel dan Tipe Data
	var angka int = 10
	teks := "Belajar Golang"
	fmt.Println("Angka:", angka, "| Teks:", teks)

	// Kondisi
	if angka > 5 {
		fmt.Println("Angka lebih besar dari 5")
	} else {
		fmt.Println("Angka kurang dari atau sama dengan 5")
	}

	// Perulangan
	fmt.Println("Perulangan angka 1-5:")
	for i := 1; i <= 5; i++ {
		fmt.Println(i)
	}

	// Fungsi
	hasilTambah := tambah(3, 7)
	fmt.Println("Hasil penjumlahan 3 + 7 =", hasilTambah)

	// Struct dan Interface
	orang := Orang{Nama: "Budi", Umur: 25}
	fmt.Println(orang.Perkenalan())

	persegi := Persegi{Sisi: 4}
	lingkaran := Lingkaran{JariJari: 7}

	cetakLuas(persegi)
	cetakLuas(lingkaran)

	// Menggunakan Pustaka Eksternal (UUID)
	id := uuid.New()
	fmt.Println("UUID unik:", id)

	// Menggunakan Pustaka Standar (time)
	fmt.Println("Waktu sekarang:", time.Now())
}

// Fungsi Tambahan
func tambah(a, b int) int {
	return a + b
}

// Fungsi dengan Parameter Interface
func cetakLuas(bentuk Bentuk) {
	fmt.Printf("Luas: %.2f\n", bentuk.Luas())
}
