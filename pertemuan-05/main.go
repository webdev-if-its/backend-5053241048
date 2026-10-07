package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// TODO: lihat SOAL.md untuk kontrak lengkap tiap fungsi di bawah.
// Ganti setiap "panic" dengan implementasi yang benar. Tambahkan import
// sendiri kalau dibutuhkan. Gunakan sync.WaitGroup dengan Add/Done klasik
// (jangan wg.Go -- baru ada di Go 1.25).

// ErrUkuranTidakValid dikembalikan (lewat Hasil.Err) untuk file berukuran <= 0.
var ErrUkuranTidakValid = errors.New("ukuran file tidak valid")

// ErrBatasTidakValid dikembalikan UnduhBatas kalau batas paralel < 1.
var ErrBatasTidakValid = errors.New("batas paralel harus minimal 1")

// File adalah satu file yang akan "diunduh". Unduhan di tugas ini
// DISIMULASIKAN dengan time.Sleep -- tidak ada jaringan sungguhan.
type File struct {
	Nama     string
	UkuranKB int
	Durasi   time.Duration // lamanya simulasi unduhan
}

// Hasil adalah hasil satu unduhan.
type Hasil struct {
	Nama     string
	UkuranKB int
	Err      error
}

// Pengunduh adalah fungsi yang mengunduh satu file. Fungsi-fungsi di bawah
// menerimanya sebagai parameter (bukan memanggil Unduh langsung) supaya test
// bisa menyuntikkan unduhan palsu dan mengamati konkurensinya.
type Pengunduh func(File) Hasil

// Jalankan menjalankan fn di goroutine baru dan LANGSUNG kembali. (Level 1)
func Jalankan(fn func()) {
	go fn()
}

// JalankanN menjalankan fn(0) ... fn(n-1), masing-masing di goroutine sendiri,
// dan LANGSUNG kembali (tidak menunggu). (Level 2)
func JalankanN(n int, fn func(i int)) {
	for i := 0; i < n; i++ {
		go fn(i)
	}
}

// Unduh mensimulasikan satu unduhan. (Level 3)
func Unduh(f File) Hasil {
	if f.UkuranKB <= 0 {
		return Hasil{Nama: f.Nama, UkuranKB: f.UkuranKB, Err: ErrUkuranTidakValid}
	}
	time.Sleep(f.Durasi)
	return Hasil{Nama: f.Nama, UkuranKB: f.UkuranKB, Err: nil}
}

// UnduhBerurutan mengunduh file satu per satu -- pembanding. (Level 3)
func UnduhBerurutan(files []File) []Hasil {
	results := make([]Hasil, len(files))
	for i, f := range files {
		results[i] = Unduh(f)
	}
	return results
}

// UnduhSemua mengunduh semua file BERSAMAAN dengan sync.WaitGroup dan
// mengembalikan hasilnya SESUAI URUTAN files. (Level 4-6)
func UnduhSemua(files []File, unduh Pengunduh) []Hasil {
	hasil := make([]Hasil, len(files))
	wg := sync.WaitGroup{}
	for i, f := range files {
		wg.Add(1)
		go func(i int, f File) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					hasil[i] = Hasil{Nama: f.Nama, UkuranKB: f.UkuranKB, Err: fmt.Errorf("%v", r)}
				}
			}()
			hasil[i] = unduh(f)
		}(i, f)
	}
	wg.Wait()
	return hasil
}

// UnduhKeChannel mengunduh semua file bersamaan dan mengirim tiap hasil ke
// channel begitu selesai; channel ditutup setelah semuanya selesai. (Level 7)
func UnduhKeChannel(files []File, unduh Pengunduh) <-chan Hasil {
	ch := make(chan Hasil, len(files))
	wg := sync.WaitGroup{}
	for _, f := range files {
		wg.Add(1)
		go func(f File) {
			defer wg.Done()
			ch <- unduh(f)
		}(f)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}

// UnduhBatas seperti UnduhSemua tapi paling banyak `maks` unduhan berjalan
// pada saat yang sama. (Level 8)
func UnduhBatas(files []File, unduh Pengunduh, maks int) ([]Hasil, error) {
	if maks < 1 {
		return nil, ErrBatasTidakValid
	}
	hasil := make([]Hasil, len(files))
	sem := make(chan struct{}, maks)
	var wg sync.WaitGroup

	for i, f := range files {
		sem <- struct{}{}
		wg.Add(1)
		go func(i int, f File) {
			defer wg.Done()
			defer func() {
				<-sem
			}()
			hasil[i] = unduh(f)
		}(i, f)
	}
	wg.Wait()
	return hasil, nil
}

// GabungChannel menggabungkan banyak channel menjadi satu (fan-in); channel
// hasil ditutup setelah SEMUA channel masukan tertutup. (Level 9)
func GabungChannel(chs ...<-chan Hasil) <-chan Hasil {
	out := make(chan Hasil)
	var wg sync.WaitGroup
	for _, c := range chs {
		wg.Add(1)
		go func(c <- chan Hasil) {
			defer wg.Done()
			for h := range c {
				out <- h
			}
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// Ringkasan adalah potret progres unduhan. (Level 10)
type Ringkasan struct {
	Total   int
	Selesai int
	Gagal   int
	TotalKB int
}

// String memformat ringkasan untuk laporan. (Level 10)
func (r Ringkasan) String() string {
	return fmt.Sprintf("%d dari %d berhasil, %d gagal, total %d KB", r.Selesai, r.Total, r.Gagal, r.TotalKB)
}

// Progres mencatat progres unduhan dan AMAN dipakai banyak goroutine
// sekaligus. (Level 10)
type Progres struct {
	mu sync.Mutex
	r Ringkasan
}

// NewProgres membuat Progres untuk `total` unduhan. (Level 10)
func NewProgres(total int) *Progres {
	return &Progres{r: Ringkasan{Total: total}}
}

// Catat mencatat satu hasil unduhan. (Level 10)
func (p *Progres) Catat(h Hasil) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if h.Err != nil {
		p.r.Gagal++
		return
	}
	p.r.Selesai++
	p.r.TotalKB += h.UkuranKB
}

// Snapshot mengembalikan potret progres saat ini. (Level 10)
func (p *Progres) Snapshot() Ringkasan {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.r
}

func main() {
	fmt.Println("Concurrent File Downloader - pertemuan 5")

	files := []File{
		{"a.zip", 120, 300 * time.Millisecond},
		{"b.pdf", 80, 200 * time.Millisecond},
		{"c.mp4", 900, 500 * time.Millisecond},
		{"rusak.bin", 0, 100 * time.Millisecond},
	}

	mulai := time.Now()
	hasil := UnduhSemua(files, Unduh)
	progres := NewProgres(len(files))
	for _, h := range hasil {
		progres.Catat(h)
	}
	fmt.Println(progres.Snapshot())
	fmt.Println("selesai dalam", time.Since(mulai).Round(10*time.Millisecond))
}
