import { styled } from "@web/theme";

export type Props = {
  onDatasetAdd: () => void;
};

const DatasetDetails: React.FC<Props> = ({ onDatasetAdd }) => {
  return (
    <Wrapper>
      <div>Some picture here</div>
      <p>buttons</p>
      <button onClick={onDatasetAdd}>Add dataset</button>
      <p>main content with title, maybe buttons and a description</p>
      <p>THIS NEEDS TO BE SCROLLABLE IF CONTENT LONG</p>
    </Wrapper>
  );
};

export default DatasetDetails;

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  padding: 24px;
`;
