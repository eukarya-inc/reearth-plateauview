import { useCallback, useEffect, useMemo, useState } from "react";

import { Root, Data, Template } from "./newTypes";
import processCatalog, { CatalogRawItem } from "./processCatalog";
import { useCurrentOverrides } from "./state";
import { ReearthApi } from "./types";
import { mergeProperty, postMsg } from "./utils";

export default () => {
  const [projectID, setProjectID] = useState<string>();
  const [inEditor, setInEditor] = useState(true);
  const [backendAccessToken, setBackendAccessToken] = useState<string>();
  const [backendURL, setBackendURL] = useState<string>();
  const [cmsURL, setCMSURL] = useState<string>();

  // ****************************************
  // Init
  useEffect(() => {
    postMsg({ action: "init" }); // Needed to trigger sending initialization data to sidebar
  }, []);
  // ****************************************

  // ****************************************
  // Override
  const [overrides, updateOverrides] = useCurrentOverrides();

  const handleOverridesUpdate = useCallback(
    (updatedProperties: Partial<ReearthApi>) => {
      updateOverrides([overrides, updatedProperties].reduce((p, v) => mergeProperty(p, v)));
    },
    [overrides, updateOverrides],
  );

  useEffect(() => {
    postMsg({ action: "updateOverrides", payload: overrides });
  }, [overrides]);
  // ****************************************

  // ****************************************
  // Minimize
  const [minimized, setMinimize] = useState(false);

  useEffect(() => {
    setTimeout(() => {
      postMsg({ action: "minimize", payload: minimized });
    }, 250);
  }, [minimized]);
  // ****************************************

  // ****************************************
  // Dataset
  const [selectedDatasets, updateDatasets] = useState<CatalogRawItem[]>([]);

  const handleDatasetAdd = useCallback((dataset: CatalogRawItem) => {
    updateDatasets(oldDatasets => [...oldDatasets, dataset]);
    postMsg({ action: "addDatasetToScene", payload: dataset });
  }, []);

  const handleDatasetRemove = useCallback(
    (id: string) => updateDatasets(oldDatasets => oldDatasets.filter(d => d.id !== id)),
    [],
  );

  const handleDatasetRemoveAll = useCallback(() => updateDatasets([]), []);
  // ****************************************

  // ****************************************
  // Catalog
  const [plateauData, setPlateauData] = useState<any[]>([]);
  const [usecaseData, setUsecaseData] = useState<any[]>([]);
  const [datasetData, setDatasetData] = useState<any[]>([]);

  useEffect(() => {
    async function fetchRawData() {
      const plateau = (await (await fetch(`${cmsURL}/plateau`)).json()).results;
      const usecase = (await (await fetch(`${cmsURL}/usecase`)).json()).results;
      const dataset = (await (await fetch(`${cmsURL}/dataset`)).json()).results;
      setPlateauData(plateau);
      setUsecaseData(usecase);
      setDatasetData(dataset);
    }
    if (cmsURL) {
      fetchRawData();
    }
  }, [cmsURL, setPlateauData, setUsecaseData, setDatasetData]);

  const rawCatalog = useMemo(() => {
    console.log("PLATEAU: ", plateauData);
    console.log("USECASE: ", usecaseData);
    console.log("DATASET: ", datasetData);
    return processCatalog(plateauData, usecaseData, datasetData);
  }, [plateauData, usecaseData, datasetData]);

  const handleModalOpen = useCallback(() => {
    const selectedIds = selectedDatasets.map(d => d.id);
    postMsg({
      action: "datacatalog-modal-open",
      payload: { addedDatasets: selectedIds, rawCatalog },
    });
  }, [rawCatalog, selectedDatasets]);
  // ****************************************

  // ****************************************
  // Templates
  const [templates, setTemplates] = useState<Template[]>([]);

  const handleTemplateAdd = useCallback(
    async (newTemplate?: Template) => {
      if (!backendURL || !backendAccessToken) return;
      const res = await fetch(`${backendURL}/viz/plateau/templates`, {
        headers: {
          authorization: `Bearer ${backendAccessToken}`,
        },
        method: "POST",
        body: JSON.stringify(newTemplate),
      });
      if (res.status !== 200) return;
      const data = await res.json();
      setTemplates(t => [...t, data.results]);
      return data.results as Template;
    },
    [backendURL, backendAccessToken],
  );

  const handleTemplateUpdate = useCallback(
    async (template: Template) => {
      if (!template.modelId || !backendURL || !backendAccessToken) return;
      const res = await fetch(`${backendURL}/viz/plateau/templates/${template.modelId}`, {
        headers: {
          authorization: `Bearer ${backendAccessToken}`,
        },
        method: "PATCH",
        body: JSON.stringify(template),
      });
      if (res.status !== 200) return;
      const updatedTemplate = (await res.json()).results;
      setTemplates(t => {
        return t.map(t2 => {
          if (t2.id === updatedTemplate.id) {
            return updatedTemplate;
          }
          return t2;
        });
      });
    },
    [backendURL, backendAccessToken],
  );

  const handleTemplateRemove = useCallback(
    async (template: Template) => {
      if (!template.modelId || !backendURL || !backendAccessToken) return;
      const res = await fetch(`${backendURL}/viz/plateau/templates/${template.modelId}`, {
        headers: {
          authorization: `Bearer ${backendAccessToken}`,
        },
        method: "DELETE",
      });
      if (res.status !== 200) return;
      setTemplates(t => t.filter(t2 => t2.modelId !== template.modelId));
    },
    [backendURL, backendAccessToken],
  );

  // ****************************************

  // ****************************************
  // Processed Data

  // MIGHT NEED FIXING. MIGHT NEED FIXING. MIGHT NEED FIXING
  const [data, setData] = useState<Data[]>();
  const processedSelectedDatasets = useMemo(() => {
    if (!data) return;
    return selectedDatasets
      .map(d => {
        if (d.modelType === "usecase") {
          // If usecase, check "data" for saved template, components, etc
          return data.filter(d3 => d3.dataId === `plateau-2022-${d.cityName}`);
        } else if (d.modelType === "plateau") {
          // Else, if PLATEAUデータ(plateau), do ....(HARDCODED TEMPLATE)
          return data;
        } else if (d.modelType === "dataset") {
          // Else, if 関連データセット(dataset), do ....(HARDCODED TEMPLATE)
          return data;
        } else {
          // SOME DEFAULT?????
          return data;
        }
      })
      .flat(1)
      .filter(p => p);
  }, [data, selectedDatasets]);
  // MIGHT NEED FIXING. MIGHT NEED FIXING. MIGHT NEED FIXING

  // ****************************************

  useEffect(() => {
    const eventListenerCallback = (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      if (e.data.type === "msgFromModal") {
        if (e.data.payload.dataset) {
          handleDatasetAdd(e.data.payload.dataset);
        }
      } else if (e.data.type === "init") {
        setProjectID(e.data.payload.projectID);
        setInEditor(e.data.payload.inEditor);
        setBackendAccessToken(e.data.payload.backendAccessToken);
        setBackendURL(e.data.payload.backendURL);
        setCMSURL(`${e.data.payload.cmsURL}/api/p/plateau-2022`);
      }
    };
    addEventListener("message", e => eventListenerCallback(e));
    return () => {
      removeEventListener("message", eventListenerCallback);
    };
  }, [handleDatasetAdd]);

  useEffect(() => {
    // Fetch from backend (if projectID)
    if (projectID) {
      // fetch
    } else {
      // Fetch from backend (if no projectID)
      const fetchBackend = async () => {
        const res = await fetch(`${cmsURL}/viz/plateau`);
        if (res.status !== 200) return;
        const results: Root = (await res.json()).results;
        setTemplates(results.templates);
        setData(results.data);
      };
      if (cmsURL) {
        fetchBackend();
      }
    }
  }, [projectID, cmsURL]);

  return {
    processedSelectedDatasets,
    overrides,
    minimized,
    inEditor,
    backendURL,
    templates,
    handleTemplateAdd,
    handleTemplateUpdate,
    handleTemplateRemove,
    setMinimize,
    handleDatasetRemove,
    handleDatasetRemoveAll,
    handleOverridesUpdate,
    handleModalOpen,
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
