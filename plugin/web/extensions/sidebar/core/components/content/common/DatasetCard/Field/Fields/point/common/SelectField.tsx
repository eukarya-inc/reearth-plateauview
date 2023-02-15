import { Select } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { ComponentProps } from "react";

import { FieldTitle, BorderlessFieldValue, FieldWrapper } from "../commonComponents";

type Props = {
  title: string;
  titleWidth?: number;
} & ComponentProps<typeof Select>;

const SelectField: React.FC<Props> = ({ title, titleWidth, ...props }) => {
  return (
    <FieldWrapper>
      <FieldTitle width={titleWidth}>{title}</FieldTitle>
      <BorderlessFieldValue>
        <StyledSelect {...props} />
      </BorderlessFieldValue>
    </FieldWrapper>
  );
};

export default SelectField;

const StyledSelect = styled(Select)`
  width: 100%;
`;
