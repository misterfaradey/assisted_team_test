package dto

type FlightAllDto struct {
	From   string `form:"from" binding:"required"`
	To     string `form:"to" binding:"required"`
	Price  int    `form:"price"` //1=maxPrice another=minPrice
	Type   string `form:"type"`
	OneWay bool   `form:"oneway"`
}

type FlightMostDto struct {
	From   string `form:"from" binding:"required"`
	To     string `form:"to" binding:"required"`
	Sort   string `form:"sort" binding:"required"`
	OneWay bool   `form:"oneway"`
}
