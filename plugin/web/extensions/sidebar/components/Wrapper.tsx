import { Content } from "@web/extensions/sharedComponents/Layout";
import Info from "@web/extensions/sidebar/components/tabs/Info";
import MapSettings from "@web/extensions/sidebar/components/tabs/MapSettings";
import Share from "@web/extensions/sidebar/components/tabs/Share";
import { styled, commonStyles } from "@web/theme";
import { memo, ReactNode } from "react";

type Props = {
  className?: string;
  header?: ReactNode;
  footer?: ReactNode;
  current: string;
  minimized?: boolean;
};

const SidebarWrapper: React.FC<Props> = ({ className, header, footer, current, minimized }) => {
  return (
    <Wrapper minimized={minimized}>
      {header}
      {!minimized && (
        <>
          <ContentWrapper className={className}>
            {
              {
                shareNprint: <Share />,
                about: <Info />,
                mapSetting: <MapSettings />,
              }[current]
            }
          </ContentWrapper>
          {footer}
        </>
      )}
    </Wrapper>
  );
};
export default memo(SidebarWrapper);

const ContentWrapper = styled(Content)`
  width: 100%;
  background: #dcdcdc;
  flex: 1;
  padding: 20px 42px 20px 12px;
  box-sizing: border-box;
  overflow: scroll;
`;

const Wrapper = styled.div<{ minimized?: boolean }>`
  display: flex;
  flex-direction: column;
  ${commonStyles.mainWrapper}
  transition: height 0.5s, width 0.5s, border-radius 0.5s;
  ${({ minimized }) => minimized && commonStyles.minimizedWrapper}
`;
