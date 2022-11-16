import { Icon, Input, Tabs } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useState } from "react";

import FileTree, { Catalog } from "./FileTree";

export type Props = {
  catalog?: Catalog;
  onOpenDetails?: (id: string) => void;
};

const DatasetTree: React.FC<Props> = ({ catalog, onOpenDetails }) => {
  const [searchTerm, setSearchTerm] = useState("");

  return (
    <Wrapper>
      <StyledInput
        placeholder="input search text"
        value={searchTerm}
        onChange={e => setSearchTerm(e.currentTarget.value)}
        addonAfter={<Icon icon="search" size={15} />}
      />
      {searchTerm.length > 0 && <p>Results</p>}
      <StyledTabs
        defaultActiveKey="prefecture"
        tabBarStyle={searchTerm.length > 0 ? { display: "none" } : undefined}>
        <Tabs.TabPane key="prefecture" tab="Prefecture">
          <FileTree filter="prefecture" catalog={catalog} onOpenDetails={onOpenDetails} />
        </Tabs.TabPane>
        <Tabs.TabPane key="type" tab="Type">
          <FileTree filter="fileType" catalog={catalog} onOpenDetails={onOpenDetails} />
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
