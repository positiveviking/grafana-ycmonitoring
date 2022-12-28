package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/endpoint"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/yandex-cloud/go-sdk/iamkey"
)

type sdk struct {
	*ycsdk.SDK
	*resty.Client

	middleware *ycsdk.IamTokenMiddleware
}

func newSDK(endpoint string, monitoringEndpoint string, keyJson string) (*sdk, error) {
	creds, err := getSDKCreds(keyJson)
	if err != nil {
		return nil, err
	}

	ycSDK, err := buildSDK(endpoint, creds)
	if err != nil {
		return nil, err
	}

	rest, err := buildRest(monitoringEndpoint, ycSDK)
	if err != nil {
		return nil, err
	}

	md := ycsdk.NewIAMTokenMiddleware(ycSDK, time.Now)

	return &sdk{
		SDK:        ycSDK,
		Client:     rest,
		middleware: md,
	}, nil
}

func (s *sdk) check(ctx context.Context) error {
	token, err := s.getToken()
	if err != nil {
		return err
	}
	if _, err = s.Client.R().SetContext(ctx).SetAuthToken(token).Head(""); err != nil {
		return fmt.Errorf("api check: %w", err)
	}
	return nil
}

func (s *sdk) read(
	ctx context.Context, folderID string, req metricsReq,
) (result metrics, err error) {
	token, err := s.getToken()
	if err != nil {
		return result, err
	}

	resp, err := s.Client.R().SetAuthToken(token).
		SetContext(ctx).
		SetQueryParam("folderId", folderID).
		SetBody(req).
		SetResult(&result).
		Post("")

	if err != nil {
		return result, fmt.Errorf("metrics read: %w", err)
	}
	if resp.StatusCode() >= http.StatusBadRequest {
		return result, fmt.Errorf("bad http status: %d", resp.StatusCode())
	}
	return result, nil
}

func (s *sdk) getToken() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	token, err := s.middleware.GetIAMToken(ctx, true)
	if err != nil {
		return "", fmt.Errorf("get token: %w", err)
	}
	return token, nil
}

func getSDKCreds(keyJson string) (ycsdk.Credentials, error) {
	if keyJson == "" {
		return ycsdk.InstanceServiceAccount(), nil
	}
	var key *iamkey.Key
	if err := json.Unmarshal([]byte(keyJson), &key); err != nil {
		return nil, fmt.Errorf("api key unmarshal: %w", err)
	}
	return ycsdk.ServiceAccountKey(key)
}

func buildSDK(apiEndpoint string, creds ycsdk.Credentials) (*ycsdk.SDK, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ycSDK, err := ycsdk.Build(ctx, ycsdk.Config{
		Credentials: creds,
		Endpoint:    apiEndpoint,
	})
	if err != nil {
		return nil, fmt.Errorf("sdk build: %w", err)
	}

	return ycSDK, nil
}

func buildRest(monitoringEndpoint string, ycSDK *ycsdk.SDK) (*resty.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if monitoringEndpoint == "" {
		epResp, err := ycSDK.ApiEndpoint().ApiEndpoint().Get(ctx, &endpoint.GetApiEndpointRequest{ApiEndpointId: "monitoring"})
		if err != nil {
			return nil, fmt.Errorf("monitoring endpoint discovery: %w", err)
		}
		monitoringEndpoint = epResp.Address
	}

	monURL := url.URL{
		Scheme: "https",
		Host:   monitoringEndpoint,
		Path:   "/monitoring/v2/data/read",
	}

	return resty.New().
		SetBaseURL(monURL.String()).
		SetAuthScheme("Bearer").
		SetTimeout(time.Second * 30).
		SetRetryCount(3), nil
}
