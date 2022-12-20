import { styled } from "@web/theme";

const Clip: React.FC = () => {
  return (
    <Wrapper>
      <Title>マップを使ってみる</Title>
      <Paragraph>
        上の動画では、空間データや地図を扱うのが初めての方が、データを追加したり地図上に表現するために必要な、基本的な機能について紹介しています。
        動画を見る時間がない方は、次の手順をお試しください。
      </Paragraph>
      <ParagraphItem>
        <span />
        <Paragraph>Add Data から使用可能なデータを表示して、マップに追加してみましょう。</Paragraph>
      </ParagraphItem>
      <ParagraphItem>
        <span />
        <Paragraph>画面の左側に現れるボタンで表示/非表示ボタンを切り替えてみましょう。</Paragraph>
      </ParagraphItem>
      <ParagraphItem>
        <span />
        <Paragraph>
          マップ上でデータをクリックして、より詳細な情報や元データについての情報を見てみましょう。
        </Paragraph>
      </ParagraphItem>
      <ParagraphItem>
        <span />
        <Paragraph>Map Settingsボタンをクリックして、背景図を変更してみましょう。</Paragraph>
      </ParagraphItem>
      <ParagraphItem>
        <span />
        <Paragraph>画面右手のズームや回転ボタンを使って、視点を変化させてみましょう。</Paragraph>
      </ParagraphItem>
    </Wrapper>
  );
};

export default Clip;
const Wrapper = styled.div`
  width: 318px;
  height: 457px;
`;
const Title = styled.p`
  margin: 0;
  font-size: 16px;
  color: inherit;
`;

const Paragraph = styled.p`
  font-size: 14px;
  line-height: 22px;
`;
const ParagraphItem = styled.div`
  display: flex;
  flex-direction: row;
  align-items: flex-start;
  padding: 0px;
  gap: 8px;

  width: 301px;
  height: 44px;
`;
