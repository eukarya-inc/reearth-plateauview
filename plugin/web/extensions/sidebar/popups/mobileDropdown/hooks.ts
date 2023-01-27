import { Project, ReearthApi } from "@web/extensions/sidebar/types";
import { mergeProperty, postMsg } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useMemo, useState } from "react";

import { Data, Root } from "../../core/newTypes";
import processCatalog, { CatalogRawItem } from "../../core/processCatalog";

export const defaultProject: Project = {
  sceneOverrides: {
    default: {
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
  },
  selectedDatasets: [],
};

export default () => {
  const [projectID, setProjectID] = useState<string>();
  //   const [backendAccessToken, setBackendAccessToken] = useState<string>();
  const [backendURL, setBackendURL] = useState<string>();
  const [cmsURL, setCMSURL] = useState<string>();
  const [reearthURL, setReearthURL] = useState<string>();

  // ****************************************
  // Init
  useEffect(() => {
    postMsg({ action: "init" }); // Needed to trigger sending initialization data to sidebar
  }, []);
  // ****************************************

  // ****************************************
  // Project
  const [project, updateProject] = useState<Project>(defaultProject);

  const handleProjectSceneUpdate = useCallback(
    (updatedProperties: Partial<ReearthApi>) => {
      updateProject(({ sceneOverrides, selectedDatasets }) => {
        const updatedProject: Project = {
          sceneOverrides: [sceneOverrides, updatedProperties].reduce((p, v) => mergeProperty(p, v)),
          selectedDatasets,
        };
        postMsg({ action: "updateProject", payload: updatedProject });
        return updatedProject;
      });
    },
    [updateProject],
  );

  const handleProjectDatasetAdd = useCallback((dataset: CatalogRawItem) => {
    updateProject(({ sceneOverrides, selectedDatasets }) => {
      const updatedProject = {
        sceneOverrides,
        selectedDatasets: [...selectedDatasets, dataset],
      };
      postMsg({ action: "updateProject", payload: updatedProject });
      return updatedProject;
    });

    postMsg({ action: "addDatasetToScene", payload: dataset }); // MIGHT NEED TO MOVE THIS ELSEWHEREEEE
  }, []);

  const handleProjectDatasetRemove = useCallback(
    (id: string) =>
      updateProject(({ sceneOverrides, selectedDatasets }) => {
        const updatedProject = {
          sceneOverrides,
          selectedDatasets: selectedDatasets.filter(d => d.id !== id),
        };
        postMsg({ action: "updateProject", payload: updatedProject });
        return updatedProject;
      }),
    [],
  );

  const handleDatasetRemoveAll = useCallback(
    () =>
      updateProject(({ sceneOverrides }) => {
        const updatedProject = {
          sceneOverrides,
          selectedDatasets: [],
        };
        postMsg({ action: "updateProject", payload: updatedProject });
        return updatedProject;
      }),
    [],
  );
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

  const rawCatalog = useMemo(
    () => processCatalog(plateauData, usecaseData, datasetData),
    [plateauData, usecaseData, datasetData],
  );

  // ****************************************
  // ****************************************
  // *****************LOOK HERE**************
  // ****************************************
  // ****************************************
  // ****************************************
  console.log("NEED TO SEND THE RAW CATALOG TO MOBILE DATA CATALOG WHEN AVAILABLE: ", rawCatalog);
  // ****************************************
  // ****************************************
  // *****************LOOK HERE**************
  // ****************************************
  // ****************************************
  // ****************************************

  // ****************************************
  // Processed Data

  const [data, setData] = useState<Data[]>();
  const processedSelectedDatasets: Data[] = useMemo(() => {
    // if (!data) return data;
    console.log("PROJECT: ", project);
    return project.selectedDatasets
      .map(d => {
        console.log("DATA: ", data);
        if (d.modelType === "usecase") {
          // If usecase, check "data" for saved template, components, etc
          // return data?.filter(d3 => d3.dataId === `plateau-2022-${d.cityName}`);
          return {
            id: d.id,
            dataId: "ASDFSDFASDFasdf", // <======= NEEDS TO BE UPDATED
            type: d.type ?? "", // maybe not needed
            name: d.cityName ?? d.name,
            public: false, //<======= NEEDS TO BE UPDATED
            // visible <=== this will come from data (or be default true)
            // template <=== this will come from data (and/or be added later from editor)
            // components: data?.filter(d=> d.),
          };
          // } else if (d.modelType === "plateau") {
          //   // Else, if PLATEAUデータ(plateau), do ....(HARDCODED TEMPLATE)
          //   return d;
          // } else if (d.modelType === "dataset") {
          //   // Else, if 関連データセット(dataset), do ....(HARDCODED TEMPLATE)
          //   return d;
        } else {
          return {
            id: d.id,
            dataId: `plateau-2022-${d.cityName ?? d.name}`,
            type: d.type ?? "", // maybe not needed
            name: d.cityName ?? d.name,
            public: false,
            visible: true,
            template: "SOME TEMPLATE NAME???????????????????????????????",
            components: [],
          };
        }
      })
      .flat(1)
      .filter(p => p);
  }, [data, project]);
  // ****************************************

  useEffect(() => {
    const eventListenerCallback = (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      if (e.data.action === "msgFromModal") {
        if (e.data.payload.dataset) {
          handleProjectDatasetAdd(e.data.payload.dataset);
        }
      } else if (e.data.action === "init") {
        setProjectID(e.data.payload.projectID);
        // setBackendAccessToken(e.data.payload.backendAccessToken);
        setBackendURL(e.data.payload.backendURL);
        setCMSURL(`${e.data.payload.cmsURL}/api/p/plateau-2022`);
        setReearthURL(`${e.data.payload.reearthURL}`);
        if (e.data.payload.draftProject) {
          updateProject(e.data.payload.draftProject);
        }
      }
    };
    addEventListener("message", e => eventListenerCallback(e));
    return () => {
      removeEventListener("message", eventListenerCallback);
    };
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (!backendURL) return;
    if (projectID) {
      (async () => {
        const res = await fetch(`${backendURL}/share/${projectID}`);
        if (res.status !== 200) return;
        const data = await res.json();
        updateProject(data);
        postMsg({ action: "updateProject", payload: data });
      })();
    } else {
      (async () => {
        const res = await fetch(`${backendURL}/viz/plateau`);
        if (res.status !== 200) return;
        const results: Root = (await res.json()).results;
        // setTemplates(results.templates);
        setData(results.data);
      })();
    }
  }, [projectID, backendURL]);

  return {
    processedSelectedDatasets,
    project,
    reearthURL,
    backendURL,
    handleProjectDatasetRemove,
    handleDatasetRemoveAll,
    handleProjectSceneUpdate,
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
