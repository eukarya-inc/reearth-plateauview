import {
  ContextValue,
  useBuildingColorContext,
} from "@web/extensions/sidebar/core/BuildingColorContext";
import { getRGBAFromString, RGBA, rgbaToString } from "@web/extensions/sidebar/utils/color";
import debounce from "lodash/debounce";
import {
  Dispatch,
  RefObject,
  SetStateAction,
  useCallback,
  useEffect,
  useMemo,
  useRef,
} from "react";

import { BaseFieldProps } from "../../types";

export const useBuildingTransparency = ({
  options,
  dataID,
  onUpdate,
}: Pick<BaseFieldProps<"buildingTransparency">, "dataID"> & {
  options: Omit<BaseFieldProps<"buildingTransparency">["value"], "id" | "group" | "type">;
  onUpdate: (property: any) => void;
}) => {
  const renderRef = useRef<() => void>();
  const debouncedRender = useMemo(
    () => debounce(() => renderRef.current?.(), 100, { maxWait: 300 }),
    [],
  );
  const [, setColor] = useBuildingColorContext();

  const onUpdateRef = useRef(onUpdate);
  useEffect(() => {
    onUpdateRef.current = onUpdate;
  }, [onUpdate]);

  const render = useCallback(() => {
    renderTileset(
      {
        dataID,
        transparency: options.transparency,
      },
      onUpdateRef,
      setColor,
    );
  }, [options.transparency, dataID, setColor]);

  useEffect(() => {
    renderRef.current = render;
    debouncedRender();
  }, [render, debouncedRender]);
};

export type State = {
  dataID: string | undefined;
  transparency: number;
};

const renderTileset = (
  state: State,
  onUpdateRef: RefObject<(property: any) => void>,
  setColor: Dispatch<SetStateAction<ContextValue>>,
) => {
  const updateTileset = async () => {
    setColor(color => {
      const transparency = (state.transparency ?? 100) / 100;

      // We can get transparency from RGBA. Because the color is defined as RGBA.
      const overriddenColor = color.color;
      const defaultRGBA = rgbaToString([255, 255, 255, transparency]);
      const expression = (() => {
        return {
          expression: {
            conditions: (overriddenColor.expression.conditions as [string, string][]).map(
              ([k, v]: [string, string]) => {
                const rgba = getRGBAFromString(v);
                if (!rgba) {
                  return [k, defaultRGBA];
                }
                const composedRGBA = [...rgba.slice(0, -1), transparency] as RGBA;
                return [k, rgbaToString(composedRGBA)];
              },
            ),
          },
        };
      })();

      onUpdateRef.current?.({
        color: expression,
      });

      return {
        ...color,
        transparency,
      };
    });
  };

  updateTileset();
};
