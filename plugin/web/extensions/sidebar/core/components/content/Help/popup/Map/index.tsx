import { styled } from "@web/theme";

const Map: React.FC = () => {
  return (
    <Wrapper>
      <Title>日影機能について</Title>
      <Paragraph>
        PLATEAU
        VIEWでは建物モデルなどが見やすいように、デフォルトでは日影の効果を「なし」に設定しています。
        日影の効果をONにするための手順は以下の通りです。
      </Paragraph>
      <ParagraphItem>
        <span />
        <Paragraph>
          建物モデルをVIEWに追加した時に、左側の一覧に表示される「影」というメニューのプルダウンから「投影のみ」、「受光のみ」、「投影と受光」のいずれかを選択します。
        </Paragraph>
      </ParagraphItem>
      <ParagraphItem>
        <span />
        <Paragraph>
          この効果と合わせて、マップ右上の「Map
          Settings」の中にある「タイムライン」の「常に表示」にチェックを入れると画面下部に時間をコントロールするスケールが表示されます（③）。
          初期状態では、現在時刻で日影が表示されます。左側の再生ボタンを押すことで時間を進めることができます。
          また、つまみを動かせば任意の時間で日影を表示することもできます。
        </Paragraph>
      </ParagraphItem>
      <Paragraph>
        PLATEAU
        VIEWではSimonら（1994）の手法で推定された太陽の位置を用いて日影を計算し、表現しています。タイムバーが出ていないときは、地図が見やすいように仮想の光源から光を当てて表示しています。
        Simon, J., Bretagnon, P., Chapront, J., & Chapront-Touze, M. (1994). Numerical expressions
        for precession formulae and mean elements for the Moon and the planets. Astronomy and
        Astrophysics (Berlin. Print), 282(2), 663–683.
      </Paragraph>
    </Wrapper>
  );
};

export default Map;

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
