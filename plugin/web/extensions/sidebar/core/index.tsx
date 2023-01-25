import Feedback from "@web/extensions/sidebar/core/components/content/Feedback";
import Help from "@web/extensions/sidebar/core/components/content/Help";
import MapSettings from "@web/extensions/sidebar/core/components/content/MapSettings";
import Selection from "@web/extensions/sidebar/core/components/content/Selection";
import Share from "@web/extensions/sidebar/core/components/content/Share";
import Templates from "@web/extensions/sidebar/core/components/content/Templates";
import Header, { Pages } from "@web/extensions/sidebar/core/components/Header";
import useHooks from "@web/extensions/sidebar/core/hooks";
import { Content } from "@web/sharedComponents";
import { styled, commonStyles } from "@web/theme";
import { memo, useCallback, useState } from "react";

import MobileSidebar from "./Mobile";

export type Props = {
  className?: string;
};

const Sidebar: React.FC<Props> = ({ className }) => {
  const {
    isMobile,
    processedSelectedDatasets,
    overrides,
    minimized,
    inEditor,
    reearthURL,
    backendURL,
    templates,
    handleTemplateAdd,
    handleTemplateUpdate,
    handleTemplateRemove,
    setMinimize,
    handleDatasetRemove,
    handleDatasetRemoveAll,
    handleOverridesUpdate,
    handleModalOpen,
  } = useHooks();

  const [current, setCurrent] = useState<Pages>("data");

  const handleClick = useCallback((p: Pages) => {
    setCurrent(p);
  }, []);

  const handleMinimize = useCallback(() => {
    const html = document.querySelector("html");
    const body = document.querySelector("body");
    const root = document.getElementById("root");
    if (!minimized) {
      html?.classList.add("minimized");
      body?.classList.add("minimized");
      root?.classList.add("minimized");
    } else {
      html?.classList.remove("minimized");
      body?.classList.remove("minimized");
      root?.classList.remove("minimized");
    }
    setMinimize(!minimized);
  }, [minimized, setMinimize]);

  return isMobile ? (
    <MobileSidebar />
  ) : (
    <Wrapper className={className} minimized={minimized}>
      <Header
        current={current}
        inEditor={inEditor}
        minimized={minimized}
        onMinimize={handleMinimize}
        onClick={handleClick}
      />
      {!minimized && (
        <ContentWrapper className={className}>
          {
            {
              data: (
                <Selection
                  inEditor={inEditor}
                  selectedDatasets={processedSelectedDatasets}
                  onDatasetRemove={handleDatasetRemove}
                  onDatasetRemoveAll={handleDatasetRemoveAll}
                  onModalOpen={handleModalOpen}
                />
              ),
              map: <MapSettings overrides={overrides} onOverridesUpdate={handleOverridesUpdate} />,
              share: (
                <Share overrides={overrides} reearthURL={reearthURL} backendURL={backendURL} />
              ),
              help: <Help />,
              feedback: <Feedback backendURL={backendURL} />,
              template: (
                <Templates
                  templates={templates}
                  onTemplateAdd={handleTemplateAdd}
                  onTemplateUpdate={handleTemplateUpdate}
                  onTemplateRemove={handleTemplateRemove}
                />
              ),
            }[current]
          }
        </ContentWrapper>
      )}
    </Wrapper>
  );
};

export default memo(Sidebar);

const Wrapper = styled.div<{ minimized?: boolean }>`
  display: flex;
  flex-direction: column;
  ${commonStyles.mainWrapper}
  transition: height 0.5s, width 0.5s, border-radius 0.5s;
  ${({ minimized }) => minimized && commonStyles.minimizedWrapper}
`;

const ContentWrapper = styled(Content)`
  flex: 1;
  background: #e7e7e7;
  box-sizing: border-box;
  overflow: auto;
`;
