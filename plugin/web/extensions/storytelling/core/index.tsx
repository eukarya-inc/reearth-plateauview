import Editor from "@web/extensions/storytelling/core/components/editor";
import Header from "@web/extensions/storytelling/core/components/header";
import Player from "@web/extensions/storytelling/core/components/player";
import useHooks, { sizes } from "@web/extensions/storytelling/core/hooks";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";

const Storytelling: React.FC = () => {
  const {
    size,
    mode,
    scenes,
    ConfigProvider,
    isMobile,
    handleMinimize,
    handleSetMode,
    captureScene,
    viewScene,
    recaptureScene,
    deleteScene,
    editScene,
    moveScene,
    clearStory,
    shareStory,
  } = useHooks();

  return (
    <ConfigProvider>
      <Wrapper size={size} mode={mode}>
        <MiniPane onClick={handleMinimize} size={size}>
          <Icon icon="cornersOut" color="#4A4A4A" size={24} />
          <MiniTitle>Story</MiniTitle>
        </MiniPane>
        <ContentPane size={size}>
          <Header
            mode={mode}
            setMode={handleSetMode}
            shareStory={shareStory}
            clearStory={clearStory}
            handleMinimize={handleMinimize}
            editable={!isMobile}
            shareable={!isMobile}
          />
          {!isMobile && mode === "editor" && (
            <Editor
              scenes={scenes}
              captureScene={captureScene}
              viewScene={viewScene}
              recaptureScene={recaptureScene}
              deleteScene={deleteScene}
              editScene={editScene}
              moveScene={moveScene}
            />
          )}
          {mode === "player" && <Player scenes={scenes} viewScene={viewScene} />}
        </ContentPane>
      </Wrapper>
    </ConfigProvider>
  );
};

const Wrapper = styled.div<{ size: keyof typeof sizes; mode?: string }>`
  position: relative;
  display: inline-block;
  border-radius: 8px;
  background: ${({ mode, size }) =>
    size === "mini" ? "#fff" : mode === "player" ? "#F4F4F4" : "#fff"};
  transition: min-width 0.5s, min-height 0.5s;
  min-width: ${({ size }) => (size === "mini" ? `${sizes.mini.width}px` : "100%")};
  min-height: ${({ size }) => `${sizes[size].height}px`};
  overflow: hidden;
`;

const MiniPane = styled.div<{ size: keyof typeof sizes }>`
  position: absolute;
  left: 0;
  top: 0;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  width: ${sizes.mini.width};
  cursor: pointer;
  pointer-events: ${({ size }) => (size === "mini" ? "all" : "none")};
  opacity: ${({ size }) => (size === "mini" ? 1 : 0)};
  transition: opacity 0.25s;
`;

const MiniTitle = styled.div`
  font-weight: 700;
  font-size: 14px;
  width: auto;
`;

const ContentPane = styled.div<{ size: keyof typeof sizes }>`
  position: absolute;
  width: 100%;
  height: 100%;
  left: 0;
  top: 0;
  display: flex;
  flex-direction: column;
  pointer-events: ${({ size }) => (size === "mini" ? "none" : "all")};
  opacity: ${({ size }) => (size === "mini" ? 0 : 1)};
  transition: opacity 0.25s;
`;

export default Storytelling;
