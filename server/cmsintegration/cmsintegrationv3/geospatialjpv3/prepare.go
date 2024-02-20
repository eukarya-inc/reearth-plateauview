package geospatialjpv3

import (
	"context"
)

func Prepare(ctx context.Context, itemID, projectID string, conf Config) error {
	if conf.BuildType == "cloudbuild" {
		return prepareOnCloudBuild(ctx, prepareOnCloudBuildConfig{
			City:                  itemID,
			Project:               projectID,
			CMSURL:                conf.CMSBase,
			CMSToken:              conf.CMSToken,
			CloudBuildImage:       conf.CloudBuildImage,
			CloudBuildMachineType: conf.CloudBuildMachineType,
			CloudBuildProject:     conf.CloudBuildProject,
			CloudBuildRegion:      conf.CloudBuildRegion,
		})
	} else {
		return prepareWithCloudRunJobs(ctx, itemID, projectID, conf.CloudRunJobsJobName)
	}
}
