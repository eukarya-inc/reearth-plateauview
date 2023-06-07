import { postMsg } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useState } from "react";

export type Tab = "plateau" | "your-data" | "custom";

export default () => {
  const [currentTab, changeTabs] = useState<Tab>();
  const [inEditor, setEditorState] = useState(false);
  const [isCustomProject, setIsCustomProject] = useState(false);

  const handleClose = useCallback(() => {
    postMsg({ action: "modalClose" });
  }, []);

  useEffect(() => {
    postMsg({ action: "initDataCatalog" }); // Needed to trigger sending selected dataset ids from Sidebar
  }, []);

  useEffect(() => {
    const eventListenerCallback = (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      if (e.data.action === "initDataCatalog") {
        setEditorState(e.data.payload.inEditor);
        const isCustomProject =
          e.data.payload.customBackendURL &&
          e.data.payload.customBackendProjectName &&
          e.data.payload.customBackendAccessToken;
        changeTabs(
          e.data.payload.currentDatasetDataSource
            ? e.data.payload.currentDatasetDataSource
            : isCustomProject
            ? "custom"
            : "plateau",
        );
        setIsCustomProject(isCustomProject);
        if (e.data.payload.currentDatasetDataSource) {
          postMsg({ action: "clearCurrentDatasetDataSource" });
        }
      }
    };
    addEventListener("message", eventListenerCallback);
    return () => {
      removeEventListener("message", eventListenerCallback);
    };
  }, []);

  return {
    currentTab,
    inEditor,
    isCustomProject,
    handleClose,
    handleTabChange: changeTabs,
  };
};
