import { Layout, MenuProps } from "antd";
import { MenuInfo } from "rc-menu/lib/interface";
import React,{ memo, useMemo, useState } from "react";

import "../../../node_modules/antd/dist/antd.less";
import "../../styles/style.less";

import { ReactComponent as DataBase } from "../UI/Icon/Icons/dataBase.svg";
import { ReactComponent as Info } from "../UI/Icon/Icons/info.svg";
import { ReactComponent as Share } from "../UI/Icon/Icons/share.svg";
import { ReactComponent as Sliders } from "../UI/Icon/Icons/sliders.svg";
import { ReactComponent as Template } from "../UI/Icon/Icons/template.svg";

import WdContent from "./WdContent";
import WdFooter from "./WdFooter";
import WdHeader from "./WdHeader";

const items: MenuProps["items"] = [
  {
    key: "mapData",
    icon: <DataBase />,
  },
  {
    key: "mapSetting",
    icon: <Sliders />,
  },
  {
    key: "shareNprint",
    icon: <Share />,
  },
  {
    key: "about",
    icon: <Info />,
  },
  {
    key: "template",
    icon: <Template />,
  },
];
export type Props = {
  className?: string;
  isInsideEditor: boolean;
};

const WdLayout: React.FC<Props> = ({ className, isInsideEditor }) => {
  const [current, setCurrent] = useState("mapData");

  const headerItems = useMemo(() => {
    return !isInsideEditor ? [...items.slice(0, -1)] : [...items];
  }, [isInsideEditor]);

  const handleClick: MenuProps["onClick"] = e => {
    console.log("click ", e.key);
    setCurrent(e.key);
  };

  return (
    <Layout className={className}>
      <WdHeader current={current} items={headerItems} onClick={(e: MenuInfo) => handleClick(e)} />
      <WdContent current={current} />
      <WdFooter />
    </Layout>
  );
};
export default memo(WdLayout);
