package gameinfo

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/morlay/genshin-level-up/pkg/httputil"
)

func ServerFromUID(uid int) string {
	switch strconv.Itoa(uid)[0] {
	case '1', '2':
		return "cn_gf01"
	case '5':
		return "cn_qd01"
	case '6':
		return "os_usa"
	case '7':
		return "os_euro"
	case '8':
		return "os_asia"
	case '9':
		return "os_cht"
	default:
		return ""
	}
}

const salt = "xV8v4Qu54lUKrEYFZkJhB8cuOh9Asafs"

func getDS(query map[string][]string, body map[string]interface{}) string {
	t := time.Now().Unix()
	// Integer between 100000 - 200000
	random := rand.Intn(100000) + 100000

	b, q := "", ""

	if len(body) > 0 {
		d, _ := json.Marshal(body)
		b = string(d)
	}

	if len(query) > 0 {
		q = url.Values(query).Encode()
	}

	check := md5.Sum([]byte(fmt.Sprintf("salt=%s&t=%d&r=%d&b=%s&q=%s", salt, t, random, b, q)))

	return fmt.Sprintf("%d,%d,%x", t, random, check)
}

const apiEndpoint = "https://api-takumi.mihoyo.com/game_record/app/genshin/api"
const hoyolabVersion = "2.11.1"

func NewClient(cookie string) *Client {
	return &Client{
		Transports: []httputil.Transport{
			NewCommonTransport(cookie),
		},
	}
}

type Client struct {
	Transports []httputil.Transport
}

func (c *Client) GetAllCharacters(ctx context.Context, uid int) ([]Character, error) {
	usrInfo, err := c.GetUserInfo(ctx, uid)
	if err != nil {
		return nil, err
	}

	characterIDs := make([]int, len(usrInfo.Avatars))

	for i := range usrInfo.Avatars {
		characterIDs[i] = usrInfo.Avatars[i].ID
	}

	list, err := c.GetCharacters(context.Background(), uid, characterIDs)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (c *Client) GetUserInfo(ctx context.Context, uid int) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", apiEndpoint+"/index", nil)
	if err != nil {
		return nil, err
	}

	query := url.Values{}

	query.Add("role_id", strconv.Itoa(uid))
	query.Add("server", ServerFromUID(194435467))

	req.Header.Set("DS", getDS(query, nil))

	req.URL.RawQuery = query.Encode()

	resp, err := httputil.GetShortConnClientContext(ctx, 5*time.Second, c.Transports...).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData := struct {
		RetCode int      `json:"retcode,omitempty"`
		Message string   `json:"message,omitempty"`
		Data    UserInfo `json:"data,omitempty"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	if respData.RetCode != 0 {
		if respData.RetCode == 10101 {
			origin := httputil.RequestFromContext(ctx)
			u := "//" + origin.Host + origin.URL.Path

			return nil, fmt.Errorf(
				"%d: %s, 请登录 https://bbs.mihoyo.com/ys/, 并在浏览器控制台执行 window.open(`%s?cookie=${btoa(document.cookie)}`)",
				respData.RetCode, respData.Message, u)
		}

		return nil, fmt.Errorf("%d: %s", respData.RetCode, respData.Message)
	}

	return &respData.Data, nil
}

func (c *Client) GetCharacters(ctx context.Context, uid int, characterIDs []int) ([]Character, error) {
	requestBody := map[string]interface{}{
		"role_id":       uid,
		"server":        ServerFromUID(194435467),
		"character_ids": characterIDs,
	}

	b, _ := json.Marshal(requestBody)

	req, err := http.NewRequestWithContext(ctx, "POST", apiEndpoint+"/character", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("DS", getDS(nil, requestBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := httputil.GetShortConnClientContext(ctx, 5*time.Second, c.Transports...).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//
	respData := struct {
		RetCode int       `json:"retcode,omitempty"`
		Message string    `json:"message,omitempty"`
		Data    *UserInfo `json:"data,omitempty"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	if respData.RetCode != 0 {
		return nil, fmt.Errorf("%d: %s", respData.RetCode, respData.Message)
	}

	return respData.Data.Avatars, nil
}

func NewCommonTransport(cookie string) httputil.Transport {
	return func(next http.RoundTripper) http.RoundTripper {
		return &commonRoundTripper{cookie: cookie, next: next}
	}
}

type commonRoundTripper struct {
	cookie string
	next   http.RoundTripper
}

func (c *commonRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,en-US;q=0.8")

	req.Header.Set("Origin", "https://webstatic.mihoyo.com")
	req.Header.Set("Referer", "https://webstatic.mihoyo.com/app/community-game-records/index.html?v=6")

	// `miHoYoBBS/<hoyolabVersion>` required
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 9; Unspecified Device) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/39.0.0.0 Mobile Safari/537.36 miHoYoBBS/"+hoyolabVersion)

	req.Header.Set("x-rpc-app_version", hoyolabVersion)
	req.Header.Set("x-rpc-client_type", "5")
	req.Header.Set("X-Requested-With", "com.mihoyo.hyperion")

	req.Header.Set("Cookie", c.cookie)

	return c.next.RoundTrip(req)
}
