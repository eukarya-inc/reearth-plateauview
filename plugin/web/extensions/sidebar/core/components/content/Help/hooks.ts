import { useCallback, useState } from "react";

import { postMsg } from "../../../utils";

type Tab = "basic" | "map" | "shadow" | "clip";

type Items = {
  label: string;
  key: Tab;
};

const items: Items[] = [
  { label: "基本操作", key: "basic" },
  { label: "マップをつくあってみる", key: "map" },
  { label: "日影機能について", key: "shadow" },
  { label: "クリップ機能", key: "clip" },
];

export default () => {
  const [selectedTab, changeTab] = useState<Tab>("basic");

  const handleItemClicked = useCallback((key: Tab) => {
    changeTab(key);
    postMsg({ action: "popup-message", payload: key });
  }, []);

  return {
    items,
    selectedTab,
    handleItemClicked,
  };
};
