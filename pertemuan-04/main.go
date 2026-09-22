package main

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

// TODO: lihat SOAL.md untuk kontrak lengkap tiap fungsi/method di bawah.
// Ganti setiap "panic" dengan implementasi yang benar. Tambahkan import
// (mis. "strings") sendiri kalau memang dibutuhkan.

// ErrTugasTidakDitemukan dikembalikan ketika ID tugas tidak ada.
var ErrTugasTidakDitemukan = errors.New("tugas tidak ditemukan")

// ErrInputKosong dikembalikan ketika judul kosong (atau hanya spasi).
var ErrInputKosong = errors.New("input tidak boleh kosong")

// ErrDaftarKosong dikembalikan oleh Max ketika slice-nya kosong.
var ErrDaftarKosong = errors.New("daftar kosong")

// Timestamps dipakai lewat embedding (composition) di Task.
type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Touch menyetel UpdatedAt ke waktu sekarang. (Level 9)
func (ts *Timestamps) Touch() {
	ts.UpdatedAt = time.Now()
}

// Task merepresentasikan satu tugas.
type Task struct {
	ID      int
	Judul   string
	Selesai bool
	Timestamps
}

// NewTask membuat Task baru dari judul. (Level 1)
func NewTask(judul string) (Task, error) {
	j := strings.TrimSpace(judul)
	if j == "" {
		return Task{}, ErrInputKosong
	}
	now := time.Now()
	time := Timestamps {
		CreatedAt: now,
		UpdatedAt: now,
	}
	return Task{
		Judul: j,
		Selesai: false,
		Timestamps: time,
	}, nil
}

// MarkDone menandai tugas selesai. (Level 2)
func (t *Task) MarkDone() {
	t.Selesai = true
	t.Touch()
}

// Rename mengganti judul tugas. (Level 3)
func (t *Task) Rename(judul string) error {
	j := strings.TrimSpace(judul)
	if j == "" {
		return ErrInputKosong
	}
	t.Judul = j
	t.Touch()
	return nil
}

// String membuat Task memenuhi fmt.Stringer. (Level 10)
func (t Task) String() string {
	panic("belum diimplementasikan")
}

// TaskStore adalah kontrak penyimpanan tugas. JANGAN diubah -- yang kalian
// tulis adalah implementasinya (MemoryStore).
type TaskStore interface {
	Add(t Task) (Task, error)
	Get(id int) (Task, error)
	List() []Task
	Delete(id int) error
}

// MemoryStore menyimpan tugas di memori.
type MemoryStore struct {
	// TODO: tambahkan field yang kalian butuhkan (mis. slice tugas dan penghitung ID)
	Tugas []Task
	NextID int
}

// Pastikan *MemoryStore memenuhi TaskStore -- kalau tidak, kode gagal
// dikompilasi di sini, bukan baru ketahuan saat dijalankan.
var _ TaskStore = (*MemoryStore)(nil)

// NewMemoryStore membuat store kosong.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		NextID: 1,
		Tugas:  []Task{},
	}
}

// Add menyimpan tugas dan memberinya ID baru. (Level 4)
func (m *MemoryStore) Add(t Task) (Task, error) {
	t.Judul = strings.TrimSpace(t.Judul)
	if t.Judul == "" {
		return Task{}, ErrInputKosong
	}
	t.ID = m.NextID + 1
	m.Tugas = append(m.Tugas, t)
	m.NextID++
	return t, nil
}

// Get mencari tugas menurut ID. (Level 4)
func (m *MemoryStore) Get(id int) (Task, error) {
	for _, v := range(m.Tugas) {
		if v.ID == id {
			return v, nil
		}
	}
	return Task{}, ErrTugasTidakDitemukan
}

// List mengembalikan seluruh tugas. (Level 5)
func (m *MemoryStore) List() []Task {
	if len(m.Tugas) == 0 {
		return []Task{}
	}
	dst := make([]Task, len(m.Tugas))
	copy(dst, m.Tugas)
	return dst
}

// Delete menghapus tugas menurut ID. (Level 6)
func (m *MemoryStore) Delete(id int) error {
	for i, v := range(m.Tugas) {
		if v.ID == id {
			m.Tugas = append(m.Tugas[:i], m.Tugas[i+1:]...)
			return nil
		}
	}
	return ErrTugasTidakDitemukan
}

// Filter mengembalikan elemen xs yang lolos pred. (Level 7)
func Filter[T any](xs []T, pred func(T) bool) []T {
	t := make([]T, 0)
	for _, v := range(xs) {
		if pred(v) {
			t = append(t, v)
		}
	}
	return t
}

// Map mengubah tiap elemen xs dengan f. (Level 7)
func Map[T, U any](xs []T, f func(T) U) []U {
	if len(xs) == 0 {
		return []U{}
	}
	var u []U
	for _, v := range(xs) {
		u = append(u, f(v))
	}
	return u
}

// Contains melaporkan apakah v ada di xs. (Level 8)
func Contains[T comparable](xs []T, v T) bool {
	if len(xs) == 0 {
		return false
	}
	for _, val := range(xs) {
		if val == v {
			return true
		}
	}
	return false
}

// Max mengembalikan elemen terbesar di xs. (Level 8)
func Max[T cmp.Ordered](xs []T) (T, error) {
	if len(xs) == 0 {
		var null T
		return null, ErrDaftarKosong
	}
	slices.Sort(xs)
	slices.Reverse(xs)
	return xs[0], nil
}

// Gabung menyambung String() tiap elemen dengan pemisah sep. (Level 10)
func Gabung[T fmt.Stringer](xs []T, sep string) string {
	panic("")
}

func main() {
	fmt.Println("Task Manager v1 - pertemuan 4")

	var store TaskStore = NewMemoryStore()
	for _, judul := range []string{"Belajar struct", "Belajar interface"} {
		t, err := NewTask(judul)
		if err != nil {
			fmt.Println("gagal membuat tugas:", err)
			continue
		}
		if _, err := store.Add(t); err != nil {
			fmt.Println("gagal menyimpan tugas:", err)
		}
	}
	for _, t := range store.List() {
		fmt.Println(t)
	}
}
