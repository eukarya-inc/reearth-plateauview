import { MenuProps } from "antd";

type MenuItem = Required<MenuProps>["items"][number];

function getItem(
  label: React.ReactNode,
  key?: React.Key | null,
  children?: MenuItem[],
  type?: "group",
): MenuItem {
  return {
    key,
    children,
    label,
    type,
  } as MenuItem;
}

const items: MenuItem[] = [
  getItem("基本操作", "basic", []),
  getItem("マップをつくあってみる", "map", []),
  getItem("日影機能について", "shadow", []),
  getItem("クリップ機能", "clip", []),
];

export default () => {
  const handleItemSelected: MenuProps["onSelect"] = e => {
    console.log("click", e.key);
  };
  return {
    items,
    handleItemSelected,
  };
};
