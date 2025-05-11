package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/LainInTheWired/ctf_backend/contest/model"
	"golang.org/x/xerrors"
)

type TeamRepository interface {
	ListTeamUsersByContest(cid int, uid *int) ([]model.Team, error)
}

type teamRepository struct {
	HTTPClient *http.Client
	URL        string
}

func NewTeamRepository(h *http.Client, url string) TeamRepository {
	return &teamRepository{
		HTTPClient: h,
		URL:        url,
	}
}

func (r *teamRepository) ListTeamUsersByContest(cid int, uid *int) ([]model.Team, error) {
	var endpoint string
	if uid == nil {
		endpoint = fmt.Sprintf("%s/team/%d/user", r.URL, cid)
	} else {
		endpoint = fmt.Sprintf("%s/team/%d/user?userid=%d", r.URL, cid, *uid)
		fmt.Println(endpoint)
	}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, xerrors.Errorf("can't create http request: %w", err)
	}
	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// リクエストの送信
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading clone response body: %v", err)
	}

	// json.Unmarshalでデコード
	var teams []model.Team
	if err := json.Unmarshal(body, &teams); err != nil {
		return nil, xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return nil, xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}

	// クローン作成のUPIDを表示
	log.Printf("teams infomation: %v\n", teams)

	return teams, nil
}
