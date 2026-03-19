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
	"yunion.io/x/pkg/errors"

	api "yunion.io/x/onecloud/pkg/apis/identity"
	"yunion.io/x/onecloud/pkg/httperrors"
	"yunion.io/x/onecloud/pkg/keystone/driver/oauth2"
)

type SCasIAMDriverFactory struct{}

func (drv SCasIAMDriverFactory) NewDriver(appId string, secret string) oauth2.IOAuth2Driver {
	return NewCasIAMOAuth2Driver(appId, secret)
}

func (drv SCasIAMDriverFactory) TemplateName() string {
	return api.IdpTemplateCasIAM
}

func (drv SCasIAMDriverFactory) IdpAttributeOptions() api.SIdpAttributeOptions {
	return api.SIdpAttributeOptions{
		UserNameAttribute:        "name",
		UserIdAttribute:          "user_id",
		UserDisplaynameAttribtue: "display_name",
	}
}

func (drv SCasIAMDriverFactory) ValidateConfig(conf api.SOAuth2IdpConfigOptions) error {
	if len(conf.AppId) == 0 {
		return errors.Wrap(httperrors.ErrInputParameter, "missing app_id")
	}
	if len(conf.Secret) == 0 {
		return errors.Wrap(httperrors.ErrInputParameter, "missing secret")
	}
	return nil
}

func init() {
	oauth2.Register(&SCasIAMDriverFactory{})
}
