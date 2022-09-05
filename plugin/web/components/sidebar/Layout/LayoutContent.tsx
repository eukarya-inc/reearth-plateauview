import { Content } from "antd/lib/layout/layout";
import React, { memo, ReactNode } from "react";

type Props = {
  className?: string;
  children?: ReactNode;
  current: string;
};

const LayoutContent: React.FC<Props> = ({ className, children }) => {
  return (
    <Content className={className}>
      {/* {
        {
          mapSetting: <MapSettingTab />,
          shareNprint: <ShareTab />,
          about: <InfoTab />,
        }[current]
      } */}
      {children}
    </Content>
  );
};
export default memo(LayoutContent);
