import { Collapse, Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";

import Editor from "./components/editor";
import Viewer from "./components/viewer";
import useHooks from "./hooks";

const Infobox: React.FC = () => {
  const { mode, dataState, feature, fields, saveFields } = useHooks();

  return (
    <Wrapper>
      {dataState === "loading" && <SimplePane>Loading...</SimplePane>}
      {dataState === "empty" && <SimplePane>Empty</SimplePane>}
      {dataState === "ready" && mode === "edit" && fields && feature && (
        <StyledCollapse
          bordered={false}
          collapsible={"header"}
          defaultActiveKey={[0]}
          expandIconPosition="end"
          expandIcon={({ isActive }: { isActive?: boolean }) => (
            <IconWrapper active={!!isActive}>
              <Icon icon="arrowDown" color="#000000" size={18} />
            </IconWrapper>
          )}>
          <Editor key={0} fields={fields} feature={feature} saveFields={saveFields} />
        </StyledCollapse>
      )}
      {dataState === "ready" && mode === "view" && feature && fields && (
        <StyledCollapse
          bordered={false}
          defaultActiveKey={[0]}
          expandIconPosition="end"
          expandIcon={({ isActive }: { isActive?: boolean }) => (
            <IconWrapper active={!!isActive}>
              <Icon icon="arrowDown" color="#000000" size={18} />
            </IconWrapper>
          )}>
          <Viewer feature={feature} fields={fields} key={0} />
        </StyledCollapse>
      )}
    </Wrapper>
  );
};

const Wrapper = styled.div`
  padding: 0 12px;
`;

const SimplePane = styled.div`
  width: 100%;
  padding: 12px 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
`;

const StyledCollapse = styled(Collapse)`
  background: none;

  .ant-collapse-header {
    font-size: 16px;
    color: #000;
    background: #fff;
    padding: 12px 40px 12px 12px;
    align-items: center !important;
  }
  .ant-collapse-content-box {
    border-top: 1px solid #e0e0e0;
    padding: 12px !important;
  }
  .ant-collapse-item {
    border-bottom: none;
  }
`;

const IconWrapper = styled.div<{ active: boolean }>`
  width: 20px;
  height: 20px;
  cursor: pointer;

  svg {
    transition: all 0.25s ease;
    transform: ${({ active }) => (active ? "rotate(0)" : "rotate(90deg)")};
  }
`;

export default Infobox;
