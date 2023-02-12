import { Project, ReearthApi } from "@web/extensions/sidebar/types";
import { mergeProperty, postMsg } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useState } from "react";

import { Data, Template } from "../newTypes";

import { Pages } from "./Header";

export type Catalog = {
  [key: string]: Catalog | CatalogItem[];
};

export type CatalogItem = {
  id?: string;
  name?: string;
  pref?: string;
  city?: string;
  city_en?: string;
  city_code?: string;
  ward?: string;
  ward_en?: string;
  ward_code?: string;
  type?: string;
  format?: string;
  layers?: string;
  url?: string;
  desc?: string;
  search_index?: string;
  year?: string;
  config?: any;
};

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
  const [inEditor, setInEditor] = useState(true);
  const [backendAccessToken, setBackendAccessToken] = useState<string>();
  const [backendURL, setBackendURL] = useState<string>();
  // const [cmsURL, setCMSURL] = useState<string>();
  const [reearthURL, setReearthURL] = useState<string>();

  const [data, setData] = useState<Data[]>();
  const [project, updateProject] = useState<Project>(defaultProject);
  const [processedSelectedDatasets, setProcessedSelectedDatasets] = useState<Data[]>([]);

  const handleBackendFetch = useCallback(async () => {
    if (!backendURL) return;
    const res = await fetch(`${backendURL}/sidebar/plateauview`);
    if (res.status !== 200) return;
    const resData = await res.json();

    setTemplates(resData.templates);
    setData(resData.data);
  }, [backendURL]);

  // ****************************************
  // Init
  useEffect(() => {
    postMsg({ action: "init" }); // Needed to trigger sending initialization data to sidebar
  }, []);
  // ****************************************

  // ****************************************
  // Project

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

  const handleProjectDatasetAdd = useCallback((dataset: any) => {
    updateProject(({ sceneOverrides, selectedDatasets }) => {
      const updatedProject: Project = {
        sceneOverrides,
        selectedDatasets: [
          ...selectedDatasets,
          {
            id: dataset.id,
            dataId: `plateau-2022-${dataset.cityName ?? dataset.name}`,
            type: dataset.type,
            name: dataset.cityName ?? dataset.name,
            visible: true,
          } as Data,
        ],
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

  const handleProjectDatasetRemoveAll = useCallback(
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

  const handleDatasetUpdate = useCallback(
    (updatedDataset: Data) => {
      if (processedSelectedDatasets.length < 1) return;

      const updatedProcessedDatasets = [...processedSelectedDatasets];
      const datasetIndex = updatedProcessedDatasets.findIndex(d2 => d2.id === updatedDataset.id);

      updatedProcessedDatasets[datasetIndex] = updatedDataset;
      setProcessedSelectedDatasets(updatedProcessedDatasets);
    },
    [processedSelectedDatasets],
  );

  const handleDatasetSave = useCallback(
    (datasetID: string) => {
      (async () => {
        if (!inEditor) return;
        const datasetToSave = processedSelectedDatasets.find(d => d.id === datasetID);
        const isNew = !data?.find(d => d.id === datasetID);

        if (!backendURL || !backendAccessToken || !datasetToSave) return;

        const fetchURL = !isNew
          ? `${backendURL}/sidebar/plateauview/data/${datasetToSave.id}`
          : `${backendURL}/sidebar/plateauview/data`;

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
        // setTemplates(t => [...t, data.results]);
        console.log("DATA JUST SAVED: ", data2);
        handleBackendFetch(); // MAYBE UPDATE THIS LATER TO JUST UPDATE THE LOCAL VALUE
      })();
    },
    [data, processedSelectedDatasets, inEditor, backendAccessToken, backendURL, handleBackendFetch],
  );

  // ****************************************

  // ****************************************
  // Catalog
  const [catalogData, setCatalog] = useState<Catalog>({});

  useEffect(() => {
    if (catalogData) {
      console.log("CATALOG DATA: ", catalogData);
    }
  }, [catalogData]);

  useEffect(() => {
    (async function fetchRawData() {
      const catalog: CatalogItem[] = await (
        await fetch("https://api.plateau.reearth.io/datacatalog")
      ).json();

      setCatalog(processCatalogByPref(catalog));
    })();
  }, []);

  const handleModalOpen = useCallback(() => {
    const selectedIds = project.selectedDatasets.map(d => d.id);
    postMsg({
      action: "catalogModalOpen",
      payload: { addedDatasets: selectedIds, catalogData },
    });
  }, [catalogData, project.selectedDatasets]);
  // ****************************************

  // ****************************************
  // Templates
  const [templates, setTemplates] = useState<Template[]>([]);

  const handleTemplateAdd = useCallback(
    async (newTemplate?: Template) => {
      if (!backendURL || !backendAccessToken) return;
      const res = await fetch(`${backendURL}/sidebar/plateauview/templates`, {
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
      const res = await fetch(`${backendURL}/sidebar/plateauview/templates/${template.modelId}`, {
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
      const res = await fetch(`${backendURL}/sidebar/plateauview/templates/${template.modelId}`, {
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
        setBackendAccessToken(e.data.payload.backendAccessToken);
        setBackendURL(e.data.payload.backendURL);
        // setCMSURL(`${e.data.payload.cmsURL}/api/p/plateau-2022`);
        setReearthURL(`${e.data.payload.reearthURL}`);
        if (e.data.payload.draftProject) {
          updateProject(e.data.payload.draftProject);
        }
      } else if (e.data.action === "triggerCatalogOpen") {
        handleModalOpen();
      } else if (e.data.action === "triggerHelpOpen") {
        handlePageChange("help");
      }
    };
    addEventListener("message", eventListenerCallback);
    return () => {
      removeEventListener("message", eventListenerCallback);
    };
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (!backendURL) return;
    if (projectID) {
      (async () => {
        const res = await fetch(`${backendURL}/share/plateauview/${projectID}`);
        if (res.status !== 200) return;
        const data = await res.json();
        if (data) {
          updateProject(data);
          postMsg({ action: "updateProject", payload: data });
        }
      })();
    }
  }, [projectID, backendURL]);

  useEffect(() => {
    if (backendURL) {
      handleBackendFetch();
    }
  }, [backendURL]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    setProcessedSelectedDatasets(
      !data
        ? project.selectedDatasets
        : project.selectedDatasets
            .map(sd => {
              const savedData = data.find(d => d.dataId === sd.dataId);
              if (savedData) {
                return {
                  ...sd,
                  ...savedData,
                };
              } else {
                return sd;
              }
            })
            .flat(1)
            .filter(p => p),
    );
  }, [data, project.selectedDatasets]);

  const [currentPage, setCurrentPage] = useState<Pages>("data");

  const handlePageChange = useCallback((p: Pages) => {
    setCurrentPage(p);
  }, []);

  return {
    catalogData,
    project,
    processedSelectedDatasets,
    inEditor,
    reearthURL,
    backendURL,
    templates,
    currentPage,
    handlePageChange,
    handleTemplateAdd,
    handleTemplateUpdate,
    handleTemplateRemove,
    handleDatasetSave,
    handleDatasetUpdate,
    handleProjectDatasetAdd,
    handleProjectDatasetRemove,
    handleProjectDatasetRemoveAll,
    handleProjectSceneUpdate,
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

// Re-order catalog by Prefecture (行政コード)
export const processCatalogByPref = (c: CatalogItem[]): Catalog => {
  const byPref: { [key: string]: CatalogItem[] } = c.reduce(
    (acc: { [key: string]: CatalogItem[] }, cur) => {
      const key = cur["pref"];
      if (key) {
        if (!acc[key]) {
          acc[key] = [];
        }
        acc[key].push(cur);
      }
      return acc;
    },
    {},
  );

  const byRest = Object.keys(byPref).map(prefKey => {
    const byCity = byPref[prefKey].reduce((acc2, cur2) => {
      const generateStructure = (keys: string[], object: CatalogItem) => {
        return keys.reduce((o, k) => {
          if (!o[k as keyof typeof o]) {
            o[k as keyof typeof o] = [];
          }
          if (k === keys[keys.length - 1]) {
            o[k as keyof typeof o].push(cur2);
          } else {
            o[k as keyof typeof o] = generateStructure(keys.slice(1), o[k as keyof typeof o]);
          }
          return o;
        }, object);
      };

      let cityKeys: string[] = [];
      let nameKeys: string[] = [];

      if (cur2["city"]) {
        cityKeys = cur2["city"].split("/");
      }
      if (cur2["name"]) {
        nameKeys = cur2["name"].split("/");
      }

      const filterKeys = nameKeys.length > 1 ? [...cityKeys, ...nameKeys] : cityKeys;

      return generateStructure(filterKeys, acc2);
    }, {});

    return { [prefKey]: byCity };
  });

  return byRest.reduce((acc, item) => {
    return { ...acc, ...item };
  }, {}) as Catalog;
};

export const prefectures = [
  "全国",
  "東京都",
  "北海道",
  "青森県",
  "岩手県",
  "宮城県",
  "秋田県",
  "山形県",
  "福島県",
  "茨城県",
  "栃木県",
  "群馬県",
  "埼玉県",
  "千葉県",
  "神奈川県",
  "新潟県",
  "富山県",
  "石川県",
  "福井県",
  "山梨県",
  "長野県",
  "岐阜県",
  "静岡県",
  "愛知県",
  "三重県",
  "滋賀県",
  "京都府",
  "大阪府",
  "兵庫県",
  "奈良県",
  "和歌山県",
  "鳥取県",
  "島根県",
  "岡山県",
  "広島県",
  "山口県",
  "徳島県",
  "香川県",
  "愛媛県",
  "高知県",
  "福岡県",
  "佐賀県",
  "長崎県",
  "熊本県",
  "大分県",
  "宮崎県",
  "鹿児島県",
  "沖縄県",
];

export const dataTypes = [
  "建築物モデル",
  "避難施設",
  "ランドマーク",
  "鉄道駅",
  "道路モデル",
  "植生",
  "都市設備",
  "土地利用",
  "土砂災害警戒区域",
  "都市計画決定情報", // Sub as-is
  "洪水浸水想定区域", // Sub by name
  "津波浸水想定区域", // Sub by name
  "高潮浸水想定区域", // Sub by name
  "内水浸水想定区域", // Sub by name
  "緊急輸送道路",
  "鉄道",
  "公園",
  "行政界",
  "ユースケース", // Sub by name
];

export type Tag = {
  type: "location" | "data-type";
  name: string;
};

// tags: [
//   { name: item.prefecture, type: "location" },
//   { name: item.city_name, type: "location" },
//   { name: item.data_format, type: "data-type" },
//   { name: item.type, type: "data-type" },
// ].filter(t => !!t.name) as Tag[],
