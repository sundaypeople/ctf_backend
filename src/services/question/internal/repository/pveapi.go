package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/LainInTheWired/ctf_backend/question/model"
	"github.com/cockroachdb/errors"
	"golang.org/x/xerrors"
)

type pveapiRepository struct {
	HTTPClient *http.Client
	URL        string
}
type PVEAPIRepository interface {
	Cloudinit(conf *model.CloudinitResponse) error
	CreateVM(conf *model.CreateVM) (string, error)
	DeleteVM(vmid int) error
	GetIPByVMID(vmid int) (*model.ResponseIPs, error)
	Template(vmid int) error
}

func NewPVEAPIRepository(h *http.Client, url string) PVEAPIRepository {
	return &pveapiRepository{
		HTTPClient: h,
		URL:        url,
	}
}

func (r *pveapiRepository) Cloudinit(conf *model.CloudinitResponse) error {
	// フォームデータの作成
	endpoint := fmt.Sprintf("http://%s:8000/cloudinit", r.URL)

	// フォームデータの作成
	jsend, err := json.Marshal(conf)
	if err != nil {
		return errors.Wrap(err, "can't change json")
	}

	// 新しいPOSTリクエストの作成
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsend))
	if err != nil {
		return xerrors.Errorf("can't create http request: %w", err)
	}

	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/json")

	// リクエストの送信
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスの読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Error reading clone response body: %w")
	}

	// json.Unmarshalでデコード
	var pveresp model.PveapiResponse[string]
	if err := json.Unmarshal(body, &pveresp); err != nil {
		return xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}
	return nil
}

func (r *pveapiRepository) CreateVM(conf *model.CreateVM) (string, error) {
	// フォームデータの作成
	endpoint := fmt.Sprintf("http://%s:8000/vm", r.URL)
	// フォームデータの作成

	jsend, err := json.Marshal(conf)
	if err != nil {
		return "", errors.Wrap(err, "can't change json")
	}

	// 新しいPOSTリクエストの作成
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsend))
	if err != nil {
		return "", xerrors.Errorf("can't create http request: %w", err)
	}

	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/json")

	// リクエストの送信
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return "", xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスの読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading clone response body: %v", err)
	}

	// json.Unmarshalでデコード
	var pveresp model.PveapiResponse[string]
	if err := json.Unmarshal(body, &pveresp); err != nil {
		return "", xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return "", xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}
	return pveresp.Data, nil
}

func (r *pveapiRepository) DeleteVM(vmid int) error {
	// フォームデータの作成
	endpoint := fmt.Sprintf("http://%s:8000/vm", r.URL)
	// フォームデータの作成

	conf := struct {
		ID int `json:"id"`
	}{
		ID: vmid,
	}

	jsend, err := json.Marshal(conf)
	if err != nil {
		return errors.Wrap(err, "can't change json")
	}

	// 新しいPOSTリクエストの作成
	req, err := http.NewRequest("DELETE", endpoint, bytes.NewBuffer(jsend))
	if err != nil {
		return xerrors.Errorf("can't create http request: %w", err)
	}

	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/json")

	// リクエストの送信
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスの読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading clone response body: %v", err)
	}

	// json.Unmarshalでデコード
	var pveresp model.PveapiResponse[string]
	if err := json.Unmarshal(body, &pveresp); err != nil {
		return xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}
	return nil
}
func (r *pveapiRepository) GetIPByVMID(vmid int) (*model.ResponseIPs, error) {
	endpoint := fmt.Sprintf("http://%s:8000/vm/%d/ips", r.URL, vmid)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't create http request: %w")
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "fail http request:")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading clone response body: %w")
	}

	var ifs model.ResponseIPs
	if err := json.Unmarshal(body, &ifs); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body: %w")
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return nil, errors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}
	return &ifs, nil
}

func (r *pveapiRepository) Template(vmid int) error {
	// フォームデータの作成
	endpoint := fmt.Sprintf("http://%s:8000/template", r.URL)
	// フォームデータの作成

	jsend, err := json.Marshal(map[string]int{"id": vmid})
	if err != nil {
		return errors.Wrap(err, "can't change json")
	}

	// 新しいPOSTリクエストの作成
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsend))
	if err != nil {
		return xerrors.Errorf("can't create http request: %w", err)
	}

	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/json")

	// リクエストの送信
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return xerrors.Errorf("fail http request: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスの読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading clone response body: %v", err)
	}

	// json.Unmarshalでデコード
	var pveresp model.PveapiResponse[string]
	if err := json.Unmarshal(body, &pveresp); err != nil {
		return xerrors.Errorf("can't unmarshal response body: %w", err)
	}

	// エラーチェック
	if resp.StatusCode >= 400 {
		return xerrors.Errorf("API Error: status code %d, response: %s", resp.StatusCode, resp.Status)
	}
	return nil
}
