import { Input, Spin, Tabs } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useMemo, useState } from "react";

import {
  DataCatalogGroup,
  DataCatalogItem,
  DataSource,
  getDataCatalogTree,
  GroupBy,
} from "../../../../api/api";
import Tags, { Tag as TagType } from "../Tags";

import FileTree from "./FileTree";

export type Tag = TagType;

export type Props = {
  addedDatasetDataIDs?: string[];
  isMobile?: boolean;
  catalog?: DataCatalogItem[];
  selectedTags?: Tag[];
  filter: GroupBy;
  selectedItem?: DataCatalogItem | DataCatalogGroup;
  expandedFolders?: { id?: string; name?: string }[];
  searchTerm: string;
  dataSource?: DataSource;
  setExpandedFolders?: React.Dispatch<React.SetStateAction<{ id?: string; name?: string }[]>>;
  onSearch: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onSelect?: (item?: DataCatalogItem | DataCatalogGroup) => void;
  addDisabled: (dataID: string) => boolean;
  onFilter: (filter: GroupBy) => void;
  onTagSelect?: (tag: Tag) => void;
  onDatasetAdd: (dataset: DataCatalogItem, keepModalOpen?: boolean) => void;
};

// function typeFilter(catalog: Catalog): DataCatalog {
//   const filteredCatalog: CatalogItem[] = prefectures.map(p => {
//     const items: CatalogItem[] = catalog.filter(i => {
//       if (i.prefecture === p) {
//         return {
//           type: "item",
//           ...i,
//         };
//       }
//     }) as CatalogItem[];

//     return {
//       type: "group",
//       name: p,
//       children: items,
//     };
//   });
//   return filteredCatalog;
// }

// function tagFilter(catalog: CatalogRawItem[], tags?: Tag[]): DataCatalog {
//   return catalog
//     .filter(item =>
//       tags?.every(selectedTag => item.tags?.some(tag => selectedTag.name === tag.name)),
//     )
//     .map(item => ({ type: "item", ...item } as CatalogItem));
// }

const DatasetTree: React.FC<Props> = ({
  addedDatasetDataIDs,
  isMobile,
  catalog,
  selectedTags,
  filter,
  selectedItem,
  expandedFolders,
  searchTerm,
  dataSource,
  setExpandedFolders,
  onSearch,
  onSelect,
  addDisabled,
  onFilter,
  onTagSelect,
  onDatasetAdd,
}) => {
  const [loading, _toggleLoading] = useState(false); // needs implementation

  const dataCatalogTree = useMemo(
    () =>
      catalog &&
      getDataCatalogTree(
        catalog,
        filter,
        dataSource === "custom",
        searchTerm.length > 0 ? searchTerm : undefined,
      ),
    [catalog, filter, dataSource, searchTerm],
  );

  const showInput = useMemo(
    () => !selectedTags?.length || searchTerm.length > 0,
    [searchTerm.length, selectedTags?.length],
  );

  const showTags = useMemo(
    () => selectedTags && selectedTags.length > 0 && searchTerm.length === 0,
    [searchTerm.length, selectedTags],
  );

  const showTabs = useMemo(
    () => searchTerm.length > 0 || selectedTags?.length,
    [searchTerm.length, selectedTags],
  );

  return (
    <Wrapper isMobile={isMobile}>
      {showInput && (
        <StyledInput placeholder="検索" value={searchTerm} onChange={onSearch} loading={loading} />
      )}
      {showTags && <Tags tags={selectedTags} onTagSelect={onTagSelect} />}
      {searchTerm.length > 0 && <p style={{ margin: "0", alignSelf: "center" }}>検索結果</p>}
      {dataSource !== "custom" ? (
        <StyledTabs
          activeKey={filter}
          tabBarStyle={showTabs ? { display: "none" } : { userSelect: "none" }}
          onChange={active => onFilter(active as GroupBy)}>
          <Tabs.TabPane key="city" tab="都道府県" style={{ position: "relative" }}>
            {dataCatalogTree ? (
              <FileTree
                addedDatasetDataIDs={addedDatasetDataIDs}
                catalog={dataCatalogTree}
                isMobile={isMobile}
                selectedItem={selectedItem}
                expandedFolders={expandedFolders}
                dataSource={dataSource}
                setExpandedFolders={setExpandedFolders}
                onSelect={onSelect}
                addDisabled={addDisabled}
                onDatasetAdd={onDatasetAdd}
              />
            ) : (
              <Loading>
                <Spin />
              </Loading>
            )}
          </Tabs.TabPane>
          <Tabs.TabPane key="type" tab="種類" style={{ position: "relative" }}>
            {dataCatalogTree ? (
              <FileTree
                addedDatasetDataIDs={addedDatasetDataIDs}
                catalog={dataCatalogTree}
                isMobile={isMobile}
                selectedItem={selectedItem}
                expandedFolders={expandedFolders}
                dataSource={dataSource}
                setExpandedFolders={setExpandedFolders}
                onSelect={onSelect}
                addDisabled={addDisabled}
                onDatasetAdd={onDatasetAdd}
              />
            ) : (
              <Loading>
                <Spin />
              </Loading>
            )}
          </Tabs.TabPane>
        </StyledTabs>
      ) : dataCatalogTree ? (
        <FileTree
          addedDatasetDataIDs={addedDatasetDataIDs}
          catalog={dataCatalogTree}
          isMobile={isMobile}
          selectedItem={selectedItem}
          expandedFolders={expandedFolders}
          dataSource={dataSource}
          setExpandedFolders={setExpandedFolders}
          onSelect={onSelect}
          addDisabled={addDisabled}
          onDatasetAdd={onDatasetAdd}
        />
      ) : (
        <Loading>
          <Spin />
        </Loading>
      )}
    </Wrapper>
  );
};

export default DatasetTree;

const Wrapper = styled.div<{ isMobile?: boolean }>`
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: ${({ isMobile }) => (isMobile ? "24px 12px" : "24px 0 24px 12px")};
  width: ${({ isMobile }) => (isMobile ? "100%" : "310px")};
`;

const StyledInput = styled(Input.Search)`
  .ant-input {
    :hover {
      border: 1px solid var(--theme-color);
    }
  }
  .ant-input-group-addon {
    width: 32px;
    padding: 0;
    :hover {
      cursor: pointer;
    }
  }
`;

const StyledTabs = styled(Tabs)`
  .ant-tabs-nav {
    border-bottom: 0.5px solid #c7c5c5;
    padding: 0 10px;
  }
  .ant-tabs-tab:hover {
    color: var(--theme-color);
  }
  .ant-tabs-ink-bar {
    background: var(--theme-color);
  }
  .ant-tabs-tab.ant-tabs-tab-active .ant-tabs-tab-btn {
    color: var(--theme-color);
  }
`;

const Loading = styled.div`
  position: absolute;
  width: 100%;
  height: 400px;
  left: 0;
  top: 0;
  display: flex;
  align-items: center;
  justify-content: center;
`;
