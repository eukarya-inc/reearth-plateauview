import { Content } from "antd/lib/layout/layout";
import React,{ memo, ReactNode } from "react";

import "../../../node_modules/antd/dist/antd.less";
import InfoTab from "../Tabs/InfoTab";
import MapSettingTab from "../Tabs/MapSettingTab";
import ShareTab from "../Tabs/ShareTab";

type Props = {
  className?: string;
  children?: ReactNode;
  current: string;
};

//  {
//     key: "mapData",
//   },
//   {
//     key: "mapSetting",
//   },
//   {
//     key: "shareNprint",
//   },
//   {
//     key: "about",
//   },
//   {
//     key: "template",
//   },

const WdContent: React.FC<Props> = ({ className, current, children }) => {
  return (
    <Content className={className}>
      {
        {
          mapSetting: <MapSettingTab />,
          shareNprint: <ShareTab />,
          about: <InfoTab />,
        }[current]
      }
      {children}
    </Content>
  );
};
export default memo(WdContent);
