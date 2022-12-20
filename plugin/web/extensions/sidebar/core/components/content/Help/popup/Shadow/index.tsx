import { styled } from "@web/theme";

const Shadow: React.FC = () => {
  return (
    <Wrapper>
      <Title>クリップ機能について</Title>
      <Paragraph>
        この動画では、3Dモデルの断面を表示するための、クリップボックスの使い方を紹介しています。
        （クリップ機能の使い方）
      </Paragraph>
      <ParagraphItem>
        <span />
        <Paragraph>
          クリップ機能が使えるデータ（建物モデルやBIMデータ）を表示すると、左側の一覧の中に「クリップ機能」というチェックボックスが表示されます。
        </Paragraph>
      </ParagraphItem>
      <ParagraphItem>
        <span />
        <Paragraph>チェックボックスにチェックを入れると、機能を有効にできます。</Paragraph>
      </ParagraphItem>
      <ParagraphItem>
        <span />
        <Paragraph>機能を有効にすると、画面中央に立方体が表示されます。</Paragraph>
      </ParagraphItem>
      <ParagraphItem>
        <span />
        <Paragraph>
          この立方体を移動して、3Dモデルに重ねると、立方体の内側、あるいは外側のみ3Dモデルを表示することができます。内側、外側の切り替えにはプルダウンの「ボックス内をクリップ」、「ボックス外をクリップ」を選択します。
        </Paragraph>
      </ParagraphItem>
      <ParagraphItem>
        <span />
        <Paragraph>
          立方体のサイズや形は頂点の赤いポイント、あるいは各面の中央に配置された青いポイントを使って変更することができます。
        </Paragraph>
      </ParagraphItem>
      <ParagraphItem>
        <span />
        <Paragraph>
          「クリップボックスを表示」のチェックボックスはデフォルトではチェックが入っていますが、このチェックを外すと、最初に表示された立方体を非表示にすることができます。立方体は表示されませんが、クリップ機能は有効です。
        </Paragraph>
      </ParagraphItem>
      <ParagraphItem>
        <span />
        <Paragraph>
          「クリップボックスを地面にスナップ」のチェックボックスはデフォルトではチェックが入っていて、地下に潜らないようになっています。地下のオブジェクトに対してもクリップ機能を使いたい場合は、このチェックを外してください。
        </Paragraph>
      </ParagraphItem>
    </Wrapper>
  );
};

export default Shadow;
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
