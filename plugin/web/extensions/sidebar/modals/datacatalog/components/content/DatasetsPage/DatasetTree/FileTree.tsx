import { Catalog, CatalogItem } from "@web/extensions/sidebar/core/components/hooks";
import { styled } from "@web/theme";
import { useCallback, useState } from "react";

import TreeBuilder from "./TreeBuilder";

export type { Catalog, CatalogItem } from "@web/extensions/sidebar/core/components/hooks";

export type Props = {
  addedDatasetIds?: string[];
  catalog: Catalog;
  isMobile?: boolean;
  onDatasetAdd: (dataset: CatalogItem) => void;
  onOpenDetails?: (data?: CatalogItem) => void;
};

const FileTree: React.FC<Props> = ({
  addedDatasetIds,
  catalog,
  isMobile,
  onDatasetAdd,
  onOpenDetails,
}) => {
  const [selectedId, select] = useState<string>();

  const handleSelect = useCallback((id?: string) => {
    select(id);
  }, []);

  return (
    <TreeWrapper isMobile={isMobile}>
      <Tree>
        <TreeBuilder
          catalog={catalog}
          addedDatasetIds={addedDatasetIds}
          selectedId={selectedId}
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
