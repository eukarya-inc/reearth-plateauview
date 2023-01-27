import { postMsg } from "@web/extensions/sidebar/utils";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useEffect, useState } from "react";

export type Tab = "data" | "detail" | "menu";

export type Props = {
  className?: string;
};

const MobileSidebar: React.FC<Props> = ({ className }) => {
  const [selected, setSelected] = useState<Tab | undefined>();

  const handleClick = useCallback(
    (tab: Tab) => {
      if (selected === tab) {
        setSelected(undefined);
      } else {
        setSelected(tab);
      }
      postMsg({ action: "mobileDropdownOpen" });
    },
    [selected],
  );

  useEffect(() => {
    if (selected) {
      postMsg({
        action: "msgToMobileDropdown",
        payload: selected,
      });
    } else {
      postMsg({ action: "popupClose" });
    }
  }, [selected]);

  useEffect(() => {
    const eventListenerCallback = (e: any) => {
      if (e.source !== parent) return;
      if (e.data.type && e.data.type === "initPopup") {
        postMsg({ action: "msgToPopup", payload: selected });
      }
    };
    (globalThis as any).addEventListener("message", (e: any) => eventListenerCallback(e));
    return () => {
      (globalThis as any).removeEventListener("message", eventListenerCallback);
    };
  });

  useEffect(() => {
    const html = document.querySelector("html");
    const body = document.querySelector("body");
    const root = document.getElementById("root");
    html?.classList.add("mobile");
    body?.classList.add("mobile");
    root?.classList.add("mobile");

    return () => {
      html?.classList.remove("mobile");
      body?.classList.remove("mobile");
      root?.classList.remove("mobile");
    };
  }, []);

  return (
    <Wrapper className={className}>
      <PlateauIcon icon="plateauLogo" size={114} wide />
      <IconGroup>
        <StyledIcon
          icon="database"
          selected={selected === "data"}
          onClick={() => handleClick("data")}
        />
        <StyledIcon
          icon="database"
          selected={selected === "detail"}
          onClick={() => handleClick("detail")}
        />
        <StyledIcon
          icon="menu"
          selected={selected === "menu"}
          onClick={() => handleClick("menu")}
        />
      </IconGroup>
    </Wrapper>
  );
};

export default MobileSidebar;

const Wrapper = styled.div`
  display: flex;
  justify-content: space-between;
  height: 56px;
  width: 100%; //NEED TO FIX WIDTH
  background: #f4f4f4;
  padding: 12px;
`;

const PlateauIcon = styled(Icon)`
  text-align: left;
`;

const IconGroup = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  height: 100%;
`;

const StyledIcon = styled(Icon)<{ selected?: boolean }>`
  background: ${({ selected }) => (selected ? "#00bebe" : "transparent")};
  color: ${({ selected }) => (selected ? "white" : "#00bebe")};
  padding: 4px;
  cursor: pointer;
  transition: background 0.3s, color 0.3s;

  :hover {
    background: #00bebe;
    color: white;
  }
`;
