import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useMemo, useState } from "react";

type XYZ = {
  x: number;
  y: number;
  z: number;
};

type Camera = {
  east: number;
  north: number;
  south: number;
  west: number;
  direction: XYZ;
  position: XYZ;
  up: XYZ;
};

type Group = {
  type: "group";
  isOpen?: boolean; // default false. Maybe not needed for v2?
  members?: Data[]; // if type is "group"
};

type Item = {
  type: "3d-tiles" | "wms-no-description"; // Many types... like "wms-no-description". So might need to make multiple possible individual Item types
  url?: string;
  description?: string;
  tags?: Tag[];
  customProperties?: {
    initialCamera?: Camera;
  };
};

export type Data = {
  id: string; // ex. "//PLATEAU データセット/東京都"
  name: string; // "東京都"
} & (Group | Item);

export type Catalog = Data[];

export type FilterType = "prefecture" | "fileType" | "tag";

export type Tag = {
  type: "location" | "data-type";
  name: string;
};

export type Props = {
  filter?: FilterType;
  catalog?: Catalog;
  onOpenDetails?: (data?: Data) => void;
};

const TreeBuilder: React.FC<{
  data: Data;
  selectedId?: string;
  nestLevel: number;
  onOpenDetails?: (data?: Data) => void;
  onSelect?: (id: string) => void;
}> = ({ data, selectedId, nestLevel, onOpenDetails, onSelect }) => {
  const [isOpen, open] = useState(data.type === "group" ? data.isOpen : undefined);
  const selected = useMemo(() => selectedId === data.id, [selectedId, data.id]);

  const handleOpenDetails = useCallback(() => {
    if (data.type === "group") return;
    onOpenDetails?.(data);

    onSelect?.(data.id);
  }, [data, onOpenDetails, onSelect]);

  return data.type === "group" ? (
    <Folder key={data.id} isOpen={isOpen}>
      <FolderItem nestLevel={nestLevel} onClick={() => open(!isOpen)}>
        <StyledIcon icon={isOpen ? "folderOpen" : "folder"} size={20} />
        <Name>{data.name}</Name>
      </FolderItem>
      {data.members?.map(m =>
        TreeBuilder({ data: m, selectedId, nestLevel: nestLevel + 1, onOpenDetails, onSelect }),
      )}
    </Folder>
  ) : (
    <FolderItem key={data.id} nestLevel={nestLevel} selected={selected} onClick={handleOpenDetails}>
      <StyledIcon icon={"file"} size={20} />
      <Name>{data.name}</Name>
    </FolderItem>
  );
};

const FileTree: React.FC<Props> = ({ filter, catalog, onOpenDetails }) => {
  const [selectedId, select] = useState<string>();
  console.log(filter, "filter");
  const plateauDatasets = (catalog?.[0] as Group).members;

  const handleSelect = useCallback((id?: string) => {
    select(id);
  }, []);

  return (
    <Tree>
      {plateauDatasets?.map(dataset =>
        TreeBuilder({
          data: dataset,
          selectedId,
          nestLevel: 1,
          onOpenDetails,
          onSelect: handleSelect,
        }),
      )}
    </Tree>
  );
};

export default FileTree;

const Tree = styled.div`
  display: flex;
  flex-direction: column;
  flex-wrap: wrap;
  width: 298px;
`;

const Folder = styled.div<{ isOpen?: boolean }>`
  width: 100%;
  ${({ isOpen }) => (isOpen ? "height: 100%;" : "overflow: hidden; height: 29px;")}
`;

const FolderItem = styled.div<{ nestLevel: number; selected?: boolean }>`
  display: flex;
  align-items: center;
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

const Name = styled.p`
  margin: 0;
  user-select: none;
`;

const StyledIcon = styled(Icon)``;
