import { Icon, Input, Tabs } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useState } from "react";

import Tags from "../Tags";

import FileTree, { Catalog, Data, FilterType, Tag } from "./FileTree";

export type Props = {
  catalog?: Catalog;
  selectedTags?: Tag[];
  onTagSelect: (tag: Tag) => void;
  onOpenDetails?: (data?: Data) => void;
};

const DatasetTree: React.FC<Props> = ({ catalog, selectedTags, onTagSelect, onOpenDetails }) => {
  const [filterType, setFilterType] = useState<FilterType>("fileType");
  const [searchTerm, setSearchTerm] = useState("");

  const handleSearch = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.currentTarget.value);
  }, []);

  const handleFilter = useCallback((filter: FilterType) => {
    setFilterType(filter);
  }, []);

  return (
    <Wrapper>
      <StyledInput
        placeholder="input search text"
        value={searchTerm}
        onChange={handleSearch}
        addonAfter={<Icon icon="search" size={15} />}
      />
      {searchTerm.length > 0 && <p>Results</p>}
      {selectedTags && <Tags tags={selectedTags} onTagSelect={onTagSelect} />}
      <StyledTabs
        defaultActiveKey="prefecture"
        tabBarStyle={searchTerm.length > 0 ? { display: "none" } : undefined}
        onChange={active => handleFilter(active as FilterType)}>
        <Tabs.TabPane key="prefecture" tab="Prefecture">
          <FileTree filter="prefecture" catalog={catalog} onOpenDetails={onOpenDetails} />
        </Tabs.TabPane>
        <Tabs.TabPane key="type" tab="Type">
          <FileTree filter={filterType} catalog={catalog} onOpenDetails={onOpenDetails} />
        </Tabs.TabPane>
      </StyledTabs>
    </Wrapper>
  );
};

export default DatasetTree;

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 24px 0 24px 12px;
  width: 310px;
`;

const StyledInput = styled(Input)`
  .ant-input {
    :hover {
      border: 1px solid #00bebe;
    }
  }
  .ant-input-group-addon {
    width: 32px;
    padding: 0;
  }
`;

const StyledTabs = styled(Tabs)`
  .ant-tabs-nav {
    border-bottom: 0.5px solid #c7c5c5;
    padding: 0 10px;
  }
  .ant-tabs-tab:hover {
    color: #00bebe;
  }
  .ant-tabs-ink-bar {
    background: #00bebe;
  }
  .ant-tabs-tab.ant-tabs-tab-active .ant-tabs-tab-btn {
    color: #00bebe;
  }
`;
