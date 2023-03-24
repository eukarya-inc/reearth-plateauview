import { postMsg, generateID } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useMemo, useState } from "react";

import { Tab } from "..";
import {
  DataCatalogItem,
  getDataCatalog,
  RawDataCatalogItem,
} from "../../../../modals/datacatalog/api/api";
import { BuildingSearch, Data, Template } from "../../../types";
import { handleDataCatalogProcessing, updateExtended } from "../../utils";

import useProjectHooks from "./projectHooks";

export default () => {
  const [inEditor, setInEditor] = useState(true);
  const [selected, setSelected] = useState<Tab | undefined>();

  const [catalogURL, setCatalogURL] = useState<string>();
  const [catalogProjectName, setCatalogProjectName] = useState<string>();
  const [reearthURL, setReearthURL] = useState<string>();
  const [backendURL, setBackendURL] = useState<string>();
  const [backendProjectName, setBackendProjectName] = useState<string>();
  const [buildingSearch, setBuildingSearch] = useState<BuildingSearch>([]);
  // const [expandedFolders, setExpandedFolders] = useState<{ id?: string; name?: string }[]>([]);

  const [fieldTemplates, setFieldTemplates] = useState<Template[]>([]);
  const [infoboxTemplates, setInfoboxTemplates] = useState<Template[]>([]);

  const [data, setData] = useState<Data[]>();

  const [catalogData, setCatalog] = useState<RawDataCatalogItem[]>([]);
  const [searchTerm, setSearchTerm] = useState("");

  const handleSearch = useCallback(({ target: { value } }: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(value);
    postMsg({ action: "saveSearchTerm", payload: { searchTerm: value } });
  }, []);

  const processedCatalog = useMemo(() => {
    const c = handleDataCatalogProcessing(catalogData, data);
    return inEditor ? c : c.filter(c => !!c.public);
  }, [catalogData, inEditor, data]);

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
  }, [backendURL, backendProjectName, setInfoboxTemplates, setFieldTemplates]);

  const {
    project,
    updateProject,
    setProjectID,
    setCleanseOverride,
    handleProjectSceneUpdate,
    handleProjectDatasetAdd,
    handleProjectDatasetRemove,
    handleProjectDatasetRemoveAll,
    handleProjectDatasetsUpdate,
    handleStorySaveData,
    handleOverride,
  } = useProjectHooks({
    fieldTemplates,
    backendURL,
    backendProjectName,
    processedCatalog,
    buildingSearch,
  });

  const handleDatasetUpdate = useCallback(
    (updatedDataset: DataCatalogItem, cleanseOverride?: any) => {
      updateProject?.(project => {
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
            setCleanseOverride?.(cleanseOverride);
          }
          updatedDatasets[datasetIndex] = updatedDataset;
        }
        const updatedProject = {
          ...project,
          datasets: updatedDatasets,
        };
        postMsg({ action: "updateProject", payload: updatedProject });
        postMsg({ action: "msgToPopup", payload: { project: updatedProject } });
        return updatedProject;
      });
    },
    [updateProject, setCleanseOverride],
  );

  // ****************************************
  // Init

  useEffect(() => {
    postMsg({ action: "init" }); // Needed to trigger sending initialization data to sidebar
  }, []);

  useEffect(() => {
    const catalogBaseUrl = catalogURL || backendURL;
    if (catalogBaseUrl) {
      getDataCatalog(catalogBaseUrl, catalogProjectName).then(res => {
        setCatalog(res);
      });
    }
  }, [backendURL, catalogProjectName, catalogURL]);

  useEffect(() => {
    if (backendURL) {
      handleBackendFetch();
    }
  }, [backendURL]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    postMsg({ action: "updateDataCatalog", payload: processedCatalog });
  }, [processedCatalog]);

  // ****************************************

  const handleInfoboxFieldsFetch = useCallback(
    (dataID: string) => {
      let fields;
      const catalogItem = processedCatalog?.find(d => d.dataID === dataID);
      if (catalogItem) {
        const name = catalogItem?.type;
        const dataType = catalogItem?.type_en;
        fields = infoboxTemplates.find(ft => ft.type === "infobox" && ft.dataType === dataType) ?? {
          id: "",
          type: "infobox",
          name,
          dataType,
          fields: [],
        };
      }

      postMsg({
        action: "infoboxFieldsFetch",
        payload: fields,
      });
    },
    [processedCatalog, infoboxTemplates],
  );

  useEffect(() => {
    const eventListenerCallback = (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      if (e.data.action === "mobileDatasetAdd") {
        handleProjectDatasetAdd(e.data.payload);
      } else if (e.data.action === "mobileDatasetUpdate") {
        handleDatasetUpdate(e.data.payload);
      } else if (e.data.action === "mobileDatasetRemove") {
        handleProjectDatasetRemove(e.data.payload);
      } else if (e.data.action === "mobileDatasetRemoveAll") {
        handleProjectDatasetRemoveAll();
      } else if (e.data.action === "mobileProjectDatasetsUpdate") {
        handleProjectDatasetsUpdate(e.data.payload);
      } else if (e.data.action === "mobileProjectSceneUpdate") {
        handleProjectSceneUpdate(e.data.payload);
      } else if (e.data.action === "mobileBuildingSearch") {
        handleBuildingSearch(e.data.payload);
      } else if (e.data.action === "init" && e.data.payload) {
        setProjectID(e.data.payload.projectID);
        setInEditor(e.data.payload.inEditor);
        setCatalogURL(e.data.payload.catalogURL);
        setCatalogProjectName(e.data.payload.catalogProjectName);
        setReearthURL(`${e.data.payload.reearthURL}`);
        setBackendURL(e.data.payload.backendURL);
        setBackendProjectName(e.data.payload.backendProjectName);
        if (e.data.payload.searchTerm) setSearchTerm(e.data.payload.searchTerm);
        if (e.data.payload.draftProject) {
          updateProject(e.data.payload.draftProject);
        }
      } else if (e.data.action === "triggerCatalogOpen") {
        handleModalOpen();
      } else if (e.data.action === "storyShare") {
        setSelected("menu");
      } else if (e.data.action === "storySaveData") {
        handleStorySaveData(e.data.payload);
      } else if (e.data.action === "infoboxFieldsFetch") {
        handleInfoboxFieldsFetch(e.data.payload);
      } else if (e.data.action === "buildingSearchOverride") {
        handleBuildingSearchOverride(e.data.payload);
      } else if (e.data.action === "buildingSearchClose") {
        handleBuildingSearchClose(e.data.payload);
      }
    };
    addEventListener("message", eventListenerCallback);
    return () => {
      removeEventListener("message", eventListenerCallback);
    };
  }, [handleInfoboxFieldsFetch]); // eslint-disable-line react-hooks/exhaustive-deps

  // ****************************************
  // Building Search
  const handleBuildingSearch = useCallback(
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

  const handleBuildingSearchOverride = useCallback(
    ({ dataID, overrides }: { dataID: string; overrides: any }) => {
      setBuildingSearch(bs => {
        const id = generateID();
        const fieldItem = {
          dataID,
          active: true,
          field: {
            id,
            type: "search",
            updatedAt: new Date(),
            override: overrides,
          },
          cleanseField: {
            id,
            type: "search",
            updatedAt: new Date(),
          },
        };
        const target = bs.find(b => b.dataID === dataID);
        if (target) {
          target.active = true;
          target.field = fieldItem.field;
          target.cleanseField = fieldItem.cleanseField;
        } else {
          bs.push(fieldItem);
        }
        return [...bs];
      });
    },
    [],
  );

  const handleBuildingSearchClose = useCallback(({ dataID }: { dataID: string }) => {
    setBuildingSearch(bs => {
      const target = bs.find(b => b.dataID === dataID);
      if (target) {
        target.active = false;
      }
      return [...bs];
    });
  }, []);

  const handleModalOpen = useCallback(() => {
    postMsg({
      action: "catalogModalOpen",
    });
  }, []);

  return {
    selected,
    catalog: processedCatalog,
    project,
    inEditor,
    reearthURL,
    backendURL,
    backendProjectName,
    templates: fieldTemplates,
    buildingSearch,
    searchTerm,
    setSelected,
    handleSearch,
    handleDatasetUpdate,
    handleProjectDatasetAdd,
    handleProjectDatasetRemove,
    handleProjectDatasetRemoveAll,
    handleProjectDatasetsUpdate,
    handleProjectSceneUpdate,
    handleModalOpen,
    handleBuildingSearch,
    handleOverride,
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
