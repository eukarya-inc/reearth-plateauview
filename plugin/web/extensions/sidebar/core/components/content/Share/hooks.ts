import { usePublishUrl } from "@web/extensions/sidebar/core/state";
import { postMsg } from "@web/extensions/sidebar/core/utils";
import { useCallback } from "react";

export default () => {
  const [publishUrl, setPublishUrl] = usePublishUrl();

  const handleScreenshotShow = useCallback(() => {
    postMsg({ action: "screenshot" });
  }, []);

  const handleScreenshotSave = useCallback(() => {
    postMsg({ action: "screenshot-save" });
  }, []);

  const handleProjectShare = useCallback(() => {
    const suffix = makeUrlSuffix();
    if (!publishUrl) {
      // To do: hit PLATEAU backend endpoint and create project, returning
      // publish URl or error.
      setPublishUrl(`https://plateauview.mlit.go.jp/${suffix}`);
    }
  }, []);

  return {
    publishUrl,
    handleProjectShare,
    handleScreenshotShow,
    handleScreenshotSave,
  };
};

function makeUrlSuffix() {
  let result = "";
  const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  const charactersLength = characters.length;
  for (let i = 0; i < 15; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}
