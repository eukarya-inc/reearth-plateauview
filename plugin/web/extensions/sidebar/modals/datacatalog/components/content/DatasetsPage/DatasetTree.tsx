import { styled } from "@web/theme";

const DatasetTree: React.FC = () => {
  return (
    <Wrapper>
      <input placeholder="input search text" />
      <p>A file system like tree for datasets</p>
    </Wrapper>
  );
};

export default DatasetTree;

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  padding: 24px 12px;
`;
