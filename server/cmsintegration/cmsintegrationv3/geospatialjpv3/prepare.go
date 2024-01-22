package geospatialjpv3

import (
	"context"

	run "cloud.google.com/go/run/apiv2"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"

	runpb "cloud.google.com/go/run/apiv2/runpb"
	"github.com/reearth/reearthx/log"
)

type PrepareConfig struct {
	gcpProjectID string
	gcpLocation  string
	wetRun       bool
}

func (c *PrepareConfig) RequestPreparing(ctx context.Context, w *cmswebhook.Payload) error {
	log.Debugfc(ctx, "geospatialjp webhook: RequestPreparing")
	client, err := run.NewJobsClient(ctx)
	if err != nil {
		log.Debugfc(ctx, "geospatialjp webhook: failed to create run client: %v", err)
		return err
	}
	defer client.Close()

	const Job = "plateauview-api-worker"

	//TODO: wet run implementation

	overrides := runpb.RunJobRequest_Overrides{
		ContainerOverrides: []*runpb.RunJobRequest_Overrides_ContainerOverride{
			{Args: []string{"prepare-gspatialjp", "--city=" + w.ItemData.Item.ID, "--project=" + w.ProjectID()}}}}

	req := &runpb.RunJobRequest{
		Name:      "projects/" + c.gcpProjectID + "/locations/" + c.gcpLocation + "/jobs/" + Job,
		Overrides: &overrides,
	}

	log.Debugfc(ctx, "geospatialjp webhook: run job: %v", req)
	op, err := client.RunJob(ctx, req)
	if err != nil {
		log.Debugfc(ctx, "geospatialjp webhook: failed to run job: %v", err)
		return err
	}

	//TODO: 実際はwaitする必要無し
	log.Debugfc(ctx, "geospatialjp webhook: waiting for job to complete")
	if _, err := op.Wait(ctx); err != nil {
		log.Debugfc(ctx, "geospatialjp webhook: failed to wait for job: %v", err)
		return err
	}

	return nil
}
