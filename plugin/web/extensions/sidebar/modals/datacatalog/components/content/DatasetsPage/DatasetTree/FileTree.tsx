import { DataCatalogItem, DataCatalogGroup } from "@web/extensions/sidebar/core/types";
import { styled } from "@web/theme";
import { useCallback, useState } from "react";

import TreeBuilder from "./TreeBuilder";

export type Props = {
  addedDatasetDataIDs?: string[];
  catalog: (DataCatalogItem | DataCatalogGroup)[];
  isMobile?: boolean;
  expandAll?: boolean;
  onDatasetAdd: (dataset: DataCatalogItem) => void;
  onOpenDetails?: (data?: DataCatalogItem) => void;
};

const FileTree: React.FC<Props> = ({
  addedDatasetDataIDs,
  catalog,
  isMobile,
  expandAll,
  onDatasetAdd,
  onOpenDetails,
}) => {
  const [selectedID, select] = useState<string>();

  const handleSelect = useCallback((dataID?: string) => {
    select(dataID);
  }, []);

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
          onDatasetAdd={onDatasetAdd}
          onOpenDetails={onOpenDetails}
          onSelect={handleSelect}
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
