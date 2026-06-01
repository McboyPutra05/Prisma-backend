package utils

import (
	"fmt"
	"prisma-laundry-backend/models"
	"strings"

	"github.com/xuri/excelize/v2"
)

// GenerateObatExcel bertugas murni hanya untuk merakit file Excel Obat
func GenerateObatExcel(obats []models.Obat) (*excelize.File, error) {
	f := excelize.NewFile()

	if len(obats) == 0 {
		f.SetSheetName("Sheet1", "Data Kosong")
		f.SetCellValue("Data Kosong", "A1", "Tanggal")
		f.SetCellValue("Data Kosong", "B1", "Jenis Obat")
		f.SetCellValue("Data Kosong", "C1", "Banyaknya")
		return f, nil
	}

	groupedData := make(map[string][]models.Obat)
	for _, obat := range obats {
		year := "Unknown"
		parts := strings.Split(obat.Tanggal, "/")
		if len(parts) == 3 {
			year = parts[2]
		} else {
			partsDash := strings.Split(obat.Tanggal, "-")
			if len(partsDash) == 3 && len(partsDash[0]) == 4 {
				year = partsDash[0]
			} else if len(partsDash) == 3 && len(partsDash[2]) == 4 {
				year = partsDash[2]
			}
		}
		groupedData[year] = append(groupedData[year], obat)
	}

	firstSheet := true
	for year, dataPerYear := range groupedData {
		if firstSheet {
			f.SetSheetName("Sheet1", year)
			firstSheet = false
		} else {
			f.NewSheet(year)
		}

		f.SetCellValue(year, "A1", "Tanggal")
		f.SetCellValue(year, "B1", "Jenis Obat")
		f.SetCellValue(year, "C1", "Banyaknya")

		for i, obat := range dataPerYear {
			baris := i + 2
			f.SetCellValue(year, fmt.Sprintf("A%d", baris), obat.Tanggal)
			f.SetCellValue(year, fmt.Sprintf("B%d", baris), obat.Name)
			f.SetCellValue(year, fmt.Sprintf("C%d", baris), obat.Qty)
		}
	}

	return f, nil
}

// GeneratePengeluaranExcel bertugas murni hanya untuk merakit file Excel Pengeluaran
func GeneratePengeluaranExcel(pengeluarans []models.Pengeluaran) (*excelize.File, error) {
	f := excelize.NewFile()

	if len(pengeluarans) == 0 {
		f.SetSheetName("Sheet1", "Data Kosong")
		f.SetCellValue("Data Kosong", "A1", "Tanggal")
		f.SetCellValue("Data Kosong", "B1", "Barang")
		f.SetCellValue("Data Kosong", "C1", "Total")
		return f, nil
	}

	groupedData := make(map[string][]models.Pengeluaran)
	for _, p := range pengeluarans {
		year := "Unknown"
		parts := strings.Split(p.Tanggal, "/")
		if len(parts) == 3 {
			year = parts[2]
		} else {
			partsDash := strings.Split(p.Tanggal, "-")
			if len(partsDash) == 3 && len(partsDash[0]) == 4 {
				year = partsDash[0]
			} else if len(partsDash) == 3 && len(partsDash[2]) == 4 {
				year = partsDash[2]
			}
		}
		groupedData[year] = append(groupedData[year], p)
	}

	firstSheet := true
	for year, dataPerYear := range groupedData {
		if firstSheet {
			f.SetSheetName("Sheet1", year)
			firstSheet = false
		} else {
			f.NewSheet(year)
		}

		f.SetCellValue(year, "A1", "Tanggal")
		f.SetCellValue(year, "B1", "Barang")
		f.SetCellValue(year, "C1", "Total")

		for i, p := range dataPerYear {
			baris := i + 2
			f.SetCellValue(year, fmt.Sprintf("A%d", baris), p.Tanggal)
			f.SetCellValue(year, fmt.Sprintf("B%d", baris), p.Barang)
			f.SetCellValue(year, fmt.Sprintf("C%d", baris), p.Total)
		}
	}

	return f, nil
}