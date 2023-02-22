import basicOperation from "@web/extensions/sidebar/core/assets/basicOperationImg.png";
import BasicOperButton from "@web/extensions/sidebar/core/assets/BasicOperButton.png";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";

const BasicOperation: React.FC = () => {
  return (
    <Wrapper>
      <Title>視点や画面移動</Title>
      <ContentWrapper>
        <img src={basicOperation} />
        <Title>データの追加</Title>
        <Paragraph>以下のボタンで建物やデータを地図に追加してください</Paragraph>
        <img src={BasicOperButton} />
        <Paragraph>
          カタログから検索するウインドウが表示されたら、
          <br /> ① 表示したいエリアに対応するフォルダをクリックして開く
          <br />
          ② 建物モデルや重ね合わせたいデータを
          <InlineIcon icon="plusCircle" size={16} />
          で選択する
        </Paragraph>
      </ContentWrapper>
    </Wrapper>
  );
};

export default BasicOperation;

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 0px 16px;
  gap: 24px;
  width: 333px;
  height: 1100px;
`;
const Title = styled.p`
  margin: 0;
  font-size: 16px;
  line-height: 24px;
  color: rgba(0, 0, 0, 0.85);
`;

const ContentWrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 0px;
  gap: 20px;
  width: 301px;
  height: 1052px;
`;

const Paragraph = styled.p`
  font-size: 14px;
  line-height: 22px;
  color: rgba(0, 0, 0, 0.45);
`;

const InlineIcon = styled(Icon)`
  display: inline-block;
  color: #00bebe;
`;
