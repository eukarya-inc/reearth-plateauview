import { Switch } from "@web/sharedComponents";
import { styled } from "@web/theme";
import ReactMarkdown from "react-markdown";
import gfm from "remark-gfm";

import { BaseField as BaseFieldProps } from ".";

type Props = BaseFieldProps<"description"> & {
  isMarkdown?: boolean;
  value?: string;
  editMode?: boolean;
};

const plugins = [gfm];

const Description: React.FC<Props> = ({ value, isMarkdown, editMode }) => {
  return editMode ? (
    <div>
      <Text>内容</Text>
      <TextBox rows={4} defaultValue={value} />
      <SwitchWrapper>
        <Switch checked={isMarkdown} size="small" />
        <Text>マークダウン</Text>
      </SwitchWrapper>
    </div>
  ) : isMarkdown && value ? (
    <ReactMarkdown remarkPlugins={plugins} linkTarget="_blank">
      {value}
    </ReactMarkdown>
  ) : (
    <div>{value}</div>
  );
};

export default Description;

const Text = styled.p`
  margin: 0;
`;

const TextBox = styled.textarea`
  width: 100%;
`;

const SwitchWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 12px;
`;
