import Musa from "mini-framework";
import { socket } from "../../..";
import ChatBox from "../../components/ChatBox";
import Button from "../../components/Button";
import UsersList from "../../components/UsersList";
import {
  sendGameStartEvent,
  sendLobbyMessage,
} from "../../events/eventHelpers";
import { Timer } from "../../components/Timer";

function LobbyPage({ lobbyMessages, users, timer }) {
  function sendMessage(event: Event) {
    event.preventDefault();
    sendLobbyMessage(socket, event.currentTarget.message.value);
  }

  return (
    <div id="LobbyPage">
      <div className="TitleWrapper">
        <h1 className="Title text-center">Waiting for players <Timer users={users}/></h1>
      </div>
      <ChatBox
        title="Lobby chat"
        messages={lobbyMessages}
        sendMessage={sendMessage}
        id="LobbyChat"
      />
      <div className="LobbyMiddle">
        <div className="UsersListBox">
          <h2>Currently in the lobby: ({users.length})</h2>
          <UsersList users={users} />
        </div>
        <Button
          onClick={() => {
            sendGameStartEvent(socket);
          }}
          id="StartGameButton">
          Start Game
        </Button>
      </div>
      
    </div>
  );
}

export default LobbyPage;
