import { Layout, MenuProps } from "antd";
import { MenuInfo } from "rc-menu/lib/interface";
import { memo, useMemo, useState } from "react";

import "../../../node_modules/antd/dist/antd.less";
import "../../styles/style.less";
import Icon from "../UI/Icon";

import WdContent from "./WdContent";
import WdFooter from "./WdFooter";
import WdHeader from "./WdHeader";

const items: MenuProps["items"] = [
  {
    key: "mapData",
    icon: <Icon icon="dataBase" />,
  },
  {
    key: "mapSetting",
    icon: <Icon icon="sliders" />,
  },
  {
    key: "shareNprint",
    icon: <Icon icon="paperPlabeTill" />,
  },
  {
    key: "about",
    icon: <Icon icon="info" />,
  },
  {
    key: "alignLeft",
    icon: <Icon icon="alignLeft" />,
  },
];
export type Props = {
  className?: string;
  isInsideEditor: boolean;
  currentTab: string;
};

const WdLayout: React.FC<Props> = ({
  className,
  isInsideEditor,
  currentTab,
}) => {
  const [current, setCurrent] = useState(currentTab);

  const headerItems = useMemo(() => {
    return !isInsideEditor ? [...items.slice(0, -1)] : [...items];
  }, [isInsideEditor]);

  const handleClick: MenuProps["onClick"] = (e) => {
    console.log("click ", e.key);
    setCurrent(e.key);
  };

  return (
    <Layout className={className}>
      <WdHeader
        current={current}
        items={headerItems}
        onClick={(e: MenuInfo) => handleClick(e)}
      />
      <WdContent current={current} />
      <WdFooter />
    </Layout>
  );
};
export default memo(WdLayout);
