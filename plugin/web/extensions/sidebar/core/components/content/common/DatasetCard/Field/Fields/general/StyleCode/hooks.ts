import { postMsg } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useState } from "react";

import { BaseFieldProps } from "../../types";

export default ({
  value,
  dataID,
  onUpdate,
}: Pick<BaseFieldProps<"styleCode">, "value" | "dataID" | "onUpdate">) => {
  const [code, editCode] = useState(value.src);

  const handleEditCode = useCallback(
    (e: React.ChangeEvent<HTMLTextAreaElement>) => {
      editCode(e.currentTarget.value);
      onUpdate({
        ...value,
        src: code,
      });
    },
    [code, onUpdate, value],
  );

  useEffect(() => {
    console.log("update", code);
    postMsg({ action: "updateLayerStyleCode", payload: { dataID, code } });
  }, [dataID, code]);

  return {
    code,
    editCode,
    handleEditCode,
  };
};
