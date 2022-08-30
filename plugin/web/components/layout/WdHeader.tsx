import { CloseOutlined } from "@ant-design/icons";
import { Button, Menu, MenuProps, Row } from "antd";
import { Header } from "antd/lib/layout/layout";
import { MenuInfo } from "rc-menu/lib/interface";
import React, { memo } from "react";

import { ReactComponent as PlateauLogo } from "../UI/Icon/Icons/plateauLogo.svg";

import "../../../node_modules/antd/dist/antd.less";

type Props = {
  className?: string;
  items: MenuProps["items"];
  onClick: (e: MenuInfo) => void;
  current: string;
};

const WdHeader: React.FC<Props> = ({ className, items, current, onClick }) => {
  return (
    <Header className={className}>
      <Row className="topHeader">
        <Button type="primary" icon={<CloseOutlined />} className="closWidgetbtn" />
      </Row>
      <Row className="bottomHeader">
        <PlateauLogo height={114.25} width={100} />
        <Menu
          selectable={true}
          className="navHeader"
          onClick={onClick}
          selectedKeys={[current]}
          mode="horizontal"
          items={items}
        />
      </Row>
    </Header>
  );
};
export default memo(WdHeader);
