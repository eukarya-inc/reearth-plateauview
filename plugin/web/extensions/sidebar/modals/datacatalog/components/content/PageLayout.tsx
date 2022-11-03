import { styled } from "@web/theme";
import { ReactNode } from "react";

export type Props = {
  left: ReactNode;
  right: ReactNode;
};

const PageLayout: React.FC<Props> = ({ left, right }) => {
  return (
    <Body>
      <FileView>{left}</FileView>
      <Divider />
      <Details>{right}</Details>
    </Body>
  );
};

export default PageLayout;

const Body = styled.div`
  display: flex;
  flex: 1;
`;

const FileView = styled.div`
  width: 402px;
`;

const Details = styled.div`
  flex: 1;
`;

const Divider = styled.div`
  border-right: 1px solid #c7c5c5;
  padding-left: 10px;
  margin: 24px 10px 24px 0;
`;
