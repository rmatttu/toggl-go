package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

// ClientConfigure is client configure
type ClientConfigure struct {
	Timeout  int
	APIToken string
}

// WorkspaceResponse Toggl workspace api response
type WorkspaceResponse struct {
	ID                          int       `json:"id"`
	Name                        string    `json:"name"`
	Profile                     int       `json:"profile"`
	Premium                     bool      `json:"premium"`
	Admin                       bool      `json:"admin"`
	DefaultHourlyRate           int       `json:"default_hourly_rate"`
	DefaultCurrency             string    `json:"default_currency"`
	OnlyAdminsMayCreateProjects bool      `json:"only_admins_may_create_projects"`
	OnlyAdminsSeeBillableRates  bool      `json:"only_admins_see_billable_rates"`
	OnlyAdminsSeeTeamDashboard  bool      `json:"only_admins_see_team_dashboard"`
	ProjectsBillableByDefault   bool      `json:"projects_billable_by_default"`
	Rounding                    int       `json:"rounding"`
	RoundingMinutes             int       `json:"rounding_minutes"`
	APIToken                    string    `json:"api_token"`
	At                          time.Time `json:"at"`
	IcalEnabled                 bool      `json:"ical_enabled"`
}

// DetailedReportRequest is detailed report request
type DetailedReportRequest struct {
	UserAgent   string
	WorkspaceID int
	Since       time.Time
	Until       time.Time
	Page        int
}

// DetailedReportResponse Toggl detailed report api response
type DetailedReportResponse struct {
	TotalGrand      int `json:"total_grand"`
	TotalBillable   int `json:"total_billable"`
	TotalCount      int `json:"total_count"`
	PerPage         int `json:"per_page"`
	TotalCurrencies []struct {
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	} `json:"total_currencies"`
	Data []struct {
		ID          int      `json:"id"`
		Pid         int      `json:"pid"`
		Tid         int      `json:"tid"`
		UID         int      `json:"uid"`
		Description string   `json:"description"`
		Start       string   `json:"start"`
		End         string   `json:"end"`
		Updated     string   `json:"updated"`
		Dur         int      `json:"dur"`
		User        string   `json:"user"`
		UseStop     bool     `json:"use_stop"`
		Client      string   `json:"client"`
		Project     string   `json:"project"`
		Task        int      `json:"task"`
		Billable    float64  `json:"billable"`
		IsBillable  bool     `json:"is_billable"`
		Cur         string   `json:"cur"`
		Tags        []string `json:"tags"`
	} `json:"data"`
}

const (
	// TogglAPIBaseURL is toggl api base url
	TogglAPIBaseURL = "https://api.track.toggl.com/api/v8"
	// TogglAPIWorkspaces is toggl workspaces api
	TogglAPIWorkspaces = "workspaces"
	// TogglReportsAPIBaseURL is toggl reports api
	TogglReportsAPIBaseURL = "https://api.track.toggl.com/reports/api/v2"
	// TogglReportsAPIDetailedReport is toggl detailed reports api
	TogglReportsAPIDetailedReport = "details"
)

func accessAPI(conf ClientConfigure, url string) ([]byte, error) {
	client := &http.Client{}
	client.Timeout = time.Second * time.Duration(conf.Timeout)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	req.SetBasicAuth(conf.APIToken, "api_token")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// FetchMainWorkspace is fetch main workspace
func FetchMainWorkspace(conf ClientConfigure) (*WorkspaceResponse, error) {
	u, _ := url.Parse(TogglAPIBaseURL)
	u.Path = path.Join(u.Path, TogglAPIWorkspaces)
	responseRaw, err := accessAPI(conf, u.String())
	if err != nil {
		return nil, err
	}
	var workspaces []WorkspaceResponse
	if err := json.Unmarshal(responseRaw, &workspaces); err != nil {
		return nil, err
	}
	return &workspaces[0], nil
}

func FetchDetailedReport(conf ClientConfigure, req DetailedReportRequest) (*DetailedReportResponse, string, error) {
	u, _ := url.Parse(TogglReportsAPIBaseURL)
	u.Path = path.Join(u.Path, TogglReportsAPIDetailedReport)

	q := u.Query()
	q.Set("user_agent", req.UserAgent)
	q.Set("workspace_id", strconv.Itoa(req.WorkspaceID))
	q.Set("since", req.Since.Format("2006-01-02")) // ISO 8601 date (YYYY-MM-DD) format. Defaults to today - 6 days.
	q.Set("until", req.Until.Format("2006-01-02")) // ISO 8601 date (YYYY-MM-DD) format. Note: Maximum date span (until - since) is one year. Defaults to today
	q.Set("page", strconv.Itoa(req.Page))          // integer, default 1
	u.RawQuery = q.Encode()

	responseRaw, err := accessAPI(conf, u.String())
	if err != nil {
		return nil, "", err
	}
	var res DetailedReportResponse
	if err := json.Unmarshal(responseRaw, &res); err != nil {
		return nil, "", err
	}
	return &res, string(responseRaw), nil
}
