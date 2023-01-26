import { UserDataItem } from "@web/extensions/sidebar/modals/datacatalog/types";
import { Tabs } from "@web/sharedComponents";
import { styled } from "@web/theme";

import LocalDataTab from "./LocalDataTab";
import WebDataTab from "./WebDataTab";

export type Props = {
  onOpenDetails?: (data?: UserDataItem) => void;
};

const FileSelectPane: React.FC<Props> = ({ onOpenDetails }) => {
  return (
    <Wrapper>
      <Tabs defaultActiveKey="local" style={{ marginBottom: "12px" }}>
        <Tabs.TabPane tab="Add Local Data" key="local">
          <LocalDataTab onOpenDetails={onOpenDetails} />
        </Tabs.TabPane>
        <Tabs.TabPane tab="Add Web Data" key="web">
          <WebDataTab onOpenDetails={onOpenDetails} />
        </Tabs.TabPane>
      </Tabs>
    </Wrapper>
  );
};

export default FileSelectPane;

const Wrapper = styled.div`
  padding: 24px 12px;
`;
