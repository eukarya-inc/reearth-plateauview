import { postMsg } from "@web/extensions/sidebar/utils";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { ReactNode, useCallback, useEffect, useState } from "react";

import PopupItem from "../sharedComponents/PopupItem";

type ItemKey = "map" | "share" | "feedback";

type MenuItem = {
  key: ItemKey;
  title: string;
  icon: ReactNode;
};

const menuItems: MenuItem[] = [
  {
    key: "map",
    title: "マップ設定",
    icon: <Icon icon="sliders" />,
  },
  {
    key: "share",
    title: "共有・印刷",
    icon: <Icon icon="share" />,
  },
  {
    key: "feedback",
    title: "ご意見・ご要望",
    icon: <Icon icon="feedback" />,
  },
];

const Menu: React.FC = () => {
  const [currentItem, changeItem] = useState<MenuItem | undefined>();

  const handleHeightUpdate = () => {
    const el = document.getElementById("menu");
    const currentHeight = el ? window.getComputedStyle(el).height : undefined;
    postMsg({
      action: "msgFromPopup",
      payload: { height: currentHeight },
    });
  };

  const handleClick = useCallback(
    (item: MenuItem) => {
      if (currentItem === item) {
        changeItem(undefined);
      } else {
        changeItem(item);
      }
      handleHeightUpdate();
    },
    [currentItem],
  );

  useEffect(() => {
    handleHeightUpdate();
  }, []);

  return (
    <Wrapper id="menu">
      {currentItem ? (
        <PopupItem onBack={() => changeItem(undefined)}>
          {currentItem.icon}
          <Title>{currentItem.title}</Title>
        </PopupItem>
      ) : (
        menuItems.map(i => (
          <PopupItem key={i.key} onClick={() => handleClick(i)}>
            {i.icon}
            <Title>{i.title}</Title>
          </PopupItem>
        ))
      )}
    </Wrapper>
  );
};

export default Menu;

const Wrapper = styled.div`
  height: 144px;
  width: 100%;
`;
const Title = styled.p`
  margin: 0;
`;
