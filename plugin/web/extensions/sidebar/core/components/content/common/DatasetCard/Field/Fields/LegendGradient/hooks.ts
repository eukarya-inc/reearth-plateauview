import { getColorsGradients } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useState } from "react";

const legendStyles: { [key: string]: string } = {
  square: "四角",
  circle: "丸",
  line: "線",
};

const Fields: { [key: string]: string } = {
  unDefiend: "-",
};
export type LegendStyleType = "square" | "circle" | "line";

export type LegendItem = {
  title: string;
  color: string;
};

export type LegendGradientType = {
  id?: string;
  style: LegendStyleType;
  items: LegendItem[];
};

export type DisplayValues = string[];

export default (value: LegendGradientType) => {
  const [currentLegendGradient, updateLegend] = useState<LegendGradientType>(value);
  const [items, setLegendItems] = useState<LegendItem[] | undefined>(value.items);
  const [colorGradients, setColorGradients] = useState<string[]>();
  const [displayValues, setDisplayValues] = useState<DisplayValues>();
  const [startColor, setStartColor] = useState("#753E13");
  const [endColor, setEndColor] = useState("#F2EEEB");
  const [step, setStep] = useState(10);

  const handleStyleChange = useCallback((style: LegendStyleType) => {
    updateLegend(l => {
      return {
        ...l,
        style,
      };
    });
  }, []);

  const handleChooseField = useCallback((field: string) => {
    console.log(field);
  }, []);

  const handleLegendItemsChange = useCallback(() => {
    if (items === undefined) return;
    if (displayValues === undefined) return;
    const result = colorGradients?.map((color, index) => {
      return {
        title: displayValues[index],
        color: color,
      };
    });
    setLegendItems(result);

    updateLegend(l => {
      return {
        ...l,
        items,
      };
    });
  }, [colorGradients, displayValues, items]);

  const handleStepChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setStep(parseInt(e.target.value));
  }, []);

  const handleStartColorChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      setStartColor(e.target.value);
    },

    [],
  );

  const handleEndColorChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setEndColor(e.target.value);
  }, []);

  useEffect(() => {
    setDisplayValues(() => {
      return [...Array(step).keys()].map(x => (x + 0).toString());
    });
  }, [step]);

  useEffect(() => {
    setColorGradients(getColorsGradients(startColor, endColor, step));
    handleLegendItemsChange();
  }, [startColor, endColor, step, handleLegendItemsChange]);

  return {
    startColor,
    endColor,
    step,
    legendStyles,
    Fields,
    currentLegendGradient,
    handleStyleChange,
    handleChooseField,
    handleStepChange,
    handleStartColorChange,
    handleEndColorChange,
  };
};
