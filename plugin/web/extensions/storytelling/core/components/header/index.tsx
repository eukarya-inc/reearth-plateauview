import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";

import type { Mode } from "../../hooks";

import Tab from "./Tab";

type Props = {
  mode: Mode;
  editable: boolean;
  shareable: boolean;
  setMode: (m: Mode) => void;
  shareStory: () => void;
  clearStory: () => void;
  handleMinimize: () => void;
};

const Header: React.FC<Props> = ({
  mode,
  editable,
  shareable,
  setMode,
  shareStory,
  clearStory,
  handleMinimize,
}) => {
  return (
    <StyledHeader>
      <HeaderMain>
        <WidgetTitle>Story</WidgetTitle>
        {editable && (
          <Tab
            mode="editor"
            icon="pencil"
            text="Editor mode"
            currentMode={mode}
            onClick={setMode}
          />
        )}
        <Tab
          mode="player"
          icon="play"
          text="Play mode"
          theme="grey"
          currentMode={mode}
          onClick={setMode}
        />
      </HeaderMain>
      <HeaderBtns>
        {mode === "editor" && (
          <IconBtn onClick={clearStory}>
            <Icon icon="eraser" size={24} />
          </IconBtn>
        )}
        {shareable && (
          <IconBtn onClick={shareStory}>
            <Icon icon="paperPlane" size={24} />
          </IconBtn>
        )}
        <IconBtn onClick={handleMinimize}>
          <Icon icon="cross" size={24} />
        </IconBtn>
      </HeaderBtns>
    </StyledHeader>
  );
};

const StyledHeader = styled.div`
  height: 40px;
  background: #dfdfdf;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
`;

const HeaderMain = styled.div`
  display: flex;
  gap: 10px;
  height: 100%;
`;
const HeaderBtns = styled.div`
  display: flex;
  gap: 2px;
  height: 100%;
`;

const WidgetTitle = styled.div`
  color: #4a4a4a;
  font-weight: 700;
  font-size: 14px;
  line-height: 20px;
  padding: 10px 12px;
`;

const IconBtn = styled.div`
  display: flex;
  width: 40px;
  height: 40px;
  align-items: center;
  justify-content: center;
  background: var(--theme-color);
  color: #fff;
  cursor: pointer;
`;

export default Header;
