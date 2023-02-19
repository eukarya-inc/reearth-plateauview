import isEqual from "lodash/isEqual";
import { useCallback, useEffect, useState } from "react";

import { BaseFieldProps } from "../../types";

import { useBuildingShadow } from "./useBuildingShadow";

type OptionsState = Omit<BaseFieldProps<"buildingShadow">["value"], "id" | "group" | "type">;

const useHooks = ({
  value,
  dataID,
  onUpdate,
}: Pick<BaseFieldProps<"buildingShadow">, "value" | "dataID" | "onUpdate">) => {
  const [options, setOptions] = useState<OptionsState>({
    shadow: value.shadow,
  });

  const handleUpdate = useCallback(
    <P extends keyof OptionsState>(prop: P, v?: OptionsState[P]) => {
      setOptions(o => {
        const next = { ...o, [prop]: v };
        onUpdate({ id: value.id, type: value.type, group: value.group, ...next });
        return next;
      });
    },
    [onUpdate, value],
  );

  const handleUpdateSelect = useCallback(
    (prop: keyof OptionsState) => (value: any) => {
      handleUpdate(prop, value as OptionsState["shadow"]);
    },
    [handleUpdate],
  );

  useEffect(() => {
    if (!isEqual(options, value)) {
      setOptions({ ...value });
    }
  }, [value, options]);

  useBuildingShadow({ value, dataID });

  return {
    options,
    handleUpdateSelect,
  };
};

export default useHooks;
