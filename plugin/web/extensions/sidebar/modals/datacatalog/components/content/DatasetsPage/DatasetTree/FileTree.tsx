import { DataCatalogItem, DataCatalogGroup } from "@web/extensions/sidebar/core/types";
import { postMsg } from "@web/extensions/sidebar/utils";
import { styled } from "@web/theme";
import { useCallback, useEffect, useState } from "react";

import TreeBuilder from "./TreeBuilder";

export type Props = {
  addedDatasetDataIDs?: string[];
  catalog: (DataCatalogItem | DataCatalogGroup)[];
  isMobile?: boolean;
  expandAll?: boolean;
  addDisabled: (dataID: string) => boolean;
  onDatasetAdd: (dataset: DataCatalogItem) => void;
  onOpenDetails?: (data?: DataCatalogItem) => void;
};

const FileTree: React.FC<Props> = ({
  addedDatasetDataIDs,
  catalog,
  isMobile,
  expandAll,
  addDisabled,
  onDatasetAdd,
  onOpenDetails,
}) => {
  const [selectedItem, selectItem] = useState<DataCatalogItem>();
  const [expandedFolders, setExpandedFolders] = useState<{ id?: string; name?: string }[]>([]);

  const handleSelect = useCallback((item?: DataCatalogItem) => {
    selectItem(item);
  }, []);

  useEffect(() => {
    postMsg({ action: "catalogModalOpen" }); // Needed to trigger sending client storage data from Sidebar
  }, []);

  useEffect(() => {
    const eventListenerCallback = (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      if (e.data.action === "catalogModalOpen") {
        if (e.data.payload.expandedFolders) setExpandedFolders(e.data.payload.expandedFolders);
        if (e.data.payload.dataset) {
          const item = e.data.payload.dataset;
          onOpenDetails?.(item);
          handleSelect(item);
          if (item.path) {
            const expandedFolders = item.path.map((item: string) => ({ name: item }));
            const index = expandedFolders.findIndex((folder: { id?: string; name?: string }) =>
              folder.id ? folder.id === item.id : folder.name === item.name,
            );
            if (index >= 0) expandedFolders.splice(index, 1);
            setExpandedFolders(expandedFolders);
          }
          postMsg({
            action: "saveDataset",
            payload: { dataset: undefined },
          });
        }
      }
    };
    addEventListener("message", eventListenerCallback);
    return () => {
      removeEventListener("message", eventListenerCallback);
    };
  }, [expandedFolders, handleSelect, onOpenDetails]);

  return (
    <TreeWrapper isMobile={isMobile}>
      <Tree>
        <TreeBuilder
          catalogItem={catalog}
          addedDatasetDataIDs={addedDatasetDataIDs}
          isMobile={isMobile}
          expandAll={expandAll}
          selectedID={selectedItem?.id}
          nestLevel={0}
          expandedFolders={expandedFolders}
          addDisabled={addDisabled}
          onDatasetAdd={onDatasetAdd}
          onOpenDetails={onOpenDetails}
          onSelect={handleSelect}
          setExpandedFolders={setExpandedFolders}
        />
      </Tree>
    </TreeWrapper>
  );
};

export default FileTree;

const TreeWrapper = styled.div<{ isMobile?: boolean }>`
  width: ${({ isMobile }) => (isMobile ? "100%" : "298px")};
  height: ${({ isMobile }) => (isMobile ? "100%" : "400px")};
  overflow-y: scroll;
`;

const Tree = styled.div`
  display: flex;
  flex-direction: column;
  flex-wrap: wrap;
`;
