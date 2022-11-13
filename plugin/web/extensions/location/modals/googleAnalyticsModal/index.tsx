import useHooks from "@web/extensions/location/hook";
import { Typography } from "@web/sharedComponents";
import { styled } from "@web/theme";

import CommonModalWrapper from "../commonModalWrapper";

const GoogleAnalyticstModal: React.FC = () => {
  const { handlegoogleModalChange } = useHooks();

  const { Link } = Typography;

  return (
    <CommonModalWrapper
      title="Google Analytics の利用について"
      onModalChange={handlegoogleModalChange}>
      <Paragraph>
        当サイトでは、サービス向上やウェブサイトの改善のためにGoogle
        Inc.の提供するアクセス分析のツールであるGoogle Analyticsを利用した計測を行っております。
        Google
        Analyticsは、当サイトが発行するCookieを利用して、個人を特定する情報を含まずにウェブサイトの利用データ（アクセス状況、トラフィック、閲覧環境など）を収集しております。
        Cookieの利用に関してはGoogleのプライバシーポリシーと規約に基づいております。
        取得したデータはウェブサイト利用状況の分析、サイト運営者へのレポートの作成、その他のサービスの提供に関わる目的に限り、これを使用します。
        Google Analyticsの利用規約及びプライバシーポリシーに関する説明については、Google
        Analyticsのサイトをご覧ください。
      </Paragraph>

      <Link href="https://marketingplatform.google.com/about/analytics/terms/jp/" target="_blank">
        Google Analytics利用規約{" "}
      </Link>
      <Link href="https://policies.google.com/privacy?hl=ja" target="_blank">
        Googleのプライバシーポリシー
      </Link>
      <Link href="https://marketingplatform.google.com/about/analytics/" target="_blank">
        Google Analyticsに関する詳細情報
      </Link>

      <Paragraph>
        Google Analytics オプトアウト アドオンを利用してGoogle
        Analyticsのトラッキングを拒否することも可能です。 Google Analytics オプトアウト
        アドオンは、JavaScriptによるデータの使用をウェブサイトのユーザーが無効にできるように開発された機能です。
        この機能を利用するには、このアドオンをダウンロードして、ご利用のブラウザにインストールしてください。
      </Paragraph>

      <Link href="https://tools.google.com/dlpage/gaoptout?hl=ja" target="_blank">
        Google Analytics オプトアウト アドオン
      </Link>
    </CommonModalWrapper>
  );
};
export default GoogleAnalyticstModal;
const Paragraph = styled.p`
  font-size: 12px;
`;
