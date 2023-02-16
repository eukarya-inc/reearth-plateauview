import { postMsg } from "@web/extensions/sidebar/utils";
import { uniq, intersection } from "lodash";
import Papa from "papaparse";
import { useCallback, useEffect, useState, useRef } from "react";

import type {
  DatasetIndexes,
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

  // Data
  const [datasetIndexes, setDatasetIndexes] = useState<DatasetIndexes>();
  const [conditions, setConditions] = useState<Condition[]>([]);
  const [results, setResults] = useState<Result[]>([]);
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
    // TODO: flyTo the selected feature
    if (selected.length === 1) {
      postMsg({
        action: "cameraFlyTo",
        payload: [
          {
            lng: selected[0].Longitude,
            lat: selected[0].Latitude,
          },
          { duration: 2 },
        ],
      });
    }
  }, [selected]);

  useEffect(() => {
    // TODO: flyTo the result (feature) if only one
    if (results.length === 1) {
      postMsg({
        action: "cameraFlyTo",
        payload: [
          {
            lng: results[0].Longitude,
            lat: results[0].Latitude,
          },
          { duration: 2 },
        ],
      });
    }
  }, [results]);

  useEffect(() => {
    // TODO: Update 3D tiles style
  }, [highlightAll, showMatchingOnly, selected, results]);

  const initDatasetData = useCallback(
    (rawDatasetData: RawDatasetData) => {
      if (!rawDatasetData) return;

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

        setDatasetIndexes({
          title: rawDatasetData.title,
          indexes: indexData,
        });
        setConditions(indexData.map(index => ({ field: index.field, values: [] })));

        // preload results data
        setTimeout(() => {
          searchIndexes.current?.forEach(si => {
            loadResultsData(si);
          });
        }, 0);
      })();

      // const datasetIndexes = ((window as any).buildingSearchInit?.data ??
      //   TEST_DATASET_INDEX_DATA) as DatasetIndexes;

      // setDatasetIndexes(datasetIndexes);
      // setConditions(datasetIndexes.indexes.map(index => ({ field: index.field, values: [] })));
    },
    [loadResultsData],
  );

  const popupClose = useCallback(() => {
    postMsg({ action: "popupClose" });
  }, []);

  useEffect(() => {
    if ((window as any).buildingSearchInit) {
      const init = (window as any).buildingSearchInit;
      if (init.viewport.isMobile) {
        setIsMobile(true);
        setSizes({
          ...sizes,
          mobile: { width: init.viewport.width * 0.9, height: sizes.mobile.height },
        });
      }
    }

    document.documentElement.style.setProperty("--theme-color", "#00BEBE");

    initDatasetData((window as any).buildingSearchInit?.data);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const onMessage = useCallback(
    (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      switch (e.data.type) {
        case "resize":
          handleResize(e.data.payload);
          break;
        default:
          break;
      }
    },
    [handleResize],
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
    datasetIndexes,
    results,
    highlightAll,
    showMatchingOnly,
    selected,
    isSearching,
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
