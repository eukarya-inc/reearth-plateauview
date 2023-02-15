import { styled } from "@web/theme";

import { BaseFieldProps } from "../types";

const Search: React.FC<BaseFieldProps<"search">> = ({ editMode }) => {
  return editMode ? (
    <Wrapper>
      <Title>Enabled</Title>
    </Wrapper>
  ) : null;
};

const Wrapper = styled.div`
  width: 100%;
  height: 32px;
  display: flex;
  justify-content: flex-start;
  align-items: center;
  gap: 4px;
`;

const Title = styled.div`
  width: 100%;
  font-size: 14px;
  display: flex;
  align-items: center;
`;

export default Search;
