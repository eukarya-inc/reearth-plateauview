import Feedback from "@web/extensions/sidebar/core/components/content/Feedback";
import MapSettings from "@web/extensions/sidebar/core/components/content/MapSettings";
import Share from "@web/extensions/sidebar/core/components/content/Share";
import { Project, ReearthApi } from "@web/extensions/sidebar/types";
import { postMsg } from "@web/extensions/sidebar/utils";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { ReactNode, useCallback, useEffect, useMemo, useState } from "react";

import PopupItem from "../sharedComponents/PopupItem";

export type Props = {
  hideFeedback: boolean;
  project: Project;
  reearthURL?: string;
  backendURL?: string;
  backendProjectName?: string;
  isCustomProject: boolean;
  customReearthURL?: string;
  customBackendURL?: string;
  customBackendProjectName?: string;
  onProjectSceneUpdate: (updatedProperties: Partial<ReearthApi>) => void;
};

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

const Menu: React.FC<Props> = ({
  hideFeedback,
  project,
  reearthURL,
  backendURL,
  backendProjectName,
  isCustomProject,
  customReearthURL,
  customBackendURL,
  customBackendProjectName,
  onProjectSceneUpdate,
}) => {
  const [currentItem, changeItem] = useState<MenuItem | undefined>();

  const menuItemsFiltered = useMemo(() => {
    return hideFeedback ? menuItems.filter(item => item.key !== "feedback") : menuItems;
  }, [hideFeedback]);

  const handleHeightUpdate = useCallback(
    (id: string) => {
      const el = document.getElementById(id);
      let currentHeight = el ? parseFloat(window.getComputedStyle(el).height) : undefined;
      if (currentItem && currentHeight) {
        currentHeight += 48;
      }
      postMsg({
        action: "msgFromPopup",
        payload: { height: currentHeight },
      });
    },
    [currentItem],
  );

  const handleClick = useCallback(
    (item: MenuItem) => {
      if (currentItem === item) {
        changeItem(undefined);
      } else {
        changeItem(item);
      }
    },
    [currentItem],
  );

  useEffect(() => {
    handleHeightUpdate(currentItem ? "content-area" : "menu");
  }, [currentItem, handleHeightUpdate]);

  return (
    <Wrapper id="menu">
      {currentItem ? (
        <>
          <PopupItem onBack={() => changeItem(undefined)}>
            {currentItem.icon}
            <Title>{currentItem.title}</Title>
          </PopupItem>
          <div id="content-area">
            {currentItem.key &&
              {
                map: (
                  <MapSettings
                    overrides={project.sceneOverrides}
                    isMobile
                    onOverridesUpdate={onProjectSceneUpdate}
                  />
                ),
                share: (
                  <Share
                    project={project}
                    reearthURL={reearthURL}
                    backendURL={backendURL}
                    backendProjectName={backendProjectName}
                    isCustomProject={isCustomProject}
                    customReearthURL={customReearthURL}
                    customBackendURL={customBackendURL}
                    customBackendProjectName={customBackendProjectName}
                    isMobile
                  />
                ),
                ...(hideFeedback
                  ? {}
                  : { feedback: <Feedback backendURL={backendURL} isMobile /> }),
              }[currentItem.key]}
          </div>
        </>
      ) : (
        menuItemsFiltered.map(i => (
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
  width: 100%;
  border-top: 1px solid #d9d9d9;
`;

const Title = styled.p`
  margin: 0;
`;
