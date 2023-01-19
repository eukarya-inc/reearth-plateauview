import welcomeScreenVideo from "@web/extensions/sidebar/core/assets/welcomeScreenVideo.png";
// import Backdrop from "@web/extensions/sidebar/modals/welcomescreen/component/Backdrop";
import useHooks from "@web/extensions/sidebar/modals/welcomescreen/hooks";
import { Checkbox, Icon } from "@web/sharedComponents";
import Video from "@web/sharedComponents/Video";
import { styled } from "@web/theme";

const WelcomeScreen: React.FC = () => {
  const { ShowVideo, handleDontShowAgain, dontShowAgain, handleShowVideo, handleClose } =
    useHooks();

  return (
    <Wrapper>
      <CloseButton>
        <Icon size={32} icon="close" onClick={handleClose} />
      </CloseButton>
      {!ShowVideo ? (
        <>
          <InnerWrapper>
            <TryMapWrapper>
              <TextWrapper width={192} height={46}>
                <Text weight={700} size={48}>
                  ようこそ
                </Text>
              </TextWrapper>
              <ImgSection>
                <TextWrapper width={215} height={32}>
                  <Text weight={500} size={20}>
                    マップを使ってみる
                  </Text>
                </TextWrapper>
                <ImgWrapper>
                  <img src={welcomeScreenVideo} onClick={handleShowVideo} />
                </ImgWrapper>
              </ImgSection>
            </TryMapWrapper>
            <BtnsWrapper>
              <ButtonWrapper>
                <TextWrapper width={84} height={21}>
                  <Text weight={500} size={14}>
                    ヘルプをみる
                  </Text>
                </TextWrapper>
              </ButtonWrapper>
              <ButtonWrapper>
                <Icon size={20} icon="plusCircle" color="#fafafa" />
                <TextWrapper width={84} height={21}>
                  <Text weight={500} size={14}>
                    ヘルプをみる
                  </Text>
                </TextWrapper>
              </ButtonWrapper>
            </BtnsWrapper>
          </InnerWrapper>
          <CheckWrapper>
            <Checkbox checked={dontShowAgain} onClick={handleDontShowAgain}>
              <TextWrapper width={192} height={46}>
                <Text weight={700} size={14}>
                  閉じて今後は表示しない
                </Text>
              </TextWrapper>
            </Checkbox>
          </CheckWrapper>
        </>
      ) : (
        <VideoWrapper>
          <Video width=" 1142" height="543" src="https://www.youtube.com/embed/pY2dM-eG5mA" />
        </VideoWrapper>
      )}
      {/* <Backdrop onClose={handleClose} /> */}
    </Wrapper>
  );
};
export default WelcomeScreen;

const Wrapper = styled.div``;

const InnerWrapper = styled.div`
  display: flex;
  flex-direction: row;
  align-items: flex-end;
  padding: 0px;
  gap: 119px;
  position: absolute;
  width: 742px;
  height: 316px;
  left: calc(50% - 742px / 2 + 0.5px);
  top: calc(50% - 316px / 2 - 30.5px);
`;
const TextWrapper = styled.div<{ width: number; height: number }>`
  width: ${({ width }) => width};
  height: ${({ height }) => height};
  margin: 0px;
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
  padding: 0px;
  gap: 55px;
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
  width: 305px;
  height: 215px;
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
  padding: 0px;
  gap: 8px;
  position: absolute;
  width: 178px;
  height: 22px;
  left: calc(50% - 178px / 2 + 0.5px);
  top: calc(50% - 22px / 2 + 178.5px);
`;

const VideoWrapper = styled.div`
  position: absolute;
  width: 1142px;
  height: 543px;
  left: calc(50% - 1142px / 2 + 0.5px);
  top: calc(50% - 543px / 2);
`;
