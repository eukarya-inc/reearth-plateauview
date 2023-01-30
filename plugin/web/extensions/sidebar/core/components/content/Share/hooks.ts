import { Project as ProjectType } from "@web/extensions/sidebar/types";
import { postMsg } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useRef, useState } from "react";

export type Project = ProjectType;

export default ({
  project,
  reearthURL,
  backendURL,
  messageApi,
}: {
  project?: Project;
  reearthURL?: string;
  backendURL?: string;
  messageApi: any;
}) => {
  const [publishedUrl, setPublishedUrl] = useState<string>();
  const [shareDisabled, setShareDisable] = useState(false);
  const timer = useRef<NodeJS.Timeout | null>(null);

  const handleScreenshotShow = useCallback(() => {
    postMsg({ action: "screenshotPreview" });
  }, []);

  const handleScreenshotSave = useCallback(() => {
    postMsg({ action: "screenshotSave" });
  }, []);

  const handleProjectShare = useCallback(async () => {
    setShareDisable(true);
    if (project) {
      if (!backendURL || !reearthURL) return;
      const resp = await fetch(`${backendURL}/share`, {
        headers: {
          "Content-Type": "application/json",
        },
        method: "POST",
        body: JSON.stringify(project),
      });
      if (resp.status !== 200) {
        messageApi.open({
          type: "error",
          content: "サバーの問題です。しばらくお待ちしてもう一回して下さい",
        });
        if (timer.current) {
          clearTimeout(timer.current);
        }
      } else {
        const project = await resp.json();
        setPublishedUrl(`${reearthURL}${reearthURL.includes("?") ? "&" : "?"}projectID=${project}`);
      }
    }
    timer.current = setTimeout(() => {
      setShareDisable(false);
    }, 3000);
  }, [messageApi, reearthURL, backendURL, project, setPublishedUrl]);

  useEffect(() => {
    return () => {
      if (timer.current) {
        clearTimeout(timer.current);
      }
    };
  }, []);

  return {
    shareDisabled,
    publishedUrl,
    handleProjectShare,
    handleScreenshotShow,
    handleScreenshotSave,
  };
};

addEventListener("message", e => {
  if (e.source !== parent) return;
  if (e.data.type) {
    if (e.data.type === "screenshotPreview") {
      generatePrintView(e.data.payload);
    } else if (e.data.type === "screenshotSave") {
      const link = document.createElement("a");
      link.download = "screenshot.png";
      link.href = e.data.payload;
      link.click();
      link.remove();
    }
  }
});

function generatePrintView(payload?: string) {
  const doc = window.open()?.document;

  if (!doc || !payload) return;

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

  const iframeDoc = iframe.contentWindow?.document;
  if (!iframeDoc) return;

  const currentDate = new Date();
  const options: Intl.DateTimeFormatOptions = {
    weekday: "long",
    year: "numeric",
    month: "long",
    day: "numeric",
  };
  const localizedDate = currentDate.toLocaleDateString("ja", options);

  iframeDoc.open();

  const iframeHTML = `
  <html>
    <body>
      <div style="display: flex; flex-direction: column; max-width: 1200px; height: 100%; margin: 0 auto; padding: 20px;">
        <div style="display: flex; justify-content: right; align-items: center; gap: 8px; height: 60px;">
          <button onclick="downloadScreenshot()" style="padding: 8px; border: none; border-radius: 4px; background: #00BEBE; color: white; cursor: pointer;">ダウンロード</button>
          <button onclick="printScreenshot()" style="padding: 9px; border: none; border-radius: 4px; background: #00BEBE; color: white; cursor: pointer;">プリント</button>
        </div>
        <div style="display: flex; justify-content: center; width: 100%;">
          <img src="${payload}" style="max-width: 100%; object-fit: contain;" />
        </div>
        <div>
          <p>この地図は${localizedDate}にhttps://plateauview.mlit.go.jpで作られました。</p>
        </div>
      </div>
    </body>
    <script>
      const downloadScreenshot = () => {
        const link = document.createElement("a");
        link.download = "screenshot.png";
        link.href = "${payload}";
        link.click();
        link.remove();
      }
      const printScreenshot = () => {
        window.print()
      }
    </script>
  </html>

`;
  iframe.contentWindow?.document.write(iframeHTML);

  const iframeHtmlStyle = iframe.contentWindow?.document.createElement("style");
  if (iframeHtmlStyle) {
    iframeHtmlStyle.appendChild(document.createTextNode(css));
    iframeHtmlStyle.setAttribute("type", "text/css");
    iframe.contentWindow?.document.head.appendChild(iframeHtmlStyle);
  }

  iframe.contentWindow?.document.close();

  return iframe;
}
