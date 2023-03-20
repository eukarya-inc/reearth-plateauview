import { DataCatalogItem } from "@web/extensions/sidebar/core/types";
import { UserDataItem } from "@web/extensions/sidebar/modals/datacatalog/types";
import { postMsg } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useState } from "react";

export type Tab = "dataset" | "your-data";

export default () => {
  const [currentTab, changeTabs] = useState<Tab>("dataset");
  const [addedDatasetDataIDs, setAddedDatasetDataIDs] = useState<string[]>();
  const [catalog, setCatalog] = useState<DataCatalogItem[]>([]);
  const [inEditor, setEditorState] = useState(false);
  const [selectedDatasetID, setDatasetID] = useState<string>();
  const [selectedItem, selectItem] = useState<DataCatalogItem>();
  const [expandedFolders, setExpandedFolders] = useState<{ id?: string; name?: string }[]>([]);

  const handleSelect = useCallback((item?: DataCatalogItem) => {
    selectItem(item);
  }, []);

  const handleOpenDetails = useCallback((data?: DataCatalogItem) => {
    setDatasetID(data?.dataID);
  }, []);

  const handleClose = useCallback(() => {
    postMsg({ action: "modalClose" });
  }, []);

  const handleDatasetAdd = useCallback(
    (dataset: DataCatalogItem | UserDataItem, keepModalOpen?: boolean) => {
      postMsg({
        action: "msgFromModal",
        payload: {
          dataset,
        },
      });
      if (!keepModalOpen) handleClose();
    },
    [handleClose],
  );

  const handleDatasetPublish = useCallback((dataID: string, publish: boolean) => {
    postMsg({ action: "updateDataset", payload: { dataID, publish } });
  }, []);

  useEffect(() => {
    postMsg({ action: "initDataCatalog" }); // Needed to trigger sending selected dataset ids from Sidebar
  }, []);

  useEffect(() => {
    const eventListenerCallback = (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      if (e.data.action === "initDataCatalog") {
        setAddedDatasetDataIDs(e.data.payload.addedDatasets);
        setCatalog(e.data.payload.catalog);
        setEditorState(e.data.payload.inEditor);
        setExpandedFolders(e.data.payload.expandedFolders);
        if (e.data.payload.dataset) {
          const item = e.data.payload.dataset;
          handleOpenDetails(item);
          handleSelect(item);
          if (item.path) {
            setExpandedFolders(
              item.path
                .map((item: string) => ({ name: item }))
                .filter((folder: { name?: string }) => folder.name !== item.name),
            );
          }
          postMsg({
            action: "saveDataset",
            payload: { dataset: undefined },
          });
        }
      } else if (e.data.action === "updateDataCatalog") {
        if (e.data.payload.updatedCatalog) {
          setCatalog(e.data.payload.updatedCatalog);
        }
        if (e.data.payload.updatedDatasetDataIDs) {
          setAddedDatasetDataIDs(e.data.payload.updatedDatasetDataIDs);
        }
      }
    };
    addEventListener("message", eventListenerCallback);
    return () => {
      removeEventListener("message", eventListenerCallback);
    };
  }, [handleOpenDetails, handleSelect]);

  return {
    currentTab,
    catalog,
    addedDatasetDataIDs,
    inEditor,
    selectedDatasetID,
    selectedItem,
    expandedFolders,
    setExpandedFolders,
    handleSelect,
    handleOpenDetails,
    handleClose,
    handleTabChange: changeTabs,
    handleDatasetAdd,
    handleDatasetPublish,
  };
};
