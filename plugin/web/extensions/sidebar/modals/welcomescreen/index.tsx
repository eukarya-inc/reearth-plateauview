import welcomeScreenVideo from "@web/extensions/sidebar/core/assets/welcomeScreenVideo.png";
import useHooks from "@web/extensions/sidebar/modals/welcomescreen/hooks";
import { Checkbox, Icon } from "@web/sharedComponents";
import Video from "@web/sharedComponents/Video";
import { styled } from "@web/theme";

const WelcomeScreen: React.FC = () => {
  const {
    ShowVideo,
    handleDontShowAgain,
    dontShowAgain,
    handleShowVideo,
    handleClose,
    handleCloseVideo,
  } = useHooks();

  return (
    <Wrapper>
      {!ShowVideo ? (
        <>
          <CloseButton>
            <Icon size={32} icon="close" onClick={handleClose} />
          </CloseButton>
          <InnerWrapper>
            <TryMapWrapper>
              <Text weight={700} size={48}>
                ようこそ
              </Text>
              <ImgSection>
                <Text weight={500} size={20}>
                  マップを使ってみる
                </Text>
                <ImgWrapper>
                  <img src={welcomeScreenVideo} onClick={handleShowVideo} />
                </ImgWrapper>
              </ImgSection>
            </TryMapWrapper>
            <BtnsWrapper>
              <ButtonWrapper>
                <Text weight={500} size={14}>
                  ヘルプをみる
                </Text>
              </ButtonWrapper>
              <ButtonWrapper>
                <Icon size={20} icon="plusCircle" color="#fafafa" />
                <Text weight={500} size={14}>
                  カタログから検索する
                </Text>
              </ButtonWrapper>
            </BtnsWrapper>
          </InnerWrapper>
          <CheckWrapper>
            <Checkbox checked={dontShowAgain} onClick={handleDontShowAgain}>
              <Text weight={700} size={14}>
                閉じて今後は表示しない
              </Text>
            </Checkbox>
          </CheckWrapper>
        </>
      ) : (
        <>
          <CloseButton>
            <Icon size={32} icon="close" onClick={handleCloseVideo} />
          </CloseButton>
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

const InnerWrapper = styled.div`
  display: flex;
  flex-direction: row;
  align-items: flex-end;
  padding: 0px;
  gap: 119px;
  width: 742px;
  height: 316px;
`;

const Text = styled.p<{ weight: number; size: number }>`
  font-weight: ${({ weight }) => weight}px;
  font-size: ${({ size }) => size}px;
  margin-bottom: 0px;
  color: #fafafa;
`;

const TryMapWrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: flex-end;
  padding: 0px;
  gap: 50px;
  width: 305px;
  height: 316px;
`;

const BtnsWrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 24px;
  width: 318px;
  height: 154px;
`;

const ImgSection = styled.div`
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 0px;
  gap: 24px;
`;

const ImgWrapper = styled.div`
  width: 305px;
  height: 159px;
  cursor: pointer;
`;

const CloseButton = styled.button`
  display: flex;
  justify-content: center;
  align-items: center;
  position: absolute;
  right: 0;
  top: 0;
  height: 32px;
  width: 32px;
  border: none;
  background: #00bebe;
  color: white;
  cursor: pointer;
`;

const ButtonWrapper = styled.div<{ selected?: boolean }>`
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  padding: 10px;
  width: 318px;
  height: 41px;
  background: ${({ selected }) => (selected ? "#d1d1d1" : "#00bebe")};
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
  padding-top: 40px;
  gap: 8px;
  width: 192px;
  height: 46px;
`;

const VideoWrapper = styled.div`
  width: 1142px;
  height: 543px;
`;
