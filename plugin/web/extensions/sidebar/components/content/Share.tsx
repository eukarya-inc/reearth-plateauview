import { Input, Row, Icon } from "@web/extensions/sharedComponents";
import CommonPage from "@web/extensions/sidebar/components/content/CommonPage";
import { usePublishUrl } from "@web/extensions/sidebar/state";
import { postMsg } from "@web/extensions/sidebar/utils";
import { styled } from "@web/theme";
import { memo, useCallback, useEffect } from "react";

function makeUrlSuffix() {
  let result = "";
  const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  const charactersLength = characters.length;
  for (let i = 0; i < 15; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}

function generatePrintView(payload?: string) {
  const doc = window.open()?.document;

  if (!doc || !payload) return;

  const iframeHTML = `
    <div style="display: flex; flex-direction: column; max-width: 1200px; height: 100%; margin: 0 auto; padding: 20px;">
      <div style="display: flex; justify-content: right; align-items: center; gap: 8px; height: 60px;">
        <button style="padding: 8px; border: none; border-radius: 4px; background: #00BEBE; color: white; cursor: pointer;">Download map</button>
        <button style="padding: 9px; border: none; border-radius: 4px; background: #00BEBE; color: white; cursor: pointer;">Print</button>
      </div>
      <div style="display: flex; justify-content: center; width: 100%;">
        <img src="${payload}" style="max-width: 100%; object-fit: contain;" />
      </div>
      <div>
        <p>This map was created using https://plateauview.mlit.go.jp on ${new Date()}</p>
      </div>
    </div>
`;

  const css = `html,body{ margin: 0; }`;

  const styleTag = doc.createElement("style");
  styleTag.appendChild(document.createTextNode(css));
  styleTag.setAttribute("type", "text/css");
  doc.head.appendChild(styleTag);

  const iframe = doc.createElement("iframe");
  iframe.style.width = "100%";
  iframe.style.height = "100%";
  iframe.style.border = "none";

  doc.body.appendChild(iframe);

  iframe.contentWindow?.document.open();
  iframe.contentWindow?.document.write(iframeHTML);

  const iframeHtmlStyle = iframe.contentWindow?.document.createElement("style");
  if (iframeHtmlStyle) {
    iframeHtmlStyle.appendChild(document.createTextNode(css));
    iframeHtmlStyle.setAttribute("type", "text/css");
    iframe.contentWindow?.document.head.appendChild(iframeHtmlStyle);
  }
  iframe.contentWindow?.document.close();
  console.log(iframe, "iframe");
  return iframe;
}

const Share: React.FC = () => {
  const [publishUrl, setPublishUrl] = usePublishUrl();

  useEffect(() => {
    const suffix = makeUrlSuffix();
    if (!publishUrl) {
      // To do: get available url from PLATEAU backend and set that here,
      // OR need to add a regen if we get an error on the copy/publish button
      setPublishUrl(`https://plateauview.mlit.go.jp/${suffix}`);
    }
  }, [publishUrl, setPublishUrl]);

  // TO DO: handle screenshot (download/show print view) functionality
  const handleScreenshotShow = useCallback(() => {
    postMsg({ action: "screenshot" });
  }, []);

  const handleScreenshotSave = useCallback(() => {
    postMsg({ action: "screenshot-save" });
  }, []);

  addEventListener("message", e => {
    if (e.source !== parent) return;
    if (e.data.type) {
      if (e.data.type === "screenshot") {
        console.log(e.data, "dataaaaaaaaaa");
        generatePrintView(e.data.payload);
      } else if (e.data.type === "screenshot-save") {
        console.log(e.data, "dataaaaaaaaaa saveeeeeee");
        const link = document.createElement("a");
        link.download = "screenshot.png";
        link.href = e.data.payload;
        link.click();
        link.remove();
      }
    }
  });

  return (
    <CommonPage title="共有・印刷">
      <>
        <Subtitle>URLで共有</Subtitle>
        <InputGroup>
          <FlexWrapper>
            <Input value={publishUrl} />
            <StyledButton>
              <Icon icon="copy" />
            </StyledButton>
          </FlexWrapper>
          <SubText>このURLを使えば誰でもこのマップにアクセスできます。</SubText>
        </InputGroup>
        <Subtitle>HTMLページへの埋め込みは下記のコードをお使いください：</Subtitle>
        <InputGroup>
          <FlexWrapper>
            <Input value={`<iframe src=${publishUrl} />`} />
            <StyledButton>
              <Icon icon="copy" />
            </StyledButton>
          </FlexWrapper>
          <SubText>このURLを使えば誰でもこのマップにアクセスできます。</SubText>
        </InputGroup>
      </>
      <>
        <Subtitle>印刷</Subtitle>
        <SectionWrapper>
          <ButtonWrapper>
            <Button onClick={handleScreenshotSave}>Download map (png)</Button>
            <Button onClick={handleScreenshotShow}>Show Print View</Button>
          </ButtonWrapper>
          <SubText>このマップを印刷できる状態で表示</SubText>
        </SectionWrapper>
      </>
    </CommonPage>
  );
};

export default memo(Share);

const Text = styled.p`
  font-size: 14px;
  margin: 0;
`;

const Subtitle = styled(Text)`
  margin-bottom: 15px;
`;

const SubText = styled.p`
  font-size: 12px;
  color: #b1b1b1;
  margin: 8px 0 16px;
`;

const SectionWrapper = styled(Row)`
  display: flex;
  flex-direction: column;
  align-items: flex-start;
`;

const InputGroup = styled(Input.Group)`
  display: flex;
  flex-direction: row;
  align-items: flex-start;
  flex-wrap: wrap;
  width: 100%;
`;

const FlexWrapper = styled.div`
  display: flex;
  width: 100%;
`;

const ButtonWrapper = styled(FlexWrapper)`
  gap: 8px;
`;

const Button = styled.button`
  height: 37px;
  width: 160px;
  border: none;
  border-radius: 3px;
  background: #ffffff;
  font-size: 14px;
  line-height: 21px;
  cursor: pointer;
`;

const StyledButton = styled.button`
  background: #00bebe;
  border: none;
  border-radius: 2px;
  width: 40px;
  cursor: pointer;

  :hover {
    background: #00bebe;
    border-color: #00bebe;
    color: white;
  }
`;
