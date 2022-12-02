import {
  MoveDirection,
  WebsocketMessage,
  WebsocketMessageType,
} from "../../websocket";

export function sendGlobalMessage(s: WebSocket, text: string) {
  let msg: WebsocketMessage = {
    type: WebsocketMessageType.GlobalMessage,
    content: text,
  };

  s.send(JSON.stringify(msg));
}

export function sendLobbyMessage(s: WebSocket, text: string) {
  let msg: WebsocketMessage = {
    type: WebsocketMessageType.LobbyMessage,
    content: text,
  };

  s.send(JSON.stringify(msg));
}

export function sendLobbyJoinEvent(s: WebSocket) {
  console.log("in sendlobbyjoinevent from frontend to backend");
  let msg: WebsocketMessage = {
    type: WebsocketMessageType.JoinLobby,
    content: null,
  };

  s.send(JSON.stringify(msg));
}

export function sendGameStartEvent(s: WebSocket) {
  let msg: WebsocketMessage = {
    type: WebsocketMessageType.GameStart,
    content: null,
  };

  s.send(JSON.stringify(msg));
}

export function sendPlayerMoveEvent(s: WebSocket, direction: MoveDirection) {
  let msg: WebsocketMessage = {
    type: WebsocketMessageType.PlayerMove,
    content: direction,
  };

  s.send(JSON.stringify(msg));
}

export function sendPlayerBombPlace(s: WebSocket, direction: MoveDirection) {
  let msg: WebsocketMessage = {
    type: WebsocketMessageType.BombPlant,
    content: direction.toString(),
  };

  s.send(JSON.stringify(msg));
}

export function sendSetUsername(s: WebSocket, username: string) {
  let msg: WebsocketMessage = {
    type: WebsocketMessageType.SetUsername,
    content: username,
  };

  s.send(JSON.stringify(msg));
}
