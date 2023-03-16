import { DataCatalogItem } from "@web/extensions/sidebar/core/types";
import { UserDataItem } from "@web/extensions/sidebar/modals/datacatalog/types";
import { postMsg } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useState } from "react";

import { GroupBy } from "../api/api";

export type Tab = "dataset" | "your-data";
export type TreeTab = "city" | "type";

export default () => {
  const [currentTab, changeTabs] = useState<Tab>("dataset");
  const [currentTreeTab, changeTreeTab] = useState<TreeTab>("city");
  const [addedDatasetDataIDs, setAddedDatasetDataIDs] = useState<string[]>();
  const [catalog, setCatalog] = useState<DataCatalogItem[]>([]);
  const [inEditor, setEditorState] = useState(false);
  const [filter, setFilter] = useState<GroupBy>("city");

  const handleClose = useCallback(() => {
    postMsg({ action: "modalClose" });
  }, []);

  const handleFilter = useCallback((filter: GroupBy) => {
    setFilter(filter);
  }, []);

  const handleDatasetAdd = useCallback(
    (dataset: DataCatalogItem | UserDataItem) => {
      postMsg({
        action: "msgFromModal",
        payload: {
          dataset,
        },
      });
      handleClose();
    },
    [handleClose],
  );

  const handleDatasetPublish = useCallback((dataID: string, publish: boolean) => {
    postMsg({ action: "updateDataset", payload: { dataID, publish } });
  }, []);

  const handleTreeTabChange = useCallback(
    (treeTab: TreeTab) => {
      handleFilter(treeTab);
      changeTreeTab(treeTab);
      postMsg({ action: "storageSave", payload: { key: "currentTreeTab", value: treeTab } });
    },
    [handleFilter],
  );

  useEffect(() => {
    postMsg({ action: "initDataCatalog" }); // Needed to trigger sending selected dataset ids from Sidebar
  }, []);

  useEffect(() => {
    const eventListenerCallback = (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      if (e.data.action === "initDataCatalog") {
        setAddedDatasetDataIDs(e.data.payload.addedDatasets);
        setCatalog(e.data.payload.dataCatalog);
        setEditorState(e.data.payload.inEditor);
        handleTreeTabChange(e.data.payload.currentTreeTab);
      } else if (e.data.action === "updateCatalog") {
        setCatalog(e.data.payload);
      }
    };
    addEventListener("message", eventListenerCallback);
    return () => {
      removeEventListener("message", eventListenerCallback);
    };
  }, [handleTreeTabChange]);

  return {
    currentTab,
    currentTreeTab,
    catalog,
    addedDatasetDataIDs,
    inEditor,
    filter,
    handleFilter,
    handleClose,
    handleTabChange: changeTabs,
    handleTreeTabChange,
    handleDatasetAdd,
    handleDatasetPublish,
  };
};
