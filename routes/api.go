package routes

import (
	"prisma-laundry-backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.POST("/login", controllers.Login)
		api.POST("/create-user", controllers.CreateUser)
		api.GET("/customers", controllers.GetCustomers)
		api.GET("/petugas", controllers.GetPetugas)
		api.GET("/admin", controllers.GetAdmin)
		api.POST("/customers", controllers.CreateCustomer)

		api.POST("/pengeluaran", controllers.CreatePengeluaran)
		api.GET("/pengeluaran", controllers.GetPengeluaran)
		api.GET("/pengeluaran/export", controllers.ExportPengeluaran)
		api.POST("/obat", controllers.CreateObat)
		api.GET("/obat", controllers.GetObat)
		api.GET("/obat/export", controllers.ExportObatBulanan)

		api.POST("/status-barang", controllers.CreateStatusBarang)
		api.GET("/status-barang", controllers.GetStatusBarang)
		api.PUT("/status-barang/:id/status", controllers.UpdateStatusBarang)
		api.PUT("/status-barang/:id/keterangan", controllers.UpdateKeteranganBarang)
		api.DELETE("/status-barang/:id", controllers.DeleteStatusBarang)

		api.POST("/tagihan", controllers.CreateTagihan)
		api.GET("/tagihan/user/:id", controllers.GetTagihanUser)
		api.PATCH("/tagihan/:id/status", controllers.UpdateStatus)
		api.GET("/tagihan", controllers.GetAllTagihanRekap)
		api.GET("/tagihan/rinci/:id", controllers.GetTagihanRinci)
		api.POST("/tagihan/rinci", controllers.CreateTagihanRinci)
		api.GET("/tagihan/export/rinci/:id", controllers.ExportTagihanRinci)
		api.GET("/tagihan/export/general", controllers.ExportTagihanGeneral)

		api.GET("/pembayaran", controllers.GetPembayaran)
        api.POST("/pembayaran", controllers.CreatePembayaran)
		api.GET("/pembayaran/customer/:id", controllers.GetPembayaranCustomer)

		api.GET("/notifications/:id", controllers.GetNotifications)
		api.PUT("/notifications/read/:id", controllers.MarkNotificationRead)
	}
}
