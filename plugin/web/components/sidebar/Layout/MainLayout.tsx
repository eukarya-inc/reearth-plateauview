// import " antd/dist/antd.less";
import { Layout, MenuProps } from "antd";
import { MenuInfo } from "rc-menu/lib/interface";
import { memo, useMemo, useState } from "react";

import { ReactComponent as DataBase } from "../../common/Icon/Icons/dataBase.svg";
import { ReactComponent as Info } from "../../common/Icon/Icons/info.svg";
import { ReactComponent as Share } from "../../common/Icon/Icons/share.svg";
import { ReactComponent as Sliders } from "../../common/Icon/Icons/sliders.svg";
import { ReactComponent as Template } from "../../common/Icon/Icons/template.svg";

import LayoutContent from "./LayoutContent";
import LayoutFooter from "./LayoutFooter";
import LayoutHeader from "./LayoutHeader";

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
    setCurrent(e.key);
  };

  return (
    <Layout className={className}>
      <LayoutHeader
        current={current}
        items={headerItems}
        onClick={(e: MenuInfo) => handleClick(e)}
      />
      <LayoutContent current={current} />
      <LayoutFooter />
    </Layout>
  );
};
export default memo(WdLayout);
