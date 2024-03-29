import { postMsg } from "@web/extensions/sidebar/utils";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useEffect, useMemo } from "react";

import useHooks from "./hooks";

export type Tab = "catalog" | "selection" | "menu";

export type Props = {
  className?: string;
};

const MobileSidebar: React.FC<Props> = ({ className }) => {
  const {
    selected,
    project,
    reearthURL,
    backendURL,
    catalogURL,
    catalogProjectName,
    backendProjectName,
    inEditor,
    hideFeedback,
    templates,
    searchTerm,
    isCustomProject,
    customReearthURL,
    customCatalogURL,
    customCatalogProjectName,
    customBackendURL,
    customBackendProjectName,
    customProjectName,
    customLogo,
    setSelected,
  } = useHooks();

  const handleTabSelect = useCallback(
    (tab: Tab) => {
      if (selected === tab) {
        setSelected(undefined);
        postMsg({ action: "popupClose" });
      } else if (selected) {
        setSelected(tab);
        postMsg({ action: "msgToPopup", payload: { selected: tab } });
      } else {
        setSelected(tab);
        postMsg({ action: "mobileDropdownOpen" });
      }
    },
    [selected, setSelected],
  );

  useEffect(() => {
    const eventListenerCallback = (e: any) => {
      if (e.source !== parent) return;
      if (e.data.action === "initPopup") {
        postMsg({
          action: "msgToPopup",
          payload: {
            selected,
            templates,
            project,
            inEditor,
            hideFeedback,
            searchTerm,
            catalogProjectName,
            catalogURL,
            reearthURL,
            backendURL,
            backendProjectName,
            isCustomProject,
            customReearthURL,
            customCatalogURL,
            customCatalogProjectName,
            customBackendURL,
            customBackendProjectName,
          },
        });
      } else if (e.data.action === "triggerCatalogOpen") {
        postMsg({ action: "modalClose" });
        handleTabSelect("catalog");
      } else if (e.data.action === "msgFromPopup") {
        setSelected(e.data.payload);
      } else if (e.data.action === "popupClose") {
        setSelected(undefined);
      }
    };
    (globalThis as any).addEventListener("message", eventListenerCallback);
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

  const useCustomProjectHeader = useMemo(() => {
    return customProjectName || customLogo;
  }, [customProjectName, customLogo]);

  return (
    <Wrapper className={className}>
      {useCustomProjectHeader ? (
        <CustomProjectHeader>
          {customLogo && <CustomLogo src={customLogo} />}
          {customProjectName && <CustomProjectName>{customProjectName}</CustomProjectName>}
        </CustomProjectHeader>
      ) : (
        <PlateauIcon icon="plateauLogo" size={114} wide />
      )}
      <IconGroup>
        <StyledIcon
          icon="database"
          selected={selected === "catalog"}
          onClick={() => handleTabSelect("catalog")}
        />
        <StyledIcon
          icon="visible"
          selected={selected === "selection"}
          onClick={() => handleTabSelect("selection")}
        />
        <StyledIcon
          icon="menu"
          selected={selected === "menu"}
          onClick={() => handleTabSelect("menu")}
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
  width: 100%;
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
  background: ${({ selected }) => (selected ? "var(--theme-color)" : "transparent")};
  color: ${({ selected }) => (selected ? "white" : "var(--theme-color)")};
  padding: 4px;
  cursor: pointer;
  transition: background 0.3s, color 0.3s;

  :hover {
    background: var(--theme-color);
    color: white;
  }
`;

const CustomProjectHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 8px;
`;

const CustomLogo = styled("img")`
  max-height: 32px;
`;

const CustomProjectName = styled.div`
  font-weight: 500;
  font-size: 14px;
  line-height: 21px;
  color: #000;
  max-height: 42px;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
`;
