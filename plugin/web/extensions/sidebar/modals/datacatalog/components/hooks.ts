import { Catalog, CatalogItem } from "@web/extensions/sidebar/core/components/hooks";
import { UserDataItem } from "@web/extensions/sidebar/modals/datacatalog/types";
import { postMsg } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useState } from "react";

export type Tab = "dataset" | "your-data";

export default () => {
  const [currentTab, changeTabs] = useState<Tab>("dataset");
  const [addedDatasetIds, setAddedDatasetIds] = useState<string[]>();
  const [rawCatalog, setRawCatalog] = useState<Catalog>({});

  const handleClose = useCallback(() => {
    postMsg({ action: "modalClose" });
  }, []);

  const handleDatasetAdd = useCallback(
    (dataset: CatalogItem | UserDataItem) => {
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

  useEffect(() => {
    // Needed to trigger sending selected dataset ids from Sidebar
    postMsg({ action: "initDataCatalog" });
  }, []);

  useEffect(() => {
    const eventListenerCallback = (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      if (e.data.type === "initDataCatalog") {
        setAddedDatasetIds(e.data.payload.addedDatasets);
        setRawCatalog(e.data.payload.catalogData);
      }
    };
    addEventListener("message", eventListenerCallback);
    return () => {
      removeEventListener("message", eventListenerCallback);
    };
  }, []);

  return {
    currentTab,
    rawCatalog,
    addedDatasetIds,
    handleClose,
    handleTabChange: changeTabs,
    handleDatasetAdd,
  };
};
