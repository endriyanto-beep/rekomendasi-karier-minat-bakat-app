package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Karir struct {
	Nama     string
	Kategori string
	Minat    []string
	Keahlian []string
	Gaji     int
}

type Rekomendasi struct {
	Nama       string
	Skor       int
	Persentase float64
	Gaji       int
}

type Pengguna struct {
	Nama     string
	Gender   string
	Lahir    string
	Minat    []string
	Keahlian []string
	Hasil    []Rekomendasi
}

func cocok(a, b []string) int {
	set := make(map[string]bool)
	for _, item := range a {
		set[strings.ToLower(strings.TrimSpace(item))] = true
	}
	hit := 0
	for _, item := range b {
		if set[strings.ToLower(strings.TrimSpace(item))] {
			hit++
		}
	}
	return hit
}

func skor(k Karir, minat, skill []string) int {
	return cocok(minat, k.Minat) + cocok(skill, k.Keahlian)
}

func persen(k Karir, minat, skill []string) float64 {
	total := len(k.Minat) + len(k.Keahlian)
	if total == 0 {
		return 0.0
	}
	match := cocok(minat, k.Minat) + cocok(skill, k.Keahlian)
	return float64(match) / float64(total) * 100
}

func urutSkor(r []Rekomendasi) {
	sort.SliceStable(r, func(i, j int) bool {
		return r[i].Skor > r[j].Skor
	})
}

func urutGaji(r []Rekomendasi) {
	sort.SliceStable(r, func(i, j int) bool {
		return r[i].Gaji > r[j].Gaji
	})
}

func showRekom(r []Rekomendasi, mode string) {
	if len(r) == 0 {
		fmt.Println("Belum ada rekomendasi.")
		return
	}

	fmt.Println("\nRekomendasi Utama:")
	fmt.Printf("- %s", r[0].Nama)
	switch mode {
	case "skor":
		fmt.Printf(" | Skor: %d\n", r[0].Skor)
	case "persentase":
		fmt.Printf(" | Kecocokan: %.2f%%\n", r[0].Persentase)
	case "gaji":
		fmt.Printf(" | Gaji: Rp%d\n", r[0].Gaji)
	}

	if len(r) > 1 {
		fmt.Println("\nAlternatif Lain:")
		for _, alt := range r[1:] {
			fmt.Printf("- %s", alt.Nama)
			switch mode {
			case "skor":
				fmt.Printf(" | Skor: %d\n", alt.Skor)
			case "persentase":
				fmt.Printf(" | Kecocokan: %.2f%%\n", alt.Persentase)
			case "gaji":
				fmt.Printf(" | Gaji: Rp%d\n", alt.Gaji)
			}
		}
	}
}

func listMinatKeahlian(k []Karir) ([]string, []string) {
	mMap, kMap := make(map[string]bool), make(map[string]bool)

	for _, karir := range k {
		for _, m := range karir.Minat {
			mMap[m] = true
		}
		for _, s := range karir.Keahlian {
			kMap[s] = true
		}
	}

	var minat, keahlian []string
	for m := range mMap {
		minat = append(minat, m)
	}
	for s := range kMap {
		keahlian = append(keahlian, s)
	}
	sort.Strings(minat)
	sort.Strings(keahlian)

	fmt.Println("\nDaftar Minat:")
	for i, m := range minat {
		fmt.Printf("%d. %s\n", i+1, m)
	}
	fmt.Println("\nDaftar Keahlian:")
	for i, s := range keahlian {
		fmt.Printf("%d. %s\n", i+1, s)
	}

	return minat, keahlian
}

func showKarir(k []Karir) {
	fmt.Println("\nDaftar Karir:")
	for i, v := range k {
		fmt.Printf("%d. %s (%s)\n", i+1, v.Nama, v.Kategori)
	}
}

func isiData(in *bufio.Reader, daftar []Karir, awal Pengguna) Pengguna {
	p := awal

	fmt.Println("Nama:", p.Nama)
	input, _ := in.ReadString('\n')
	if s := strings.TrimSpace(input); s != "" {
		p.Nama = s
	}

	fmt.Println("Jenis Kelamin:", p.Gender)
	input, _ = in.ReadString('\n')
	if s := strings.TrimSpace(input); s != "" {
		p.Gender = s
	}

	fmt.Println("Tanggal Lahir (DD-MM-YYYY):", p.Lahir)
	input, _ = in.ReadString('\n')
	if s := strings.TrimSpace(input); s != "" {
		p.Lahir = s
	}

	minatList, skillList := listMinatKeahlian(daftar)
	p.Minat, p.Keahlian = nil, nil

	fmt.Println("\nNomor minat dipilih (pisahkan dengan koma):")
	input, _ = in.ReadString('\n')
	for _, s := range strings.Split(strings.TrimSpace(input), ",") {
		if idx, err := strconv.Atoi(strings.TrimSpace(s)); err == nil && idx >= 1 && idx <= len(minatList) {
			p.Minat = append(p.Minat, minatList[idx-1])
		}
	}

	fmt.Println("Nomor keahlian dipilih (pisahkan dengan koma):")
	input, _ = in.ReadString('\n')
	for _, s := range strings.Split(strings.TrimSpace(input), ",") {
		if idx, err := strconv.Atoi(strings.TrimSpace(s)); err == nil && idx >= 1 && idx <= len(skillList) {
			p.Keahlian = append(p.Keahlian, skillList[idx-1])
		}
	}

	for _, k := range daftar {
		nilai := skor(k, p.Minat, p.Keahlian)
		if nilai > 0 {
			p.Hasil = append(p.Hasil, Rekomendasi{
				Nama:       k.Nama,
				Skor:       nilai,
				Persentase: persen(k, p.Minat, p.Keahlian),
				Gaji:       k.Gaji,
			})
		}
	}

	urutSkor(p.Hasil)
	return p
}
func updateRekomendasi(p *Pengguna, daftar []Karir) {
	p.Hasil = nil

	for _, k := range daftar {
		nilai := skor(k, p.Minat, p.Keahlian)
		if nilai > 0 {
			p.Hasil = append(p.Hasil, Rekomendasi{
				Nama:       k.Nama,
				Skor:       nilai,
				Persentase: persen(k, p.Minat, p.Keahlian),
				Gaji:       k.Gaji,
			})
		}
	}

	urutSkor(p.Hasil)
}

func main() {
	karirList := []Karir{
		{"Software Engineer", "Teknologi", []string{"pemrograman", "logika"}, []string{"golang", "java", "oop"}, 17000000},
		{"Data Analyst", "Teknologi", []string{"data", "analisis"}, []string{"sql", "excel", "python"}, 15000000},
		{"Graphic Designer", "Desain", []string{"desain", "visual"}, []string{"photoshop", "illustrator"}, 10000000},
		{"Public Relations Specialist", "Komunikasi", []string{"komunikasi", "hubungan publik"}, []string{"public speaking", "menulis"}, 11000000},
		{"Project Manager", "Manajemen", []string{"planning", "organisasi"}, []string{"scrum", "timeline"}, 14000000},
		{"Sales Executive", "Penjualan", []string{"negosiasi", "target"}, []string{"komunikasi", "presentasi"}, 13000000},
		{"Content Writer", "Media", []string{"menulis", "konten"}, []string{"seo", "blogging"}, 9000000},
		{"Research Scientist", "Riset", []string{"penelitian", "analisa"}, []string{"statistik", "laboratorium"}, 15500000},
		{"Nurse", "Kesehatan", []string{"merawat", "kesehatan"}, []string{"perawatan", "medis"}, 9500000},
		{"Teacher", "Pendidikan", []string{"mengajar", "berbagi ilmu"}, []string{"presentasi", "kurikulum"}, 8500000},
		{"Lawyer", "Hukum", []string{"regulasi", "analisa"}, []string{"dokumen hukum", "negosiasi"}, 16000000},
		{"Politician", "Politik", []string{"politik", "strategi"}, []string{"public speaking", "kampanye"}, 15500000},
		{"Environmental Scientist", "Lingkungan", []string{"alam", "riset"}, []string{"pengamatan", "data lingkungan"}, 13500000},
		{"Farmer", "Pertanian", []string{"berkebun", "tanaman"}, []string{"pengolahan lahan", "irigasi"}, 7000000},
		{"Logistic Coordinator", "Logistik", []string{"logistik", "manajemen"}, []string{"supply chain", "pengiriman"}, 11500000},
		{"Tour Guide", "Pariwisata", []string{"jalan-jalan", "komunikasi"}, []string{"bahasa", "interpersonal"}, 8500000},
		{"Hotel Manager", "Pariwisata", []string{"hospitality", "manajemen"}, []string{"koordinasi", "operasional"}, 12000000},
		{"Cybersecurity Analyst", "Teknologi", []string{"keamanan", "jaringan"}, []string{"pentest", "firewall"}, 16500000},
		{"AI Specialist", "Teknologi", []string{"machine learning", "logika"}, []string{"python", "tensorflow"}, 18000000},
		{"Digital Marketer", "Pemasaran", []string{"konten", "strategi"}, []string{"ads", "analytics"}, 12000000},
		{"Financial Analyst", "Keuangan", []string{"investasi", "data"}, []string{"excel", "akuntansi"}, 14000000},
		{"UX/UI Designer", "Desain", []string{"pengalaman pengguna", "desain"}, []string{"figma", "wireframe"}, 11000000},
		{"Mobile App Developer", "Teknologi", []string{"app", "programming"}, []string{"kotlin", "swift"}, 15000000},
		{"Biomedical Engineer", "Kesehatan", []string{"teknologi", "biomedis"}, []string{"alat kesehatan", "rancang bangun"}, 14500000},
		{"Mechanical Engineer", "Teknik", []string{"mesin", "rancang bangun"}, []string{"autocad", "solidworks"}, 13500000},
		{"Civil Engineer", "Teknik", []string{"struktur", "konstruksi"}, []string{"autocad", "sketsa"}, 14000000},
		{"Architect", "Teknik", []string{"desain", "bangunan"}, []string{"sketchup", "autocad"}, 15000000},
		{"Chef", "Kuliner", []string{"memasak", "kreativitas"}, []string{"resep", "plating"}, 10000000},
		{"Pharmacist", "Kesehatan", []string{"obat", "kesehatan"}, []string{"farmasi", "resep"}, 13000000},
		{"Psychologist", "Kesehatan", []string{"empati", "observasi"}, []string{"konseling", "psikotes"}, 12500000},
	}

	var user Pengguna
	riwayat := []Pengguna{
		{
			Nama:     "Dina",
			Gender:   "Perempuan",
			Lahir:    "12-08-2000",
			Minat:    []string{"pemrograman", "logika"},
			Keahlian: []string{"golang", "oop"},
			Hasil: []Rekomendasi{
				{"Software Engineer", 3, 75.0, 17000000},
				{"Cybersecurity Analyst", 2, 50.0, 16500000},
			},
		},
		{
			Nama:     "Rizky",
			Gender:   "Laki-laki",
			Lahir:    "04-04-1999",
			Minat:    []string{"desain", "visual"},
			Keahlian: []string{"photoshop", "illustrator"},
			Hasil: []Rekomendasi{
				{"Graphic Designer", 4, 100.0, 10000000},
			},
		},
		{
			Nama:     "Sari",
			Gender:   "Perempuan",
			Lahir:    "21-11-2001",
			Minat:    []string{"menulis", "konten"},
			Keahlian: []string{"seo", "blogging"},
			Hasil: []Rekomendasi{
				{"Content Writer", 4, 100.0, 9000000},
				{"Digital Marketer", 3, 60.0, 12000000},
			},
		},
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== APLIKASI REKOMENDASI KARIR ===")
		fmt.Println("1. Masukkan Data Pengguna Baru")
		fmt.Println("2. Lihat Rekomendasi Karir")
		fmt.Println("3. Rekomendasi Berdasarkan Persentase")
		fmt.Println("4. Rekomendasi Berdasarkan Gaji")
		fmt.Println("5. Tampilkan Semua Karir")
		fmt.Println("6. Lihat/Kelola Riwayat Pengguna")
		fmt.Println("7. Keluar")
		fmt.Print(">> ")

		pilihan, _ := reader.ReadString('\n')
		pilihan = strings.TrimSpace(pilihan)

		switch pilihan {
		case "1":
			user = isiData(reader, karirList, Pengguna{})
			updateRekomendasi(&user, karirList)
			riwayat = append(riwayat, user)
			fmt.Println("Data pengguna berhasil disimpan")

		case "2":
			if user.Nama == "" {
				fmt.Println("Belum ada data pengguna. Pilih menu 1 terlebih dahulu.")
				continue
			}
			updateRekomendasi(&user, karirList)
			urutSkor(user.Hasil)
			showRekom(user.Hasil, "skor")

		case "3":
			if user.Nama == "" {
				fmt.Println("Belum ada data pengguna. Pilih menu 1 terlebih dahulu.")
				continue
			}
			updateRekomendasi(&user, karirList)
			sort.SliceStable(user.Hasil, func(i, j int) bool {
				return user.Hasil[i].Persentase > user.Hasil[j].Persentase
			})
			showRekom(user.Hasil, "persentase")

		case "4":
			if user.Nama == "" {
				fmt.Println("Belum ada data pengguna. Pilih menu 1 terlebih dahulu.")
				continue
			}
			updateRekomendasi(&user, karirList)
			urutGaji(user.Hasil)
			showRekom(user.Hasil, "gaji")

		case "5":
			showKarir(karirList)

		case "6":
			if len(riwayat) == 0 {
				fmt.Println("Belum ada riwayat pengguna.")
				continue
			}
			for {
				fmt.Println("\nRiwayat Pengguna:")
				for i, p := range riwayat {
					fmt.Printf("%d. Nama: %s\n", i+1, p.Nama)
					fmt.Printf("   Jenis Kelamin: %s\n", p.Gender)
					fmt.Printf("   Tanggal Lahir: %s\n", p.Lahir)
					fmt.Printf("   Minat: %s\n", strings.Join(p.Minat, ", "))
					fmt.Printf("   Keahlian: %s\n", strings.Join(p.Keahlian, ", "))
					if len(p.Hasil) > 0 {
						fmt.Printf("   Rekomendasi Utama: %s (Skor: %d)\n", p.Hasil[0].Nama, p.Hasil[0].Skor)
					} else {
						fmt.Printf("   Rekomendasi Utama: Tidak tersedia\n")
					}
				}

				fmt.Println("1. Edit Pengguna")
				fmt.Println("2. Hapus Pengguna")
				fmt.Println("3. Kembali")
				fmt.Print(">> ")

				sub, _ := reader.ReadString('\n')
				sub = strings.TrimSpace(sub)

				if sub == "1" {
					fmt.Print("Masukkan nomor pengguna yang ingin diedit: ")
					nomorStr, _ := reader.ReadString('\n')
					nomor, err := strconv.Atoi(strings.TrimSpace(nomorStr))
					if err != nil || nomor < 1 || nomor > len(riwayat) {
						fmt.Println("Nomor tidak valid.")
						continue
					}

					namaSebelum := riwayat[nomor-1].Nama

					riwayat[nomor-1] = isiData(reader, karirList, riwayat[nomor-1])
					updateRekomendasi(&riwayat[nomor-1], karirList)

					if user.Nama == namaSebelum {
						user = riwayat[nomor-1]
					}

					fmt.Println("Data berhasil diperbarui.")
					break

				} else if sub == "2" {
					fmt.Print("Masukkan nomor pengguna yang ingin dihapus: ")
					nomorStr, _ := reader.ReadString('\n')
					nomor, err := strconv.Atoi(strings.TrimSpace(nomorStr))
					if err != nil || nomor < 1 || nomor > len(riwayat) {
						fmt.Println("Nomor tidak valid.")
						continue
					}

					nama := riwayat[nomor-1].Nama
					riwayat = append(riwayat[:nomor-1], riwayat[nomor:]...)
					fmt.Printf("Data pengguna '%s' berhasil dihapus.\n", nama)
					break

				} else if sub == "3" {
					break

				} else {
					fmt.Println("Pilihan tidak dikenali.")
				}
			}

		case "7":
			fmt.Println("Terima kasih telah menggunakan aplikasi.")
			return

		default:
			fmt.Println("Pilihan tidak dikenali.")
		}
	}
}
