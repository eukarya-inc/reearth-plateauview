import CommonPage from "@web/extensions/sidebar/core/components/content/CommonPage";
import Menu from "@web/sharedComponents/Menu";
import { styled } from "@web/theme";

import useHooks from "./hooks";

const Help: React.FC = () => {
  const { items, handleItemClicked } = useHooks();
  const { SubMenu } = Menu;
  return (
    <CommonPage title="使い方">
      <MenuWrapper onClick={handleItemClicked} mode="vertical">
        {items.map(item => (
          <SubMenu
            key={item?.key}
            onTitleClick={item?.onclick}
            title={
              <span>
                <span>{item?.label}</span>
              </span>
            }
          />
        ))}
      </MenuWrapper>
    </CommonPage>
  );
};

export default Help;

const MenuWrapper = styled(Menu)`
  display: flex;
  flex-direction: column;
  align-items: space-between;
  padding: 0px;
  gap: 12px;
  background: #e7e7e7;
  width: 326px;
  .ant-menu-submenu-open,
  .ant-menu-submenu-active,
  .ant-menu-submenu-selected {
    background: #00bebe !important;
    color: #ffffff;
  }
  .ant-menu-submenu-vertical:hover,
  .ant-menu-submenu:hover,
  .ant-menu-submenu-title:hover {
    color: #e7e7e7 !important;
  }
  .ant-menu-submenu-arrow::before,
  .ant-menu-submenu-arrow::after {
    background: #e7e7e7;
  }
`;
