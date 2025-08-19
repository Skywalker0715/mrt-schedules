package station

import (
	"net/http"

	"github.com/Skywalker0715/mrt-schedules/common/response"

	"github.com/gin-gonic/gin"
)

func InitializeRouter(router *gin.RouterGroup) {
	// Initialize the service
	stationService := NewService()
	station := router.Group("stations")

	station.GET("/", func(c *gin.Context) {
		GetAllStations(c, stationService)
	})

	station.GET("/:id", func(c *gin.Context) {
		CheckSchedulesByStation(c, stationService)
	})
}

func GetAllStations(c *gin.Context, service Service) {
	datas, err := service.GetAllStation()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			response.Response{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		response.Response{
			Success: true,
			Message: "Successfully get all station",
			Data:    datas,
		},
	)
}

func CheckSchedulesByStation(c *gin.Context, service Service) {
	id := c.Param("id")

	datas, err := service.CheckSchedulesByStation(id)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			response.Response{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		response.Response{
			Success: true,
			Message: "Successfully get schedules by station",
			Data:    datas,
		},
	)
}
