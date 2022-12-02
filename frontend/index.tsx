import Musa from "mini-framework";
import GamePage from "./src/pages/GamePage";
import { Events, WebsocketMessageType } from "./websocket";
import {
  sendGameStartEvent,
  sendPlayerMoveEvent,
  sendPlayerBombPlace,
} from "./src/events/eventHelpers";
import LobbyPage from "./src/pages/LobbyPage";
import ChatMessage from "./src/models/ChatMessage";
import { User } from "./src/models/User";
import { MapInfo, mapObject } from "./src/models/Maps";
import { Page } from "./src/models/utils";
import { toI } from "./src/functions/coords";
import InitialPage from "./src/pages/InitialPage";
import { MoveDirection } from "./websocket";
import { resetTimer, stopTimer } from "./src/components/Timer";

export { sendGameStartEvent };

export let socket;

const websocketURL = "ws://localhost:8080/game";

let throttleFunc;
const throttle = (wait: number) => {
  let first = true;
  return function ex(...s: any[]) {
    if (first) {
      throttleFunc.apply(this, s);
      first = false;

      setTimeout(() => {
        first = true;
      }, wait);
    }
  };
};

let timer2 = 30;
let currentSpeed = 500;
let throttledMover = throttle(currentSpeed);
let previousMap: MapInfo;
function App() {
  // State tracking and stuff should be here
  const [globalMessages, setGlobalMessages] = Musa.useState<ChatMessage[]>([]);
  const [users, setUsers] = Musa.useState<User[]>([]);
  const [usersID, setUsersID] = Musa.useState<User[]>([]);
  const [winner, setWinner] = Musa.useState<String | number>(new String());

  const [lobbyMessages, setLobbyMessages] = Musa.useState<ChatMessage[]>([]);
  const [gameMap, setGameMap] = Musa.useState<MapInfo>([]);
  const [dir, setDir] = Musa.useState(["stand", "stand", "stand", "stand"]);
  const [dirID, setDirID] = Musa.useState(new Set());
  const [timer, setTimer] = Musa.useState(30);

  const [currentPage, setCurrentPage] = Musa.useState<Page>(Page.Initial);

  const handleSetGameMap = (map: MapInfo) => {
    setGameMap(map);
    previousMap = map;
  };

  if (!socket) {
    socket = new WebSocket(websocketURL);
  }

  socket.onmessage = (event) => {
    let msg: Events = JSON.parse(event.data);
    let copyMap = previousMap ? previousMap : gameMap.slice();

    switch (msg.type) {
      case WebsocketMessageType.SetUsername:
        {
          let id = msg.content.ID;
          let newName = msg.content.Username;
          setUsers(
            users.map((user: User) => {
              if (user.id === id) {
                user.setName(newName);
              }

              return user;
            })
          );
        }
        break;
      case WebsocketMessageType.GlobalMessage:
        {
          let tempMessages: ChatMessage[] = globalMessages;
          tempMessages.push(msg.content);
          setGlobalMessages(tempMessages);
        }
        break;
      case WebsocketMessageType.LobbyMessage:
        {
          let tempMessages: ChatMessage[] = lobbyMessages;
          let dirCopy = dir.slice();
          dirCopy = dirCopy.map((i) =>
            !i.includes("stand") ? "stand" + i : i
          );
          setDir(dirCopy);
          tempMessages.push(msg.content);
          setLobbyMessages(tempMessages);
        }
        break;
      case WebsocketMessageType.GameStart:
        {
          stopTimer();
          setCurrentPage(Page.Game);
        }
        break;
      case WebsocketMessageType.BombPlant:
        {
          let i = toI(msg.content.X, msg.content.Y);
          copyMap[i] = mapObject.Bomb;
          let dirCopy = dir.slice();
          dirCopy = dirCopy.map((i) =>
            !i.includes("stand") ? "stand" + i : i
          );
          let currUser = users.find((u) => u.getId() == msg.content.UserID);
          console.log(currUser.getBombsLeft());
          if (currUser.getBombsLeft() == 1) {
            currUser.setBombsLeft(0);
          } else {
            currUser.setBombsLeft(currUser.getBombsLeft() - 1);
          }
          setDir(dirCopy);
          handleSetGameMap(copyMap);
        }
        break;
      case WebsocketMessageType.BombExplode:
        {
          if (msg.content.AffectedBlocks) {
            msg.content.AffectedBlocks.forEach((b) => {
              let i = toI(b.X, b.Y);
              copyMap[i] = mapObject.Explosion;
            });
          }
          let dirCopy = dir.slice();
          dirCopy = dirCopy.map((i) =>
            !i.includes("stand") ? "stand" + i : i
          );
          let currUser = users.find((u) => u.getId() == msg.content.UserID);
          if (currUser.getBombsLeft() == 0) {
            currUser.setBombsLeft(1);
          } else {
            currUser.setBombsLeft(currUser.getBombsLeft() + 1);
          }
          setDir(dirCopy);
          if (msg.content.AffectedBlocks) {
            setTimeout(() => {
              msg.content.AffectedBlocks.forEach((b) => {
                let i = toI(b.X, b.Y);
                copyMap[i] = mapObject.Empty;
              });
              handleSetGameMap(copyMap);
            }, 300);
          }
        }
        break;
      case WebsocketMessageType.PlayerMove:
        {
          // TODO: handle player moving
          // Alex's code
          let newDirID = dirID.has(msg.content.ID);
          let userIndex = usersID.indexOf(msg.content.ID);
          let dirCopy = dir.slice();
          let tempUsers: User[] = users.slice();
          let userID = usersID.slice();
          let copyMap = new Set(dirID);

          if (!newDirID) {
            copyMap = new Set(dirID);
            copyMap.add(msg.content.ID);
            userID.push(msg.content.ID);
            dirCopy = Array.from({ length: users.length }, (a, v) => {
              if (dir[v].includes("stand")) {
                return dir[v];
              } else {
                return "stand" + dir[v];
              }
            });
            setDirID((i) => {
              let x = new Set();
              i.forEach((v) => {
                x.add(v);
              });
              x.add(msg.content.ID);
              return x;
            });
            setUsersID((i) => {
              i.push(msg.content.ID);
              tempUsers = tempUsers.sort(
                (a, b) => i.indexOf(a.id) - i.indexOf(b.id)
              );
              return i;
            });
          } else {
            tempUsers = tempUsers.sort(
              (a, b) => userID.indexOf(a.id) - userID.indexOf(b.id)
            );
            userIndex = userID.indexOf(msg.content.ID);
            dirCopy = Array.from({ length: users.length }, (a, v) => {
              if (dir[v].includes("stand")) {
                return dir[v];
              } else {
                return "stand" + dir[v];
              }
            });
            switch (msg.content.Direction) {
              case 1: {
                dirCopy[userIndex] = "up";
                break;
              }
              case 2: {
                dirCopy[userIndex] = "down";
                break;
              }
              case 3: {
                dirCopy[userIndex] = "right";
                break;
              }
              case 4: {
                dirCopy[userIndex] = "left";
                break;
              }
            }
          }

          // Set player's coordinates
          var id: string = msg.content.ID;
          var X: number = msg.content.X;
          var Y: number = msg.content.Y;
          var foundIndex = tempUsers.findIndex(
            (user: User) => user.getId() === id
          );

          if (foundIndex >= 0) {
            tempUsers[foundIndex].setX(X);
            tempUsers[foundIndex].setY(Y);
          }
          setUsers(tempUsers);
          setDir(dirCopy);
        }
        break;
      case WebsocketMessageType.JoinLobby:
        {
          resetTimer();
          let id = msg.content.ID;
          let username = msg.content.Username;

          let tempUsers: User[] = users;
          var user: User = users.find((user: User) => user.getId() === id);
          if (!user) {
            let newUser: User = new User(id, false, true, username, 3);
            tempUsers.push(newUser);
          } else {
            var index = users.findIndex(
              (findUser: User) => findUser.getId() === user.getId()
            );
            user.setInCurrentLobby(true);
            user.setName(username);
            tempUsers[index] = user;
          }
          setUsers(tempUsers);
        }
        break;
      case WebsocketMessageType.ItemDrop:
        {
          setTimeout(() => {
            let dropX = msg.content.X;
            let dropY = msg.content.Y;

            if (msg.content.Item == 4) {
              throttledMover = throttle(currentSpeed - 100);
              currentSpeed -= 100;
            }
            let powerUp = msg.content.Item + 2;
            let i = toI(dropX, dropY);
            copyMap[i] = powerUp;
            handleSetGameMap(copyMap);
          }, 325);
        }
        break;
      case WebsocketMessageType.ItemPickUp:
        {
          let id = msg.content.UserID;

          let tempUsers: User[] = users.slice();
          let user: User = tempUsers.find((user: User) => {
            return user.getId() === id;
          });

          if (user) {
            let index = tempUsers.findIndex(
              (findUser: User) => findUser.getId() === user.getId()
            );

            switch (msg.content.Item) {
              case 1: {
                if (!tempUsers[index].powerups.includes("bomb")) {
                  tempUsers[index].powerups = [...user.powerups, "bomb"];
                }
                break;
              }
              case 2: {
                if (!tempUsers[index].powerups.includes("bombpower")) {
                  tempUsers[index].powerups = [...user.powerups, "bombpower"];
                }

                if (tempUsers[index].getBombsLeft() == 0) {
                  tempUsers[index].setBombsLeft(1);
                } else {
                  tempUsers[index].setBombsLeft(
                    tempUsers[index].getBombsLeft() + 1
                  );
                }

                tempUsers[index].setBombsTotal(
                  tempUsers[index].getBombsTotal() + 1
                );
                break;
              }
              case 3: {
                if (!tempUsers[index].powerups.includes("flame")) {
                  tempUsers[index].powerups = [...user.powerups, "flame"];
                }
                break;
              }
              case 4: {
                if (!tempUsers[index].powerups.includes("speed")) {
                  tempUsers[index].powerups = [...user.powerups, "speed"];
                }
                break;
              }
            }
          }

          let i = toI(msg.content.X, msg.content.Y);
          copyMap[i] = mapObject.Empty;
          setUsers(tempUsers);
          handleSetGameMap(copyMap);
        }
        break;
      case WebsocketMessageType.SendGameMap:
        {
          handleSetGameMap(msg.content);
        }
        break;
      case WebsocketMessageType.GameStartTimer:
        {
          resetTimer(msg.content.Timer);
        }
        break;
      case WebsocketMessageType.SetOwnID:
        {
          let id = msg.content.ID;
          if (users.find((user: User) => user.id === id)) return; // combats same user being added twice

          let tempUsers: User[] = users.slice();

          let toTempUsers = new User(id, true);
          tempUsers.push(toTempUsers);

          setUsers(tempUsers);
        }
        break;
      case WebsocketMessageType.Disconnect:
        {
          let id = msg.content.ID;
          setUsers(users.filter((user: User) => user.id !== id));
        }
        break;
      case WebsocketMessageType.Connect:
        {
          let id = msg.content.ID;
          if (users.find((user: User) => user.id === id)) return; // combats same user being added twice

          let tempUsers: User[] = users;
          let toTempUsers = new User(id, false, false);

          tempUsers.push(toTempUsers);

          setUsers(tempUsers);
        }
        break;
      case WebsocketMessageType.LoseLife:
        {
          let id = msg.content.ID;
          let tempUsers = users.slice();
          tempUsers = tempUsers.map((user: User) => {
            if (user.id === id) {
              user.lives -= 1;
              if (user.lives === 0) {
                user.is_dead = true;
              }
              return user;
            } else {
              return user;
            }
          });
          setUsers(tempUsers);
        }
        break;
      case WebsocketMessageType.AnnounceWinner:
        {
          try {
            let tempWinner = users.find(
              (user: User) => user.getId() == msg.content.ID
            );
            setWinner(tempWinner.name ? tempWinner.name : "Nameless");
          } catch {
            setWinner(42);
          }
        }
        break;
      default:
        console.log("Unknown event from server: ", msg);
    }
  };
  const mover = (e) => {
    document.removeEventListener("keydown", throttled);
    let copyMap = gameMap.slice();
    switch (e.key) {
      case "ArrowDown": {
        e.preventDefault();
        sendPlayerMoveEvent(socket, MoveDirection.Down);
        break;
      }
      case "ArrowUp": {
        e.preventDefault();
        sendPlayerMoveEvent(socket, MoveDirection.Up);
        break;
      }
      case "ArrowRight": {
        e.preventDefault();
        sendPlayerMoveEvent(socket, MoveDirection.Right);
        break;
      }
      case "ArrowLeft": {
        e.preventDefault();
        sendPlayerMoveEvent(socket, MoveDirection.Left);
        break;
      }
      case " ": {
        e.preventDefault();
        let x = copyMap.indexOf(6);
        if (x < 0) {
          x = copyMap.indexOf(11);
        }
        sendPlayerBombPlace(socket, MoveDirection.Down);
        break;
      }
    }
  };
  throttleFunc = mover;

  let throttled = (e: Event) => {
    return throttledMover(e);
  };
  document.addEventListener("keydown", throttled);

  // Somekind of routing system
  function CurrentPage({ page }) {
    switch (page) {
      case Page.Initial:
        return (
          <InitialPage
            globalMessages={globalMessages}
            setCurrentPage={setCurrentPage}
          />
        );
      case Page.Lobby:
        return (
          <LobbyPage
            lobbyMessages={lobbyMessages}
            users={users.filter((user: User) => user.isInCurrentLobby())}
            timer={timer2}
          />
        );
      case Page.Game:
        return (
          <GamePage
            dir={dir}
            gameMap={gameMap}
            lobbyMessages={lobbyMessages}
            players={users.filter((user: User) => user.isInCurrentLobby())}
            winner={winner}
          />
        );
    }
  }

  return <CurrentPage page={currentPage} />;
}

const container = document.getElementById("root");
Musa.render(<App />, container);
