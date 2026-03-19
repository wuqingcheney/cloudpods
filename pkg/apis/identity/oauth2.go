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

package identity

// OAuth2.0
type SOAuth2IdpConfigOptions struct {
	AppId  string `json:"app_id"`
	Secret string `json:"secret"`

	SIdpAttributeOptions
}

// SCasIAMOAuth2ConfigOptions 统一身份认证系统 OAuth2.0 授权码模式配置
type SCasIAMOAuth2ConfigOptions struct {
	// 应用系统编码 (client_id)
	AppId string `json:"app_id" help:"Client ID assigned by the unified identity auth system" required:"true"`
	// 应用系统密钥 (client_secret)
	Secret string `json:"secret" help:"Client secret assigned by the unified identity auth system" required:"true"`
	// 统一身份认证服务地址，例如 http://10.25.0.17
	SsoEndpoint string `json:"sso_endpoint" help:"Base URL of the unified identity auth SSO service, e.g. http://10.25.0.17" required:"true"`

	SIdpAttributeOptions
}
