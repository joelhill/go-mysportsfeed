package msf

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	logrus "github.com/sirupsen/logrus"
)

// SeasonalTeamGamelogsOptions - Are the options to hit the seasonal team gamelogs endpoint
type SeasonalTeamGamelogsOptions struct {
	// URL Parts
	URL     string
	Version string
	Sport   string
	Season  string
	Format  string

	// Optional URL Params
	Team   string
	Game   string
	Date   string
	Stats  string
	Sort   string
	Offset string
	Limit  string
	Force  string
}

// NewSeasonalTeamGamelogsOptions - Returns the options with most url parts already set to hit the seasonal games endpoint
func (s *Service) NewSeasonalTeamGamelogsOptions() *SeasonalTeamGamelogsOptions {
	return &SeasonalTeamGamelogsOptions{
		URL:     s.Config.BaseURL,
		Version: s.Config.Version,
		Sport:   s.Config.Sport,
		Format:  s.Config.Format,
		Season:  s.Config.Season,
	}
}

// SeasonalTeamGamelogs - hits the https://api.mysportsfeeds.com/{version}/pull/{sport}/{season}/team_gamelogs.{format} endoint
func (s *Service) SeasonalTeamGamelogs(options *SeasonalTeamGamelogsOptions) (GameLogIO, int, error) {

	mapping := GameLogIO{}

	// make sure we have all the required elements to build the full required url string.
	err := validateSeasonalTeamGamelogsURI(options)
	if err != nil {
		return mapping, 0, err
	}

	t := time.Now()
	cacheBuster := t.Format("20060102150405")

	uri := fmt.Sprintf("%s/%s/pull/%s/%s/team_gamelogs.%s?cachebuster=%s", options.URL, options.Version, options.Sport, options.Season, options.Format, cacheBuster)

	if len(options.Team) > 0 {
		uri = fmt.Sprintf("%s&team=%s", uri, options.Team)
	}

	if len(options.Game) > 0 {
		uri = fmt.Sprintf("%s&game=%s", uri, options.Game)
	}

	if len(options.Date) > 0 {
		uri = fmt.Sprintf("%s&date=%s", uri, options.Date)
	}

	if len(options.Stats) > 0 {
		uri = fmt.Sprintf("%s&stats=%s", uri, options.Stats)
	}

	if len(options.Sort) > 0 {
		uri = fmt.Sprintf("%s&sort=%s", uri, options.Sort)
	}

	if len(options.Offset) > 0 {
		uri = fmt.Sprintf("%s&offset=%s", uri, options.Offset)
	}

	if len(options.Limit) > 0 {
		uri = fmt.Sprintf("%s&limit=%s", uri, options.Limit)
	}

	if len(options.Force) > 0 {
		uri = fmt.Sprintf("%s&force=%s", uri, options.Force)
	}

	s.Logger = s.Logger.WithFields(logrus.Fields{
		"URI": uri,
	})
	s.Logger.Debug("SeasonalTeamGamelogs API Call")

	// make you a client
	client := retryablehttp.NewClient()

	req, err := retryablehttp.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		s.Logger.Errorf("client: could not create request: %s", err.Error())
		return mapping, 0, err
	}
	req.Header.Add("Authorization", CompressionHeaderGzip)
	req.Header.Add("Authorization", s.Config.Authorization)

	response, err := client.Do(req)
	if err != nil {
		s.Logger.Errorf("client: error making http request: %s", err.Error())
		return mapping, 0, err
	}

	if response.StatusCode < 200 || response.StatusCode > 300 {
		s.Logger.Errorf("client: something went wrong making the get request for SeasonalTeamGamelogs: %s", err.Error())
		return mapping, response.StatusCode, err
	}

	s.Logger.Infof("SeasonalTeamGamelogs Status Code: %d", response.StatusCode)

	if jErr := json.NewDecoder(response.Body).Decode(&mapping); jErr != nil {
		s.Logger.Errorf("client: error decoding response for SeasonalTeamGamelogs: %s", err.Error())
		return mapping, response.StatusCode, err
	}

	return mapping, response.StatusCode, nil
}

func validateSeasonalTeamGamelogsURI(options *SeasonalTeamGamelogsOptions) error {
	if len(options.URL) == 0 {
		return errors.New("missing required option to build the url: URL")
	}
	if len(options.Version) == 0 {
		return errors.New("missing required option to build the url: Version")
	}
	if len(options.Sport) == 0 {
		return errors.New("missing required option to build the url: Sport")
	}
	if len(options.Season) == 0 {
		return errors.New("missing required option to build the url: Season")
	}
	if len(options.Format) == 0 {
		return errors.New("missing required option to build the url: Format")
	}
	return nil
}
