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

type DataPengguna struct {
	Nama         string
	JenisKelamin string
	TanggalLahir string
	Minat        []string
	Keahlian     []string
	HasilKarir   []Rekomendasi
}

func hitungKecocokan(daftarA, daftarB []string) int {
	set := make(map[string]bool)
	for i := 0; i < len(daftarA); i++ {
		item := daftarA[i]
		set[strings.ToLower(strings.TrimSpace(item))] = true
	}

	jumlah := 0
	for i := 0; i < len(daftarB); i++ {
		item := daftarB[i]
		if set[strings.ToLower(strings.TrimSpace(item))] {
			jumlah++
		}
	}
	return jumlah
}

func hitungSkor(karir Karir, minat, keahlian []string) int {
	return hitungKecocokan(minat, karir.Minat) + hitungKecocokan(keahlian, karir.Keahlian)
}

func hitungPersentase(karir Karir, minat, keahlian []string) float64 {
	total := len(karir.Minat) + len(karir.Keahlian)
	if total == 0 {
		return 0.0
	}
	terpenuhi := hitungKecocokan(minat, karir.Minat) + hitungKecocokan(keahlian, karir.Keahlian)
	return float64(terpenuhi) / float64(total) * 100
}

func urutkanBerdasarkanSkor(rek []Rekomendasi) {
	sort.SliceStable(rek, func(i, j int) bool {
		return rek[i].Skor > rek[j].Skor
	})
}

func urutkanBerdasarkanGaji(rek []Rekomendasi) {
	sort.SliceStable(rek, func(i, j int) bool {
		return rek[i].Gaji > rek[j].Gaji
	})
}

func tampilkanRekomendasi(rek []Rekomendasi, mode string) {
	if len(rek) == 0 {
		fmt.Println("Belum ada rekomendasi.")
		return
	}

	fmt.Println("\nRekomendasi Utama:")
	r := rek[0]
	switch mode {
	case "skor":
		fmt.Printf("- %s | Skor: %d\n", r.Nama, r.Skor)
	case "persentase":
		fmt.Printf("- %s | Kecocokan: %.2f%%\n", r.Nama, r.Persentase)
	case "gaji":
		fmt.Printf("- %s | Gaji: Rp%d\n", r.Nama, r.Gaji)
	}

	if len(rek) > 1 {
		fmt.Println("\nAlternatif Lain:")
		alternatif := rek[1:]
		for i := 0; i < len(alternatif); i++ {
			alt := alternatif[i]
			switch mode {
			case "skor":
				fmt.Printf("- %s | Skor: %d\n", alt.Nama, alt.Skor)
			case "persentase":
				fmt.Printf("- %s | Kecocokan: %.2f%%\n", alt.Nama, alt.Persentase)
			case "gaji":
				fmt.Printf("- %s | Gaji: Rp%d\n", alt.Nama, alt.Gaji)
			}
		}
	}
}

func tampilkanDaftarMinatKeahlian(daftar []Karir) ([]string, []string) {
	minatMap := make(map[string]bool)
	keahlianMap := make(map[string]bool)

	for i := 0; i < len(daftar); i++ {
		k := daftar[i]
		for j := 0; j < len(k.Minat); j++ {
			m := k.Minat[j]
			minatMap[m] = true
		}
		for j := 0; j < len(k.Keahlian); j++ {
			keahlian := k.Keahlian[j]
			keahlianMap[keahlian] = true
		}
	}

	var minatList, keahlianList []string
	for m := range minatMap {
		minatList = append(minatList, m)
	}
	for k := range keahlianMap {
		keahlianList = append(keahlianList, k)
	}
	sort.Strings(minatList)
	sort.Strings(keahlianList)

	fmt.Println("\n📌 Daftar Minat Tersedia:")
	for i, m := range minatList {
		fmt.Printf("%d. %s\n", i+1, m)
	}
	fmt.Println("\n🛠 Daftar Keahlian Tersedia:")
	for i, k := range keahlianList {
		fmt.Printf("%d. %s\n", i+1, k)
	}

	return minatList, keahlianList
}

func tampilkanListKarir(daftar []Karir) {
	fmt.Println("\n📋 Daftar Karir Tersedia:")
	for i, k := range daftar {
		fmt.Printf("%d. %s (%s)\n", i+1, k.Nama, k.Kategori)
	}
}

func prosesInputPengguna(reader *bufio.Reader, daftarKarir []Karir, dataAwal DataPengguna) DataPengguna {
	penggunaBaru := dataAwal

	fmt.Printf("Nama [%s]: ", penggunaBaru.Nama)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		penggunaBaru.Nama = input
	}

	fmt.Printf("Jenis Kelamin [%s]: ", penggunaBaru.JenisKelamin)
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		penggunaBaru.JenisKelamin = input
	}

	fmt.Printf("Tanggal Lahir (DD-MM-YYYY) [%s]: ", penggunaBaru.TanggalLahir)
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		penggunaBaru.TanggalLahir = input
	}

	minatList, keahlianList := tampilkanDaftarMinatKeahlian(daftarKarir)
	penggunaBaru.Minat = []string{}
	penggunaBaru.Keahlian = []string{}

	fmt.Println("\nMasukkan nomor minat yang dipilih (pisahkan dengan koma): ")
	minatInput, _ := reader.ReadString('\n')
	minatSplit := strings.Split(strings.TrimSpace(minatInput), ",")
	for i := 0; i < len(minatSplit); i++ {
		idx := minatSplit[i]
		if i, err := strconv.Atoi(strings.TrimSpace(idx)); err == nil && i >= 1 && i <= len(minatList) {
			penggunaBaru.Minat = append(penggunaBaru.Minat, minatList[i-1])
		}
	}

	fmt.Println("Masukkan nomor keahlian yang dipilih (pisahkan dengan koma): ")
	keahlianInput, _ := reader.ReadString('\n')
	keahlianSplit := strings.Split(strings.TrimSpace(keahlianInput), ",")
	for i := 0; i < len(keahlianSplit); i++ {
		idx := keahlianSplit[i]
		if i, err := strconv.Atoi(strings.TrimSpace(idx)); err == nil && i >= 1 && i <= len(keahlianList) {
			penggunaBaru.Keahlian = append(penggunaBaru.Keahlian, keahlianList[i-1])
		}
	}

	var hasil []Rekomendasi
	for i := 0; i < len(daftarKarir); i++ {
		karir := daftarKarir[i]
		skor := hitungSkor(karir, penggunaBaru.Minat, penggunaBaru.Keahlian)
		if skor > 0 {
			persen := hitungPersentase(karir, penggunaBaru.Minat, penggunaBaru.Keahlian)
			hasil = append(hasil, Rekomendasi{karir.Nama, skor, persen, karir.Gaji})
		}
	}
	urutkanBerdasarkanSkor(hasil)
	penggunaBaru.HasilKarir = hasil

	return penggunaBaru
}

var riwayatPengguna []DataPengguna

func main() {
	daftarKarir := []Karir{
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

	var pengguna DataPengguna
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== APLIKASI REKOMENDASI KARIR ===")
		fmt.Println("1. Masukkan Data Pengguna Baru")
		fmt.Println("2. Lihat Rekomendasi Karir (Pengguna Terakhir)")
		fmt.Println("3. Lihat Rekomendasi Berdasarkan Persentase")
		fmt.Println("4. Lihat Rekomendasi Berdasarkan Gaji")
		fmt.Println("5. Tampilkan Semua Pilihan Karir")
		fmt.Println("6. Lihat & Kelola Riwayat Pengguna")
		fmt.Println("7. Keluar")
		fmt.Print(">> ")

		pilihan, _ := reader.ReadString('\n')
		pilihan = strings.TrimSpace(pilihan)

		switch pilihan {
		case "1":
			pengguna = prosesInputPengguna(reader, daftarKarir, DataPengguna{})
			riwayatPengguna = append(riwayatPengguna, pengguna)
			fmt.Println("✅ Data pengguna baru berhasil disimpan!")

		case "2":
			if pengguna.Nama == "" {
				fmt.Println("❗ Belum ada data pengguna. Silakan pilih menu 1 terlebih dahulu.")
				continue
			}
			urutkanBerdasarkanSkor(pengguna.HasilKarir)
			fmt.Println("👤 Data Pengguna Terakhir:")
			fmt.Printf("Nama: %s\n", pengguna.Nama)
			fmt.Printf("Jenis Kelamin: %s\n", pengguna.JenisKelamin)
			fmt.Printf("Tanggal Lahir: %s\n", pengguna.TanggalLahir)
			fmt.Printf("Minat: %s\n", strings.Join(pengguna.Minat, ", "))
			fmt.Printf("Keahlian: %s\n", strings.Join(pengguna.Keahlian, ", "))
			tampilkanRekomendasi(pengguna.HasilKarir, "skor")

		case "3":
			if pengguna.Nama == "" {
				fmt.Println("❗ Belum ada data pengguna. Silakan pilih menu 1 terlebih dahulu.")
				continue
			}
			sort.SliceStable(pengguna.HasilKarir, func(i, j int) bool {
				return pengguna.HasilKarir[i].Persentase > pengguna.HasilKarir[j].Persentase
			})
			tampilkanRekomendasi(pengguna.HasilKarir, "persentase")

		case "4":
			if pengguna.Nama == "" {
				fmt.Println("❗ Belum ada data pengguna. Silakan pilih menu 1 terlebih dahulu.")
				continue
			}
			urutkanBerdasarkanGaji(pengguna.HasilKarir)
			tampilkanRekomendasi(pengguna.HasilKarir, "gaji")

		case "5":
			tampilkanListKarir(daftarKarir)

		case "7":
			fmt.Println("Terima kasih telah menggunakan aplikasi.")
			return

		case "6":
			if len(riwayatPengguna) == 0 {
				fmt.Println("📭 Belum ada riwayat pengguna.")
				continue
			}

			for {
				fmt.Println("\n📖 Riwayat Pengguna Tersimpan:")
				fmt.Println(strings.Repeat("=", 40))

				for i, p := range riwayatPengguna {
					fmt.Printf("%d. Nama: %s\n", i+1, p.Nama)
					fmt.Printf("   Jenis Kelamin: %s\n", p.JenisKelamin)
					fmt.Printf("   Tanggal Lahir: %s\n", p.TanggalLahir)
					fmt.Printf("   Minat: %s\n", strings.Join(p.Minat, ", "))
					fmt.Printf("   Keahlian: %s\n", strings.Join(p.Keahlian, ", "))

					rekomendasiStr := "Belum ada rekomendasi."
					if len(p.HasilKarir) > 0 {
						rekomendasiStr = fmt.Sprintf("%s (Skor: %d)", p.HasilKarir[0].Nama, p.HasilKarir[0].Skor)
					}
					fmt.Printf("   Rekomendasi Teratas: %s\n", rekomendasiStr)
					fmt.Println(strings.Repeat("-", 40))
				}

				fmt.Println("Sub-Menu Riwayat:")
				fmt.Println("1. Edit Data Pengguna")
				fmt.Println("2. Hapus Data Pengguna")
				fmt.Println("3. Kembali ke Menu Utama")
				fmt.Print(">> ")

				pilihanSubMenu, _ := reader.ReadString('\n')
				pilihanSubMenu = strings.TrimSpace(pilihanSubMenu)

				keluarSubMenu := false
				switch pilihanSubMenu {
				case "1":
					fmt.Print("Masukkan nomor pengguna yang ingin diedit: ")
					input, _ := reader.ReadString('\n')
					nomor, err := strconv.Atoi(strings.TrimSpace(input))
					if err != nil || nomor < 1 || nomor > len(riwayatPengguna) {
						fmt.Println("❗ Nomor tidak valid.")
						continue
					}

					indeks := nomor - 1
					fmt.Println("\n--- Mengedit Data untuk:", riwayatPengguna[indeks].Nama, "---")
					dataTeredit := prosesInputPengguna(reader, daftarKarir, riwayatPengguna[indeks])
					riwayatPengguna[indeks] = dataTeredit
					fmt.Println("✅ Data berhasil diperbarui!")
					keluarSubMenu = true

				case "2":
					fmt.Print("Masukkan nomor pengguna yang ingin dihapus: ")
					input, _ := reader.ReadString('\n')
					nomor, err := strconv.Atoi(strings.TrimSpace(input))
					if err != nil || nomor < 1 || nomor > len(riwayatPengguna) {
						fmt.Println("❗ Nomor tidak valid.")
						continue
					}

					indeks := nomor - 1
					namaDihapus := riwayatPengguna[indeks].Nama

					riwayatPengguna = append(riwayatPengguna[:indeks], riwayatPengguna[indeks+1:]...)

					fmt.Printf("✅ Data untuk '%s' berhasil dihapus.\n", namaDihapus)

					if len(riwayatPengguna) == 0 {
						keluarSubMenu = true
					}

				case "3":
					keluarSubMenu = true

				default:
					fmt.Println("❗ Pilihan sub-menu tidak dikenal.")
				}

				if keluarSubMenu {
					break
				}
			}

		default:
			fmt.Println("❗ Pilihan tidak dikenal.")
		}
	}
}
