import { CloseOutlined } from "@ant-design/icons";
import { Button, Menu, MenuProps, Row } from "antd";
import { Header } from "antd/lib/layout/layout";
import { MenuInfo } from "rc-menu/lib/interface";
import React, { memo } from "react";

import "../../../node_modules/antd/dist/antd.less";
import "../../styles/style.less";
import Icon from "../UI/Icon";

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
        <Button
          type="primary"
          icon={<CloseOutlined />}
          className="closWidgetbtn"
        />
      </Row>
      <Row className="bottomHeader">
        <Icon icon="plateauLogo" height={114.25} width={100} />
        <Menu
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
