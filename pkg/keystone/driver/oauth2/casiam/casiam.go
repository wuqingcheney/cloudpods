// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package casiam

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/httputils"

	"yunion.io/x/onecloud/pkg/apis/identity"
	"yunion.io/x/onecloud/pkg/keystone/driver/oauth2"
	"yunion.io/x/onecloud/pkg/keystone/models"
)

// SCasIAMOAuth2Driver 实现统一身份认证系统的 OAuth2.0 授权码模式接入
// 文档参考: Oauth2.0授权码方式.docx / 统一身份认证系统-OAUTH2.0规范
type SCasIAMOAuth2Driver struct {
	oauth2.SOAuth2BaseDriver
}

func NewCasIAMOAuth2Driver(appId string, secret string) oauth2.IOAuth2Driver {
	return &SCasIAMOAuth2Driver{
		SOAuth2BaseDriver: oauth2.SOAuth2BaseDriver{
			AppId:  appId,
			Secret: secret,
		},
	}
}

// getSSOEndpoint 从 context 中读取配置的统一身份认证服务地址
func (drv *SCasIAMOAuth2Driver) getSSOEndpoint(ctx context.Context) string {
	if config, ok := ctx.Value("config").(identity.TConfigs); ok {
		if v, exists := config["oauth2"]["sso_endpoint"]; exists {
			return strings.TrimRight(strings.Trim(v.String(), `"`), "/")
		}
	}
	return ""
}

// GetSsoRedirectUri 构造授权码模式的跳转地址
// 访问: /bit-msa-sso/oauth/authorize?client_id=xxx&response_type=code
func (drv *SCasIAMOAuth2Driver) GetSsoRedirectUri(ctx context.Context, callbackUrl, state string) (string, error) {
	endpoint := drv.getSSOEndpoint(ctx)
	if endpoint == "" {
		return "", errors.Error("casiam sso_endpoint is not configured")
	}
	params := map[string]string{
		"client_id":     drv.AppId,
		"response_type": "code",
		"redirect_uri":  callbackUrl,
		"state":         state,
	}
	urlStr := fmt.Sprintf("%s/bit-msa-sso/oauth/authorize?%s", endpoint, jsonutils.Marshal(params).QueryString())
	return urlStr, nil
}

// sTokenResponse 授权码换取 token 的响应结构
type sTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

// fetchAccessToken 使用授权码换取 access_token
// POST /bit-msa-sso/oauth/zkyToken?client_id=xxx&client_secret=xxx&grant_type=authorization_code&code=xxx
func (drv *SCasIAMOAuth2Driver) fetchAccessToken(ctx context.Context, code string) (string, error) {
	endpoint := drv.getSSOEndpoint(ctx)
	if endpoint == "" {
		return "", errors.Error("casiam sso_endpoint is not configured")
	}
	tokenUrl := fmt.Sprintf("%s/bit-msa-sso/oauth/zkyToken", endpoint)
	params := map[string]string{
		"client_id":     drv.AppId,
		"client_secret": drv.Secret,
		"grant_type":    "authorization_code",
		"code":          code,
	}
	urlWithParams := fmt.Sprintf("%s?%s", tokenUrl, jsonutils.Marshal(params).QueryString())

	httpclient := httputils.GetDefaultClient()
	_, resp, err := httputils.JSONRequest(httpclient, ctx, httputils.POST, urlWithParams, nil, nil, true)
	if err != nil {
		return "", errors.Wrap(err, "fetch access token")
	}
	var tokenResp sTokenResponse
	if err := resp.Unmarshal(&tokenResp); err != nil {
		return "", errors.Wrap(err, "unmarshal token response")
	}
	if tokenResp.AccessToken == "" {
		return "", errors.Error("empty access_token in response")
	}
	return tokenResp.AccessToken, nil
}

// sCheckTokenResponse check_token 接口的响应结构
type sCheckTokenResponse struct {
	Active   bool   `json:"active"`
	Exp      int64  `json:"exp"`
	UserId   string `json:"user_id"`   // 稳定唯一用户ID（若对方支持）
	UserName string `json:"user_name"` // 用户名（可能随改名变化）
	Jti      string `json:"jti"`
	ClientId string `json:"client_id"`
}

// fetchUserInfo 验证 access_token 并获取登录用户信息
// POST /bit-msa-sso/oauth/check_token?token=xxx
func (drv *SCasIAMOAuth2Driver) fetchUserInfo(ctx context.Context, accessToken string) (*sCheckTokenResponse, error) {
	endpoint := drv.getSSOEndpoint(ctx)
	checkUrl := fmt.Sprintf("%s/bit-msa-sso/oauth/check_token?token=%s", endpoint, accessToken)

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	httpclient := httputils.GetDefaultClient()
	_, resp, err := httputils.JSONRequest(httpclient, ctx, httputils.POST, checkUrl, headers, nil, true)
	if err != nil {
		return nil, errors.Wrap(err, "check token")
	}
	var info sCheckTokenResponse
	if err := resp.Unmarshal(&info); err != nil {
		return nil, errors.Wrap(err, "unmarshal check_token response")
	}
	if !info.Active {
		return nil, errors.Error("token is not active")
	}
	if info.UserName == "" {
		return nil, errors.Error("empty user_name in check_token response")
	}
	return &info, nil
}

// Authenticate 授权码模式完整认证流程:
// 1. 用授权码换取 access_token
// 2. 验证 token 并获取用户名
func (drv *SCasIAMOAuth2Driver) Authenticate(ctx context.Context, code string) (map[string][]string, error) {
	accessToken, err := drv.fetchAccessToken(ctx, code)
	if err != nil {
		return nil, errors.Wrap(err, "fetchAccessToken")
	}
	userInfo, err := drv.fetchUserInfo(ctx, accessToken)
	if err != nil {
		return nil, errors.Wrap(err, "fetchUserInfo")
	}
	attrs := make(map[string][]string)
	// 优先使用稳定的 user_id 作为唯一标识；若对方接口不返回则降级用 user_name
	uid := userInfo.UserId
	if uid == "" {
		uid = userInfo.UserName
	}
	attrs["name"] = []string{userInfo.UserName}
	attrs["user_id"] = []string{uid}
	attrs["display_name"] = []string{userInfo.UserName}
	return attrs, nil
}

// ClearUserProjectRoles 在登录时清除用户在指定项目下的所有旧角色，
// 使外部系统撤销角色后本地能同步生效。
// 实现 oauth2.IOAuth2RoleSyncer 接口。
func (drv *SCasIAMOAuth2Driver) ClearUserProjectRoles(ctx context.Context, userId, projectId string) error {
	roles, err := models.AssignmentManager.FetchUserProjectRoles(userId, projectId)
	if err != nil {
		return errors.Wrap(err, "FetchUserProjectRoles")
	}
	for i := range roles {
		if err := models.AssignmentManager.RemoveUserProjectRole(userId, projectId, roles[i].Id); err != nil {
			log.Warningf("casiam: remove user %s role %s from project %s fail: %s", userId, roles[i].Name, projectId, err)
		}
	}
	return nil
}
