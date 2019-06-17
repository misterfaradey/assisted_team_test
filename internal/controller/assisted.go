package controller

import (
	"errors"
	"local/test-tasks/assisted_team/internal/dto"
	"local/test-tasks/assisted_team/internal/server"
	"local/test-tasks/assisted_team/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AssistedController struct {
	service service.AssistedService
	actions []server.Action
}

func (c *AssistedController) Actions() []server.Action {
	return c.actions
}

func NewController(service service.AssistedService) server.Controller {

	c := AssistedController{
		service: service,
	}

	c.actions = append(c.actions,
		server.Action{
			HttpMethod:   "GET",
			RelativePath: "/api/flights/all",
			ActionExec:   c.flights,
		},
	)

	c.actions = append(c.actions,
		server.Action{
			HttpMethod:   "GET",
			RelativePath: "/api/flight",
			ActionExec:   c.flight,
		},
	)

	return &c
}

func (c *AssistedController) flights(ctx *gin.Context) {
	var req dto.FlightAllDto

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, err)
		return
	}

	out := c.service.GetFlights(req)

	ctx.JSON(http.StatusOK, out)
}

func (c *AssistedController) flight(ctx *gin.Context) {
	var req dto.FlightMostDto

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, err)
		return
	}

	var out service.Flights

	switch req.Sort {
	case "min_price":
		out = c.service.MinPriceFlight(req)
	case "max_price":
		out = c.service.MaxPriceFlight(req)
	case "min_time":
		out = c.service.MinTimeFlight(req)
	case "max_time":
		out = c.service.MaxTimeFlight(req)
	case "optimal":
		out = c.service.OptimalFlight(req)
	default:
		ctx.JSON(400, errors.New("set valid sort_type"))
		return
	}

	ctx.JSON(http.StatusOK, out)
}
