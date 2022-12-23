const Popup: React.FC = () => {
  // Initially open Basic
  addEventListener("message", e => {
    if (e.source !== parent) return;
    if (e.data.type) {
      if (e.data.type === "msgFromHelp") {
        console.log("POPUP MESSAGE: ", e.data.message);
      }
    }
  });
  // Handle receiving message from sidebar (selectedTab)
  return <div>I AM A POOPPP UP</div>;
};

export default Popup;
