import { postMsg } from "@web/extensions/sidebar/utils";
import { Radio } from "@web/sharedComponents";
import isEqual from "lodash/isEqual";
import { ComponentProps, useCallback, useEffect, useState } from "react";

import { BaseFieldProps } from "../../types";

import { useBuildingColor } from "./useBuildingColor";

type OptionsState = Omit<BaseFieldProps<"buildingColor">["value"], "id" | "group" | "type">;

const useHooks = ({
  value,
  dataID,
  onUpdate,
}: Pick<BaseFieldProps<"buildingColor">, "value" | "dataID" | "onUpdate">) => {
  const [options, setOptions] = useState<OptionsState>({
    colorType: value.colorType,
  });
  const [floods, setFloods] = useState<
    { id: string; label: string; featurePropertyName: string }[]
  >([]);
  const [initialized, setInitialized] = useState(false);

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

  const handleUpdateColorType: Exclude<ComponentProps<typeof Radio>["onChange"], undefined> =
    useCallback(
      e => {
        handleUpdate("colorType", e.target.value);
      },
      [handleUpdate],
    );

  useEffect(() => {
    if (!isEqual(options, value)) {
      setOptions({ ...value });
    }
  }, [value, options]);

  useEffect(() => {
    const MARK_TEXT = ["洪水浸水想定区域", "浸水予想区域"];
    const EXCLUDING_RANGE_NAME = ["改定"];
    const waitReturnedPostMsg = async (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      if (e.data.action === "findTileset") {
        const layer = e.data.payload.layer;
        const url = layer?.data?.url;
        if (!url) {
          return;
        }
        const data = await (async () => {
          try {
            return await fetch(url).then(r => r.json());
          } catch (e) {
            console.error(e);
          }
        })();

        Object.entries(data?.properties || {}).forEach(([k, v]) => {
          if (
            k.endsWith("_浸水ランク") &&
            v &&
            typeof v === "object" &&
            Object.keys(v).length > 0
          ) {
            const marker = MARK_TEXT.find(s => k.includes(s));
            if (!marker) {
              return;
            }
            const [label, rest] = k.split(`${marker}`);
            const matches = rest.match(/(（.*）)?_(.*)_浸水ランク/);
            if (!matches) {
              return;
            }
            const range = matches[1]?.match(/（(.*)）/)?.[1];
            const scale = matches[2];
            const level = scale === "想定最大規模" ? "L2" : "L1";
            setFloods(v => [
              ...v,
              {
                id: `floods-${v.length}`,
                label: `${level || ""}${scale ? `(${scale})` : ""}_浸水ランク(${label}${
                  range && !EXCLUDING_RANGE_NAME.includes(range) ? `: ${range}` : ""
                })`,
                featurePropertyName: k,
              },
            ]);
          }
        });
        removeEventListener("message", waitReturnedPostMsg);
        setInitialized(true);
      }
    };
    addEventListener("message", waitReturnedPostMsg);
    postMsg({
      action: "findTileset",
      payload: {
        dataID,
      },
    });
  }, [dataID]);

  useBuildingColor({ value, dataID, floods, initialized });

  return {
    options,
    floods,
    handleUpdateColorType,
  };
};

export default useHooks;