# backend-nrp

Repo tugas mata kuliah **Pengembangan Backend Dasar**, dibuat dari template [`webdev-if-its/backend-template`](https://github.com/webdev-if-its/backend-template). Ganti judul di atas jadi nama repo kalian sendiri (`backend-nrp`, contoh: `backend-5025201012`).

## Aturan Umum

- Tugas tiap pertemuan disimpan di folder `pertemuan-XX/` pada repo ini.
- Commit message wajib menyebut level yang dicapai: `pertemuan-XX: level N selesai`.
- Deadline push: sebelum pertemuan berikutnya dimulai.
- Semua level dicek otomatis lewat `go test` — baca `pertemuan-XX/SOAL.md` tiap minggu untuk detail levelnya.

## Mengambil Pertemuan Baru Tiap Minggu

Repo ini **tidak otomatis sinkron** dengan template dosen. Begitu ada pertemuan baru, jalankan (ganti `pertemuan-02` sesuai minggu berjalan):

```bash
git fetch https://github.com/webdev-if-its/backend-template.git main
git checkout FETCH_HEAD -- pertemuan-02
```

Perintah ini **aman dijalankan kapan pun** — tidak akan menimpa folder pertemuan lain yang sudah kalian kerjakan, karena hanya mengambil folder yang disebutkan. Setelah itu, commit folder barunya seperti biasa.

Kalau dosen memperbaiki sesuatu di pertemuan yang sudah dirilis (mis. ada bug di test), biasanya cukup ambil ulang file yang diperbaiki saja, bukan seluruh folder — akan diumumkan file mana yang berubah.

---

Bagian di bawah ini **isi bertahap** sesuai level yang sedang kalian kerjakan (lihat `pertemuan-01/SOAL.md`) — heading-nya dicek otomatis, jangan diganti namanya.

## Identitas
- Nama: Alfarel Sandriano Subektiansyah
- NRP: 5053241048
- Kelas: M

## Commit vs Push
Commit merupakan operasi untuk menyimpan perubahan yang telah dibuat ke db git lokal (`.git`).
Sedangkan, push merupakan operasi untuk mengirim perubahan ke remote repository (misalnya GitHub).

## Reproducibility
Reproducible berarti apabila rekan tim menjalankan program Go di perangkatnya sendiri, hasilnya harus identik dengan apabila kita menjalankan program Go di perangkat kita sendiri. Entah itu hasil binary, test, maupun runtime behavior.

Apabila rekan tim memiliki versi Go yang berbeda, hasilnya tergantung perbedaan versi yang dimiliki oleh rekan tim. Apabila versi yang berbeda merupakan versi minor, maka tidak akan ada masalah serius yang muncul. Masalah serius akan muncul apabila perbedaan versi Go merupakan major. Contoh masalah yang dihasilkan antara lain:
- perbedaan semantik loop variable (Go 1.21 vs Go 1.22)
- ada versi minimum di `go.mod` tetapi toolchain lokal lebih rendah

## Catatan Merge Conflict
<<<<<<< HEAD
	Sapa("test")
=======
	Sapa(nama)
>>>>>>> staging

merge conflict terjadi karena ada line yang sama dengan versi yang berbeda di kedua branch.

Versi yang benar harusnya yang staging. Untuk memperbaiki conflict, hapus semua line yang tidak diperlukan termasuk `Sapa("test")` dan conflict marker.

## Kenapa .gitignore Penting
`.gitignore` itu penting karena baik di environment testing maupun production, akan ada hal-hal yang harus kita sembunyikan ketika kita melakukan push ke remote server. Misal menyembunyikan `.env` yang terdapat API key third party service tertentu. Dan apabila file build/IDE ikut ter-push ke remote server, dan rekan tim kita mendapatkan copy-an file tersebut, maka bisa saja:
1. build error ketika rekan tim melakukan build program Go
2. konfigurasi IDE rekan tim ter-override oleh file IDE kita

## Refleksi
Sejauh ini, soal-soal yang diberikan cukup menjelaskan hal-hal yang awalnya saya tidak memahami. Dan tidak ada kesulitan yang signifikan ketika mengerjakan soal-soal pertemuan 1.
