# 统一身份认证系统 OAuth2.0 接入文档

## 概述

本文档描述如何将统一身份认证系统（CasIAM）通过 OAuth2.0 授权码模式接入 Cloudpods。

接入实现位于：`pkg/keystone/driver/oauth2/casiam/`

---

## 接入流程

```
用户访问 Cloudpods
    ↓
Cloudpods 构造授权跳转地址
    → GET /bit-msa-sso/oauth/authorize?client_id=xxx&response_type=code&redirect_uri=xxx&state=xxx
    ↓
用户在统一身份认证系统登录
    ↓
统一身份认证系统回调 redirect_uri，携带 authorizationCode=xxx
    ↓
Cloudpods 用授权码换取 access_token
    → POST /bit-msa-sso/oauth/zkyToken?client_id=xxx&client_secret=xxx&grant_type=authorization_code&code=xxx
    ↓
Cloudpods 验证 token 获取用户信息
    → POST /bit-msa-sso/oauth/check_token?token=xxx
    ↓
创建/同步本地用户，分配项目和角色，完成登录
```

---

## 对接前提

向统一身份认证系统管理员申请以下信息：

| 信息 | 说明 | 示例 |
|------|------|------|
| `client_id` | 应用系统编码 | `BIT-MSA` |
| `client_secret` | 应用系统密钥 | `xxxxxx` |
| 服务地址 | SSO 服务基础地址 | `http://10.25.0.17` |
| `redirect_uri` | 回调地址（提供给对方配置） | `https://your-cloudpods/api/v1/idp/xxx/callback` |

### 环境地址

| 环境 | 地址 |
|------|------|
| 测试环境 | `http://10.25.0.17` |
| 科研外网生产 | `http://10.21.255.19` |
| 科研内网生产 | `http://10.54.255.19` |

---

## 配置方法

### 方式一：通过 climc 命令行创建

```bash
# 基础创建（最简配置）
climc idp-create-casiam-oauth2 \
  --name "统一身份认证" \
  --app-id "BIT-MSA" \
  --secret "your-client-secret" \
  --sso-endpoint "http://10.25.0.17"

# 完整配置（推荐，含权限映射）
climc idp-create-casiam-oauth2 \
  --name "统一身份认证" \
  --app-id "BIT-MSA" \
  --secret "your-client-secret" \
  --sso-endpoint "http://10.25.0.17" \
  --default-project-id "<cloudpods-project-id>" \
  --default-role-id "<cloudpods-role-id>" \
  --auto-create-user \
  --auto-create-project
```

### 方式二：通过 climc 更新已有 IDP 配置

```bash
climc idp-config-casiam-oauth2 <idp-id-or-name> \
  --app-id "BIT-MSA" \
  --secret "your-client-secret" \
  --sso-endpoint "http://10.21.255.19" \
  --default-project-id "<cloudpods-project-id>" \
  --default-role-id "<cloudpods-role-id>"
```

### 方式三：通过 YAML 编辑器配置

```bash
# 先禁用 IDP
climc idp-disable <idp-id>

# 打开 YAML 编辑器
climc idp-config-edit <idp-id>
```

YAML 格式示例：
```yaml
oauth2:
  app_id: BIT-MSA
  secret: your-client-secret
  sso_endpoint: http://10.25.0.17
  default_project_id: <project-id>
  default_role_id: <role-id>
  user_name_attribute: name
  user_id_attribute: user_id
  user_displayname_attribute: display_name
```

---

## 配置参数说明

### 必填参数

| 参数 | 说明 |
|------|------|
| `app_id` | 统一身份认证分配的 client_id |
| `secret` | 统一身份认证分配的 client_secret |
| `sso_endpoint` | 统一身份认证服务基础地址（不含路径，不含末尾斜杠） |

### 权限映射参数（可选，推荐配置）

| 参数 | 说明 |
|------|------|
| `default_project_id` | 用户登录后默认加入的 Cloudpods 项目 ID |
| `default_role_id` | 用户在默认项目中的角色 ID（如 member、admin）|
| `auto_create_user` | 首次登录时自动创建本地用户（默认 true）|
| `auto_create_project` | 项目不存在时自动创建（默认 false）|

### 查询项目和角色 ID

```bash
# 查询项目列表
climc project-list

# 查询角色列表
climc role-list
```

---

## 常用管理命令

```bash
# 查看所有 IDP
climc idp-list

# 查看 IDP 详情
climc idp-show <idp-id>

# 查看 IDP 配置
climc idp-config-show <idp-id>

# 启用/禁用
climc idp-enable <idp-id>
climc idp-disable <idp-id>

# 手动同步
climc idp-sync <idp-id>

# 获取 SSO 跳转地址（用于测试）
climc idp-sso-url <idp-id>

# 获取回调地址（提供给统一身份认证系统配置）
climc idp-sso-callback-url <idp-id>

# 设置为默认 SSO
climc idp-default-sso <idp-id> --enable

# 删除 IDP（需先禁用）
climc idp-disable <idp-id>
climc idp-delete <idp-id>
```

---

## 用户同步机制

### 同步时机

CasIAM 驱动使用 **登录时同步**（SyncOnAuth）模式，不做定时全量同步。

### 各场景处理

| 场景 | 处理方式 |
|------|----------|
| 首次登录 | 自动创建本地用户，加入默认项目和角色 |
| 再次登录 | 同步更新显示名、邮箱、手机号等属性 |
| 外部修改显示名/邮箱 | 下次登录自动同步 ✅ |
| 外部删除用户 | 登录时 check_token 返回 active=false，登录失败 ✅ |
| 外部撤销角色 | 登录时先清除旧角色，再重新分配 ✅ |
| 外部修改用户名 | 若对方返回稳定 user_id 则正常复用；否则创建新用户 ⚠️ |

### 关于用户唯一 ID

- 若统一身份认证的 `check_token` 接口返回 `user_id` 字段，驱动会优先使用该字段作为唯一标识
- 若不返回 `user_id`，则降级使用 `user_name` 作为标识（用户名修改后会创建新用户）
- **建议**：联系统一身份认证管理员确认 `check_token` 是否支持返回稳定的 `user_id`

---

## 接口参考

### 1. 获取授权码
```
GET /bit-msa-sso/oauth/authorize
  ?client_id=BIT-MSA
  &response_type=code
  &redirect_uri=<回调地址>
  &state=<随机状态值>
```

### 2. 授权码换 Token
```
POST /bit-msa-sso/oauth/zkyToken
  ?client_id=BIT-MSA
  &client_secret=xxx
  &grant_type=authorization_code
  &code=<授权码>

响应：
{
  "access_token": "eyJ...",
  "refresh_token": "eyJ...",
  "token_type": "bearer",
  "expires_in": 7199,
  "scope": "all"
}
```

### 3. 验证 Token 获取用户信息
```
POST /bit-msa-sso/oauth/check_token?token=<access_token>

响应：
{
  "active": true,
  "exp": 1739507817,
  "user_name": "zhangsan",
  "user_id": "U-10086",    // 若支持
  "jti": "806108f2-...",
  "client_id": "BIT-MSA",
  "scope": ["all"]
}
```

### 4. 刷新 Token
```
POST /bit-msa-sso/oauth/token
  ?grant_type=refresh_token
  &client_id=BIT-MSA
  &client_secret=xxx
  &refresh_token=<refresh_token>
```

---

## 代码位置

| 文件 | 说明 |
|------|------|
| `pkg/keystone/driver/oauth2/casiam/casiam.go` | 驱动主体，认证逻辑 |
| `pkg/keystone/driver/oauth2/casiam/factory.go` | 工厂类，注册驱动 |
| `pkg/keystone/driver/oauth2/casiam/doc.go` | 包声明 |
| `pkg/apis/identity/oauth2.go` | SCasIAMOAuth2ConfigOptions 配置类型 |
| `pkg/apis/identity/config.go` | IdpTemplateCasIAM 模板常量 |
| `pkg/keystone/driver/oauth2/types.go` | IOAuth2RoleSyncer 接口定义 |
| `pkg/keystone/driver/oauth2/oauth2.go` | 登录时角色清理逻辑 |
| `pkg/keystone/models/assignments.go` | RemoveUserProjectRole 公开方法 |
| `pkg/keystone/service/drivers.go` | 驱动注册 import |
| `cmd/climc/shell/identity/identityproviders.go` | climc 命令 |
