package utils

import (
	"fmt"
	"prisma-laundry-backend/models"
	"time"

	"github.com/xuri/excelize/v2"
)

func GenerateTagihanRinciExcel(customerName string, tagihans []models.TagihanRinci) (*excelize.File, error) {
	f := excelize.NewFile()
	sheet := "Sheet1"

	styleJudul, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 18},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	styleTengah, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	styleBoldTengah, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	styleAngka, _ := f.NewStyle(&excelize.Style{
		NumFmt: 3, 
	})
	
	formatTanggal := "dd/mm/yyyy"
	styleTanggal, _ := f.NewStyle(&excelize.Style{
		CustomNumFmt: &formatTanggal, 
		Alignment:    &excelize.Alignment{Horizontal: "center"},
	})

	f.MergeCell(sheet, "A1", "G1")
	f.SetCellValue(sheet, "A1", "PRISMA LAUNDRY")
	f.SetCellStyle(sheet, "A1", "A1", styleJudul)

	f.MergeCell(sheet, "A2", "G2")
	f.SetCellValue(sheet, "A2", "jl. pahlawan nomor xx, sukabumi selatan, jakarta barat")
	f.SetCellStyle(sheet, "A2", "A2", styleTengah)

	f.MergeCell(sheet, "A3", "G3")
	f.SetCellValue(sheet, "A3", "No. Telp 0812198719823")
	f.SetCellStyle(sheet, "A3", "A3", styleTengah)

	f.SetCellValue(sheet, "A5", "Kpd.")
	f.SetCellValue(sheet, "B5", customerName)

	baris := 7
	nomorUrut := 1

	for i := 0; i < len(tagihans); i++ {
		t := tagihans[i]

		isFirstOfGroup := false
		if i == 0 {
			isFirstOfGroup = true
		} else {
			prev := tagihans[i-1]
			if t.Tanggal != prev.Tanggal || t.Merk != prev.Merk || t.KodePo != prev.KodePo {
				isFirstOfGroup = true
			}
		}

		if isFirstOfGroup {
			f.SetCellValue(sheet, fmt.Sprintf("A%d", baris), nomorUrut)
			
			parsedDate, err := time.Parse("2006-01-02", t.Tanggal)
			cellTgl := fmt.Sprintf("B%d", baris)
			
			if err == nil {
				f.SetCellValue(sheet, cellTgl, parsedDate)
				f.SetCellStyle(sheet, cellTgl, cellTgl, styleTanggal)
			} else {
				f.SetCellValue(sheet, cellTgl, t.Tanggal)
			}
			
			nomorUrut++
		}

		f.SetCellValue(sheet, fmt.Sprintf("C%d", baris), t.Qty)
		itemText := t.Item
		if t.NoCelana != "" {
			itemText = fmt.Sprintf("%s (No. %s)", t.Item, t.NoCelana)
		}
		f.SetCellValue(sheet, fmt.Sprintf("D%d", baris), itemText) 
		f.SetCellValue(sheet, fmt.Sprintf("E%d", baris), t.JenisCucian)
		
		f.SetCellValue(sheet, fmt.Sprintf("F%d", baris), t.Harga)
		f.SetCellStyle(sheet, fmt.Sprintf("F%d", baris), fmt.Sprintf("F%d", baris), styleAngka)
		
		f.SetCellValue(sheet, fmt.Sprintf("G%d", baris), t.Total)
		f.SetCellStyle(sheet, fmt.Sprintf("G%d", baris), fmt.Sprintf("G%d", baris), styleAngka)

		baris++
		isLastOfGroup := false
		if i == len(tagihans)-1 {
			isLastOfGroup = true
		} else {
			next := tagihans[i+1]
			if t.Tanggal != next.Tanggal || t.Merk != next.Merk || t.KodePo != next.KodePo {
				isLastOfGroup = true
			}
		}

		if isLastOfGroup {
			if t.KodePo != "" {
				cellPO := fmt.Sprintf("B%d", baris)
				f.SetCellValue(sheet, cellPO, t.KodePo)
				f.SetCellStyle(sheet, cellPO, cellPO, styleBoldTengah)
			}
			
			if t.Merk != "" {
				cellMerk := fmt.Sprintf("D%d", baris)
				f.SetCellValue(sheet, cellMerk, t.Merk)
				f.SetCellStyle(sheet, cellMerk, cellMerk, styleBoldTengah)
			}

			if t.KodePo != "" || t.Merk != "" {
				baris++ 
			}
		}
	}

	f.SetColWidth(sheet, "A", "A", 5) 
	f.SetColWidth(sheet, "B", "B", 15)
	f.SetColWidth(sheet, "C", "C", 8) 
	f.SetColWidth(sheet, "D", "E", 25)
	f.SetColWidth(sheet, "F", "G", 15) 

	return f, nil
}