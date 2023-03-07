import JSON5 from "json5";
import { useCallback, useEffect, useState } from "react";

import { BaseFieldProps } from "../../types";

export default ({ value, onUpdate }: Pick<BaseFieldProps<"styleCode">, "value" | "onUpdate">) => {
  const [code, editCode] = useState(value.src);

  const onApply = useCallback(() => {
    try {
      const styleObject = JSON5.parse(code);
      onUpdate({
        ...value,
        src: code,
        override: styleObject,
      });
      // eslint-disable-next-line no-empty
    } catch (error) {}
  }, [onUpdate, code, value]);

  useEffect(() => {
    onApply();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const onEdit = useCallback((e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newValue = e.target.value;
    editCode(newValue);
  }, []);

  return {
    code,
    editCode,
    onApply,
    onEdit,
  };
};
