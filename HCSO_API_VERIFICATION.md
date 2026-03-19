# HCSO Client 接口路径验证报告

验证时间：2026-03-20
验证依据：HCS 8.6.0 API 文档

---

## 验证结论总览

| 服务 | Module | 代码路径 | 文档路径 | 状态 |
|------|--------|---------|---------|------|
| ECS 云服务器 | mod_servers.go | `v1/cloudservers` | `v1/{tenant_id}/cloudservers` | ✅ 一致 |
| ECS 云服务器(Nova) | mod_servers.go | `v2.1/servers` | `v2.1/{tenant_id}/servers` | ✅ 一致 |
| ECS 规格 | mod_flavors.go | `v1/flavors` | `v1/{tenant_id}/cloudservers/flavors` | ✅ 一致 |
| ECS 网卡 | mod_interface.go | `v2/os-interface` | `v2/{tenant_id}/servers/{id}/os-interface` | ✅ 一致 |
| ECS 密钥对 | mod_keypairs.go | `v2/os-keypairs` | `v2.1/{tenant_id}/os-keypairs` | ⚠️ v2 vs v2.1（兼容）|
| ECS 安全组 | mod_secgroups.go | `v2.1/os-security-groups` | `v2.1/{tenant_id}/os-security-groups` | ✅ 一致 |
| ECS 可用区 | mod_zones.go | `v2/os-availability-zone` | `v2.1/{tenant_id}/os-availability-zone` | ⚠️ v2 vs v2.1（兼容）|
| EVS 云硬盘 | mod_disks.go | `v2/cloudvolumes` | `v2/{tenant_id}/cloudvolumes` | ✅ 一致 |
| EVS 快照 | mod_snapshots.go | `v2/snapshots` | `v2/{tenant_id}/snapshots` | ✅ 一致 |
| EVS 配额 | mod_quotas.go | `v1/quotas` | `v1/{tenant_id}/quotas` | ✅ 一致 |
| IMS 镜像 | mod_images.go | `v2/cloudimages` | `v2/cloudimages` | ✅ 一致 |
| VPC | mod_vpcs.go | `v1/vpcs` | `v1/{tenant_id}/vpcs` | ✅ 一致 |
| 子网 | mod_subnets.go | `v1/subnets` | `v1/{tenant_id}/subnets` | ✅ 一致 |
| EIP | mod_eips.go | `v1/publicips` | `v1/{tenant_id}/publicips` | ✅ 一致 |
| 带宽 | mod_bandwidths.go | `v1/bandwidths` | `v1/{tenant_id}/bandwidths` | ✅ 一致 |
| 安全组 | mod_secgroups.go | `v1/security-groups` | `v1/{tenant_id}/security-groups` | ✅ 一致 |
| 安全组规则 | mod_secgroup_rules.go | `v1/security-group-rules` | `v1/{tenant_id}/security-group-rules` | ✅ 一致 |
| 端口 | mod_port.go | `v1/ports` | `v1/{tenant_id}/ports` | ✅ 一致 |
| VPC对等连接 | mod_vpc_peerings.go | `v2.0/vpc/peerings` | `v2.0/vpc/peerings/{id}` | ✅ 一致 |
| VPC路由 | mod_vpc_routes.go | `v2.0/vpc/routes` | `v2.0/vpc/routes/{id}` | ✅ 一致 |
| ELB 负载均衡器 | mod_loadbalancers.go | `v2.0/lbaas/loadbalancers` | `v2.0/lbaas/loadbalancers`（旧）/ `v3/elb/loadbalancers`（新） | ✅ 旧版仍支持 |
| ELB 监听器 | mod_loadbalancer_listeners.go | `v2.0/lbaas/listeners` | `v2.0/lbaas/listeners`（旧）/ `v3/elb/listeners`（新） | ✅ 旧版仍支持 |
| ELB 后端组 | mod_loadbalancer_backend_group.go | `v2.0/lbaas/pools` | `v2.0/lbaas/pools`（旧）/ `v3/elb/pools`（新） | ✅ 旧版仍支持 |
| ELB 后端服务器 | mod_loadbalancer_backend.go | `v2.0/lbaas/members` | `v2.0/lbaas/pools/{id}/members`（旧） | ✅ 旧版仍支持 |
| ELB 健康检查 | mod_loadbalancer_healthcheck.go | `v2.0/lbaas/healthmonitors` | `v2.0/lbaas/healthmonitors`（旧） | ✅ 旧版仍支持 |
| ELB 证书 | mod_loadbalancer_certificates.go | `v2.0/lbaas/certificates` | `v2.0/lbaas/certificates`（旧） | ✅ 旧版仍支持 |
| ELB 转发策略 | mod_loadbalancer_policies.go | `v2.0/lbaas/l7policies` | `v2.0/lbaas/l7policies`（旧） | ✅ 旧版仍支持 |
| ELB 白名单 | mod_loadbalancer_whitelists.go | `v2.0/lbaas/whitelists` | `v2.0/lbaas/whitelists`（旧） | ✅ 旧版仍支持 |
| NAT 网关 | mod_natgateway.go | `v2.0/nat_gateways` | `v2.0/nat_gateways` | ✅ 一致 |
| SNAT 规则 | mod_snat_rules.go | `v2.0/snat_rules` | `v2.0/snat_rules` | ✅ 一致 |
| DNAT 规则 | mod_dnat_rules.go | `v2.0/dnat_rules` | `v2.0/dnat_rules` | ✅ 一致 |
| RDS 数据库实例 | mod_dbinstance.go | `v3/instances` | `v3/{project_id}/instances` | ✅ 一致 |
| RDS 备份 | mod_dbinstance_backup.go | `v3/backups` | `v3/{project_id}/backups` | ✅ 一致 |
| RDS 数据存储 | mod_dbinstance_datastore.go | `v3/datastores` | `v3/{project_id}/datastores` | ✅ 一致 |
| RDS 规格 | mod_dbinstance_flavor.go | `v3/flavors` | `v3/{project_id}/flavors` | ✅ 一致 |
| RDS 存储类型 | mod_dbinstance_storage.go | `v3/storage-type` | `v3/{project_id}/storage-type` | ✅ 一致 |
| DCS 缓存 | mod_elasticcache.go | `v1.0/instances` | `v1.0/{project_id}/instances` | ✅ 一致 |
| CES 监控 | mod_ces.go | `V1.0/metrics` | `V1.0/{project_id}/metrics` | ✅ 一致 |
| CTS 审计 | mod_traces.go | `v2.0/system/trace` | `v2.0/{project_id}/system/trace` | ✅ 一致 |
| IAM 用户 | mod_users.go | `v3/users` | `v3/users` | ✅ 已修正（原 v3.0/OS-USER）|
| IAM 用户组 | mod_groups.go | `v3/groups` | `v3/groups` | ✅ 一致 |
| IAM 角色 | mod_roles.go | `v3/roles` | `v3/roles` | ✅ 一致 |
| IAM 项目 | mod_projects.go | `v3/projects` | `v3/projects` | ✅ 一致 |
| IAM 域 | mod_domains.go | `v3/auth/domains` | `v3/auth/domains` | ✅ 一致 |
| IAM 区域 | mod_regions.go | `v3/regions` | `v3/regions` | ✅ 一致 |
| IAM 服务 | mod_services.go | `v3/services` | `v3/services` | ✅ 一致 |
| IAM 端点 | mod_endpoints.go | `v3/endpoints` | `v3/endpoints` | ✅ 一致 |
| IAM 凭据 | mod_credential.go | `v3.0/OS-CREDENTIAL/credentials` | `v3.0/OS-CREDENTIAL/securitytokens` | ⚠️ 路径不同（功能不同）|
| IAM SAML | mod_saml_provider.go | `v3/OS-FEDERATION/identity_providers` | `v3/OS-FEDERATION/identity_providers` | ✅ 一致 |
| IAM 映射 | mod_mapping.go | `v3/OS-FEDERATION/mappings` | `v3/OS-FEDERATION/mappings` | ✅ 一致 |
| EPS 企业项目 | mod_enterpriceprojects.go | `v1.0/enterprise-projects` | `v1.0/enterprise-projects` | ✅ 一致 |
| SFS Turbo | mod_sfs.go | `v1/sfs-turbo/shares` | `v1/{project_id}/sfs-turbo/shares` | ✅ 一致 |
| Jobs 任务 | mod_jobs.go | `v1/jobs` | `v1/{tenant_id}/jobs` | ✅ 一致 |

---

## 问题汇总

### 已修复
| 文件 | 问题 | 修复内容 |
|------|------|----------|
| mod_users.go | version `v3.0/OS-USER` 与文档不符 | 改为 `v3`，同时清理冗余的 `SetVersion("v3")` 调用 |

### 轻微差异（向下兼容，无需修改）
| 文件 | 问题 | 说明 |
|------|------|------|
| mod_keypairs.go | 代码用 `v2`，文档标注 `v2.1` | OpenStack v2/v2.1 向下兼容，HCS 同时支持 |
| mod_zones.go | 代码用 `v2`，文档标注 `v2.1` | 同上 |

### 可升级项（非强制，v2.0 旧版仍支持）
| 文件 | 当前版本 | 新版本 | 说明 |
|------|---------|--------|------|
| mod_loadbalancers.go 等所有 ELB module | `v2.0/lbaas/*` | `v3/elb/*` | HCS 8.6.0 新增 v3 接口，v2.0 旧版接口仍保留，功能等价但 v3 有更多新特性 |

### 注意项
| 文件 | 说明 |
|------|------|
| mod_credential.go | 代码路径 `v3.0/OS-CREDENTIAL/credentials` 用于访问密钥管理，文档中 `v3.0/OS-CREDENTIAL/securitytokens` 是获取临时 AK/SK，两者功能不同，代码路径正确 |

---

## ELB v3 接口说明

HCS 8.6.0 ELB 同时提供两套接口：

| 版本 | 路径格式 | 状态 | 说明 |
|------|---------|------|------|
| v2.0（旧） | `/v2.0/lbaas/{resource}` | 仍支持 | 代码当前使用此版本 |
| v3（新） | `/v3/{project_id}/elb/{resource}` | 推荐 | 新增功能更丰富，如安全策略、日志等 |

**建议：** 当前 v2.0 接口可正常使用，如需使用 v3 新特性（安全策略、logtank 等），可在后续版本中逐步迁移至 v3。
