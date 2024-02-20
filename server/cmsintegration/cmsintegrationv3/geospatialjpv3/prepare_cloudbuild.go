package geospatialjpv3

import (
	"context"
	"path"

	"github.com/reearth/reearthx/log"
	"github.com/reearth/reearthx/rerror"
	"google.golang.org/api/cloudbuild/v1"
)

type prepareOnCloudBuildConfig struct {
	City                  string
	Project               string
	CMSURL                string
	CMSToken              string
	CloudBuildImage       string
	CloudBuildMachineType string
	CloudBuildProject     string
	CloudBuildRegion      string
}

const defaultDockerImage = "eukarya/plateauview2-sidecar-worker:latest"

func prepareOnCloudBuild(ctx context.Context, conf prepareOnCloudBuildConfig) error {
	if conf.CloudBuildImage == "" {
		conf.CloudBuildImage = defaultDockerImage
	}

	log.Debugfc(ctx, "geospatialjp webhook: prepare (cloud build): %s", ppp.Sprint(conf))

	return runCloudBuild(ctx, CloudBuildConfig{
		Image: conf.CloudBuildImage,
		Args: []string{
			"--city=" + conf.City,
			"--project=" + conf.Project,
			"--wetrun",
		},
		Env: []string{
			"REEARTH_CMS_URL=" + conf.CMSURL,
			"REEARTH_CMS_TOKEN=" + conf.CMSToken,
		},
		MachineType: conf.CloudBuildMachineType,
		Project:     conf.CloudBuildProject,
		Region:      conf.CloudBuildRegion,
	})
}

type CloudBuildConfig struct {
	Image       string
	Args        []string
	Env         []string
	MachineType string
	Region      string
	Project     string
}

func runCloudBuild(ctx context.Context, conf CloudBuildConfig) error {
	cb, err := cloudbuild.NewService(ctx)
	if err != nil {
		return rerror.ErrInternalBy(err)
	}

	machineType := ""
	if v := conf.MachineType; v != "default" {
		machineType = v
	}

	build := &cloudbuild.Build{
		Timeout:  "86400s", // 1 day
		QueueTtl: "86400s", // 1 day
		Steps: []*cloudbuild.BuildStep{
			{
				Name: conf.Image,
				Args: conf.Args,
				Env:  conf.Env,
			},
		},
		Options: &cloudbuild.BuildOptions{
			MachineType: machineType,
		},
	}

	if conf.Region != "" {
		call := cb.Projects.Locations.Builds.Create(
			path.Join("projects", conf.Project, "locations", conf.Region),
			build,
		)
		_, err = call.Do()
	} else {
		call := cb.Projects.Builds.Create(conf.Project, build)
		_, err = call.Do()
	}
	if err != nil {
		return rerror.ErrInternalBy(err)
	}
	return nil
}
