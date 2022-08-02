import { Content } from "antd/lib/layout/layout";
import { memo, ReactNode } from "react";

import "../../../node_modules/antd/dist/antd.less";
import MapSettingTab from "../Tabs/MapSettingTab";

type Props = {
  className?: string;
  children?: ReactNode;
  current?: string;
};
const WdContent: React.FC<Props> = ({ className, current, children }) => {
  console.log(current);
  return (
    <Content className={className}>
      {current == "mapSetting" ? <MapSettingTab /> : null}
      {children}
    </Content>
  );
};
export default memo(WdContent);
