import { postMsg } from "@web/extensions/sidebar/utils";
import { uniq, intersection } from "lodash";
import Papa from "papaparse";
import { useCallback, useEffect, useState, useRef } from "react";

import type {
  InitData,
  Dataset,
  Condition,
  Result,
  Viewport,
  RawDatasetData,
  IndexData,
  SearchIndex,
  SearchResults,
} from "../types";

// import { TEST_RESULT_DATA } from "./TEST_DATA";

export type Size = {
  width: number;
  height: number;
};

type DataRowIds = number[];

export default () => {
  // UI
  const [minimized, setMinimized] = useState<boolean>(false);
  const [sizes, setSizes] = useState<{ default: Size; mobile: Size; mini: Size }>({
    default: {
      width: 347,
      height: 524,
    },
    mobile: {
      width: 0,
      height: 524,
    },
    mini: {
      width: 238,
      height: 82,
    },
  });
  const [isMobile, setIsMobile] = useState<boolean>(false);
  const [size, setSize] = useState<Size>(sizes.default);
  const prevSizeRef = useRef<Size>(size);

  const toggleMinimize = useCallback(() => {
    prevSizeRef.current = size;
    setMinimized(!minimized);
  }, [minimized, size]);

  const setHtmlSize = useCallback((size: Size) => {
    document.documentElement.style.width = `${size.width}px`;
    document.documentElement.style.height = `${size.height}px`;
  }, []);

  useEffect(() => {
    const targetSize = minimized ? sizes.mini : isMobile ? sizes.mobile : sizes.default;
    setSize(targetSize);
    if (
      targetSize.width >= prevSizeRef.current.width ||
      targetSize.height >= prevSizeRef.current.height
    ) {
      setHtmlSize(targetSize);
    } else {
      setTimeout(() => {
        setHtmlSize(targetSize);
      }, 500);
    }
  }, [minimized, sizes, isMobile, setHtmlSize]);

  const handleResize = useCallback(
    (viewport: Viewport) => {
      if (viewport.isMobile) {
        setIsMobile(true);
        setSizes({
          ...sizes,
          mobile: { width: viewport.width * 0.9, height: sizes.mobile.height },
        });
      } else if (isMobile) {
        setIsMobile(false);
      }
    },
    [sizes, isMobile],
  );

  const [activeTab, setActiveTab] = useState<"condition" | "result">("condition");
  const onClickCondition = useCallback(() => {
    setActiveTab("condition");
  }, []);
  const onClickResult = useCallback(() => {
    setActiveTab("result");
  }, []);
  const [conditionsState, setConditionsState] = useState<"loading" | "empty" | "ready">("loading");

  // Data
  const [dataset, setDataset] = useState<Dataset>();
  const [conditions, setConditions] = useState<Condition[]>([]);
  const [results, setResults] = useState<Result[]>([]);
  const [resultStyleCondition, setResultStyleCondition] = useState<string>();
  const [highlightAll, setHighlightAll] = useState<boolean>(true);
  const [showMatchingOnly, setShowMatchingOnly] = useState<boolean>(false);
  const [selected, setSelected] = useState<Result[]>([]);
  const [isSearching, setIsSearching] = useState<boolean>(false);

  const searchIndexes = useRef<SearchIndex[]>();
  const searchResults = useRef<SearchResults[]>();

  const loadDetailData = useCallback(async (url: string) => {
    let results: { dataRowId: number }[] = [];
    await fetch(url)
      .then(response => response.text())
      .then(v => {
        results = Papa.parse(v, { header: true, skipEmptyLines: true }).data as typeof results;
      });
    return results;
  }, []);

  const loadResultsData = useCallback(async (si: SearchIndex) => {
    if (si.resultsData) return;
    await fetch(`${si.baseURL}/resultsData.csv`)
      .then(response => response.text())
      .then(v => {
        si.resultsData = Papa.parse(v, {
          header: true,
          skipEmptyLines: true,
          fastMode: true,
        }).data;
      });
  }, []);

  useEffect(() => {
    if (!dataset?.dataID) return;

    const colorConditions: [string, string][] = [];

    if (highlightAll && resultStyleCondition) {
      colorConditions.push([resultStyleCondition, "color('red')"]);
    } else {
      if (selected) {
        let selectedConditon = "";
        selected.forEach(bldg => {
          if (selectedConditon) selectedConditon += " || ";
          selectedConditon += "${gml_id} === '" + bldg.gml_id + "'";
        });
        if (selectedConditon) {
          colorConditions.push([selectedConditon, "color('red')"]);
        }
      }
    }

    colorConditions.push(["true", "color()"]);

    const showConditions: [string, string][] = [];

    if (showMatchingOnly) {
      if (resultStyleCondition) {
        showConditions.push([resultStyleCondition, "true"]);
      }
      showConditions.push(["true", "false"]);
    } else {
      showConditions.push(["true", "true"]);
    }

    const styles = {
      "3dtiles": {
        color: {
          expression: {
            conditions: colorConditions,
          },
        },
        show: {
          expression: {
            conditions: showConditions,
          },
        },
      },
    };

    postMsg({
      action: "updateDatasetInScene",
      payload: {
        dataID: dataset.dataID,
        update: styles,
      },
    });
  }, [selected, resultStyleCondition, dataset?.dataID, highlightAll, showMatchingOnly]);

  const conditionApply = useCallback(() => {
    searchResults.current = [];
    const combinedResults: Result[] = [];

    (async () => {
      if (searchIndexes.current) {
        await Promise.all(
          searchIndexes.current.map(async si => {
            // get all conditions groups for current search index
            const threeDTilesId = si.baseURL.split("/").pop() ?? "";
            const condGroups: DataRowIds[] = [];

            await Promise.all(
              conditions.map(async cond => {
                if (cond.values.length > 0) {
                  const condGroup: DataRowIds = [];
                  await Promise.all(
                    cond.values.map(async v => {
                      await loadDetailData(
                        `${si.baseURL}/${si.indexRoot.indexes[cond.field].values[v].url}`,
                      ).then(v => {
                        condGroup.push(...v.map(item => item.dataRowId));
                      });
                    }),
                  );
                  condGroups.push(condGroup);
                }
              }),
            );

            let results: Result[] = [];

            const rowIds = intersection(...condGroups);
            if (rowIds) {
              await loadResultsData(si);
              results = rowIds.map(rowId => si.resultsData?.[rowId]);
            }

            searchResults.current?.push({
              threeDTilesId,
              results,
            });

            combinedResults.push(...results);
          }),
        );

        // combine results from different search index
        setResults(combinedResults);
        setIsSearching(false);
        setResultStyleCondition(() => {
          let resultCondition = "";
          conditions.map(c => {
            if (c.values.length > 0) {
              if (resultCondition) resultCondition += " && ";
              let currentCondition = "(";
              c.values.map(value => {
                if (currentCondition !== "(") currentCondition += " || ";
                currentCondition += "${" + c.field + "} === '" + value + "'";
              });
              currentCondition += ")";
              resultCondition += currentCondition;
            }
          });
          return resultCondition;
        });
      }
    })();

    setResults([]);
    setIsSearching(true);
    setActiveTab("result");
    setSelected([]);
    setHighlightAll(true);
    setShowMatchingOnly(false);
  }, [conditions, loadDetailData, loadResultsData]);

  useEffect(() => {
    if (selected.length === 1) {
      postMsg({
        action: "cameraLookAt",
        payload: [
          {
            lng: Number(selected[0].Longitude),
            lat: Number(selected[0].Latitude),
            height: Number(selected[0].Height) + 100,
            range: 200,
          },
          { duration: 2 },
        ],
      });
    }
  }, [selected]);

  useEffect(() => {
    if (results.length > 0) {
      postMsg({
        action: "cameraLookAt",
        payload: [
          {
            lng: Number(results[0].Longitude),
            lat: Number(results[0].Latitude),
            height: Number(results[0].Height) + 100,
            range: 200,
          },
          { duration: 2 },
        ],
      });
    }
  }, [results]);

  const initDatasetData = useCallback(
    (rawDatasetData: RawDatasetData) => {
      setResults([]);
      setSelected([]);
      setActiveTab("condition");
      setConditionsState("loading");

      if (!rawDatasetData.searchIndex) {
        setConditionsState("empty");
      } else {
        searchIndexes.current = [];

        const indexData: IndexData[] = [];

        const allIndexes =
          typeof rawDatasetData.searchIndex === "string"
            ? [{ url: rawDatasetData.searchIndex }]
            : rawDatasetData.searchIndex;

        (async () => {
          await Promise.all(
            allIndexes.map(async si => {
              const baseURL = si.url.replace("/indexRoot.json", "").replace(".zip", "");

              const indexRootRes = await fetch(`${baseURL}/indexRoot.json`);
              if (indexRootRes.status !== 200) return;

              const indexRoot = await indexRootRes.json();
              if (indexRoot && searchIndexes.current) {
                searchIndexes.current.push({
                  baseURL,
                  indexRoot,
                });

                Object.keys(indexRoot.indexes).forEach(field => {
                  const f = indexData.find(mi => mi.field === field);
                  if (!f) {
                    indexData.push({
                      field,
                      values: Object.keys(indexRoot.indexes[field].values),
                    });
                  } else {
                    f.values.push(...Object.keys(indexRoot.indexes[field].values));
                  }
                });
              }
            }),
          );

          indexData.forEach(indexDataItem => {
            indexDataItem.values = uniq(indexDataItem.values);
          });

          setDataset({
            title: rawDatasetData.title,
            dataID: rawDatasetData.dataID,
            indexes: indexData,
          });
          setConditions(indexData.map(index => ({ field: index.field, values: [] })));
          setConditionsState("ready");

          // preload results data
          setTimeout(() => {
            searchIndexes.current?.forEach(si => {
              loadResultsData(si);
            });
          }, 0);
        })();
      }
    },
    [loadResultsData],
  );

  const popupClose = useCallback(() => {
    postMsg({ action: "popupClose" });
  }, []);

  const onInit = useCallback(
    (initData: InitData | undefined) => {
      if (!initData) return;
      if (initData.viewport.isMobile) {
        setIsMobile(true);
        setSizes({
          ...sizes,
          mobile: { width: initData.viewport.width * 0.9, height: sizes.mobile.height },
        });
      }
      setMinimized(false);
      initDatasetData((window as any).buildingSearchInit?.data);
    },
    [initDatasetData, sizes],
  );

  useEffect(() => {
    document.documentElement.style.setProperty("--theme-color", "#00BEBE");
    onInit((window as any).buildingSearchInit);

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const onMessage = useCallback(
    (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      switch (e.data.type) {
        case "resize":
          handleResize(e.data.payload);
          break;
        case "buildingSearchInit":
          onInit(e.data.payload);
          break;
        default:
          break;
      }
    },
    [handleResize, onInit],
  );

  useEffect(() => {
    addEventListener("message", onMessage);
    return () => {
      removeEventListener("message", onMessage);
    };
  }, [onMessage]);

  return {
    size,
    minimized,
    activeTab,
    dataset,
    results,
    highlightAll,
    showMatchingOnly,
    selected,
    isSearching,
    conditionsState,
    onClickCondition,
    onClickResult,
    toggleMinimize,
    popupClose,
    setConditions,
    conditionApply,
    setHighlightAll,
    setShowMatchingOnly,
    setSelected,
  };
};
