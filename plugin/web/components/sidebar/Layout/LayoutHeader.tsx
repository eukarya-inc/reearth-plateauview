import { CloseOutlined } from "@ant-design/icons";
import { Button, Menu, MenuProps, Row } from "antd";
import { Header } from "antd/lib/layout/layout";
import { MenuInfo } from "rc-menu/lib/interface";
import React, { memo } from "react";

import { styled } from "../../../theme";
import { ReactComponent as PlateauLogo } from "../../common/Icon/Icons/plateauLogo.svg";

type Props = {
  className?: string;
  items: MenuProps["items"];
  onClick: (e: MenuInfo) => void;
  current: string;
};

const LayoutHeader: React.FC<Props> = ({ className, items, current, onClick }) => {
  return (
    <Header className={className}>
      <TopHeader>
        <ClosWidgetbtn type="primary" icon={<CloseOutlined />} />
      </TopHeader>
      <BottomHeader>
        <PlateauLogo height={114.25} width={100} />
        <NavHeader
          selectable={true}
          onClick={onClick}
          selectedKeys={[current]}
          mode="horizontal"
          items={items}
        />
      </BottomHeader>
    </Header>
  );
};
export default memo(LayoutHeader);
const TopHeader = styled(Row)`
  direction: rtl;
  align-items: flex-start;
  height: 32px;
  width: 100%;
`;
const BottomHeader = styled(Row)`
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  padding: 0px 0px 10px 10px;
  height: 50px;
`;
const NavHeader = styled(Menu)`
  height: 40px;
  width: 100%;
`;
const ClosWidgetbtn = styled(Button)`
  border-radius: 0%;
  height: 32px;
  width: 32px;
`;
