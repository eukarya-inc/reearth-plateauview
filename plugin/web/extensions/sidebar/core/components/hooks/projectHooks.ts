import { UserDataItem } from "@web/extensions/sidebar/modals/datacatalog/types";
import { Project, ReearthApi } from "@web/extensions/sidebar/types";
import { generateID, mergeProperty, postMsg } from "@web/extensions/sidebar/utils";
import { useCallback, useState } from "react";

import { Data, DataCatalogItem, Template } from "../../types";
import { FieldComponent } from "../content/common/DatasetCard/Field/Fields/types";

import { mergeOverrides } from ".";

export const defaultProject: Project = {
  sceneOverrides: {
    default: {
      camera: {
        lat: 35.65075152248653,
        lng: 139.7617718208305,
        altitude: 2219.7187259974316,
        heading: 6.132702058010316,
        pitch: -0.5672459184621266,
        roll: 0.00019776785897196447,
        fov: 1.0471975511965976,
        height: 2219.7187259974316,
      },
      sceneMode: "3d",
      depthTestAgainstTerrain: false,
    },
    terrain: {
      terrain: true,
      terrainType: "cesiumion",
      terrainCesiumIonAccessToken:
        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiI3NGI5ZDM0Mi1jZDIzLTRmMzEtOTkwYi0zZTk4Yzk3ODZlNzQiLCJpZCI6NDA2NDYsImlhdCI6MTYwODk4MzAwOH0.3rco62ErML11TMSEflsMqeUTCDbIH6o4n4l5sssuedE",
      terrainCesiumIonAsset: "286503",
    },
    tiles: [
      {
        id: "tokyo",
        tile_url: "https://cyberjapandata.gsi.go.jp/xyz/seamlessphoto/{z}/{x}/{y}.jpg",
        tile_type: "url",
      },
    ],
    atmosphere: { shadows: true },
  },
  datasets: [],
  userStory: undefined,
};

export default ({
  data,
  fieldTemplates,
  backendURL,
  backendProjectName,
  backendAccessToken,
  inEditor,
  processedCatalog,
  setLoading,
  handleBackendFetch,
}: {
  data?: Data[];
  fieldTemplates?: Template[];
  backendURL?: string;
  backendProjectName?: string;
  backendAccessToken?: string;
  inEditor?: boolean;
  processedCatalog: DataCatalogItem[];
  setLoading?: React.Dispatch<React.SetStateAction<boolean>>;
  handleBackendFetch: () => Promise<void>;
}) => {
  const [project, updateProject] = useState<Project>(defaultProject);
  const [cleanseOverride, setCleanseOverride] = useState<string>();

  const processOverrides = useCallback(
    (dataset: DataCatalogItem, activeIDs?: string[]) => {
      if (!activeIDs) return undefined;
      let overrides = undefined;

      const inactivefields = dataset?.components?.filter(c => !activeIDs.find(id => id === c.id));
      const inactiveTemplates = inactivefields?.filter(af => af.type === "template");
      if (inactiveTemplates) {
        const inactiveTemplateFields = inactiveTemplates
          .map(
            at =>
              fieldTemplates?.find(ft => at.type === "template" && at.templateID === ft.id)
                ?.components,
          )
          .reduce((acc, field) => [...(acc ?? []), ...(field ?? [])], []);

        if (inactiveTemplateFields) {
          inactivefields?.push(...inactiveTemplateFields);
        }
      }

      const activeFields: FieldComponent[] | undefined = dataset?.components
        ?.filter(c => !!activeIDs.find(id => id === c.id))
        .map(c2 => {
          if (c2.type === "template") {
            return [
              c2,
              ...(c2.components?.length
                ? c2.components
                : fieldTemplates?.find(ft => ft.id === c2.templateID)?.components ?? []),
            ];
          }
          return c2;
        })
        .reduce((acc: FieldComponent[], field: FieldComponent | FieldComponent[] | undefined) => {
          if (!field) return acc;
          return [...acc, ...(Array.isArray(field) ? field : [field])];
        }, []);

      const cleanseOverrides = mergeOverrides("cleanse", inactivefields, cleanseOverride);
      overrides = mergeOverrides("update", activeFields, cleanseOverrides);

      setCleanseOverride(undefined);

      return overrides;
    },
    [fieldTemplates, cleanseOverride],
  );

  const handleProjectSceneUpdate = useCallback(
    (updatedProperties: Partial<ReearthApi>) => {
      updateProject(({ sceneOverrides, datasets }) => {
        const updatedProject: Project = {
          sceneOverrides: [sceneOverrides, updatedProperties].reduce((p, v) => mergeProperty(p, v)),
          datasets,
        };
        postMsg({ action: "updateProject", payload: updatedProject });
        return updatedProject;
      });
    },
    [updateProject],
  );

  const handleProjectDatasetAdd = useCallback(
    (dataset: DataCatalogItem | UserDataItem) => {
      const datasetToAdd = { ...dataset } as DataCatalogItem;

      updateProject(project => {
        if (!dataset.components?.length) {
          const defaultTemplate = fieldTemplates?.find(ft =>
            dataset.type2
              ? ft.name.includes(dataset.type2)
              : dataset.type
              ? ft.name.includes(dataset.type)
              : undefined,
          );
          if (defaultTemplate && !datasetToAdd.components) {
            datasetToAdd.components = [
              {
                id: generateID(),
                type: "template",
                templateID: defaultTemplate.id,
                components: defaultTemplate.components,
              },
            ];
          }
        }

        const updatedProject: Project = {
          ...project,
          datasets: [...project.datasets, datasetToAdd],
        };

        postMsg({ action: "updateProject", payload: updatedProject });

        return updatedProject;
      });

      const activeIDs = (
        !datasetToAdd.components?.find(c => c.type === "switchGroup") || !datasetToAdd.fieldGroups
          ? datasetToAdd.components
          : datasetToAdd.components.filter(
              c =>
                (c.group && c.group === datasetToAdd.fieldGroups?.[0].id) ||
                c.type === "switchGroup",
            )
      )
        ?.filter(c => !(!datasetToAdd.config?.data && c.type === "switchDataset"))
        ?.map(c => c.id);

      const overrides = processOverrides(datasetToAdd, activeIDs);

      postMsg({
        action: "addDatasetToScene",
        payload: {
          dataset: datasetToAdd,
          overrides,
        },
      });
    },
    [fieldTemplates, processOverrides],
  );

  const handleProjectDatasetRemove = useCallback((dataID: string) => {
    updateProject(({ sceneOverrides, datasets }) => {
      const updatedProject = {
        sceneOverrides,
        datasets: datasets.filter(d => d.dataID !== dataID),
      };
      postMsg({ action: "updateProject", payload: updatedProject });
      return updatedProject;
    });
    postMsg({ action: "removeDatasetFromScene", payload: dataID });
  }, []);

  const handleProjectDatasetRemoveAll = useCallback(() => {
    updateProject(({ sceneOverrides }) => {
      const updatedProject = {
        sceneOverrides,
        datasets: [],
      };
      postMsg({ action: "updateProject", payload: updatedProject });
      return updatedProject;
    });
    postMsg({ action: "removeAllDatasetsFromScene" });
  }, []);

  const handleDatasetUpdate = useCallback(
    (updatedDataset: DataCatalogItem, cleanseOverride?: any) => {
      updateProject(project => {
        const updatedDatasets = [...project.datasets];
        const datasetIndex = updatedDatasets.findIndex(d2 => d2.dataID === updatedDataset.dataID);
        if (datasetIndex >= 0) {
          if (updatedDatasets[datasetIndex].visible !== updatedDataset.visible) {
            postMsg({
              action: "updateDatasetVisibility",
              payload: { dataID: updatedDataset.dataID, hide: !updatedDataset.visible },
            });
          }
          if (cleanseOverride) {
            setCleanseOverride(cleanseOverride);
          }
          updatedDatasets[datasetIndex] = updatedDataset;
        }
        const updatedProject = {
          ...project,
          datasets: updatedDatasets,
        };
        postMsg({ action: "updateProject", payload: updatedProject });
        return updatedProject;
      });
    },
    [],
  );

  const handleDataRequest = useCallback(
    async (dataset?: DataCatalogItem) => {
      if (!backendURL || !backendAccessToken || !dataset) return;
      const datasetToSave = convertToData(dataset);

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
    [data, backendAccessToken, backendURL, backendProjectName, handleBackendFetch],
  );

  const handleDatasetSave = useCallback(
    (dataID: string) => {
      (async () => {
        if (!inEditor) return;
        setLoading?.(true);
        const selectedDataset = project.datasets.find(d => d.dataID === dataID);

        await handleDataRequest(selectedDataset);
        setLoading?.(false);
      })();
    },
    [inEditor, project.datasets, setLoading, handleDataRequest],
  );

  const handleDatasetPublish = useCallback(
    (dataID: string, publish: boolean) => {
      (async () => {
        if (!inEditor || !processedCatalog) return;
        const dataset = processedCatalog.find(item => item.dataID === dataID);

        if (!dataset) return;

        dataset.public = publish;

        updateProject(project => {
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

        await handleDataRequest(dataset);
      })();
    },
    [processedCatalog, inEditor, handleDataRequest],
  );

  const handleOverride = useCallback(
    (dataID: string, activeIDs?: string[]) => {
      const dataset = project.datasets.find(d => d.dataID === dataID);
      if (dataset) {
        const overrides = processOverrides(dataset, activeIDs);

        postMsg({
          action: "updateDatasetInScene",
          payload: { dataID, overrides },
        });
      }
    },
    [project.datasets, processOverrides],
  );

  return {
    project,
    updateProject,
    handleProjectSceneUpdate,
    handleProjectDatasetAdd,
    handleProjectDatasetRemove,
    handleProjectDatasetRemoveAll,
    handleDatasetUpdate,
    handleDatasetSave,
    handleDatasetPublish,
    handleOverride,
  };
};

const convertToData = (item: DataCatalogItem): Data => {
  return {
    dataID: item.dataID,
    public: item.public,
    visible: item.visible ?? true,
    template: item.template,
    components: item.components,
    fieldGroups: item.fieldGroups,
  };
};
