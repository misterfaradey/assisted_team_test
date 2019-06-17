package service

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"local/test-tasks/assisted_team/internal/dto"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

const defaultCargeType = "TotalAmount"
const defaultSortType = "SingleAdult"
const defaultCurrency = "SGD"

func openFiles(files []string) ([]Flights, error) {

	if len(files) == 0 {
		return nil, errors.New("no files")
	}

	data := []Flights{}

	for _, file := range files {

		xmlFile, err := os.Open(file)
		if err != nil {
			return nil, err
		}

		byteValue, err := ioutil.ReadAll(xmlFile)

		xmlFile.Close()

		if err != nil {
			return nil, err
		}

		log.Println("file", file, "read")

		var via AirFareSearch

		err = xml.Unmarshal(byteValue, &via)
		if err != nil {
			return nil, err
		}

		data = append(data, via.PricedItineraries.Flights...)

	}

	return data, nil

}

type AssistedService interface {
	Update(fileReturn, fileOneWay []string) error

	GetFlights(dto dto.FlightAllDto) []Flights

	MinPriceFlight(dto dto.FlightMostDto) Flights
	MaxPriceFlight(dto dto.FlightMostDto) Flights
	MinTimeFlight(dto dto.FlightMostDto) Flights
	MaxTimeFlight(dto dto.FlightMostDto) Flights
	OptimalFlight(dto dto.FlightMostDto) Flights
}

type assistedService struct {
	dataOneWay []Flights
	dataReturn []Flights
	mut        sync.RWMutex
}

func NewAssistedService() AssistedService {
	a := assistedService{
		dataOneWay: make([]Flights, 0),
		dataReturn: make([]Flights, 0),
	}
	return &a
}

func (a *assistedService) MinPriceFlight(req dto.FlightMostDto) Flights {
	a.mut.RLock()
	defer a.mut.RUnlock()

	var out []Flights
	if req.OneWay {
		out = a.getFlightsOneWay(req.From, req.To)
	} else {
		out = a.getFlightsReturn(req.From, req.To)
	}

	if len(out) == 0 {
		return Flights{}
	}

	return out[0]

}

func (a *assistedService) MaxPriceFlight(req dto.FlightMostDto) Flights {
	a.mut.RLock()
	defer a.mut.RUnlock()

	var out []Flights
	if req.OneWay {
		out = a.getFlightsOneWay(req.From, req.To)
	} else {
		out = a.getFlightsReturn(req.From, req.To)
	}

	if len(out) == 0 {
		return Flights{}
	}

	return out[len(out)-1]
}

func timeSort(out []Flights) {
	sort.Slice(out, func(i, j int) bool {

		time1 := out[i].OnwardPricedItinerary.Flights.Flight[len(out[i].OnwardPricedItinerary.Flights.Flight)-1].ArrivalTimeStamp.Time.Unix() -
			out[i].OnwardPricedItinerary.Flights.Flight[0].DepartureTimeStamp.Time.Unix()

		time2 := out[j].OnwardPricedItinerary.Flights.Flight[len(out[j].OnwardPricedItinerary.Flights.Flight)-1].ArrivalTimeStamp.Time.Unix() -
			out[j].OnwardPricedItinerary.Flights.Flight[0].DepartureTimeStamp.Time.Unix()

		return time1 < time2
	})
}

//примем за формулу оптимальности стоимость часа перелета как 5 SGD (1h = 5SGD ≈ 235RUB)
func optimalSort(out []Flights) {
	sort.Slice(out, func(i, j int) bool {

		fl1, _ := out[i].Pricing.GetChargeTypeIndex(defaultSortType)

		fl2, _ := out[j].Pricing.GetChargeTypeIndex(defaultSortType)

		if fl1 == -1 || fl2 == -1 {
			return false
		}

		time1 := out[i].OnwardPricedItinerary.Flights.Flight[len(out[i].OnwardPricedItinerary.Flights.Flight)-1].ArrivalTimeStamp.Time.Unix() -
			out[i].OnwardPricedItinerary.Flights.Flight[0].DepartureTimeStamp.Time.Unix() +
			int64(time.Hour.Seconds())*int64(out[i].Pricing.ServiceCharges[fl1].Price)/5

		time2 := out[j].OnwardPricedItinerary.Flights.Flight[len(out[j].OnwardPricedItinerary.Flights.Flight)-1].ArrivalTimeStamp.Time.Unix() -
			out[j].OnwardPricedItinerary.Flights.Flight[0].DepartureTimeStamp.Time.Unix() +
			int64(time.Hour.Seconds())*int64(out[j].Pricing.ServiceCharges[fl2].Price)/5

		return time1 < time2
	})
}

func (a *assistedService) MaxTimeFlight(req dto.FlightMostDto) Flights {
	a.mut.RLock()
	defer a.mut.RUnlock()

	var out []Flights
	if req.OneWay {
		out = a.getFlightsOneWay(req.From, req.To)
	} else {
		out = a.getFlightsReturn(req.From, req.To)
	}

	if len(out) == 0 {
		return Flights{}
	}

	timeSort(out)

	return out[len(out)-1]
}

func (a *assistedService) MinTimeFlight(req dto.FlightMostDto) Flights {
	a.mut.RLock()
	defer a.mut.RUnlock()

	var out []Flights
	if req.OneWay {
		out = a.getFlightsOneWay(req.From, req.To)
	} else {
		out = a.getFlightsReturn(req.From, req.To)
	}

	if len(out) == 0 {
		return Flights{}
	}

	timeSort(out)

	return out[0]
}

func (a *assistedService) OptimalFlight(req dto.FlightMostDto) Flights {
	a.mut.RLock()
	defer a.mut.RUnlock()

	var out []Flights
	if req.OneWay {
		out = a.getFlightsOneWay(req.From, req.To)
	} else {
		out = a.getFlightsReturn(req.From, req.To)
	}

	if len(out) == 0 {
		return Flights{}
	}

	optimalSort(out)

	return out[0]
}

//собрать подходящие рейсы в обе стороны
func (a *assistedService) getFlightsReturn(from, to string) []Flights {
	out := make([]Flights, 0)

	for _, flight := range a.dataReturn {
		_len := len(flight.OnwardPricedItinerary.Flights.Flight)
		if _len == 0 {
			continue
		}
		if flight.OnwardPricedItinerary.Flights.Flight[0].Source == from &&
			flight.OnwardPricedItinerary.Flights.Flight[_len-1].Destination == to {
			out = append(out, flight)
		}
	}
	return out
}

//собрать подходящие рейсы в одну сторону
func (a *assistedService) getFlightsOneWay(from, to string) []Flights {
	out := make([]Flights, 0)

	for _, flight := range a.dataOneWay {
		_len := len(flight.OnwardPricedItinerary.Flights.Flight)
		if _len == 0 {
			continue
		}
		if flight.OnwardPricedItinerary.Flights.Flight[0].Source == from &&
			flight.OnwardPricedItinerary.Flights.Flight[_len-1].Destination == to {
			out = append(out, flight)
		}
	}
	return out
}

func (a *assistedService) GetFlights(req dto.FlightAllDto) []Flights {
	a.mut.RLock()
	defer a.mut.RUnlock()

	var out []Flights
	if req.OneWay {
		out = a.getFlightsOneWay(req.From, req.To)
	} else {
		out = a.getFlightsReturn(req.From, req.To)
	}

	//сортировать не нужно
	if req.Price == 0 && req.Type == defaultSortType {
		return out

	}

	sort.Slice(out, func(i, j int) bool {

		fl1, _ := out[i].Pricing.GetChargeTypeIndex(req.Type)

		fl2, _ := out[j].Pricing.GetChargeTypeIndex(req.Type)

		if fl1 == -1 || fl2 == -1 {
			return false
		}

		if req.Price == 1 {
			return out[i].Pricing.ServiceCharges[fl1].Price > out[j].Pricing.ServiceCharges[fl2].Price
		}

		return out[i].Pricing.ServiceCharges[fl1].Price < out[j].Pricing.ServiceCharges[fl2].Price
	})

	return out

}

func (a *assistedService) Update(fileReturn, fileOneWay []string) error {
	a.mut.Lock()
	defer a.mut.Unlock()

	dReturn, err := updateData(fileReturn)
	if err != nil {
		return err
	}

	dOneWay, err := updateData(fileOneWay)
	if err != nil {
		return err
	}

	a.dataReturn = dReturn
	a.dataOneWay = dOneWay
	return nil
}

func updateData(files []string) ([]Flights, error) {
	data, err := openFiles(files)
	if err != nil {
		return nil, err
	}

	//отсортируем сразу по возрастанию цены,
	//так как по мне это самый распространенный кейс
	sort.Slice(data, func(i, j int) bool {

		//находим элементы, что будемм сравнивать. так как тут мало элементов, то слайс будет быстрее чем map
		fl1, err := data[i].Pricing.GetChargeTypeIndex(defaultSortType)
		if err != nil {
			log.Println(err)
			return false
		}

		fl2, err := data[j].Pricing.GetChargeTypeIndex(defaultSortType)
		if err != nil {
			log.Println(err)
			return false
		}

		return data[i].Pricing.ServiceCharges[fl1].Price < data[j].Pricing.ServiceCharges[fl2].Price
	})

	//сортируем совокупность рейсов по дате отправлениия,
	//чтобы первым рейс был тот, откуда мы отправляемся,
	//а последним был тот, куда мы прилетаем. А то мало ли
	for k := range data {
		sort.Slice(data[k].OnwardPricedItinerary.Flights.Flight, func(i, j int) bool {

			return data[k].OnwardPricedItinerary.Flights.Flight[i].DepartureTimeStamp.Unix() <
				data[k].OnwardPricedItinerary.Flights.Flight[j].DepartureTimeStamp.Unix()

		})
	}
	return data, nil
}
