package service

import (
	"encoding/xml"
	"errors"
	"fmt"
	"time"
)

type AirFareSearch struct {
	PricedItineraries struct {
		Flights []Flights `xml:"Flights"`
	} `xml:"PricedItineraries"`
}

type Flights struct {
	OnwardPricedItinerary Itinerary `xml:"OnwardPricedItinerary"`
	ReturnPricedItinerary Itinerary `xml:"ReturnPricedItinerary" json:"ReturnPricedItinerary,omitempty"`
	Pricing               Pricing   `xml:"Pricing"`
}

type Itinerary struct {
	Flights struct {
		Flight []Flight `xml:"Flight"  json:"Flight,omitempty"`
	} `xml:"Flights" json:"Flights,omitempty"`
}

type Pricing struct {
	Currency       string `xml:"currency,attr"`
	ServiceCharges []struct {
		Price      float32 `xml:",chardata"`
		Type       string  `xml:"type,attr"`
		ChargeType string  `xml:"ChargeType,attr"`
	} `xml:"ServiceCharges"`
}

func (p *Pricing) GetChargeTypeIndex(sortType string) (int, error) {
	for k := range p.ServiceCharges {
		if p.ServiceCharges[k].ChargeType == defaultCargeType &&
			p.ServiceCharges[k].Type == sortType &&
			p.Currency == defaultCurrency {
			return k, nil
		}
	}
	return -1, errors.New(fmt.Sprint("Err:\n sortType: ", sortType, "Data: ", p))
}

type Carrier struct {
	Text string `xml:",chardata"`
	ID   string `xml:"id,attr"`
}

type Flight struct {
	Carrier            Carrier    `xml:"Carrier"`
	FlightNumber       int        `xml:"FlightNumber"`
	Source             string     `xml:"Source"`
	Destination        string     `xml:"Destination"`
	DepartureTimeStamp customTime `xml:"DepartureTimeStamp"`
	ArrivalTimeStamp   customTime `xml:"ArrivalTimeStamp" `
	Class              string     `xml:"Class"`
	NumberOfStops      int        `xml:"NumberOfStops"`
	FareBasis          string     `xml:"FareBasis"`
	WarningText        string     `xml:"WarningText"`
	TicketType         string     `xml:"TicketType"`
}

type customTime struct {
	time.Time
}

func (c *customTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortForm = "2006-01-02T1504" // yyyymmdd date format
	var v string
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	parse, err := time.Parse(shortForm, v)
	if err != nil {
		return err
	}
	*c = customTime{parse}
	return nil
}
