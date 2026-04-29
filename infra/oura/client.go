package oura

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://api.ouraring.com"

type Client struct {
	token      string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// --- Oura API v2 レスポンス型 ---

type listResponse[T any] struct {
	Data      []T     `json:"data"`
	NextToken *string `json:"next_token"`
}

// DailyReadiness は /v2/usercollection/daily_readiness の data 要素
type DailyReadiness struct {
	ID           string                `json:"id"`
	Day          string                `json:"day"`
	Score        *int                  `json:"score"`
	Timestamp    string                `json:"timestamp"`
	Contributors ReadinessContributors `json:"contributors"`
}

type ReadinessContributors struct {
	ActivityBalance     *int `json:"activity_balance"`
	BodyTemperature     *int `json:"body_temperature"`
	HRVBalance          *int `json:"hrv_balance"`
	PreviousDayActivity *int `json:"previous_day_activity"`
	PreviousNight       *int `json:"previous_night"`
	RecoveryIndex       *int `json:"recovery_index"`
	RestingHeartRate    *int `json:"resting_heart_rate"`
	SleepBalance        *int `json:"sleep_balance"`
}

// DailySleep は /v2/usercollection/daily_sleep の data 要素
type DailySleep struct {
	ID           string            `json:"id"`
	Day          string            `json:"day"`
	Score        *int              `json:"score"`
	Timestamp    string            `json:"timestamp"`
	Contributors SleepContributors `json:"contributors"`
}

type SleepContributors struct {
	DeepSleep   *int `json:"deep_sleep"`
	Efficiency  *int `json:"efficiency"`
	Latency     *int `json:"latency"`
	REMSleep    *int `json:"rem_sleep"`
	Restfulness *int `json:"restfulness"`
	Timing      *int `json:"timing"`
	TotalSleep  *int `json:"total_sleep"`
}

// InterbeatInterval は /v2/usercollection/interbeat_interval の data 要素
type InterbeatInterval struct {
	ID        string    `json:"id"`
	Day       string    `json:"day"`
	Interval  float64   `json:"interval"`   // サンプリング間隔（秒）
	Items     []float64 `json:"items"`      // RR間隔の配列（秒）
	Timestamp string    `json:"timestamp"`
}

// --- 公開メソッド ---

func (c *Client) GetDailyReadiness(ctx context.Context, date string) (*DailyReadiness, error) {
	url := fmt.Sprintf("%s/v2/usercollection/daily_readiness?start_date=%s&end_date=%s", baseURL, date, date)
	var resp listResponse[DailyReadiness]
	if err := c.get(ctx, url, &resp); err != nil {
		return nil, fmt.Errorf("oura: get daily_readiness: %w", err)
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("oura: no readiness data for %s", date)
	}
	return &resp.Data[0], nil
}

func (c *Client) GetDailySleep(ctx context.Context, date string) (*DailySleep, error) {
	url := fmt.Sprintf("%s/v2/usercollection/daily_sleep?start_date=%s&end_date=%s", baseURL, date, date)
	var resp listResponse[DailySleep]
	if err := c.get(ctx, url, &resp); err != nil {
		return nil, fmt.Errorf("oura: get daily_sleep: %w", err)
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("oura: no sleep data for %s", date)
	}
	return &resp.Data[0], nil
}

func (c *Client) GetInterbeatInterval(ctx context.Context, date string) (*InterbeatInterval, error) {
	url := fmt.Sprintf(
		"%s/v2/usercollection/interbeat_interval?start_datetime=%sT00:00:00&end_datetime=%sT23:59:59",
		baseURL, date, date,
	)
	var resp listResponse[InterbeatInterval]
	if err := c.get(ctx, url, &resp); err != nil {
		return nil, fmt.Errorf("oura: get interbeat_interval: %w", err)
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("oura: no ibi data for %s", date)
	}
	return &resp.Data[0], nil
}

// get は指数バックオフ付きの GET リクエスト共通処理。
// 429 Too Many Requests の場合は最大5回リトライする（CLAUDE.md 注意事項）。
func (c *Client) get(ctx context.Context, url string, out any) error {
	var lastErr error
	for attempt := range 5 {
		if attempt > 0 {
			wait := time.Duration(1<<attempt) * time.Second
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(wait):
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.token)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = fmt.Errorf("rate limited (429)")
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("oura: unexpected status %d for %s", resp.StatusCode, url)
		}

		return json.NewDecoder(resp.Body).Decode(out)
	}
	return fmt.Errorf("oura: max retries exceeded: %w", lastErr)
}
