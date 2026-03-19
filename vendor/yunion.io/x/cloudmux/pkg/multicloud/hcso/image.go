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

package hcso

import (
	"context"
	"fmt"
	"strings"
	"time"

	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/imagetools"

	"yunion.io/x/cloudmux/pkg/apis"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/cloudmux/pkg/multicloud/huawei"
)

type TImageOwnerType string

const (
	ImageOwnerPublic TImageOwnerType = "gold"    // 公共镜像：gold
	ImageOwnerSelf   TImageOwnerType = "private" // 私有镜像：private
	ImageOwnerShared TImageOwnerType = "shared"  // 共享镜像：shared

	EnvFusionCompute = "FusionCompute"
	EnvIronic        = "Ironic"
)

const (
	ImageStatusQueued  = "queued"  // queued：表示镜像元数据已经创建成功，等待上传镜像文件。
	ImageStatusSaving  = "saving"  // saving：表示镜像正在上传文件到后端存储。
	ImageStatusDeleted = "deleted" // deleted：表示镜像已经删除。
	ImageStatusKilled  = "killed"  // killed：表示镜像上传错误。
	ImageStatusActive  = "active"  // active：表示镜像可以正常使用
)

// https://support.huaweicloud.com/api-ims/zh-cn_topic_0020091565.html
type SImage struct {
	multicloud.SImageBase
	huawei.HuaweiTags
	storageCache *SStoragecache

	// normalized image info
	imgInfo *imagetools.ImageInfo

	Schema          string    `json:"schema"`
	MinDiskGB       int64     `json:"min_disk"`
	CreatedAt       time.Time `json:"created_at"`
	ImageSourceType string    `json:"__image_source_type"`
	ContainerFormat string    `json:"container_format"`
	File            string    `json:"file"`
	UpdatedAt       time.Time `json:"updated_at"`
	Protected       bool      `json:"protected"`
	Checksum        string    `json:"checksum"`
	ID              string    `json:"id"`
	Isregistered    string    `json:"__isregistered"`
	MinRamMB        int       `json:"min_ram"`
	Lazyloading     string    `json:"__lazyloading"`
	Owner           string    `json:"owner"`
	OSType          string    `json:"__os_type"`
	Imagetype       string    `json:"__imagetype"`
	Visibility      string    `json:"visibility"`
	VirtualEnvType  string    `json:"virtual_env_type"`
	Platform        string    `json:"__platform"`
	SizeGB          int       `json:"size"`
	ImageSize       int64     `json:"__image_size"`
	OSBit           string    `json:"__os_bit"`
	OSVersion       string    `json:"__os_version"`
	Name            string    `json:"name"`
	Self            string    `json:"self"`
	DiskFormat      string    `json:"disk_format"`
	Status          string    `json:"status"`
	// [CHANGED] HCS 8.6.0 新增字段：CPU架构，取值为 x86_64 或 aarch64
	// [ORIGIN] 原始通过 __support_arm 字段（取值 "true"/"false"）间接判断是否为 arm 架构
	Architecture string `json:"architecture"`
	// [CHANGED] HCS 8.6.0 新增字段：固件类型，取值为 bios 或 uefi
	// [ORIGIN] 原始通过 normalize 工具从镜像名称/平台信息推断固件类型
	HwFirmwareType         string `json:"hw_firmware_type"`
	SupportKVMFPGAType     string `json:"__support_kvm_fpga_type"`
	SupportKVMNVMEHIGHIO   string `json:"__support_nvme_highio"`
	SupportLargeMemory     string `json:"__support_largememory"`
	SupportDiskIntensive   string `json:"__support_diskintensive"`
	SupportHighPerformance string `json:"__support_highperformance"`
	SupportXENGPUType      string `json:"__support_xen_gpu_type"`
	SupportKVMGPUType      string `json:"__support_kvm_gpu_type"`
	SupportGPUT4           string `json:"__support_gpu_t4"`
	SupportKVMAscend310    string `json:"__support_kvm_ascend_310"`
	SupportArm             string `json:"__support_arm"` // [ORIGIN] 旧版架构判断字段，HCS 8.6.0 由 architecture 字段替代，保留用于兼容
	// [CHANGED] HCS 8.6.0 新增字段：KVM Infiniband 网卡支持标识
	SupportKVMInfiniband string `json:"__support_kvm_infiniband"`
}

func (self *SImage) GetMinRamSizeMb() int {
	return self.MinRamMB
}

func (self *SImage) GetId() string {
	return self.ID
}

func (self *SImage) GetName() string {
	return self.Name
}

func (self *SImage) GetGlobalId() string {
	return self.ID
}

func (self *SImage) GetStatus() string {
	switch self.Status {
	case ImageStatusQueued:
		return api.CACHED_IMAGE_STATUS_CACHING
	case ImageStatusActive:
		return api.CACHED_IMAGE_STATUS_ACTIVE
	case ImageStatusKilled:
		return api.CACHED_IMAGE_STATUS_CACHE_FAILED
	default:
		return api.CACHED_IMAGE_STATUS_CACHE_FAILED
	}
}

func (self *SImage) GetImageStatus() string {
	switch self.Status {
	case ImageStatusQueued:
		return cloudprovider.IMAGE_STATUS_QUEUED
	case ImageStatusActive:
		return cloudprovider.IMAGE_STATUS_ACTIVE
	case ImageStatusKilled:
		return cloudprovider.IMAGE_STATUS_KILLED
	default:
		return cloudprovider.IMAGE_STATUS_KILLED
	}
}

func (self *SImage) Refresh() error {
	new, err := self.storageCache.region.GetImage(self.GetId())
	if err != nil {
		return err
	}
	return jsonutils.Update(self, new)
}

func (self *SImage) GetImageType() cloudprovider.TImageType {
	switch self.Imagetype {
	case "gold":
		return cloudprovider.ImageTypeSystem
	case "private":
		return cloudprovider.ImageTypeCustomized
	case "shared":
		return cloudprovider.ImageTypeShared
	default:
		return cloudprovider.ImageTypeCustomized
	}
}

func (self *SImage) GetSizeByte() int64 {
	return int64(self.MinDiskGB) * 1024 * 1024 * 1024
}

func (self *SImage) getNormalizedImageInfo() *imagetools.ImageInfo {
	if self.imgInfo == nil {
		arch := "x86"
		// [ORIGIN] 原始逻辑：仅通过 __support_arm 字段判断架构
		// if strings.ToLower(self.SupportArm) == "true" {
		// 	arch = "arm"
		// }
		// [CHANGED] HCS 8.6.0 新增 architecture 字段（取值 x86_64 或 aarch64），优先使用；
		// 同时保留对旧字段 __support_arm 的兼容，以适配无 architecture 字段的旧版本环境
		if self.Architecture == "aarch64" || strings.ToLower(self.SupportArm) == "true" {
			arch = "arm"
		}
		imgInfo := imagetools.NormalizeImageInfo(self.ImageSourceType, arch, self.OSType, self.Platform, "")
		self.imgInfo = &imgInfo
	}

	return self.imgInfo
}

func (self *SImage) GetFullOsName() string {
	return self.ImageSourceType
}

func (self *SImage) GetOsType() cloudprovider.TOsType {
	return cloudprovider.TOsType(self.getNormalizedImageInfo().OsType)
}

func (self *SImage) GetOsDist() string {
	return self.getNormalizedImageInfo().OsDistro
}

func (self *SImage) GetOsVersion() string {
	return self.getNormalizedImageInfo().OsVersion
}

func (self *SImage) GetOsLang() string {
	return self.getNormalizedImageInfo().OsLang
}

func (self *SImage) GetOsArch() string {
	// [ORIGIN] 原始逻辑：通过 normalize 工具从镜像元数据推断架构
	// return self.getNormalizedImageInfo().OsArch
	// [CHANGED] HCS 8.6.0 响应中新增 architecture 字段（取值 x86_64 或 aarch64），直接使用；
	// 无此字段时（旧版本环境）回退到原始 normalize 逻辑
	if len(self.Architecture) > 0 {
		return self.Architecture
	}
	return self.getNormalizedImageInfo().OsArch
}

func (i *SImage) GetBios() cloudprovider.TBiosType {
	// [ORIGIN] 原始逻辑：通过 normalize 工具从镜像元数据推断 BIOS 类型
	// return cloudprovider.ToBiosType(i.getNormalizedImageInfo().OsBios)
	// [CHANGED] HCS 8.6.0 响应中新增 hw_firmware_type 字段（取值 bios 或 uefi），优先使用；
	// 无此字段时（旧版本环境）回退到原始 normalize 逻辑
	if len(i.HwFirmwareType) > 0 {
		return cloudprovider.ToBiosType(i.HwFirmwareType)
	}
	return cloudprovider.ToBiosType(i.getNormalizedImageInfo().OsBios)
}

func (self *SImage) GetMinOsDiskSizeGb() int {
	return int(self.MinDiskGB)
}

func (self *SImage) GetImageFormat() string {
	return self.DiskFormat
}

func (self *SImage) GetCreatedAt() time.Time {
	return self.CreatedAt
}

func (self *SImage) IsEmulated() bool {
	return false
}

func (self *SImage) Delete(ctx context.Context) error {
	return self.storageCache.region.DeleteImage(self.GetId())
}

func (self *SImage) GetIStoragecache() cloudprovider.ICloudStoragecache {
	return self.storageCache
}

func (self *SRegion) GetImage(imageId string) (*SImage, error) {
	image := &SImage{}
	err := DoGet(self.ecsClient.Images.Get, imageId, nil, image)
	if err != nil {
		return nil, errors.Wrap(err, "DoGet")
	}
	return image, nil
}

func excludeImage(image SImage) bool {
	if image.VirtualEnvType == "Ironic" {
		return true
	}

	if len(image.SupportDiskIntensive) > 0 {
		return true
	}

	if len(image.SupportKVMFPGAType) > 0 || len(image.SupportKVMAscend310) > 0 {
		return true
	}

	if len(image.SupportKVMGPUType) > 0 {
		return true
	}

	if len(image.SupportKVMNVMEHIGHIO) > 0 {
		return true
	}

	if len(image.SupportGPUT4) > 0 {
		return true
	}

	if len(image.SupportXENGPUType) > 0 {
		return true
	}

	if len(image.SupportHighPerformance) > 0 {
		return true
	}

	return false
}

// https://support.huaweicloud.com/api-ims/zh-cn_topic_0060804959.html
func (self *SRegion) GetImages(status string, imagetype TImageOwnerType, name string, envType string) ([]SImage, error) {
	queries := map[string]string{}
	if len(status) > 0 {
		queries["status"] = status
	}

	if len(imagetype) > 0 {
		queries["__imagetype"] = string(imagetype)
		if imagetype == ImageOwnerPublic {
			queries["protected"] = "True"
		}
	}
	if len(envType) > 0 {
		queries["virtual_env_type"] = envType
	}

	if len(name) > 0 {
		queries["name"] = name
	}

	images := make([]SImage, 0)
	err := doListAllWithMarker(self.ecsClient.Images.List, queries, &images)

	// 排除掉需要特定镜像才能创建的实例类型
	// https://support.huaweicloud.com/eu-west-0-api-ims/zh-cn_topic_0031617666.html#ZH-CN_TOPIC_0031617666__table48545918250
	// https://support.huaweicloud.com/productdesc-ecs/zh-cn_topic_0088142947.html
	filtedImages := make([]SImage, 0)
	for i := range images {
		if !excludeImage(images[i]) {
			filtedImages = append(filtedImages, images[i])
		}
	}

	return filtedImages, err
}

func (self *SRegion) DeleteImage(imageId string) error {
	return DoDelete(self.ecsClient.OpenStackImages.Delete, imageId, nil, nil)
}

func (self *SRegion) GetImageByName(name string) (*SImage, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("image name should not be empty")
	}

	images, err := self.GetImages("", TImageOwnerType(""), name, "")
	if err != nil {
		return nil, err
	}
	if len(images) == 0 {
		return nil, cloudprovider.ErrNotFound
	}

	log.Debugf("%d image found match name %s", len(images), name)
	return &images[0], nil
}

/*
https://support.huaweicloud.com/api-ims/zh-cn_topic_0020092109.html

	os version 取值范围： https://support.huaweicloud.com/api-ims/zh-cn_topic_0031617666.html
	用于创建私有镜像的源云服务器系统盘大小大于等于40GB且不超过1024GB。
	目前支持vhd，zvhd、raw，qcow2
	todo: 考虑使用镜像快速导入。 https://support.huaweicloud.com/api-ims/zh-cn_topic_0133188204.html
	使用OBS文件创建镜像

	* openstack原生接口支持的格式：https://support.huaweicloud.com/api-ims/zh-cn_topic_0031615566.html
*/
func (self *SRegion) ImportImageJob(name string, osDist string, osVersion string, osArch string, bucket string, key string, minDiskGB int64) (string, error) {
	os_version, err := stdVersion(osDist, osVersion, osArch)
	log.Debugf("%s %s %s: %s.min_disk %d GB", osDist, osVersion, osArch, os_version, minDiskGB)
	if err != nil {
		log.Debugln(err)
	}

	params := jsonutils.NewDict()
	params.Add(jsonutils.NewString(name), "name")
	image_url := fmt.Sprintf("%s:%s", bucket, key)
	params.Add(jsonutils.NewString(image_url), "image_url")
	if len(os_version) > 0 {
		params.Add(jsonutils.NewString(os_version), "os_version")
	}
	params.Add(jsonutils.NewBool(true), "is_config_init")
	params.Add(jsonutils.NewBool(true), "is_config")
	params.Add(jsonutils.NewInt(minDiskGB), "min_disk")
	// [CHANGED] HCS 8.6.0 制作镜像接口新增以下参数：
	// virtual_env_type: 明确指定镜像环境类型为 FusionCompute（弹性云服务器镜像）
	params.Add(jsonutils.NewString("FusionCompute"), "virtual_env_type")
	// architecture: HCS 8.6.0 新增字段，指定 CPU 架构（x86_64 或 aarch64）
	// [ORIGIN] 原始代码未传入 virtual_env_type 和 architecture 参数
	hcsArch := archToHCS(osArch)
	if len(hcsArch) > 0 {
		params.Add(jsonutils.NewString(hcsArch), "architecture")
	}

	ret, err := self.ecsClient.Images.PerformAction2("action", "", params, "")
	if err != nil {
		return "", err
	}

	return ret.GetString("job_id")
}

// archToHCS 将内部架构标识转换为 HCS 8.6.0 IMS 接口所需的架构字段值（x86_64 或 aarch64）
func archToHCS(osArch string) string {
	switch osArch {
	case "64", apis.OS_ARCH_X86_64, "x86":
		return "x86_64"
	case apis.OS_ARCH_AARCH64, apis.OS_ARCH_ARM, "arm64":
		return "aarch64"
	}
	return ""
}

func formatVersion(osDist string, osVersion string) (string, error) {
	err := fmt.Errorf("unsupport version %s.reference: https://support.huaweicloud.com/api-ims/zh-cn_topic_0031617666.html", osVersion)
	dist := strings.ToLower(osDist)
	if dist == "ubuntu" || dist == "redhat" || dist == "centos" || dist == "oracle" || dist == "euleros" {
		parts := strings.Split(osVersion, ".")
		if len(parts) < 2 {
			return "", err
		}

		return parts[0] + "." + parts[1], nil
	}

	if dist == "debian" {
		parts := strings.Split(osVersion, ".")
		if len(parts) < 3 {
			return "", err
		}

		return parts[0] + "." + parts[1] + "." + parts[2], nil
	}

	if dist == "fedora" || dist == "windows" || dist == "suse" {
		parts := strings.Split(osVersion, ".")
		if len(parts) < 1 {
			return "", err
		}

		return parts[0], nil
	}

	if dist == "opensuse" {
		parts := strings.Split(osVersion, ".")
		if len(parts) == 0 {
			return "", err
		}

		if len(parts) == 1 {
			return parts[0], nil
		}

		if len(parts) >= 2 {
			return parts[0] + "." + parts[1], nil
		}
	}

	return "", err
}

// https://support.huaweicloud.com/api-ims/zh-cn_topic_0031617666.html
func stdVersion(osDist string, osVersion string, osArch string) (string, error) {
	// 架构
	arch := ""
	switch osArch {
	case "64", apis.OS_ARCH_X86_64, apis.OS_ARCH_AARCH64, apis.OS_ARCH_ARM:
		arch = "64bit"
	case "32", apis.OS_ARCH_X86_32, apis.OS_ARCH_AARCH32:
		arch = "32bit"
	default:
		return "", fmt.Errorf("unsupported arch %s.reference: https://support.huaweicloud.com/api-ims/zh-cn_topic_0031617666.html", osArch)
	}

	_dist := strings.Split(strings.TrimSpace(osDist), " ")[0]
	_dist = strings.ToLower(_dist)
	// 版本
	ver, err := formatVersion(_dist, osVersion)
	if err != nil {
		return "", err
	}

	//  操作系统
	dist := ""

	switch _dist {
	case "ubuntu":
		return fmt.Sprintf("Ubuntu %s server %s", ver, arch), nil
	case "redhat":
		dist = "Redhat Linux Enterprise"
	case "centos":
		dist = "CentOS"
	case "fedora":
		dist = "Fedora"
	case "debian":
		dist = "Debian GNU/Linux"
	case "windows":
		dist = "Windows Server"
	case "oracle":
		dist = "Oracle Linux Server release"
	case "suse":
		dist = "SUSE Linux Enterprise Server"
	case "opensuse":
		dist = "OpenSUSE"
	case "euleros":
		dist = "EulerOS"
	default:
		return "", fmt.Errorf("unsupported os %s. reference: https://support.huaweicloud.com/api-ims/zh-cn_topic_0031617666.html", dist)
	}

	return fmt.Sprintf("%s %s %s", dist, ver, arch), nil
}
