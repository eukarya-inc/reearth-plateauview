import { Project, ReearthApi } from "@web/extensions/sidebar/types";
import { generateID, mergeProperty, postMsg } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";

import { getDataCatalog, RawDataCatalogItem } from "../../modals/datacatalog/api/api";
import { UserDataItem } from "../../modals/datacatalog/types";
import { Data, DataCatalogItem, Template } from "../types";

import { Story as FieldStory, StoryItem } from "./content/common/DatasetCard/Field/Fields/types";
import { Pages } from "./Header";

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

export default () => {
  const [projectID, setProjectID] = useState<string>();
  const [inEditor, setInEditor] = useState(true);

  const [catalogURL, setCatalogURL] = useState<string>();
  const [reearthURL, setReearthURL] = useState<string>();
  const [backendURL, setBackendURL] = useState<string>();
  const [backendProjectName, setBackendProjectName] = useState<string>();
  const [backendAccessToken, setBackendAccessToken] = useState<string>();

  const [data, setData] = useState<Data[]>();
  const [project, updateProject] = useState<Project>(defaultProject);
  const [selectedDatasets, setSelectedDatasets] = useState<DataCatalogItem[]>([]);
  const [savingTemplate, setSaveTemplate] = useState<boolean>(false);
  const [savingDataset, setSaveDataset] = useState<boolean>(false);

  const handleBackendFetch = useCallback(async () => {
    if (!backendURL) return;
    const res = await fetch(`${backendURL}/sidebar/${backendProjectName}`);
    if (res.status !== 200) return;
    const resData = await res.json();

    if (resData.templates) {
      setFieldTemplates(resData.templates.filter((t: Template) => t.type === "field"));
      setInfoboxTemplates(resData.templates.filter((t: Template) => t.type === "infobox"));
    }
    setData(resData.data);
  }, [backendURL, backendProjectName]);

  // ****************************************
  // Init
  const [catalogData, setCatalog] = useState<RawDataCatalogItem[]>([]);

  useEffect(() => {
    postMsg({ action: "init" }); // Needed to trigger sending initialization data to sidebar
  }, []);

  useEffect(() => {
    if (catalogURL) {
      getDataCatalog(catalogURL).then(res => {
        setCatalog(res);
      });
    }
  }, [catalogURL]);

  useEffect(() => {
    if (backendURL) {
      handleBackendFetch();
    }
  }, [backendURL]); // eslint-disable-line react-hooks/exhaustive-deps

  const processedCatalog = useMemo(() => {
    const c = handleDataCatalogProcessing(catalogData, data);
    return inEditor ? c : c.filter(c => !!c.public);
  }, [catalogData, inEditor, data]);

  useEffect(() => {
    postMsg({ action: "updateCatalog", payload: processedCatalog });
  }, [processedCatalog]);

  // ****************************************

  // ****************************************
  // Project

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
      updateProject(project => {
        let dataToAdd = data?.find(d => d.dataID === dataset.dataID);

        if (!dataToAdd) {
          dataToAdd = convertToData(dataset as DataCatalogItem);
        }

        const updatedProject: Project = {
          ...project,
          datasets: [...project.datasets, dataToAdd],
        };

        postMsg({ action: "updateProject", payload: updatedProject });
        setSelectedDatasets(sds => [...sds, dataset as DataCatalogItem]);

        return updatedProject;
      });

      // const options = data?.find(d => d.id === dataset.id)?.components;
      postMsg({ action: "addDatasetToScene", payload: { dataset } });
    },
    [data],
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
    setSelectedDatasets(sds => sds.filter(sd => sd.dataID !== dataID));
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
    setSelectedDatasets([]);
    postMsg({ action: "removeAllDatasetsFromScene" });
  }, []);

  const handleDatasetUpdate = useCallback((updatedDataset: DataCatalogItem) => {
    setSelectedDatasets(selectedDatasets => {
      const updatedDatasets = [...selectedDatasets];
      const datasetIndex = updatedDatasets.findIndex(d2 => d2.dataID === updatedDataset.dataID);
      if (datasetIndex >= 0) {
        if (updatedDatasets[datasetIndex].visible !== updatedDataset.visible) {
          postMsg({
            action: "updateDatasetVisibility",
            payload: { dataID: updatedDataset.dataID, hide: !updatedDataset.visible },
          });
        }
        updatedDatasets[datasetIndex] = updatedDataset;
      }
      updateProject(project => {
        return {
          ...project,
          datasets: updatedDatasets.map(ud => convertToData(ud)),
        };
      });
      return updatedDatasets;
    });
  }, []);

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
      setSaveDataset(true);
      (async () => {
        if (!inEditor) return;
        const selectedDataset = selectedDatasets.find(d => d.dataID === dataID);

        await handleDataRequest(selectedDataset);
      })();
      setSaveDataset(false);
    },
    [setSaveDataset, inEditor, selectedDatasets, handleDataRequest],
  );

  const handleDatasetPublish = useCallback(
    (dataID: string, publish: boolean) => {
      (async () => {
        if (!inEditor || !processedCatalog) return;
        const dataset = processedCatalog.find(item => item.dataID === dataID);

        if (!dataset) return;

        dataset.public = publish;

        await handleDataRequest(dataset);
      })();
    },
    [processedCatalog, inEditor, handleDataRequest],
  );

  // ****************************************

  // ****************************************
  // Templates
  const [fieldTemplates, setFieldTemplates] = useState<Template[]>([]);

  const handleTemplateAdd = useCallback(async () => {
    if (!backendURL || !backendProjectName || !backendAccessToken) return;
    const res = await fetch(`${backendURL}/sidebar/${backendProjectName}/templates`, {
      headers: {
        authorization: `Bearer ${backendAccessToken}`,
      },
      method: "POST",
      body: JSON.stringify({ type: "field", name: "新しいテンプレート" }),
    });
    if (res.status !== 200) return;
    const newTemplate = await res.json();
    setFieldTemplates(t => [...t, newTemplate]);
    return newTemplate as Template;
  }, [backendURL, backendProjectName, backendAccessToken]);

  const handleTemplateSave = useCallback(
    async (template: Template) => {
      if (!backendURL || !backendProjectName || !backendAccessToken) return;
      setSaveTemplate(true);
      const res = await fetch(
        `${backendURL}/sidebar/${backendProjectName}/templates/${template.id}`,
        {
          headers: {
            authorization: `Bearer ${backendAccessToken}`,
          },
          method: "PATCH",
          body: JSON.stringify(template),
        },
      );
      if (res.status !== 200) return;
      const updatedTemplate = await res.json();
      setFieldTemplates(t => {
        return t.map(t2 => {
          if (t2.id === updatedTemplate.id) {
            return updatedTemplate;
          }
          return t2;
        });
      });
      setSaveTemplate(false);
    },
    [backendURL, backendProjectName, backendAccessToken],
  );

  const handleTemplateRemove = useCallback(
    async (id: string) => {
      if (!backendURL || !backendProjectName || !backendAccessToken) return;
      const res = await fetch(`${backendURL}/sidebar/${backendProjectName}/templates/${id}`, {
        headers: {
          authorization: `Bearer ${backendAccessToken}`,
        },
        method: "DELETE",
      });
      if (res.status !== 200) return;
      setFieldTemplates(t => t.filter(t2 => t2.id !== id));
    },
    [backendURL, backendProjectName, backendAccessToken],
  );

  // ****************************************
  // Story
  const handleStorySaveData = useCallback((story: StoryItem & { dataID?: string }) => {
    if (story.id && story.dataID) {
      // save database story
      setSelectedDatasets(sd => {
        const tarStory = (
          sd
            .find(s => s.dataID === story.dataID)
            ?.components?.find(c => c.type === "story") as FieldStory
        )?.stories?.find((st: StoryItem) => st.id === story.id);
        if (tarStory) {
          tarStory.scenes = story.scenes;
        }
        return sd;
      });
    }

    // save user story
    updateProject(project => {
      const updatedProject: Project = {
        ...project,
        userStory: {
          scenes: story.scenes,
        },
      };
      postMsg({ action: "updateProject", payload: updatedProject });
      return updatedProject;
    });
  }, []);

  const handleInitUserStory = useCallback((story: StoryItem) => {
    postMsg({ action: "storyPlay", payload: story });
  }, []);

  // ****************************************

  // Infobox
  const [infoboxTemplates, setInfoboxTemplates] = useState<Template[]>([]);

  const handleInfoboxTemplateAdd = useCallback(
    async (template: Omit<Template, "id">) => {
      if (!backendURL || !backendProjectName || !backendAccessToken) return;
      const res = await fetch(`${backendURL}/sidebar/${backendProjectName}/templates`, {
        headers: {
          authorization: `Bearer ${backendAccessToken}`,
        },
        method: "POST",
        body: JSON.stringify(template),
      });
      if (res.status !== 200) return;
      const newTemplate = await res.json();
      setInfoboxTemplates(t => [...t, newTemplate]);
      return newTemplate as Template;
    },
    [backendURL, backendProjectName, backendAccessToken],
  );

  const handleInfoboxTemplateSave = useCallback(
    async (template: Template) => {
      if (!backendURL || backendProjectName || !backendAccessToken) return;
      const res = await fetch(
        `${backendURL}/sidebar/${backendProjectName}/templates/${template.id}`,
        {
          headers: {
            authorization: `Bearer ${backendAccessToken}`,
          },
          method: "PATCH",
          body: JSON.stringify(template),
        },
      );
      if (res.status !== 200) return;
      const updatedTemplate = await res.json();
      setInfoboxTemplates(t => {
        return t.map(t2 => {
          if (t2.id === updatedTemplate.id) {
            return updatedTemplate;
          }
          return t2;
        });
      });
      postMsg({
        action: "infoboxFieldsSaved",
      });
    },
    [backendURL, backendProjectName, backendAccessToken],
  );

  const handleInfoboxFieldsFetch = useCallback(
    (dataID: string) => {
      const name = catalogData?.find(d => d.id === dataID)?.type ?? "";
      const fields = infoboxTemplates.find(ft => ft.type === "infobox" && ft.name === name) ?? {
        id: "",
        type: "infobox",
        name,
        fields: [],
      };
      postMsg({
        action: "infoboxFieldsFetch",
        payload: fields,
      });
    },
    [catalogData, infoboxTemplates],
  );
  const handleInfoboxFieldsFetchRef = useRef<any>();
  handleInfoboxFieldsFetchRef.current = handleInfoboxFieldsFetch;

  const handleInfoboxFieldsSave = useCallback(
    async (template: Template) => {
      if (template.id) {
        handleInfoboxTemplateSave(template);
      } else {
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        const { id, ...templateData } = template;
        handleInfoboxTemplateAdd(templateData);
      }
    },
    [handleInfoboxTemplateAdd, handleInfoboxTemplateSave],
  );
  const handleInfoboxFieldsSaveRef = useRef<any>();
  handleInfoboxFieldsSaveRef.current = handleInfoboxFieldsSave;

  // ****************************************

  useEffect(() => {
    const eventListenerCallback = (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      if (e.data.action === "msgFromModal") {
        if (e.data.payload.dataset) {
          handleProjectDatasetAdd(e.data.payload.dataset);
        }
      } else if (e.data.action === "init" && e.data.payload) {
        setProjectID(e.data.payload.projectID);
        setInEditor(e.data.payload.inEditor);
        setCatalogURL(e.data.payload.catalogURL);
        setReearthURL(`${e.data.payload.reearthURL}`);
        setBackendURL(e.data.payload.backendURL);
        setBackendProjectName(e.data.payload.backendProjectName);
        setBackendAccessToken(e.data.payload.backendAccessToken);
        if (e.data.payload.draftProject) {
          updateProject(e.data.payload.draftProject);
        }
      } else if (e.data.action === "updateDataset") {
        handleDatasetPublish(e.data.payload.dataID, e.data.payload.publish);
      } else if (e.data.action === "triggerCatalogOpen") {
        handleModalOpen();
      } else if (e.data.action === "triggerHelpOpen") {
        handlePageChange("help");
      } else if (e.data.action === "storyShare") {
        setCurrentPage("share");
      } else if (e.data.action === "storySaveData") {
        handleStorySaveData(e.data.payload);
      } else if (e.data.action === "infoboxFieldsFetch") {
        handleInfoboxFieldsFetchRef.current(e.data.payload);
      } else if (e.data.action === "infoboxFieldsSave") {
        handleInfoboxFieldsSaveRef.current(e.data.payload);
      }
    };
    addEventListener("message", eventListenerCallback);
    return () => {
      removeEventListener("message", eventListenerCallback);
    };
  }, [handleDatasetPublish]); // eslint-disable-line react-hooks/exhaustive-deps

  const fetchedSharedProject = useRef(false);

  useEffect(() => {
    if (!backendURL || !backendProjectName || fetchedSharedProject.current) return;
    if (projectID && processedCatalog.length) {
      (async () => {
        const res = await fetch(`${backendURL}/share/${backendProjectName}/${projectID}`);
        if (res.status !== 200) return;
        const data = await res.json();
        if (data) {
          updateProject(data);
          postMsg({ action: "updateProject", payload: data });
          (data.datasets as Data[]).forEach(d => {
            const dataset = processedCatalog.find(item => item.dataID === d.dataID);
            if (dataset) {
              setSelectedDatasets(sds => [...sds, dataset]);
              postMsg({ action: "addDatasetToScene", payload: { dataset } });
            }
          });
          if (data.userStory) {
            handleInitUserStory(data.userStory);
          }
        }
        fetchedSharedProject.current = true;
      })();
    }
  }, [projectID, backendURL, backendProjectName, processedCatalog, handleInitUserStory]);

  const [currentPage, setCurrentPage] = useState<Pages>("data");

  const handlePageChange = useCallback((p: Pages) => {
    setCurrentPage(p);
  }, []);

  // ThreeDTilesSearch
  const handleThreeDTilesSearch = useCallback(
    (dataID: string) => {
      const plateauItem = catalogData.find(pd => pd.id === dataID);
      const searchIndex = plateauItem?.["search_index"];

      postMsg({
        action: "buildingSearchOpen",
        payload: {
          title: plateauItem?.["name"] ?? "",
          dataID,
          searchIndex,
        },
      });
    },
    [catalogData],
  );

  const handleModalOpen = useCallback(() => {
    postMsg({
      action: "catalogModalOpen",
    });
  }, []);

  return {
    catalog: processedCatalog,
    project,
    selectedDatasets,
    inEditor,
    reearthURL,
    backendURL,
    backendProjectName,
    templates: fieldTemplates,
    currentPage,
    savingTemplate,
    savingDataset,
    handlePageChange,
    handleTemplateAdd,
    handleTemplateSave,
    handleTemplateRemove,
    handleDatasetSave,
    handleDatasetUpdate,
    handleProjectDatasetAdd,
    handleProjectDatasetRemove,
    handleProjectDatasetRemoveAll,
    handleProjectSceneUpdate,
    handleModalOpen,
    handleThreeDTilesSearch,
  };
};

addEventListener("message", e => {
  if (e.source !== parent) return;
  if (e.data.type) {
    if (e.data.type === "extended") {
      updateExtended(e.data.payload);
    }
  }
});

function updateExtended(e: { vertically: boolean }) {
  const html = document.querySelector("html");
  const body = document.querySelector("body");
  const root = document.getElementById("root");

  if (e?.vertically) {
    html?.classList.add("extended");
    body?.classList.add("extended");
    root?.classList.add("extended");
  } else {
    html?.classList.remove("extended");
    body?.classList.remove("extended");
    root?.classList.remove("extended");
  }
}

const newItem = (ri: RawDataCatalogItem): DataCatalogItem => {
  return {
    ...ri,
    dataID: ri.id,
    public: false,
    visible: true,
    fieldGroups: [{ id: generateID(), name: "グループ1" }],
  };
};

const handleDataCatalogProcessing = (
  catalog: (DataCatalogItem | RawDataCatalogItem)[],
  savedData?: Data[],
): DataCatalogItem[] =>
  catalog.map(item => {
    if (!savedData) return newItem(item);

    const savedData2 = savedData.find(d => d.dataID === ("dataID" in item ? item.dataID : item.id));
    if (savedData2) {
      return {
        ...item,
        ...savedData2,
      };
    } else {
      return newItem(item);
    }
  });

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
