import Musa from "mini-framework";
import GameWindow from "../../components/GameWindow";
import Dashboard from "../../components/Dashboard";
import ChatBox from "../../components/ChatBox";
import { sendLobbyMessage } from "../../events/eventHelpers";
import { socket } from "../../..";
import Helper from "../../components/Helper";

const downKey =
  "https://icons.iconarchive.com/icons/chromatix/keyboard-keys/128/arrow-down-icon.png";
const leftKey =
  "https://icons.iconarchive.com/icons/chromatix/keyboard-keys/128/arrow-left-icon.png";
const rightKey =
  "https://icons.iconarchive.com/icons/chromatix/keyboard-keys/128/arrow-right-icon.png";
const upKey =
  "https://icons.iconarchive.com/icons/chromatix/keyboard-keys/128/arrow-up-icon.png";
const spaceKey =
  "https://raw.githubusercontent.com/q2apro/keyboard-keys-speedflips/master/single-keys-blank/200dpi/spacebar.png";

function GamePage({ gameMap, dir, lobbyMessages, players, winner }) {
  function sendMessage(event: Event) {
    event.preventDefault();
    sendLobbyMessage(socket, event.currentTarget.message.value);
  }

  return (
    <div id="GamePage">
      <ChatBox
        title="Lobby chat"
        messages={lobbyMessages}
        sendMessage={sendMessage}
      />
      <div className="GameMiddle">
        <div className="MainGame">
          <GameWindow
            dir={dir}
            gameMap={gameMap}
            players={players}
            winner={winner}
          />
          <Helper />
        </div>
        <Dashboard players={players} dir={dir} />
      </div>
    </div>
  );
}

export default GamePage;
