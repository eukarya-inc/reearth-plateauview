import { Icon, Button, Col } from "@web/extensions/sharedComponents";
import { styled } from "@web/theme";
import { memo, ReactNode, useMemo } from "react";

export type Pages = "mapData" | "mapSetting" | "shareNprint" | "about" | "template";

export type TabProps = {
  key: Pages;
  icon: ReactNode;
};

type Props = {
  className?: string;
  current: string;
  isInsideEditor?: boolean;
  minimized?: boolean;
  onMinimize: () => void;
  onClick: (p: Pages) => void;
};

const plateauWebsiteUrl = "https://www.mlit.go.jp/plateau/";

const Header: React.FC<Props> = ({
  className,
  current,
  isInsideEditor,
  minimized,
  onMinimize,
  onClick,
}) => {
  const headerItems = useMemo(() => {
    const items: TabProps[] = [
      {
        key: "mapData",
        icon: <StyledIcon icon="dataBase" />,
      },
      {
        key: "mapSetting",
        icon: <StyledIcon icon="sliders" />,
      },
      {
        key: "shareNprint",
        icon: <StyledIcon icon="share" />,
      },
      {
        key: "about",
        icon: <StyledIcon icon="info" />,
      },
      {
        key: "template",
        icon: <StyledIcon icon="template" />,
      },
    ];

    return !isInsideEditor ? [...items.slice(0, -1)] : [...items];
  }, [isInsideEditor]);

  return (
    <HeaderWrapper className={className}>
      <TopSection minimized={minimized}>
        <PlateauIcon
          icon="plateauLogo"
          size={114}
          wide
          onClick={() => window.open(plateauWebsiteUrl, "_blank", "noopener")}
        />
        <MinimizeButton
          icon={<Icon icon={minimized ? "menu" : "close"} onClick={onMinimize} />}
          minimized={minimized}
        />
      </TopSection>
      {!minimized && (
        <Nav>
          {headerItems?.map((i, idx) => (
            <IconWrapper key={idx} current={current === i.key} onClick={() => onClick(i.key)}>
              {i.icon}
            </IconWrapper>
          ))}
        </Nav>
      )}
    </HeaderWrapper>
  );
};

export default memo(Header);

const HeaderWrapper = styled(Col)`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  height: 164px;
  width: 100%;
  background-color: #ffff;
  border-radius: 0 8px 8px 0;
  transition: height 0.5s, width 0.5s, border-radius 0.5s;
`;

const Nav = styled.div`
  display: flex;
  justify-content: space-around;
  width: 100%;
  padding: 0 53px;
`;

const TopSection = styled.div<{ minimized?: boolean }>`
  display: flex;
  position: relative;
  width: 100%;
  height: 100%;
  padding: 24px;
  border-radius: 0 8px 8px 0;
  transition: all 1s;
`;

const MinimizeButton = styled(Button)<{ minimized?: boolean }>`
  position: absolute;
  right: 0;
  top: 0;
  border: none;
  border-radius: 0;
  height: 32px;
  width: 32px;
  background: #00bebe;
  color: white;

  :focus {
    background: #00bebe;
  }

  :hover,
  :active,
  :checked {
    background: #c7c5c5;
    color: white;
  }

  ${({ minimized }) => minimized && "position: static;"}
`;

const PlateauIcon = styled(Icon)<{ minimized?: boolean }>`
  cursor: pointer;
  width: 100%;
  margin: auto;
  text-align: center;

  ${({ minimized }) => minimized && "text-align: left;"}
`;

const IconWrapper = styled.div<{ current?: boolean }>`
  cursor: pointer;
  padding: 6px;
  border-bottom-style: solid;
  border-bottom-width: 1px;
  border-bottom-color: ${({ current }) => (current ? "#00bebe" : "transparent")};
  color: ${({ current }) => (current ? "#00bebe" : "#C7C5C5")};
  transition: border-bottom-color 0.5s, color 0.5s;

  :hover {
    border-bottom-color: #00bebe;
    color: #00bebe;
  }
`;

const StyledIcon = styled(Icon)`
  width: 100%;
`;
