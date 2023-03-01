import {
  createContext,
  Dispatch,
  FC,
  PropsWithChildren,
  SetStateAction,
  useContext,
  useState,
} from "react";

export type ContextValue = {
  color: {
    expression: {
      conditions: string[][];
    };
  };
  transparency: number;
};

const defaultValue: ContextValue = {
  color: {
    expression: {
      conditions: [["true", "rgba(255, 255, 255, 1)"]],
    },
  },
  transparency: 1.0,
};

const context = createContext<[ContextValue, Dispatch<SetStateAction<ContextValue>>] | undefined>(
  undefined,
);

export const BuildingColorProvider: FC<PropsWithChildren> = ({ children }) => {
  const value = useState(defaultValue);
  return <context.Provider value={value}>{children}</context.Provider>;
};

export const useBuildingColorContext = () => {
  const ctx = useContext(context);
  if (!ctx) {
    throw new Error("BuildingColorContext is not initialized");
  }
  return ctx;
};
