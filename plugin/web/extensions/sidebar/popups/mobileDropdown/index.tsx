import type { Tab } from "@web/extensions/sidebar/core/Mobile";
import { postMsg } from "@web/extensions/sidebar/utils";
import { styled } from "@web/theme";
import { useEffect, useState } from "react";

import Data from "./Data";
import Detail from "./Detail";
import Menu from "./Menu";

const MobileDropdown: React.FC = () => {
  const [currentTab, setCurrentTab] = useState<Tab>("menu");

  useEffect(() => {
    postMsg({ action: "initPopup" });
  }, []);

  useEffect(() => {
    const eventListenerCallback = (e: any) => {
      if (e.source !== parent) return null;
      if (e.data.type) {
        if (e.data.type === "msgToPopup" && e.data.message) {
          setCurrentTab(e.data.message);
        }
      }
    };
    (globalThis as any).addEventListener("message", (e: any) => eventListenerCallback(e));
    return () => {
      (globalThis as any).removeEventListener("message", eventListenerCallback);
    };
  });

  return (
    <Wrapper>
      {currentTab &&
        {
          data: <Data />,
          detail: <Detail />,
          menu: <Menu />,
        }[currentTab]}
    </Wrapper>
  );
};

export default MobileDropdown;

const Wrapper = styled.div`
  width: 100%;
  height: 100%;
`;
