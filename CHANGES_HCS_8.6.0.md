# HCS 8.6.0 API 适配变更记录

日期：2026-03-20

---

## 一、CI/CD 流水线修改

**修改目的：** 将镜像构建推送目标从阿里云 ACR 改为华为云 SWR，仅构建 linux/amd64 架构镜像。

**影响文件（共12个）：**

| 文件 | 变更内容 |
|------|----------|
| `.github/workflows/docker_glance.yml` | 触发分支改为 `releases/**`/`main`/`master`，登录华为云 SWR，ARCH=amd64，REGISTRY=swr.cn-north-4.myhuaweicloud.com/iiecas-ziheng |
| `.github/workflows/docker_apigateway.yml` | 同上 |
| `.github/workflows/docker_climc.yml` | 同上 |
| `.github/workflows/docker_cloudmon.yml` | 同上 |
| `.github/workflows/docker_esxi_agent.yml` | 同上 |
| `.github/workflows/docker_host.yml` | 同上 |
| `.github/workflows/docker_keystone.yml` | 同上 |
| `.github/workflows/docker_logger.yml` | 同上 |
| `.github/workflows/docker_notify.yml` | 同上 |
| `.github/workflows/docker_region.yml` | 同上 |
| `.github/workflows/docker_scheduler.yml` | 同上 |
| `.github/workflows/docker_webconsole.yml` | 同上 |

**GitHub Secrets 配置要求：**
- `SWR_USERNAME`：华为云 SWR 用户名（格式：`cn-north-4@<AK>`）
- `SWR_PASSWORD`：华为云 SWR 登录密码（对应 SK）

---

## 二、HCS 8.6.0 IMS 镜像服务 API 适配

**文件：** `vendor/yunion.io/x/cloudmux/pkg/multicloud/hcso/image.go`

### 2.1 SImage 结构体新增字段

| 字段名 | JSON Key | 类型 | 说明 |
|--------|----------|------|------|
| `Architecture` | `architecture` | string | HCS 8.6.0 新增，CPU 架构，取值 `x86_64` 或 `aarch64`，替代旧 `__support_arm` |
| `HwFirmwareType` | `hw_firmware_type` | string | HCS 8.6.0 新增，固件类型，取值 `bios` 或 `uefi` |
| `SupportKVMInfiniband` | `__support_kvm_infiniband` | string | HCS 8.6.0 新增，KVM Infiniband 网卡支持标识 |

### 2.2 方法逻辑变更

| 方法 | 原始逻辑 | 变更后逻辑 |
|------|----------|------------|
| `getNormalizedImageInfo()` | 通过 `__support_arm=="true"` 判断 arm | 优先用 `architecture=="aarch64"` 判断，兼容旧 `__support_arm` |
| `GetOsArch()` | 通过 normalize 工具推断架构 | 优先返回 `architecture` 字段，字段为空时回退 normalize |
| `GetBios()` | 通过 normalize 工具推断 BIOS 类型 | 优先使用 `hw_firmware_type`，字段为空时回退 normalize |
| `ImportImageJob()` | 未传入 `virtual_env_type` 和 `architecture` | 新增传入 `virtual_env_type=FusionCompute` 和 `architecture` 参数 |

### 2.3 新增函数

- `archToHCS(osArch string) string`：将内部架构标识转换为 HCS IMS 接口所需的 `x86_64` 或 `aarch64`

---

## 三、HCS 8.6.0 ECS 弹性云服务器 API 适配

### 3.1 instance.go

**文件：** `vendor/yunion.io/x/cloudmux/pkg/multicloud/hcso/instance.go`

#### IpAddress 结构体新增字段

| 字段名 | JSON Key | 类型 | 说明 |
|--------|----------|------|------|
| `NetworkTags` | `network_tags` | []string | HCS 8.6.0 新增，网络标签列表 |

#### OSExtendedVolumesVolumesAttached 结构体新增字段

| 字段名 | JSON Key | 类型 | 说明 |
|--------|----------|------|------|
| `HwPassthrough` | `hw:passthrough` | string | 磁盘模式：`false`=VBD，`true`=SCSI |
| `Multiattach` | `multiattach` | string | 是否为共享磁盘 |

#### 新增结构体

- `VirtualCipherCardDevice`：虚拟密码卡设备信息（含 `uuid`、`name`、`mac` 字段）

#### SInstance 结构体新增字段

| 字段名 | JSON Key | 类型 | 说明 |
|--------|----------|------|------|
| `OSEXTSRVATTROsHostname` | `OS-EXT-SRV-ATTR:os_hostname` | string | 云服务器内部 hostname |
| `SystemSerialNumber` | `system_serial_number` | string | 云服务器 serial number |
| `VirtualCipherCardDevice` | `virtual_cipher_card_device` | []VirtualCipherCardDevice | 虚拟密码卡设备信息 |

### 3.2 instancetype.go

**文件：** `vendor/yunion.io/x/cloudmux/pkg/multicloud/hcso/instancetype.go`

#### OSExtraSpecs 结构体新增字段

| 字段名 | JSON Key | 类型 | 说明 |
|--------|----------|------|------|
| `ResourceType` | `resource_type` | string | 资源类型 |
| `EcsVirtualizationEnvTypes` | `ecs:virtualization_env_types` | string | 虚拟化类型：FusionCompute(XEN) 或 CloudCompute(KVM) |
| `CapabilitiesCpuInfoArch` | `capabilities:cpu_info:arch` | string | CPU 架构，如 `x86_64`、`aarch64` |
| `HwCpuMode` | `hw:cpu_mode` | string | CPU 模式 |
| `HwCpuModel` | `hw:cpu_model` | string | CPU 型号 |
| `PciPassthroughAlias` | `pci_passthrough:alias` | string | 本地直通 GPU 型号和数量 |

#### GetCpuArch() 方法变更

| 原始逻辑 | 变更后逻辑 |
|----------|------------|
| 通过 `ecs:instance_architecture` 和 flavor ID 前缀（`k` 开头）判断 | 优先使用 `capabilities:cpu_info:arch`，其次 `ecs:instance_architecture`，最后 flavor ID 前缀 |

---

## 四、HCS 8.6.0 EVS 云硬盘 API 适配

**文件：** `vendor/yunion.io/x/cloudmux/pkg/multicloud/hcso/disk.go`

### 4.1 VolumeImageMetadata 结构体新增字段

| 字段名 | JSON Key | 类型 | 说明 |
|--------|----------|------|------|
| `Architecture` | `architecture` | string | HCS 8.6.0 新增，镜像 CPU 架构，取值 `x86_64` 或 `aarch64` |
| `HwFirmwareType` | `hw_firmware_type` | string | HCS 8.6.0 新增，固件类型，取值 `bios` 或 `uefi` |

### 4.2 GetDriver() 方法变更

| 原始逻辑 | 变更后逻辑 |
|----------|------------|
| 硬编码返回 `"scsi"` | 根据 `Metadata.AttachedMode` 动态判断，`scsi` 时返回 `"scsi"`，否则返回 `"virtio"` |

---

## 五、未适配说明

| 服务 | 说明 |
|------|------|
| AS 弹性伸缩 | hcso 无对应实现文件，属预期缺失，cloudpods 不直接管理云厂商 AS 服务 |
| BMS 裸金属服务器 | HCS 8.6.0 中 BMS 复用 ECS 接口（通过 `tags=__type_baremetal` 区分），host.go 无需修改 |
| HCEOS | 为操作系统级别 API，与 cloudpods 云管平台无直接对应关系，无需适配 |

---

## 七、IAM 用户管理接口路径修正

**文件：** `vendor/yunion.io/x/cloudmux/pkg/multicloud/hcso/client/modules/mod_users.go`

| 项目 | 原始值 | 修改后 | 说明 |
|------|--------|--------|------|
| `version` | `"v3.0/OS-USER"` | `"v3"` | HCS 8.6.0 IAM 2.0 文档用户接口路径为 `/v3/users` |
| `List()` | 临时 `SetVersion("v3")` | 直接调用，无需切换 | 默认已是 v3 |
| `Delete()` | 临时 `SetVersion("v3")` | 直接调用，无需切换 | 默认已是 v3 |
| `ListGroups()` | 临时 `SetVersion("v3")` | 直接调用，无需切换 | 默认已是 v3 |
| `ResetPassword()` | 使用默认 `v3.0/OS-USER` | 使用默认 `v3` | 修复原本路径不一致问题 |

**接口路径变化：**
- 原：`/v3.0/OS-USER/users`
- 改：`/v3/users`

---

## 八、EVS 云硬盘 API 适配（补充）

**文件：** `vendor/yunion.io/x/cloudmux/pkg/multicloud/hcso/disk.go`

### 8.1 DiskMeta 结构体新增字段

| 字段名 | JSON Key | 类型 | 说明 |
|--------|----------|------|------|
| `HwPassthrough` | `hw:passthrough` | string | HCS 8.6.0 磁盘设备类型：`true`=SCSI，`false`/缺省=VBD |

### 8.2 SDisk 结构体新增字段

| 字段名 | JSON Key | 类型 | 说明 |
|--------|----------|------|------|
| `Shareable` | `shareable` | string | 是否为可共享云硬盘 |
| `OsVolHostAttrHost` | `os-vol-host-attr:host` | string | 云硬盘所在主机 |
| `OsVolTenantAttrTenantId` | `os-vol-tenant-attr:tenant_id` | string | 所属租户ID |
| `OsVolMigStatusAttrMigstat` | `os-vol-mig-status-attr:migstat` | string | 云硬盘迁移状态 |
| `EncryptionInfo` | `encryption_info` | interface{} | 云硬盘主机侧加密信息 |

### 8.3 GetDriver() 方法优化

| 原始逻辑 | 变更后逻辑 |
|----------|------------|
| 通过 `Metadata.AttachedMode` 判断 | 优先通过 `Metadata.HwPassthrough == "true"` 判断，回退到 `AttachedMode` |

### 8.4 snapshot.go

经核查，`SSnapshot` 结构体字段已完整覆盖 HCS 8.6.0 EVS 快照响应字段，无需修改。

---

## 九、VPC 虚拟私有云 API 适配

### 9.1 vpc.go - SVpc 新增字段

| 字段名 | JSON Key | 说明 |
|--------|----------|------|
| `Description` | `description` | VPC 描述信息 |
| `ExtendCidr` | `extend_cidr` | VPC 扩展网段 |

### 9.2 network.go - SNetwork 新增字段

| 字段名 | JSON Key | 说明 |
|--------|----------|------|
| `Description` | `description` | 子网描述信息 |
| `Ipv6CidrBlock` | `cidr_v6` | IPv6 网段 |
| `Ipv6GatewayIP` | `gateway_ip_v6` | IPv6 网关地址 |

### 9.3 securitygroup.go - SSecurityGroup 新增字段

| 字段名 | JSON Key | 说明 |
|--------|----------|------|
| `ProjectID` | `project_id` | 资源所属项目ID |
| `CreatedAt` | `created_at` | 创建时间 |
| `UpdatedAt` | `updated_at` | 更新时间 |
| `SysTags` | `sys_tags` | 系统标签列表 |

### 9.4 vpc_peering.go - SVpcPeering 新增字段

| 字段名 | JSON Key | 说明 |
|--------|----------|------|
| `Description` | `description` | 对等连接描述 |
| `CreatedAt` | `created_at` | 创建时间 |
| `UpdatedAt` | `updated_at` | 更新时间 |

---

## 十、RDS 云数据库 API 适配

**文件：** `vendor/yunion.io/x/cloudmux/pkg/multicloud/hcso/dbinstance.go`

### SDBInstance 结构体新增字段

| 字段名 | JSON Key | 类型 | 说明 |
|--------|----------|------|------|
| `Alias` | `alias` | string | 实例备注信息 |
| `PrivateDnsNames` | `private_dns_names` | []string | 实例内网域名列表 |
| `EnableSsl` | `enable_ssl` | bool | 是否开启SSL |
| `Cpu` | `cpu` | string | CPU大小 |
| `Mem` | `mem` | string | 内存大小（GB）|
| `ReadOnlyByUser` | `read_only_by_user` | bool | 用户设置的只读状态 |
| `ChargeInfo` | `charge_info` | interface{} | 计费信息 |
| `AssociatedWithDdm` | `associated_with_ddm` | bool | 是否已被DDM实例关联 |

---

## 十一、EIP 弹性公网IP API 适配

**文件：** `vendor/yunion.io/x/cloudmux/pkg/multicloud/hcso/eip.go`

### SEipAddress 结构体新增字段

| 字段名 | JSON Key | 说明 |
|--------|----------|------|
| `ProjectID` | `project_id` | 项目ID |
| `CreatedAt` | `created_at` | 创建时间（新格式，与 create_time 并存）|
| `UpdatedAt` | `updated_at` | 更新时间 |
| `BandwidthDirection` | `bandwidth_direction` | 带宽限速模式：egress/bidirectional |
| `RouterID` | `router_id` | 关联路由ID |
| `ExternalNetID` | `external_net_id` | 外部网络ID |
| `OpStatus` | `op_status` | 运营状态 |
| `Expiry` | `expiry` | 到期时间 |

### Bandwidth 结构体新增字段

| 字段名 | JSON Key | 说明 |
|--------|----------|------|
| `Description` | `description` | 带宽描述 |
| `Shared` | `shared` | 是否租户级别资源 |
| `RuleType` | `rule_type` | 带宽类型 |
| `Direction` | `direction` | 限速模式：egress/bidirectional |
| `Status` | `status` | 带宽状态 |

---

## 十二、SFS Turbo 高性能弹性文件服务 API 适配

**文件：** `vendor/yunion.io/x/cloudmux/pkg/multicloud/hcso/sfs-turbo.go`

### SfsTurbo 结构体新增字段

| 字段名 | JSON Key | 说明 |
|--------|----------|------|
| `Version` | `version` | 文件系统版本号 |
| `Bandwidth` | `bandwidth` | 文件系统总带宽（MB/s）|
| `OptionalEndpoint` | `optional_endpoint` | 可选的挂载IP地址 |
| `Scenario` | `scenario` | 文件系统场景：HCS/HCSStandard |
| `ExactShareType` | `exact_share_type` | 文件系统详细规格 |
| `Ownership` | `ownership` | 文件系统所属云服务名称 |
| `InstanceId` | `instanceId` | 节点ID（预留字段）|
| `InstanceType` | `instanceType` | 节点类型（预留字段）|
| `StatusDetail` | `statusDetail` | 请求ID（预留字段）|


1. 所有代码修改处均保留原始逻辑注释（`[ORIGIN]`）和变更说明（`[CHANGED]`），便于回溯
2. 新增字段均为可选字段，旧版 HCS 环境返回为空时不影响现有逻辑
3. 架构判断逻辑均保留了向下兼容，支持无 `architecture` 字段的旧版本环境
