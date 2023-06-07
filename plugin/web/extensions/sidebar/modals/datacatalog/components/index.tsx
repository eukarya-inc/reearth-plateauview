import DatasetsPage from "@web/extensions/sidebar/modals/datacatalog/components/content/DatasetsPage";
import YourDataPage from "@web/extensions/sidebar/modals/datacatalog/components/content/YourDataPage";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";

import useHooks from "./hooks";
import useDataset from "./useDataset";

const DataCatalog: React.FC = () => {
  const { currentTab, inEditor, isCustomProject, handleClose, handleTabChange } = useHooks();

  const {
    catalog: customCatalog,
    addedDatasetDataIDs: customAddedDatasetDataIDs,
    selectedItem: customSelectedItem,
    expandedFolders: customExpandedFolders,
    searchTerm: customSearchTerm,
    filter: customFilter,
    setExpandedFolders: setCustomExpandedFolders,
    handleSearch: handleCustomSearch,
    handleSelect: handleCustomSelect,
    handleFilter: handleCustomFilter,
    handleDatasetAdd: handleCustomDatasetAdd,
    handleDatasetPublish: handleCustomDatasetPublish,
  } = useDataset({ inEditor, dataSource: "custom" });

  const {
    catalog,
    addedDatasetDataIDs,
    selectedItem,
    expandedFolders,
    searchTerm,
    filter,
    setExpandedFolders,
    handleSearch,
    handleSelect,
    handleFilter,
    handleDatasetAdd,
    handleDatasetPublish,
  } = useDataset({ inEditor, dataSource: "plateau" });

  return (
    <Wrapper>
      <Header>
        <Title>データカタログ</Title>
        <TabsWrapper>
          {isCustomProject && (
            <Tab selected={currentTab === "custom"} onClick={() => handleTabChange("custom")}>
              <TabName>Custom Dataset</TabName>
            </Tab>
          )}
          <Tab selected={currentTab === "plateau"} onClick={() => handleTabChange("plateau")}>
            <Logo icon="plateauLogoPart" selected={currentTab === "plateau"} />
            <TabName>PLATEAUデータセット</TabName>
          </Tab>
          <Tab selected={currentTab === "your-data"} onClick={() => handleTabChange("your-data")}>
            <Icon icon="user" />
            <TabName>Myデータ</TabName>
          </Tab>
        </TabsWrapper>
        <CloseButton>
          <Icon size={32} icon="close" onClick={handleClose} />
        </CloseButton>
      </Header>
      {currentTab === "your-data" ? (
        <YourDataPage onDatasetAdd={handleDatasetAdd} />
      ) : currentTab === "plateau" ? (
        <DatasetsPage
          catalog={catalog}
          addedDatasetDataIDs={addedDatasetDataIDs}
          inEditor={inEditor}
          selectedItem={selectedItem}
          expandedFolders={expandedFolders}
          searchTerm={searchTerm}
          filter={filter}
          dataSource={"plateau"}
          editable={!isCustomProject}
          setExpandedFolders={setExpandedFolders}
          onSearch={handleSearch}
          onSelect={handleSelect}
          onFilter={handleFilter}
          onDatasetAdd={handleDatasetAdd}
          onDatasetPublish={handleDatasetPublish}
        />
      ) : (
        <DatasetsPage
          catalog={customCatalog}
          addedDatasetDataIDs={customAddedDatasetDataIDs}
          inEditor={inEditor}
          selectedItem={customSelectedItem}
          expandedFolders={customExpandedFolders}
          searchTerm={customSearchTerm}
          filter={customFilter}
          dataSource={"custom"}
          editable={isCustomProject}
          setExpandedFolders={setCustomExpandedFolders}
          onSearch={handleCustomSearch}
          onSelect={handleCustomSelect}
          onFilter={handleCustomFilter}
          onDatasetAdd={handleCustomDatasetAdd}
          onDatasetPublish={handleCustomDatasetPublish}
        />
      )}
    </Wrapper>
  );
};

export default DataCatalog;

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  width: 905px;
  height: 590px;
  background: #f4f4f4;
  box-shadow: 0px 3px 6px -4px rgba(0, 0, 0, 0.12), 0px 6px 16px rgba(0, 0, 0, 0.08),
    0px 9px 28px 8px rgba(0, 0, 0, 0.05);
`;

const Header = styled.div`
  display: flex;
  background: #dcdcdc;
  height: 48px;
  position: relative;
`;

const Title = styled.p`
  align-self: center;
  font-size: 14px;
  font-weight: 700;
  margin: 0 12px;
  color: #4a4a4a;
  user-select: none;
`;

const TabsWrapper = styled.div`
  display: flex;
  margin: 8px 10px 0 10px;
  gap: 10px;
`;

const Tab = styled.div<{ selected?: boolean }>`
  display: flex;
  gap: 8px;
  border-width: 1px 1px 0px 1px;
  border-style: solid;
  border-color: ${({ selected }) => (selected ? "#f4f4f4" : "#c8c8c8")};
  border-radius: 2px 2px 0px 0px;
  background: ${({ selected }) => (selected ? "#f4f4f4" : "#c8c8c8")};
  color: ${({ selected }) => (selected ? "var(--theme-color)" : "#898989")};
  padding: 8px 12px;
  cursor: pointer;
`;

const Logo = styled(Icon)<{ selected?: boolean }>`
  opacity: ${({ selected }) => (selected ? 1 : 0.4)};
`;

const TabName = styled.p`
  margin: 0;
  user-select: none;
`;

const CloseButton = styled.button`
  display: flex;
  justify-content: center;
  align-items: center;
  position: absolute;
  right: 0;
  height: 48px;
  width: 48px;
  border: none;
  background: var(--theme-color);
  color: white;
  cursor: pointer;
`;
