import { postMsg } from "@web/extensions/sidebar/utils";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useEffect, useState } from "react";

import useHooks from "./hooks";

export type Tab = "catalog" | "selection" | "menu";

export type Props = {
  className?: string;
};

const MobileSidebar: React.FC<Props> = ({ className }) => {
  const [selected, setSelected] = useState<Tab | undefined>();

  const {
    project,
    catalog,
    // inEditor,
    // reearthURL,
    // backendURL,
    // backendProjectName,
    templates,
    // currentPage,
    // loading,
    buildingSearch,
    // handlePageChange,
    // handleDatasetSave,
    // handleProjectDatasetRemove,
    // handleDatasetUpdate,
    // handleProjectDatasetAdd,
    // handleProjectDatasetRemoveAll,
    // handleProjectDatasetsUpdate,
    // handleProjectSceneUpdate,
    // handleModalOpen,
    // handleBuildingSearch,
    // handleOverride,
  } = useHooks();

  // const {
  //   catalog,
  //   project,
  //   loading,
  //   templates,
  //   buildingSearch,
  //   reearthURL,
  //   backendURL,
  //   backendProjectName,
  //   searchTerm,
  //   handleSearch,
  //   handleDatasetSave,
  //   handleDatasetUpdate,
  //   handleProjectDatasetAdd,
  //   handleProjectDatasetRemove,
  //   handleProjectDatasetRemoveAll,
  //   handleProjectDatasetsUpdate,
  //   handleProjectSceneUpdate,
  //   handleBuildingSearch,
  // } = useHooks();

  const handleTabSelect = useCallback(
    (tab: Tab) => {
      if (selected === tab) {
        setSelected(undefined);
        postMsg({ action: "popupClose" });
      } else if (selected) {
        setSelected(tab);
        postMsg({ action: "msgToPopup", payload: tab });
      } else {
        setSelected(tab);
        postMsg({ action: "mobileDropdownOpen" });
      }
    },
    [selected],
  );

  useEffect(() => {
    const eventListenerCallback = (e: any) => {
      if (e.source !== parent) return;
      if (e.data.action === "initPopup") {
        postMsg({
          action: "msgToPopup",
          payload: {
            selected,
            project,
            templates,
            catalog,
            buildingSearch,
          },
        });
      } else if (e.data.action === "triggerCatalogOpen") {
        postMsg({ action: "modalClose" });
        handleTabSelect("catalog");
      } else if (e.data.action === "msgFromPopup") {
        setSelected(e.data.payload);
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

  return (
    <Wrapper className={className}>
      <PlateauIcon icon="plateauLogo" size={114} wide />
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
