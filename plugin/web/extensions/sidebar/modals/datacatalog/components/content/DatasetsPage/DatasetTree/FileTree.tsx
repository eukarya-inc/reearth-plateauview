import { DataCatalogItem, DataCatalogGroup } from "@web/extensions/sidebar/core/types";
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
  const [selectedID, select] = useState<string>();
  const [selectedKey, setSelectedKey] = useState("");
  const [expandedKeys, setExpandedKeys] = useState<string[]>([]);

  const handleSelect = useCallback((dataID?: string) => {
    select(dataID);
  }, []);

  const expandAllParentKeys = useCallback((key: string) => {
    const keyArr = key.split("-");
    while (keyArr.length > 1) {
      keyArr.pop();
      const parent = keyArr.join("-");
      setExpandedKeys((prevState: string[]) => {
        const newExpandedKeys = [...prevState];
        if (!prevState.includes(parent)) newExpandedKeys.push(parent);
        return newExpandedKeys;
      });
    }
  }, []);

  useEffect(() => {
    const { selectedDataset } = window as any;
    if (selectedDataset) {
      onOpenDetails?.(selectedDataset);
      handleSelect(selectedDataset.dataID);
      if (selectedKey) expandAllParentKeys(selectedKey);
      setTimeout(() => {
        (window as any).selectedDataset = undefined;
      }, 500);
    }
  }, [expandAllParentKeys, handleSelect, onOpenDetails, selectedKey]);

  return (
    <TreeWrapper isMobile={isMobile}>
      <Tree>
        <TreeBuilder
          catalogItem={catalog}
          addedDatasetDataIDs={addedDatasetDataIDs}
          isMobile={isMobile}
          expandAll={expandAll}
          selectedID={selectedID}
          nestLevel={0}
          nodeKey="0"
          expandedKeys={expandedKeys}
          addDisabled={addDisabled}
          onDatasetAdd={onDatasetAdd}
          onOpenDetails={onOpenDetails}
          onSelect={handleSelect}
          setExpandedKeys={setExpandedKeys}
          setSelectedKey={setSelectedKey}
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
