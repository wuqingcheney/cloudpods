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

package modules

import (
	"fmt"

	"yunion.io/x/jsonutils"

	"yunion.io/x/cloudmux/pkg/multicloud/hcso/client/manager"
	"yunion.io/x/cloudmux/pkg/multicloud/hcso/client/responses"
)

type SUserManager struct {
	SResourceManager
}

func NewUserManager(cfg manager.IManagerConfig) *SUserManager {
	user := &SUserManager{SResourceManager: SResourceManager{
		SBaseManager:  NewBaseManager(cfg),
		ServiceName:   ServiceNameIAM,
		Region:        cfg.GetRegionId(),
		ProjectId:     "",
		// [ORIGIN] 原始版本为 "v3.0/OS-USER"，与 HCS 8.6.0 IAM 2.0 文档不符
		// [CHANGED] HCS 8.6.0 IAM 用户管理接口路径为 /v3/users，统一改为 v3
		version:       "v3",
		Keyword:       "user",
		KeywordPlural: "users",

		ResourceKeyword: "users",
	}}
	user.SetDomainId(cfg.GetDomainId())
	return user
}

func (self *SUserManager) List(querys map[string]string) (*responses.ListResult, error) {
	// [ORIGIN] 原需临时 SetVersion("v3")，现默认已是 v3，无需重复设置
	return self.SResourceManager.List(querys)
}

func (self *SUserManager) Delete(id string) (jsonutils.JSONObject, error) {
	// [ORIGIN] 原需临时 SetVersion("v3")，现默认已是 v3，无需重复设置
	return self.SResourceManager.Delete(id, nil)
}

func (self *SUserManager) ResetPassword(id, password string) error {
	params := map[string]interface{}{
		"user": map[string]string{
			"password": password,
		},
	}
	_, err := self.SResourceManager.Update(id, jsonutils.Marshal(params))
	return err
}

func (self *SUserManager) ListGroups(userId string) (*responses.ListResult, error) {
	// [ORIGIN] 原需临时 SetVersion("v3")，现默认已是 v3，无需重复设置
	return self.SResourceManager.ListInContextWithSpec(nil, fmt.Sprintf("%s/groups", userId), nil, "groups")
}
