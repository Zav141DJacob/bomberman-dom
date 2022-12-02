import Musa from "mini-framework";
import { socket } from "../../..";
import ChatBox from "../../components/ChatBox";
import Button from "../../components/Button";
import {
  sendGlobalMessage,
  sendLobbyJoinEvent,
  sendSetUsername,
} from "../../events/eventHelpers";
import { Page } from "../../models/utils";

function InitialPage({ globalMessages, setCurrentPage }) {
  function sendMessage(event: Event) {
    event.preventDefault();
    sendGlobalMessage(socket, event.currentTarget.message.value);
  }

  return (
    <div id="InitialPage">
      <div className="TitleWrapper">
        <div className="Title TitleAnimation text-center">Bomberman!</div>
      </div>

      <div className="margin-bottom InitialLeft">
        <ChatBox
          title="Global chat"
          messages={globalMessages}
          sendMessage={sendMessage}
          id="GlobalChat"
        />
        
      </div>

      

      <div id="InitialMiddle">
        <div id="NamePrompt" className="margin-bottom">
          <div>Enter your name</div>
          <input type="text" id="username" placeholder="Player"></input>
        </div>
        <Button
          onClick={() => {
            sendSetUsername(socket, document.getElementById("username").value);
            sendLobbyJoinEvent(socket);
            setCurrentPage(Page.Lobby);
          }}
        >
          Join lobby
        </Button>
      </div>
    </div>
  );
}

export default InitialPage;
