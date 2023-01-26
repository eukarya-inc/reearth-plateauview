import { UserDataItem } from "@web/extensions/sidebar/modals/datacatalog/types";
import { Tabs } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useState } from "react";

import LocalDataTab from "./LocalDataTab";
import WebDataTab from "./WebDataTab";

export type Props = {
  onOpenDetails?: (data?: UserDataItem) => void;
};

const FileSelectPane: React.FC<Props> = ({ onOpenDetails }) => {
  const [selectedLocalItem, setSelectedLocalItem] = useState<UserDataItem>();
  const [selectedWebItem, setSelectedWebItem] = useState<UserDataItem>();

  const handleTabChange = (activeKey: string) => {
    switch (activeKey) {
      case "local":
        if (onOpenDetails) onOpenDetails(selectedLocalItem);
        break;
      case "web":
        if (onOpenDetails) onOpenDetails(selectedWebItem);
        break;
    }
  };

  return (
    <Wrapper>
      <StyledTabs defaultActiveKey="local" onChange={handleTabChange}>
        <Tabs.TabPane tab="Add Local Data" key="local">
          <LocalDataTab onOpenDetails={onOpenDetails} setSelectedLocalItem={setSelectedLocalItem} />
        </Tabs.TabPane>
        <Tabs.TabPane tab="Add Web Data" key="web">
          <WebDataTab onOpenDetails={onOpenDetails} setSelectedWebItem={setSelectedWebItem} />
        </Tabs.TabPane>
      </StyledTabs>
    </Wrapper>
  );
};

export default FileSelectPane;

const Wrapper = styled.div`
  padding: 24px 12px;
`;

const StyledTabs = styled(Tabs)`
  margin-bottom: 12px;
`;
