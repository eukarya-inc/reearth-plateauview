import welcomeScreenVideo from "@web/extensions/sidebar/core/assets/welcomeScreenVideo.png";
import useHooks from "@web/extensions/sidebar/modals/welcomescreen/hooks";
import { Checkbox, Icon } from "@web/sharedComponents";
import Video from "@web/sharedComponents/Video";
import { styled } from "@web/theme";

const WelcomeScreen: React.FC = () => {
  const {
    isMobile,
    showVideo,
    dontShowAgain,
    handleDontShowAgain,
    handleShowVideo,
    handleCloseVideo,
    handleClose,
    handleOpenHelp,
    handleOpenCatalog,
  } = useHooks();

  return (
    <Wrapper>
      {!showVideo ? (
        <>
          <CloseButton size={40} icon="close" onClick={handleClose} />
          <InnerWrapper isMobile={isMobile}>
            <TextWrapper isMobile={isMobile}>
              <Title weight={700} size={isMobile ? 24 : 48}>
                ようこそ
              </Title>
              <Text weight={500} size={isMobile ? 16 : 20}>
                {isMobile ? "データがお好きですか？" : "マップを使ってみる"}
              </Text>
            </TextWrapper>
            <ContentWrapper isMobile={isMobile}>
              {!isMobile && (
                <ImgWrapper>
                  <img src={welcomeScreenVideo} onClick={handleShowVideo} />
                </ImgWrapper>
              )}
              <BtnsWrapper isMobile={isMobile}>
                {!isMobile && (
                  <ButtonWrapper onClick={handleOpenHelp}>
                    <Text weight={500} size={14}>
                      ヘルプをみる
                    </Text>
                  </ButtonWrapper>
                )}
                <ButtonWrapper onClick={handleOpenCatalog}>
                  <Icon size={20} icon="plusCircle" color="#fafafa" />
                  <Text weight={500} size={14}>
                    カタログから検索する
                  </Text>
                </ButtonWrapper>
              </BtnsWrapper>
            </ContentWrapper>
            <CheckWrapper>
              <Checkbox checked={dontShowAgain} onClick={handleDontShowAgain} />
              <Text weight={700} size={14}>
                閉じて今後は表示しない
              </Text>
            </CheckWrapper>
          </InnerWrapper>
        </>
      ) : (
        <>
          <CloseButton size={40} icon="close" onClick={handleCloseVideo} />
          <VideoWrapper>
            <Video width=" 1142" height="543" src="https://www.youtube.com/embed/pY2dM-eG5mA" />
          </VideoWrapper>
        </>
      )}
    </Wrapper>
  );
};

export default WelcomeScreen;

const Wrapper = styled.div`
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: relative;
  background: rgba(0, 0, 0, 0.7);
`;

const InnerWrapper = styled.div<{ isMobile?: boolean }>`
  display: flex;
  flex-direction: column;
  width: ${({ isMobile }) => (isMobile ? "318px" : "742px")};
`;

const Text = styled.p<{ weight: number; size: number }>`
  font-weight: ${({ weight }) => weight}px;
  font-size: ${({ size }) => size}px;
  margin: 0px;
  color: #fafafa;
`;

const Title = styled(Text)<{ isMobile?: boolean }>`
  ${({ isMobile }) => !isMobile && `margin-bottom: 24px;`}
`;

const TextWrapper = styled.div<{ isMobile?: boolean }>`
  display: flex;
  flex-direction: column;
  align-items: ${({ isMobile }) => (isMobile ? "center" : "flex-start")};
  justify-content: flex-end;
  margin-bottom: 24px;
`;

const ContentWrapper = styled.div<{ isMobile?: boolean }>`
  display: flex;
  ${({ isMobile }) =>
    isMobile &&
    `
  flex-direction: column;
  align-items: center;
  `};
  justify-content: space-between;
`;

const BtnsWrapper = styled.div<{ isMobile?: boolean }>`
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 24px;
  width: ${({ isMobile }) => (isMobile ? "100%" : "318px")};
`;

const ImgWrapper = styled.div`
  width: 305px;
  height: 159px;
  cursor: pointer;
`;

const CloseButton = styled(Icon)`
  display: flex;
  justify-content: center;
  align-items: center;
  position: absolute;
  right: 0;
  top: 0;
  height: 48px;
  width: 48px;
  border: none;
  background: var(--theme-color);
  color: white;
  cursor: pointer;
`;

const ButtonWrapper = styled.div<{ selected?: boolean }>`
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  padding: 10px;
  width: 100%;
  height: 41px;
  background: ${({ selected }) => (selected ? "#d1d1d1" : "var(--theme-color)")};
  border-radius: 4px;
  border: none;
  gap: 8px;
  cursor: pointer;
  transition: background 0.3s;
  :hover {
    background: #d1d1d1;
  }
`;

const CheckWrapper = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  align-self: center;
  gap: 8px;
  margin-top: 50px;
`;

const VideoWrapper = styled.div`
  width: 1142px;
  height: 543px;
`;
