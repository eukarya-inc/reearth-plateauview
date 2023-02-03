import useHooks from "@web/extensions/sidebar/core/components/hooks";
import { postMsg } from "@web/extensions/sidebar/utils";
import { styled } from "@web/theme";
import { useEffect, useMemo, useState } from "react";

import { Tab } from "../../core/components/Mobile";

import Catalog from "./Catalog";
import Menu from "./Menu";
import Selection from "./Selection";

type Props = {
  isMobile?: boolean;
};

const MobileDropdown: React.FC<Props> = ({ isMobile }) => {
  const [currentTab, setCurrentTab] = useState<Tab>("catalog");

  const {
    rawCatalog,
    project,
    reearthURL,
    backendURL,
    handleDatasetAdd,
    handleDatasetSave,
    handleDatasetUpdate,
    handleProjectDatasetRemove,
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
    (globalThis as any).addEventListener("message", eventListenerCallback);
    return () => {
      (globalThis as any).removeEventListener("message", eventListenerCallback);
    };
  });

  const addedDatasetIds = useMemo(
    () => project.selectedDatasets.map(dataset => dataset.id),
    [project.selectedDatasets],
  );

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
              onDatasetSave={handleDatasetSave}
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
