import Editor from "@web/extensions/storytelling/core/components/editor";
import Header from "@web/extensions/storytelling/core/components/header";
import useHooks, { size } from "@web/extensions/storytelling/core/hooks";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";

const Storytelling: React.FC = () => {
  const { minimized, handleMinimize, mode, setMode, contentTheme } = useHooks();

  return (
    <Wrapper minimized={minimized} theme={contentTheme}>
      <MiniPane onClick={handleMinimize} minimized={minimized}>
        <Icon icon="cornersOut" color="#4A4A4A" size={24} />
        <MiniTitle>Story</MiniTitle>
      </MiniPane>
      <ContentPane minimized={minimized}>
        <Header mode={mode} setMode={setMode} handleMinimize={handleMinimize} />
        {mode === "editor" && <Editor />}
      </ContentPane>
    </Wrapper>
  );
};

const Wrapper = styled.div<{ minimized: boolean; theme?: string }>`
  position: relative;
  display: inline-block;
  border-radius: 8px;
  background: ${({ theme }) => (theme === "grey" ? "#F4F4F4" : "#fff")};
  transition: min-width 0.5s, min-height 0.5s;
  min-width: ${({ minimized }) => (minimized ? `${size.mini.width}px` : "100%")};
  min-height: ${({ minimized }) =>
    minimized ? `${size.mini.height}px` : `${size.extend.height}px`};
  overflow: hidden;
`;

const MiniPane = styled.div<{ minimized: boolean }>`
  position: absolute;
  left: 0;
  top: 0;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  width: ${size.mini.width};
  cursor: pointer;
  pointer-events: ${({ minimized }) => (minimized ? "all" : "none")};
  opacity: ${({ minimized }) => (minimized ? 1 : 0)};
  transition: opacity 0.25s;
`;

const MiniTitle = styled.div`
  font-family: "Noto Sans";
  font-weight: 700;
  font-size: 14px;
  width: auto;
`;

const ContentPane = styled.div<{ minimized: boolean }>`
  position: absolute;
  width: 100%;
  height: 100%;
  left: 0;
  top: 0;
  display: flex;
  flex-direction: column;
  pointer-events: ${({ minimized }) => (minimized ? "none" : "all")};
  opacity: ${({ minimized }) => (minimized ? 0 : 1)};
  transition: opacity 0.25s;
`;

export default Storytelling;
