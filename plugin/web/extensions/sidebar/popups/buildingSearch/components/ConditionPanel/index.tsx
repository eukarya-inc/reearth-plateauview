import { Empty, Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";

import type { DatasetIndexes, Condition as ConditionType } from "../../types";

import Condition from "./Condition";

type Props = {
  active: boolean;
  datasetIndexes?: DatasetIndexes;
  conditionsState: "loading" | "empty" | "ready";
  conditionApply: () => void;
  setConditions: React.Dispatch<React.SetStateAction<ConditionType[]>>;
};

const ConditionPanel: React.FC<Props> = ({
  active,
  datasetIndexes,
  conditionsState,
  conditionApply,
  setConditions,
}) => {
  return (
    <Wrapper active={active}>
      {conditionsState === "loading" && <Loading>Loading...</Loading>}
      {conditionsState === "ready" && (
        <ConditionWrapper>
          <DatasetInfo>
            <Icon icon="database" size={24} />
            <DatasetName>{datasetIndexes?.title}</DatasetName>
          </DatasetInfo>
          <Conditions>
            {datasetIndexes?.indexes.map(indexItem => (
              <Condition
                key={indexItem.field}
                indexItem={indexItem}
                setConditions={setConditions}
              />
            ))}
          </Conditions>
          <ButtonWrapper>
            <Button onClick={conditionApply}>Apply</Button>
          </ButtonWrapper>
        </ConditionWrapper>
      )}
      {conditionsState === "empty" && (
        <EmptyWrapper>
          <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />
        </EmptyWrapper>
      )}
    </Wrapper>
  );
};

const Wrapper = styled.div<{ active: boolean }>`
  display: ${({ active }) => (active ? "flex" : "none")};
  padding: 8px 0;
  align-items: center;
  height: 100%;
  overflow: hidden;
`;

const ConditionWrapper = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
`;

const Loading = styled.div`
  width: 100%;
  display: flex;
  justify-content: center;
  color: #8c8c8c;
`;

const EmptyWrapper = styled.div`
  width: 100%;
  display: flex;
  justify-content: center;
`;

const DatasetInfo = styled.div`
  height: 40px;
  display: flex;
  align-items: flex-start;
  justify-content: flex-start;
  flex-shrink: 0;
  width: 100%;
  gap: 8px;
  padding: 9px 20px;
`;

const DatasetName = styled.div`
  font-size: 16px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
`;

const Conditions = styled.div`
  height: 350px;
  overflow: auto;
`;

const ButtonWrapper = styled.div`
  padding: 6px 20px;
`;

const Button = styled.div`
  width: 67px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  background: var(--theme-color);
  border-radius: 4px;
  cursor: pointer;
`;

export default ConditionPanel;
