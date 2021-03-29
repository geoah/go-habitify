package habitify

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	BaseURL = "https://api.habitify.me"

	JournalEndpoint     = BaseURL + "/journal"
	HabitsEndpoint      = BaseURL + "/habits"
	HabitEndpoint       = BaseURL + "/habits/%s"
	HabitStatusEndpoint = BaseURL + "/habits/%s/status"
	HabitLogsEndpoint   = BaseURL + "/habits/%s/logs"
	HabitNotesEndpoint  = BaseURL + "/habits/%s/notes"
	AreasEndpoint       = BaseURL + "/areas"
)

type (
	Habit struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		IsArchived  bool      `json:"is_archived"`
		StartDate   time.Time `json:"start_date"`
		TimeOfDay   []string  `json:"time_of_day"`
		AreaID      string    `json:"area_id"`
		Recurrence  string    `json:"recurrence"`
		CreatedDate time.Time `json:"created_date"`
		Goal        struct {
			UnitType    string `json:"unit_type"`
			Value       int    `json:"value"`
			Periodicity string `json:"periodicity"`
		} `json:"goal"`
		LogMethod string `json:"log_method"`
		Status    string `json:"status"`
	}
	HabitStatus struct {
		Status   string `json:"status"`
		Progress struct {
			CurrentValue  int       `json:"current_value"`
			TargetValue   int       `json:"target_value"`
			UnitType      string    `json:"unit_type"`
			Periodicity   string    `json:"periodicity"`
			ReferenceDate time.Time `json:"reference_date"`
		} `json:"progress"`
	}
	HabitLog struct {
		ID          string    `json:"id"`
		Value       float64   `json:"value"`
		CreatedDate time.Time `json:"created_date"`
		TargetDate  time.Time `json:"targe_date"`
		UnitType    string    `json:"unit_type"`
		HabitID     string    `json:"habit_id"`
	}
	HabitNote struct {
		ID          string    `json:"id"`
		Content     string    `json:"content"`
		CreatedDate time.Time `json:"created_date"`
		HabitID     string    `json:"habit_id"`
	}
	Area struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		CreatedDate time.Time `json:"created_date"`
	}
	Error struct {
		Reason         string `json:"reason"`
		AdditionalInfo struct {
			ActualUnitCategory   string `json:"actual_unit_category"`
			ExpectedUnitCategory string `json:"expected_unit_category"`
		} `json:"additional_info"`
	}
)

type Client struct {
	apiKey     string
	httpClient *resty.Client
}

func New(apiKey string) *Client {
	return &Client{
		httpClient: resty.New(),
		apiKey:     apiKey,
	}
}

func (c *Client) GetJournal(
	targetDate time.Time,
) ([]Habit, error) {
	res := []Habit{}
	httpRes, err := c.httpClient.R().
		SetQueryParams(
			map[string]string{
				"target_date": targetDate.Format(time.RFC3339),
			},
		).
		SetHeader("Authorization", c.apiKey).
		SetResult(&res).
		Get(JournalEndpoint)
	if httpRes.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status, %d", httpRes.StatusCode())
	}
	return res, err
}

func (c *Client) GetHabits() ([]Habit, error) {
	res := []Habit{}
	httpRes, err := c.httpClient.R().
		SetHeader("Authorization", c.apiKey).
		SetResult(&res).
		Get(HabitsEndpoint)
	if httpRes.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status, %d", httpRes.StatusCode())
	}
	return res, err
}

func (c *Client) GetHabitLogs(
	habitID string,
	from time.Time,
	to time.Time,
) ([]HabitLog, error) {
	res := []HabitLog{}
	httpRes, err := c.httpClient.R().
		SetHeader("Authorization", c.apiKey).
		SetQueryParams(
			map[string]string{
				"from": from.UTC().Format(time.RFC3339),
				"to":   to.UTC().Format(time.RFC3339),
			},
		).
		SetResult(&res).
		Get(fmt.Sprintf(HabitLogsEndpoint, habitID))
	if httpRes.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status, %d", httpRes.StatusCode())
	}
	return res, err
}

func (c *Client) AddHabitLogs(
	habitID string,
	targetDate time.Time,
	unitType string,
	value string,
) (*HabitLog, error) {
	if unitType == "" {
		unitType = "rep"
	}
	res := &HabitLog{}
	req := map[string]string{
		"value":       value,
		"unit_type":   unitType,
		"target_date": targetDate.UTC().Format(time.RFC3339),
	}
	errRes := Error{}
	httpRes, err := c.httpClient.R().
		SetHeader("Authorization", c.apiKey).
		SetBody(req).
		SetError(&errRes).
		SetResult(&res).
		Post(fmt.Sprintf(HabitLogsEndpoint, habitID))
	if httpRes.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status, %d", httpRes.StatusCode())
	}
	return res, err
}
