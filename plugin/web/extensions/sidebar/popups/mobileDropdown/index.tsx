import { postMsg } from "@web/extensions/sidebar/utils";
import { styled } from "@web/theme";
import { useEffect, useState } from "react";

import { Tab } from "../../core/components/Mobile";

import Catalog from "./Catalog";
import useHooks from "./hooks";
import Menu from "./Menu";
import Selection from "./Selection";

type Props = {
  isMobile?: boolean;
};

const MobileDropdown: React.FC<Props> = ({ isMobile }) => {
  const [currentTab, setCurrentTab] = useState<Tab>();

  const {
    addedDatasetIds,
    rawCatalog,
    project,
    reearthURL,
    backendURL,
    handleDatasetAdd,
    handleProjectDatasetRemove,
    handleDatasetUpdate,
    handleDatasetRemoveAll,
    handleProjectSceneUpdate,
  } = useHooks();

  useEffect(() => {
    postMsg({ action: "initPopup" });
  }, []);

  useEffect(() => {
    const eventListenerCallback = (e: any) => {
      if (e.source !== parent) return null;
      if (e.data.action) {
        if (e.data.action === "msgToPopup" && e.data.payload) {
          setCurrentTab(e.data.payload);
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
          catalog: (
            <Catalog
              addedDatasetIds={addedDatasetIds}
              isMobile={isMobile}
              rawCatalog={rawCatalog}
              onDatasetAdd={handleDatasetAdd}
            />
          ),
          selection: (
            <Selection
              selectedDatasets={project.selectedDatasets}
              onDatasetUpdate={handleDatasetUpdate}
              onDatasetRemove={handleProjectDatasetRemove}
              onDatasetRemoveAll={handleDatasetRemoveAll}
            />
          ),
          menu: (
            <Menu
              project={project}
              backendURL={backendURL}
              reearthURL={reearthURL}
              onProjectSceneUpdate={handleProjectSceneUpdate}
            />
          ),
        }[currentTab]}
    </Wrapper>
  );
};

export default MobileDropdown;

const Wrapper = styled.div`
  width: 100%;
  height: 100%;
`;
