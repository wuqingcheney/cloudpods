# HCSO 接口汇总文档

生成时间：2026-03-20

> 基于 `vendor/yunion.io/x/cloudmux/pkg/multicloud/hcso/client/modules/` 目录下共 **53 个 module 文件**整理。

---

## 一、接口统计概览

| 服务 | 服务名常量 | Module 数量 | 说明 |
|------|-----------|------------|------|
| ECS 弹性云服务器 | ServiceNameECS | 7 | 云服务器、规格、网卡、密钥对、安全组、可用区 |
| EVS 云硬盘 | ServiceNameEVS | 3 | 云硬盘、快照、配额 |
| IMS 镜像服务 | ServiceNameIMS | 2 | 镜像（标准+OpenStack兼容） |
| VPC 虚拟私有云 | ServiceNameVPC | 8 | VPC、子网、EIP、安全组、端口、带宽、对等连接、路由 |
| ELB 弹性负载均衡 | ServiceNameELB | 8 | 负载均衡器、监听器、后端组、后端、健康检查、证书、策略、白名单 |
| IAM 身份认证 | ServiceNameIAM | 10 | 用户、用户组、角色、项目、域、区域、服务、端点、凭据、SAML |
| RDS 关系型数据库 | ServiceNameRDS | 6 | 实例、备份、数据存储、规格、存储类型、任务 |
| NAT 网关 | ServiceNameNAT | 3 | NAT网关、SNAT规则、DNAT规则 |
| DCS 分布式缓存 | ServiceNameDCS | 2 | 弹性缓存实例、可用区 |
| CES 云监控 | ServiceNameCES | 1 | 监控指标 |
| CTS 云审计 | ServiceNameCTS | 1 | 操作追踪 |
| EPS 企业项目 | ServiceNameEPS | 1 | 企业项目 |
| SFS 文件存储 | ServiceNameSFSTurbo | 1 | SFS Turbo 共享文件系统 |
| Jobs 任务 | （动态） | 1 | 异步任务查询 |

**合计：54 个接口模块（含部分服务有多个 Manager 变体）**

---

## 二、各服务接口详情

### 2.1 ECS 弹性云服务器

| Module 文件 | 接口路径 | 版本 | 说明 |
|------------|---------|------|------|
| mod_servers.go | `/{tenant_id}/cloudservers` | v1 | 云服务器管理（HCS原生）|
| mod_servers.go | `/{tenant_id}/cloudservers` | v1.1 | 创建包年包月云服务器 |
| mod_servers.go | `/{tenant_id}/servers` | v2.1 | 云服务器管理（OpenStack兼容） |
| mod_servers.go | `/{tenant_id}/cloudservers` | v2 | 批量操作云服务器 |
| mod_flavors.go | `/{tenant_id}/flavors` | v1 | 云服务器规格查询 |
| mod_interface.go | `/{tenant_id}/servers/{id}/os-interface` | v2 | 网卡管理 |
| mod_keypairs.go | `/{tenant_id}/os-keypairs` | v2 | SSH 密钥对管理 |
| mod_secgroups.go | `/{tenant_id}/os-security-groups` | v2.1 | 安全组（OpenStack兼容） |
| mod_zones.go | `/{tenant_id}/os-availability-zone` | v2 | 可用区查询 |

### 2.2 EVS 云硬盘

| Module 文件 | 接口路径 | 版本 | 说明 |
|------------|---------|------|------|
| mod_disks.go | `/{tenant_id}/cloudvolumes` | v2 | 云硬盘 CRUD |
| mod_disks.go | `/{tenant_id}/cloudvolumes` | v2.1 | 云硬盘扩容等操作 |
| mod_snapshots.go | `/{tenant_id}/snapshots` | v2 | 快照管理 |
| mod_snapshots.go | `/{tenant_id}/os-vendor-snapshots` | v2 | 快照回滚（OpenStack兼容） |
| mod_quotas.go | `/{tenant_id}/quotas` | v1 | 配额查询 |

### 2.3 IMS 镜像服务

| Module 文件 | 接口路径 | 版本 | 说明 |
|------------|---------|------|------|
| mod_images.go | `/cloudimages` | v2 | 镜像 CRUD（HCS原生），不含 project_id |
| mod_images.go | `/images` | v2 | 镜像删除（OpenStack兼容） |

### 2.4 VPC 虚拟私有云

| Module 文件 | 接口路径 | 版本 | 说明 |
|------------|---------|------|------|
| mod_vpcs.go | `/{tenant_id}/vpcs` | v1 | VPC 管理 |
| mod_subnets.go | `/{tenant_id}/subnets` | v1 | 子网管理 |
| mod_eips.go | `/{tenant_id}/publicips` | v1 | 弹性公网 IP |
| mod_secgroup_rules.go | `/{tenant_id}/security-group-rules` | v1 | 安全组规则 |
| mod_secgroups.go | `/{tenant_id}/security-groups` | v1 | 安全组 |
| mod_port.go | `/{tenant_id}/ports` | v1 | 端口管理 |
| mod_bandwidths.go | `/{tenant_id}/bandwidths` | v1 | 带宽管理 |
| mod_vpc_peerings.go | `/vpc/peerings` | v2.0 | VPC 对等连接 |
| mod_vpc_routes.go | `/vpc/routes` | v2.0 | VPC 路由 |

### 2.5 ELB 弹性负载均衡

| Module 文件 | 接口路径 | 版本 | 说明 |
|------------|---------|------|------|
| mod_loadbalancers.go | `/lbaas/loadbalancers` | v2.0 | 负载均衡器 |
| mod_loadbalancer_listeners.go | `/lbaas/listeners` | v2.0 | 监听器 |
| mod_loadbalancer_backend_group.go | `/lbaas/pools` | v2.0 | 后端服务器组 |
| mod_loadbalancer_backend.go | `/lbaas/{pool_id}/members` | v2.0 | 后端服务器 |
| mod_loadbalancer_healthcheck.go | `/lbaas/healthmonitors` | v2.0 | 健康检查 |
| mod_loadbalancer_certificates.go | `/lbaas/certificates` | v2.0 | SSL 证书 |
| mod_loadbalancer_policies.go | `/lbaas/l7policies` | v2.0 | 转发策略 |
| mod_loadbalancer_rules.go | `/lbaas/{l7policy_id}/rules` | v2.0 | 转发规则 |
| mod_loadbalancer_whitelists.go | `/lbaas/whitelists` | v2.0 | 白名单 |

### 2.6 IAM 身份认证

| Module 文件 | 接口路径 | 版本 | 说明 |
|------------|---------|------|------|
| mod_users.go | `/OS-USER/users` | v3.0 | 用户管理 |
| mod_groups.go | `/groups` | v3 | 用户组管理 |
| mod_roles.go | `/roles` | v3 | 角色管理 |
| mod_projects.go | `/projects` | v3 | 项目管理 |
| mod_domains.go | `/auth/domains` | v3 | 域管理 |
| mod_regions.go | `/regions` | v3 | 区域查询 |
| mod_services.go | `/services` | v3 | 服务查询 |
| mod_endpoints.go | `/endpoints` | v3 | 端点查询 |
| mod_credential.go | `/OS-CREDENTIAL/credentials` | v3.0 | 访问凭据 |
| mod_saml_provider.go | `/OS-FEDERATION/identity_providers` | v3 | SAML 身份提供商 |
| mod_mapping.go | `/OS-FEDERATION/mappings` | v3 | SAML 映射 |

### 2.7 RDS 关系型数据库

| Module 文件 | 接口路径 | 版本 | 说明 |
|------------|---------|------|------|
| mod_dbinstance.go | `/{tenant_id}/instances` | v3 | 数据库实例 CRUD |
| mod_dbinstance_backup.go | `/{tenant_id}/backups` | v3 | 备份管理 |
| mod_dbinstance_datastore.go | `/{tenant_id}/datastores` | v3 | 数据库版本 |
| mod_dbinstance_flavor.go | `/{tenant_id}/flavors` | v3 | 数据库规格 |
| mod_dbinstance_storage.go | `/{tenant_id}/storage-type` | v3 | 存储类型 |
| mod_dbinstance_job.go | `/{tenant_id}/jobs` | v3 | 异步任务 |

### 2.8 NAT 网关

| Module 文件 | 接口路径 | 版本 | 说明 |
|------------|---------|------|------|
| mod_natgateway.go | `/nat_gateways` | v2.0 | NAT 网关 |
| mod_snat_rules.go | `/snat_rules` | v2.0 | SNAT 规则 |
| mod_dnat_rules.go | `/dnat_rules` | v2.0 | DNAT 规则 |

### 2.9 DCS 分布式缓存

| Module 文件 | 接口路径 | 版本 | 说明 |
|------------|---------|------|------|
| mod_elasticcache.go | `/{tenant_id}/instances` | v1.0 | 缓存实例管理 |
| mod_elasticcache.go | `/availableZones` | v1.0 | 可用区查询 |

### 2.10 其他服务

| Module 文件 | 服务 | 接口路径 | 版本 | 说明 |
|------------|------|---------|------|------|
| mod_ces.go | CES 云监控 | `/{tenant_id}/metrics` | V1.0 | 监控指标查询 |
| mod_traces.go | CTS 云审计 | `/{tenant_id}/system/trace` | v2.0 | 操作追踪 |
| mod_enterpriceprojects.go | EPS 企业项目 | `/enterprise-projects` | v1.0 | 企业项目管理 |
| mod_sfs.go | SFS Turbo | `/{tenant_id}/sfs-turbo/shares` | v1 | 共享文件系统 |
| mod_jobs.go | Jobs | `/{tenant_id}/jobs` | v1 | 异步任务查询（通用） |

---

## 三、HCS 8.6.0 适配变更的接口

| 服务 | 接口 | 变更内容 |
|------|------|----------|
| IMS | `GET /v2/cloudimages` | 响应新增 `architecture`、`hw_firmware_type`、`__support_kvm_infiniband` 字段 |
| IMS | `POST /v2/cloudimages/action` | 请求新增 `virtual_env_type`、`architecture` 参数 |
| ECS | `GET /v1/{tenant_id}/cloudservers/{id}` | 响应新增 `system_serial_number`、`virtual_cipher_card_device`、`OS-EXT-SRV-ATTR:os_hostname` 字段 |
| ECS | `GET /v1/{tenant_id}/cloudservers/detail` | addresses 新增 `network_tags`，volumes_attached 新增 `hw:passthrough`、`multiattach` |
| ECS | `GET /v1/{tenant_id}/flavors` | os_extra_specs 新增 `capabilities:cpu_info:arch`、`ecs:virtualization_env_types` 等字段 |
| EVS | `GET /v2/{tenant_id}/cloudvolumes/{id}` | volume_image_metadata 新增 `architecture`、`hw_firmware_type` 字段 |

---

## 四、接口版本一致性确认

经核查，所有 Module 中使用的 API 版本与 HCS 8.6.0 文档定义一致，**无需修改接口路径或版本号**，本次适配仅涉及请求/响应数据字段层面的补充。
