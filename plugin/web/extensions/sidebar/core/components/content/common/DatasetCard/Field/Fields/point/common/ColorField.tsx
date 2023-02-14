import { styled } from "@web/theme";

import { TextInput } from "./styled";

const ColorField: React.FC<{ color: string; value: string }> = ({ color, value }) => {
  return (
    <>
      <ColorBlock color={color} />
      <TextInput value={value} />
    </>
  );
};

export default ColorField;

const ColorBlock = styled.div<{ color: string; legendStyle?: "circle" | "square" | "line" }>`
  width: 30px;
  height: ${({ legendStyle }) => (legendStyle === "line" ? "3px" : "30px")};
  background: ${({ color }) => color ?? "#d9d9d9"};
  border-radius: ${({ legendStyle }) =>
    legendStyle
      ? legendStyle === "circle"
        ? "50%"
        : legendStyle === "line"
        ? "5px"
        : "2px"
      : "1px 0 0 1px"};
`;
