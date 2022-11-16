import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useState } from "react";

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
  customProperties?: {
    initialCamera?: Camera;
  };
};

export type Data = {
  id: string; // ex. "//PLATEAU データセット/東京都"
  name: string; // "東京都"
} & (Group | Item);

export type Catalog = Data[];

export type Props = {
  filter?: "prefecture" | "fileType";
  catalog?: Catalog;
  onOpenDetails?: (id: string) => void;
};

const TreeBuilder: React.FC<{
  data: Data;
  nestLevel: number;
  onOpenDetails?: (id: string) => void;
}> = ({ data, nestLevel, onOpenDetails }) => {
  const [isOpen, open] = useState(data.type === "group" ? data.isOpen : undefined);

  return data.type === "group" ? (
    <Folder key={data.id} isOpen={isOpen}>
      <FolderItem nestLevel={nestLevel} onClick={() => open(!isOpen)}>
        <StyledIcon icon={isOpen ? "folderOpen" : "folder"} size={20} />
        <Name>{data.name}</Name>
      </FolderItem>
      {data.members?.map(m => TreeBuilder({ data: m, nestLevel: nestLevel + 1 }))}
    </Folder>
  ) : (
    <FolderItem key={data.id} nestLevel={nestLevel} onClick={() => onOpenDetails?.(data.id)}>
      <StyledIcon icon={"file"} size={20} />
      <Name>{data.name}</Name>
    </FolderItem>
  );
};

const FileTree: React.FC<Props> = ({ filter, catalog, onOpenDetails }) => {
  console.log(filter, "filter");
  const plateauDatasets = (catalog?.[0] as Group).members;
  return (
    <Tree>
      {plateauDatasets?.map(dataset => TreeBuilder({ data: dataset, nestLevel: 1, onOpenDetails }))}
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

const FolderItem = styled.div<{ nestLevel: number }>`
  display: flex;
  align-items: center;
  box-sizing: border-box;
  gap: 8px;
  min-height: 29px;

  padding-left: ${({ nestLevel }) => (nestLevel ? `${nestLevel * 8}px` : "8px")};
  padding-right: 8px;
  cursor: pointer;

  :hover {
    background: #dcdcdc;
  }
`;

const Name = styled.p`
  margin: 0;
`;

const StyledIcon = styled(Icon)``;
