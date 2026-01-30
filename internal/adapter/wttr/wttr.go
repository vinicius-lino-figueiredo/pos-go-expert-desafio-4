// Package wttr TODO
package wttr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/vinicius-lino-figueiredo/pos-go-expert-desafio-4/domain"
)

const baseURL = "https://wttr.in/"

// ErrNoConditionFound TODO
var ErrNoConditionFound = errors.New("no condition found")

// ErrStatusCode TODO
type ErrStatusCode struct {
	Status int
}

// Error implements [error].
func (e ErrStatusCode) Error() string {
	return fmt.Sprintf("unexpected status code %d", e.Status)
}

// Wttr TODO
type Wttr struct {
	cl *http.Client
}

// NewTemperatureGetter TODO
func NewTemperatureGetter(cl *http.Client) domain.TemperatureGetter {
	return &Wttr{cl: cl}
}

// GetTemperature implements [domain.TemperatureGetter].
func (w *Wttr) GetTemperature(ctx context.Context, location string) (float64, error) {
	u, err := w.getURL(location)
	if err != nil {
		return 0, fmt.Errorf("mounting url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return 0, fmt.Errorf("creating request: %w", err)
	}

	res, err := w.cl.Do(req)
	if err != nil {
		return 0, fmt.Errorf("doing request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, ErrStatusCode{Status: res.StatusCode}
	}

	var body wttr
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return 0, fmt.Errorf("decoding response: %w", err)
	}

	if len(body.CurrentCondition) == 0 {
		return 0, ErrNoConditionFound
	}

	c, err := strconv.ParseFloat(body.CurrentCondition[0].TempC, 64)
	if err != nil {
		return 0, fmt.Errorf("converting temperature number: %w", err)
	}

	return c, nil

}

func (w *Wttr) getURL(location string) (string, error) {
	b, err := url.JoinPath(baseURL, location)
	if err != nil {
		return "", fmt.Errorf("joining location: %w", err)
	}

	u, err := url.Parse(b)
	if err != nil {
		return "", fmt.Errorf("parsing URL: %w", err)
	}

	u.RawQuery = url.Values{"format": {"j1"}}.Encode()

	// u.Query().Add("json", "j1")

	return u.String(), nil

}

type wttr struct {
	CurrentCondition []struct {
		FeelsLikeC string `json:"FeelsLikeC"`
		FeelsLikeF string `json:"FeelsLikeF"`
		Cloudcover string `json:"cloudcover"`
		Humidity   string `json:"humidity"`
		LangPt     []struct {
			Value string `json:"value"`
		} `json:"lang_pt"`
		LocalObsDateTime string `json:"localObsDateTime"`
		ObservationTime  string `json:"observation_time"`
		PrecipInches     string `json:"precipInches"`
		PrecipMm         string `json:"precipMM"`
		Pressure         string `json:"pressure"`
		PressureInches   string `json:"pressureInches"`
		TempC            string `json:"temp_C"`
		TempF            string `json:"temp_F"`
		UvIndex          string `json:"uvIndex"`
		Visibility       string `json:"visibility"`
		VisibilityMiles  string `json:"visibilityMiles"`
		WeatherCode      string `json:"weatherCode"`
		WeatherDesc      []struct {
			Value string `json:"value"`
		} `json:"weatherDesc"`
		WeatherIconURL []struct {
			Value string `json:"value"`
		} `json:"weatherIconUrl"`
		Winddir16Point string `json:"winddir16Point"`
		WinddirDegree  string `json:"winddirDegree"`
		WindspeedKmph  string `json:"windspeedKmph"`
		WindspeedMiles string `json:"windspeedMiles"`
	} `json:"current_condition"`
	NearestArea []struct {
		AreaName []struct {
			Value string `json:"value"`
		} `json:"areaName"`
		Country []struct {
			Value string `json:"value"`
		} `json:"country"`
		Latitude   string `json:"latitude"`
		Longitude  string `json:"longitude"`
		Population string `json:"population"`
		Region     []struct {
			Value string `json:"value"`
		} `json:"region"`
		WeatherURL []struct {
			Value string `json:"value"`
		} `json:"weatherUrl"`
	} `json:"nearest_area"`
	Request []struct {
		Query string `json:"query"`
		Type  string `json:"type"`
	} `json:"request"`
	Weather []struct {
		Astronomy []struct {
			MoonIllumination string `json:"moon_illumination"`
			MoonPhase        string `json:"moon_phase"`
			Moonrise         string `json:"moonrise"`
			Moonset          string `json:"moonset"`
			Sunrise          string `json:"sunrise"`
			Sunset           string `json:"sunset"`
		} `json:"astronomy"`
		AvgtempC string `json:"avgtempC"`
		AvgtempF string `json:"avgtempF"`
		Date     string `json:"date"`
		Hourly   []struct {
			DewPointC        string `json:"DewPointC"`
			DewPointF        string `json:"DewPointF"`
			FeelsLikeC       string `json:"FeelsLikeC"`
			FeelsLikeF       string `json:"FeelsLikeF"`
			HeatIndexC       string `json:"HeatIndexC"`
			HeatIndexF       string `json:"HeatIndexF"`
			WindChillC       string `json:"WindChillC"`
			WindChillF       string `json:"WindChillF"`
			WindGustKmph     string `json:"WindGustKmph"`
			WindGustMiles    string `json:"WindGustMiles"`
			Chanceoffog      string `json:"chanceoffog"`
			Chanceoffrost    string `json:"chanceoffrost"`
			Chanceofhightemp string `json:"chanceofhightemp"`
			Chanceofovercast string `json:"chanceofovercast"`
			Chanceofrain     string `json:"chanceofrain"`
			Chanceofremdry   string `json:"chanceofremdry"`
			Chanceofsnow     string `json:"chanceofsnow"`
			Chanceofsunshine string `json:"chanceofsunshine"`
			Chanceofthunder  string `json:"chanceofthunder"`
			Chanceofwindy    string `json:"chanceofwindy"`
			Cloudcover       string `json:"cloudcover"`
			DiffRad          string `json:"diffRad"`
			Humidity         string `json:"humidity"`
			LangPt           []struct {
				Value string `json:"value"`
			} `json:"lang_pt"`
			PrecipInches    string `json:"precipInches"`
			PrecipMm        string `json:"precipMM"`
			Pressure        string `json:"pressure"`
			PressureInches  string `json:"pressureInches"`
			ShortRad        string `json:"shortRad"`
			TempC           string `json:"tempC"`
			TempF           string `json:"tempF"`
			Time            string `json:"time"`
			UvIndex         string `json:"uvIndex"`
			Visibility      string `json:"visibility"`
			VisibilityMiles string `json:"visibilityMiles"`
			WeatherCode     string `json:"weatherCode"`
			WeatherDesc     []struct {
				Value string `json:"value"`
			} `json:"weatherDesc"`
			WeatherIconURL []struct {
				Value string `json:"value"`
			} `json:"weatherIconUrl"`
			Winddir16Point string `json:"winddir16Point"`
			WinddirDegree  string `json:"winddirDegree"`
			WindspeedKmph  string `json:"windspeedKmph"`
			WindspeedMiles string `json:"windspeedMiles"`
		} `json:"hourly"`
		MaxtempC    string `json:"maxtempC"`
		MaxtempF    string `json:"maxtempF"`
		MintempC    string `json:"mintempC"`
		MintempF    string `json:"mintempF"`
		SunHour     string `json:"sunHour"`
		TotalSnowCm string `json:"totalSnow_cm"`
		UvIndex     string `json:"uvIndex"`
	} `json:"weather"`
}
