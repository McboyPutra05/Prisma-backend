package controllers

import (
	"fmt"
	"net/http"
	"prisma-laundry-backend/config"
	"prisma-laundry-backend/models"
	"prisma-laundry-backend/services"
	"prisma-laundry-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateTagihan(c *gin.Context) {
	var input models.Tagihan

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.CreateTagihan(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, input)
}

func GetTagihanUser(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	data, err := services.GetTagihanByUser(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func UpdateStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var body struct {
		Status string `json:"status"`
	}

	c.ShouldBindJSON(&body)

	err := services.UpdateStatusTagihan(uint(id), body.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

func GetAllTagihanRekap(c *gin.Context) {
	data, err := services.GetRekapTagihanSemuaCustomer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil rekap tagihan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func GetTagihanRinci(c *gin.Context) {
	customerId := c.Param("id")
	bulan := c.Query("bulan")

	var customer struct {
		Name string
	}
	config.DB.Table("users").Select("name").Where("id = ? AND role = 'customer'", customerId).Scan(&customer)

	var totalBulanIni float64
	var totalKeseluruhan float64
	var totalPembayaran float64

	queryTotal := config.DB.Table("tagihan_rincis").Where("customer_id = ?", customerId)
	if bulan != "" && bulan != "Semua Bulan" {
		queryTotal = queryTotal.Where("bulan = ?", bulan)
	}
	queryTotal.Select("COALESCE(SUM(total), 0)").Scan(&totalBulanIni)
	config.DB.Table("tagihan_rincis").Where("customer_id = ?", customerId).Select("COALESCE(SUM(total), 0)").Scan(&totalKeseluruhan)
	config.DB.Table("pembayarans").Where("customer_id = ?", customerId).Select("COALESCE(SUM(nominal), 0)").Scan(&totalPembayaran)

	totalBelumDibayar := totalKeseluruhan - totalPembayaran
	if totalBelumDibayar < 0 {
		totalBelumDibayar = 0 
	}

	var details []models.TagihanRinci
	query := config.DB.Where("customer_id = ?", customerId)
	if bulan != "" && bulan != "Semua Bulan" {
		query = query.Where("bulan = ?", bulan)
	}
	query.Order("tanggal asc").Find(&details)

	var jumlahPO int64 = 0
	var formattedDetails []map[string]interface{}

	for i, d := range details {
		isFirstOfGroup := false
		if i == 0 {
			isFirstOfGroup = true
		} else {
			prev := details[i-1]
			if d.Tanggal != prev.Tanggal || d.Merk != prev.Merk || d.KodePo != prev.KodePo {
				isFirstOfGroup = true
			}
		}

		if isFirstOfGroup {
			jumlahPO++
		}

		formattedDetails = append(formattedDetails, map[string]interface{}{
			"tgl":     d.Tanggal,
			"qty":     d.Qty,
			"item":    d.Item,
			"trt":     d.JenisCucian,
			"harga":   d.Harga,
			"total":   d.Total,
			"merk":    d.Merk,
			"kode_po": d.KodePo,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"customer_name": customer.Name,
			"summary": gin.H{
				"total_bulan_ini":     totalBulanIni,
				"total_keseluruhan":   totalKeseluruhan,
				"total_belum_dibayar": totalBelumDibayar,
				"total_sudah_dibayar": totalPembayaran,
				"jumlah_po":           jumlahPO,
			},
			"details": formattedDetails,
		},
	})
}

func CreateTagihanRinci(c *gin.Context) {
	var input struct {
		CustomerID string `json:"customer_id"`
		Tanggal    string `json:"tanggal"`
		Merk       string `json:"merk"`
		KodePO     string `json:"kode_po"`
		Rincian    []struct {
			Qty         string `json:"qty"`
			Item        string `json:"item"`
			JenisCucian string `json:"jenisCucian"`
			Harga       string `json:"harga"`
		} `json:"rincian"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bulan := "Januari"
	if len(input.Tanggal) >= 7 {
		monthStr := input.Tanggal[5:7]
		monthInt, _ := strconv.Atoi(monthStr)
		daftarBulan := []string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
		if monthInt >= 1 && monthInt <= 12 {
			bulan = daftarBulan[monthInt]
		}
	}

	for _, r := range input.Rincian {
		hargaFloat, _ := strconv.ParseFloat(r.Harga, 64)
		qtyFloat, _ := strconv.ParseFloat(r.Qty, 64)
		total := hargaFloat * qtyFloat

		tagihan := models.TagihanRinci{
			CustomerID:  input.CustomerID,
			Bulan:       bulan,
			Tanggal:     input.Tanggal,
			Merk:        input.Merk,
			KodePo:      input.KodePO,
			Qty:         r.Qty,
			Item:        r.Item,
			JenisCucian: r.JenisCucian,
			Harga:       hargaFloat,
			Total:       total,
			Status:      "Belum Lunas",
		}
		config.DB.Create(&tagihan)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Tagihan Rinci berhasil disimpan ke Database!"})
}

func ExportTagihanRinci(c *gin.Context) {
	customerId := c.Param("id")
	bulan := c.Query("bulan")

	var customer struct {
		Name string
	}
	config.DB.Table("users").Select("name").Where("id = ? AND role = 'customer'", customerId).Scan(&customer)

	if customer.Name == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Data Customer tidak ditemukan atau role bukan customer"})
		return
	}

	var tagihans []models.TagihanRinci
	query := config.DB.Where("customer_id = ?", customerId)
	if bulan != "" {
		query = query.Where("bulan = ?", bulan)
	}
	query.Order("tanggal asc").Find(&tagihans)

	fileExcel, err := utils.GenerateTagihanRinciExcel(customer.Name, tagihans)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal merakit Excel Nota Tagihan"})
		return
	}

	filename := fmt.Sprintf("Nota_Tagihan_%s_%s.xlsx", customer.Name, bulan)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	if err := fileExcel.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim file Excel"})
		return
	}
}

func ExportTagihanGeneral(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Fitur export rekap general sedang dalam pengembangan."})
}
