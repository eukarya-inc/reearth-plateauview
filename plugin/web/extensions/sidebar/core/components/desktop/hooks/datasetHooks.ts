import { Project } from "@web/extensions/sidebar/types";
import { mergeProperty, postMsg } from "@web/extensions/sidebar/utils";
import { formatDateTime } from "@web/extensions/sidebar/utils/date";
import { findLast } from "lodash";
import { useCallback } from "react";

import { Data, DataCatalogItem, Template } from "../../../types";
import { convertToData } from "../../utils";

export default ({
  data,
  templates,
  project,
  backendURL,
  backendProjectName,
  backendAccessToken,
  publishToGeospatial,
  inEditor,
  processedCatalog,
  setCleanseOverride,
  setLoading,
  updateProject,
  handleBackendFetch,
}: {
  data?: Data[];
  templates?: Template[];
  project?: Project;
  backendURL?: string;
  backendProjectName?: string;
  backendAccessToken?: string;
  publishToGeospatial?: boolean;
  inEditor?: boolean;
  processedCatalog: DataCatalogItem[];
  setCleanseOverride?: React.Dispatch<React.SetStateAction<string | undefined>>;
  setLoading?: React.Dispatch<React.SetStateAction<boolean>>;
  updateProject?: React.Dispatch<React.SetStateAction<Project>>;
  handleBackendFetch: () => Promise<void>;
}) => {
  const handleDataRequest = useCallback(
    async (dataset?: DataCatalogItem) => {
      if (!backendURL || !backendAccessToken || !dataset) return;
      const datasetToSave = convertToData(dataset, templates);

      const isNew = !data?.find(d => d.dataID === dataset.dataID);

      const fetchURL = !isNew
        ? `${backendURL}/sidebar/${backendProjectName}/data/${dataset.id}` // should be id and not dataID because id here is the CMS item's id
        : `${backendURL}/sidebar/${backendProjectName}/data`;

      const method = !isNew ? "PATCH" : "POST";

      const res = await fetch(fetchURL, {
        headers: {
          authorization: `Bearer ${backendAccessToken}`,
        },
        method,
        body: JSON.stringify(datasetToSave),
      });
      if (res.status !== 200) {
        handleBackendFetch();
        return;
      }
      const data2 = await res.json();
      console.log("DATA JUST SAVED: ", data2);
      handleBackendFetch(); // MAYBE UPDATE THIS LATER TO JUST UPDATE THE LOCAL VALUE
    },
    [data, templates, backendAccessToken, backendURL, backendProjectName, handleBackendFetch],
  );

  const handleDatasetUpdate = useCallback(
    (updatedDataset: DataCatalogItem, cleanseOverride?: any) => {
      updateProject?.(project => {
        const updatedDatasets = [...project.datasets];
        let updatedSceneOverrides = { ...project.sceneOverrides };
        const datasetIndex = updatedDatasets.findIndex(d2 => d2.dataID === updatedDataset.dataID);
        if (datasetIndex >= 0) {
          if (updatedDatasets[datasetIndex].visible !== updatedDataset.visible) {
            postMsg({
              action: "updateDatasetVisibility",
              payload: { dataID: updatedDataset.dataID, hide: !updatedDataset.visible },
            });
          }
          if (cleanseOverride) {
            setCleanseOverride?.(cleanseOverride);
          }

          const item = findLast(updatedDataset.components, item => item.type === "currentTime");
          if (item && item?.type === "currentTime") {
            const updatedProperty = {
              timeline: {
                current: formatDateTime(item.currentDate, item.currentTime),
                start: formatDateTime(item.startDate, item.startTime),
                stop: formatDateTime(item.stopDate, item.stopTime),
              },
            };
            updatedSceneOverrides = [updatedSceneOverrides, updatedProperty].reduce((p, v) =>
              mergeProperty(p, v),
            );
          }

          updatedDatasets[datasetIndex] = updatedDataset;
        }
        const updatedProject = {
          ...project,
          datasets: updatedDatasets,
          sceneOverrides: updatedSceneOverrides,
        };
        postMsg({ action: "updateProject", payload: updatedProject });
        return updatedProject;
      });
    },
    [updateProject, setCleanseOverride],
  );

  const handleDatasetSave = useCallback(
    (dataID: string) => {
      (async () => {
        if (!inEditor) return;
        setLoading?.(true);
        const selectedDataset = project?.datasets.find(d => d.dataID === dataID);

        await handleDataRequest(selectedDataset);
        setLoading?.(false);
      })();
    },
    [inEditor, project?.datasets, setLoading, handleDataRequest],
  );

  const handleDatasetPublish = useCallback(
    (dataID: string, publish: boolean) => {
      if (!inEditor || !processedCatalog) return;
      const dataset = processedCatalog.find(item => item.dataID === dataID);

      if (!dataset) return;

      dataset.public = publish;

      updateProject?.(project => {
        const updatedDatasets = [...project.datasets];
        const datasetIndex = updatedDatasets.findIndex(d2 => d2.dataID === dataID);
        if (datasetIndex >= 0) {
          updatedDatasets[datasetIndex] = dataset;
        }
        return {
          ...project,
          datasets: updatedDatasets,
        };
      });

      handleDataRequest(dataset);

      if (publish && publishToGeospatial && dataset.itemId && backendURL && backendAccessToken) {
        fetch(`${backendURL}/publish_to_geospatialjp`, {
          headers: {
            authorization: `Bearer ${backendAccessToken}`,
            "Content-Type": "application/json",
          },
          method: "POST",
          body: JSON.stringify({ id: dataset.itemId }),
        })
          .then(r => {
            if (!r.ok)
              throw `failed to publish the data on gspatial.jp: status code is ${r.statusText}`;
          })
          .catch(console.error);
      }
    },
    [
      processedCatalog,
      inEditor,
      backendAccessToken,
      backendURL,
      publishToGeospatial,
      updateProject,
      handleDataRequest,
    ],
  );

  return {
    handleDatasetUpdate,
    handleDatasetSave,
    handleDatasetPublish,
  };
};
