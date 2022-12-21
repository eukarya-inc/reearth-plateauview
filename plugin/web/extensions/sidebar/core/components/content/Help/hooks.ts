import { useCallback } from "react";

import { postMsg } from "../../../utils";

export default () => {
  const handleItemClicked = useCallback((e: any) => {
    console.log("click", e);
    postMsg({ action: "show-popup", payload: e.key });
  }, []);

  const items = [
    { label: "基本操作", key: "basic", onclick: handleItemClicked },
    { label: "マップをつくあってみる", key: "map", onclick: handleItemClicked },
    { label: "日影機能について", key: "shadow", onclick: handleItemClicked },
    { label: "クリップ機能", key: "clip", onclick: handleItemClicked },
  ];

  return {
    items,
    handleItemClicked,
  };
};
