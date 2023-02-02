import {
  CatalogItem,
  DataCatalog as DataCatalogType,
} from "@web/extensions/sidebar/core/processCatalog";
import { Button, Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useMemo, useState } from "react";

export type DataCatalog = DataCatalogType;

export type Props = {
  catalog: DataCatalog;
  isMobile?: boolean;
  onDatasetAdd: (dataset: CatalogItem) => void;
  onOpenDetails?: (data?: CatalogItem) => void;
};

const TreeBuilder: React.FC<{
  item: CatalogItem;
  selectedId?: string;
  nestLevel: number;
  onDatasetAdd: (dataset: CatalogItem) => void;
  onOpenDetails?: (item?: CatalogItem) => void;
  onSelect?: (id: string) => void;
}> = ({ item, selectedId, nestLevel, onDatasetAdd, onOpenDetails, onSelect }) => {
  const [isOpen, open] = useState(false);

  const selected = useMemo(
    () => (item.type !== "group" ? selectedId === item.id : false),
    [selectedId, item],
  );

  const handleOpenDetails = useCallback(() => {
    if (item.type === "group") return;
    onOpenDetails?.(item);
    onSelect?.(item.id);
  }, [item, onOpenDetails, onSelect]);

  const handleClick = useCallback(() => {
    onDatasetAdd(item);
  }, [item, onDatasetAdd]);

  return item.type === "group" ? (
    <Folder key={item.name} isOpen={isOpen}>
      <FolderItem nestLevel={nestLevel} onClick={() => open(!isOpen)}>
        <NameWrapper>
          <Icon icon={isOpen ? "folderOpen" : "folder"} size={20} />
          <Name>{item.name}</Name>
        </NameWrapper>
      </FolderItem>
      {item.children.map(m =>
        TreeBuilder({
          item: m,
          selectedId,
          nestLevel: nestLevel + 1,
          onDatasetAdd,
          onOpenDetails,
          onSelect,
        }),
      )}
    </Folder>
  ) : (
    <FolderItem key={item.id} nestLevel={nestLevel} selected={selected}>
      <NameWrapper onClick={handleOpenDetails}>
        <Icon icon={"file"} size={20} />
        <Name>{item.cityName ?? item.name}</Name>
      </NameWrapper>
      <Button
        type="link"
        icon={<StyledIcon icon="plusCircle" selected={selected} />}
        onClick={handleClick}
      />
    </FolderItem>
  );
};

const FileTree: React.FC<Props> = ({ catalog, isMobile, onDatasetAdd, onOpenDetails }) => {
  const [selectedId, select] = useState<string>();

  const handleSelect = useCallback((id?: string) => {
    select(id);
  }, []);

  return (
    <TreeWrapper isMobile={isMobile}>
      <Tree>
        {catalog.map(item =>
          TreeBuilder({
            item,
            selectedId,
            nestLevel: 1,
            onDatasetAdd,
            onOpenDetails,
            onSelect: handleSelect,
          }),
        )}
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

const Folder = styled.div<{ isOpen?: boolean }>`
  width: 100%;
  ${({ isOpen }) =>
    isOpen
      ? "height: 100%;"
      : `
  height: 29px; 
  overflow: hidden;
  `}
`;

const FolderItem = styled.div<{ nestLevel: number; selected?: boolean }>`
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-sizing: border-box;
  gap: 8px;
  min-height: 29px;
  ${({ selected }) =>
    selected &&
    `
  background: #00BEBE;
  color: white;
  `}

  padding-left: ${({ nestLevel }) => (nestLevel ? `${nestLevel * 8}px` : "8px")};
  padding-right: 8px;
  cursor: pointer;

  :hover {
    background: #00bebe;
    color: white;
  }
`;

const NameWrapper = styled.div`
  display: flex;
`;

const Name = styled.p`
  margin: 0 0 0 8px;
  user-select: none;
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
  width: 200px;
`;

const StyledIcon = styled(Icon)<{ selected: boolean }>`
  color: ${({ selected }) => (selected ? "#ffffff" : "#00bebe")};
  ${FolderItem}:hover & {
    color: #ffffff;
  }
`;
