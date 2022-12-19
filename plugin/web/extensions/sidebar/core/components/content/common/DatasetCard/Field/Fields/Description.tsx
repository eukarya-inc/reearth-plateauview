import { styled } from "@web/theme";

import { BaseField as BaseFieldProps } from ".";

type Props = BaseFieldProps<"description"> & {
  editMode?: boolean;
};

const Description: React.FC<Props> = ({ editMode }) => {
  return editMode ? (
    <div>
      <TextBoxTitle>asdf</TextBoxTitle>
      <TextBox rows={4} />
      {/* markdowntoggle */}
    </div>
  ) : (
    <div>PUBLISH MODE</div>
  );
};

export default Description;

const TextBoxTitle = styled.p`
  margin: 0;
`;

const TextBox = styled.textarea`
  width: 100%;
`;
