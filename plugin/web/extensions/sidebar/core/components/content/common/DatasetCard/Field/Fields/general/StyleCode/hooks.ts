import JSON5 from "json5";
import { useCallback, useState } from "react";

import { BaseFieldProps } from "../../types";

export default ({ value, onUpdate }: Pick<BaseFieldProps<"styleCode">, "value" | "onUpdate">) => {
  const [code, editCode] = useState(value.src);

  const onApply = useCallback(() => {
    if (code) {
      try {
        const styleObject = JSON5.parse(code);
        onUpdate({
          ...value,
          src: code,
          override: styleObject,
        });
        // eslint-disable-next-line no-empty
      } catch (error) {}
    }
  }, [onUpdate, code, value]);

  const onEdit = useCallback(
    (e: React.ChangeEvent<HTMLTextAreaElement>) => {
      editCode(e.target.value);
      onUpdate({
        ...value,
        src: e.target.value,
      });
    },
    [onUpdate, value],
  );

  return {
    code,
    editCode,
    onApply,
    onEdit,
  };
};
