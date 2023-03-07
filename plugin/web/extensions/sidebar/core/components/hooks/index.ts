import useDatasetHooks from "@web/extensions/sidebar/core/components/hooks/datasetHooks";
import useProjectHooks from "@web/extensions/sidebar/core/components/hooks/projectHooks";
import useTemplateHooks from "@web/extensions/sidebar/core/components/hooks/templateHooks";
import { Project } from "@web/extensions/sidebar/types";
import { postMsg } from "@web/extensions/sidebar/utils";
import { merge } from "lodash";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";

import { getDataCatalog, RawDataCatalogItem } from "../../../modals/datacatalog/api/api";
import { Data, DataCatalogItem, Template } from "../../types";
import { Story as FieldStory, StoryItem } from "../content/common/DatasetCard/Field/Fields/types";
import { Pages } from "../Header";

import { handleDataCatalogProcessing, updateExtended } from "./utils";

export default () => {
  const [projectID, setProjectID] = useState<string>();
  const [inEditor, setInEditor] = useState(true);

  const [catalogURL, setCatalogURL] = useState<string>();
  const [catalogProjectName, setCatalogProjectName] = useState<string>();
  const [reearthURL, setReearthURL] = useState<string>();
  const [backendURL, setBackendURL] = useState<string>();
  const [backendProjectName, setBackendProjectName] = useState<string>();
  const [backendAccessToken, setBackendAccessToken] = useState<string>();

  const [data, setData] = useState<Data[]>();

  const [loading, setLoading] = useState<boolean>(false);

  const [catalogData, setCatalog] = useState<RawDataCatalogItem[]>([]);

  const processedCatalog = useMemo(() => {
    const c = handleDataCatalogProcessing(catalogData, data);
    return inEditor ? c : c.filter(c => !!c.public);
  }, [catalogData, inEditor, data]);

  const {
    fieldTemplates,
    handleInfoboxFieldsFetchRef,
    handleInfoboxFieldsSaveRef,
    setFieldTemplates,
    setInfoboxTemplates,
    handleTemplateAdd,
    handleTemplateSave,
    handleTemplateRemove,
  } = useTemplateHooks({
    backendURL,
    backendProjectName,
    backendAccessToken,
    processedCatalog,
    setLoading,
  });

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
    setCleanseOverride,
    handleProjectSceneUpdate,
    handleProjectDatasetAdd,
    handleProjectDatasetRemove,
    handleProjectDatasetRemoveAll,
    handleOverride,
  } = useProjectHooks({
    fieldTemplates,
  });

  const { handleDatasetUpdate, handleDatasetSave, handleDatasetPublish } = useDatasetHooks({
    data,
    project,
    backendURL,
    backendProjectName,
    backendAccessToken,
    inEditor,
    processedCatalog,
    setCleanseOverride,
    setLoading,
    updateProject,
    handleBackendFetch,
  });

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
    postMsg({ action: "updateCatalog", payload: processedCatalog });
  }, [processedCatalog]);

  // ****************************************

  // ****************************************
  // Story
  const handleStorySaveData = useCallback(
    (story: StoryItem & { dataID?: string }) => {
      if (story.id && story.dataID) {
        // save database story
        updateProject(project => {
          const tarStory = (
            project.datasets
              .find(d => d.dataID === story.dataID)
              ?.components?.find(c => c.type === "story") as FieldStory
          )?.stories?.find((st: StoryItem) => st.id === story.id);
          if (tarStory) {
            tarStory.scenes = story.scenes;
          }
          return project;
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
    },
    [updateProject],
  );

  const handleInitUserStory = useCallback((story: StoryItem) => {
    postMsg({ action: "storyPlay", payload: story });
  }, []);

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
        setCatalogProjectName(e.data.payload.catalogProjectName);
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
          (data.datasets as Data[]).forEach(d => {
            const dataset = processedCatalog.find(item => item.dataID === d.dataID);
            const mergedDataset: DataCatalogItem = merge(dataset, d, {});
            if (mergedDataset) {
              handleProjectDatasetAdd(mergedDataset);
            }
          });
          if (data.userStory.length > 0) {
            handleInitUserStory(data.userStory);
          }
        }
        fetchedSharedProject.current = true;
      })();
    }
  }, [
    projectID,
    backendURL,
    backendProjectName,
    processedCatalog,
    handleProjectDatasetAdd,
    handleInitUserStory,
  ]);

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
    inEditor,
    reearthURL,
    backendURL,
    backendProjectName,
    templates: fieldTemplates,
    currentPage,
    loading,
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
