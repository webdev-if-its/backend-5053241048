package main

import (
	"fmt"
	"slices"
)

// TODO(Level 1): lihat SOAL.md untuk kontrak lengkap tiap fungsi di bawah.
// Ganti setiap "panic" dengan implementasi yang benar.

func HitungSubtotal(qty int, hargaSatuan float64) float64 {
	return float64(qty) * hargaSatuan
}

func HitungTotalPesanan(qty []int, hargaSatuan []float64) float64 {
	if len(qty) != len(hargaSatuan) {
		return 0
	}
	var total_price float64
	for i := range qty {
		total_price += HitungSubtotal(qty[i], hargaSatuan[i])
	}
	return total_price
}

func TerapkanPajak(total float64, tarifPajak float64) float64 {
	return total + (total * tarifPajak)
}

func HitungDiskon(total float64) float64 {
	if total >= 1000000 {
		return total * .10
	} else if total >= 500000 && total < 1000000 {
		return total * .05
	}
	return 0
}

func TotalSetelahDiskon(qty []int, hargaSatuan []float64, tarifPajak float64) float64 {
	total := HitungTotalPesanan(qty, hargaSatuan)
	diskon := HitungDiskon(total)
	pajak := TerapkanPajak(total-diskon, tarifPajak)
	return pajak
}

func ValidasiPesanan(qty []int, hargaSatuan []float64) (bool, string) {
	slices.Sort(hargaSatuan[:])
	slices.Sort(qty[:])
	if len(qty) != len(hargaSatuan) || hargaSatuan[0] <= 0 || qty[0] <= 0 {
		return false, "pesanan tidak valid"
	}
	return true, ""
}

func TentukanStatus(total float64) string {
	if total > 1000000 {
		return "Prioritas"
	} else if total > 100000 && total <= 1000000 {
		return "Reguler"
	}
	return "Hemat"
}

func RingkasanPesanan(qty []int, hargaSatuan []float64, tarifPajak float64) string {
	total := HitungTotalPesanan(qty, hargaSatuan)
	return fmt.Sprintf("Subtotal pesanan: %f\nDiskon: %f\nTotal akhir: %f\nStatus:%s\n", total, HitungDiskon(total), TotalSetelahDiskon(qty, hargaSatuan, tarifPajak), TentukanStatus(total))
}

// TODO(Level 9): signature ini SUDAH benar (cari tahu sendiri kenapa
// bentuknya begini - lihat SOAL.md) - tinggal implementasikan isinya.
func Total(harga ...float64) float64 {
	var total float64
	for _, v := range harga {
		total += v
	}
	return total
}

// TODO(Level 10, bonus): signature ini SUDAH benar (cari tahu sendiri
// kenapa ada dua nilai balik - lihat SOAL.md) - tinggal implementasikan isinya.
func HitungOngkosKirim(beratKg float64, jarakKm float64) (float64, error) {
	if beratKg <= 0 || jarakKm <= 0 {
		return 0, fmt.Errorf("input tidak valid")
	}
	return (beratKg * 2000) + (jarakKm * 3000), nil
}

func main() {
	fmt.Println("Sales Order Processor - pertemuan 2")
}
