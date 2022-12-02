import ChatMessage from "./src/models/ChatMessage";
import { MapInfo } from "./src/models/Maps";

export interface WebsocketMessage {
  type: WebsocketMessageType;
  content: any;
}

export enum WebsocketMessageType {
  SetUsername = 1, // Sets username
  GlobalMessage, // Global chat message

  SetOwnID, // Gets own ID from server
  Disconnect, // User disconnects
  Connect, // User connects

  JoinLobby, // User requests to join lobby
  LobbyMessage, // Lobby message
  GameStart, // If there are 2 members users can request to start game
  GameStartTimer, // Notifies users about game starting early
  AnnounceWinner, // sent to users when game ends (1 player alive)

  PlayerMove,
  ThrottledPlayerMove, // Sent to user when PlayerMove is called too fast

  BombPlant, // sent to server when user plants a bomb
  BombExplode, // sent to users when bomb explodes

  ItemDrop, // sent to user when something drops after breaking block
  ItemPickUp, // TODO: sent to users when somebody picks stuff up

  SendGameMap, // sends inital map to players
  LoseLife, // sent to users when a player loses a life (from explosion for example)
}

export type Events =
  | EventGameStartTimer
  | EventThrottledPlayerMove
  | EventLoseLife
  | EventSendGameMap
  | EventItemPickUp
  | EventItemDrop
  | EventBombExplode
  | EventBombPlant
  | EventPlayerMove
  | EventGameStart
  | EventAnnounceWinner
  | EventJoinLobby
  | EventLobbyMessage
  | EventGlobalMessage
  | EventConnectResponse
  | EventSetOwnID
  | EventSetUsername
  | EventDisconnectResponse;

export interface EventGameStartTimer {
  type: WebsocketMessageType.GameStartTimer;
  content: {
    Timer: number; // Number of seconds until game starts
  };
}

export interface EventThrottledPlayerMove {
  type: WebsocketMessageType.ThrottledPlayerMove;
  content: {
    SecondsLeft: number; // Seconds left until throttle is over
  };
}

export interface EventLoseLife {
  type: WebsocketMessageType.LoseLife;
  content: {
    ID: string;
    LivesLeft: number;
  };
}

export interface EventSendGameMap {
  type: WebsocketMessageType.SendGameMap;
  content: MapInfo;
}

export interface EventItemPickUp {
  type: WebsocketMessageType.ItemPickUp;
  content: {
    UserID: string;
    Item: Items;
    X: number;
    Y: number;
  };
}

export interface EventItemDrop {
  type: WebsocketMessageType.ItemDrop;
  content: {
    Item: Items;
    X: number;
    Y: number;
  };
}

export enum Items {
  Bomb = 1,
  BombPowerup,
  FlamePowerup,
  SpeedPowerup,
}

export interface Block {
  X: number;
  Y: number;
}

export interface EventBombExplode {
  type: WebsocketMessageType.BombExplode;
  content: {
    X: number;
    Y: number;
    IsFlamable: boolean;
    AffectedBlocks: Block[];
  };
}

export interface EventBombPlant {
  type: WebsocketMessageType.BombPlant;
  content: {
    X: number;
    Y: number;
  };
}

export interface EventPlayerMove {
  type: WebsocketMessageType.PlayerMove;
  content: {
    ID: string;
    X: number;
    Y: number;
    Direction: MoveDirection;
  };
}

export enum MoveDirection {
  Up = 1,
  Down,
  Right,
  Left,
}

export interface EventGameStart {
  type: WebsocketMessageType.GameStart;
  content: MapInfo;
}

export interface EventAnnounceWinner {
  type: WebsocketMessageType.AnnounceWinner;
  content: {
    ID: string;
  };
}

export interface EventJoinLobby {
  type: WebsocketMessageType.JoinLobby;
  content: {
    ID: string;
    Username: string;
  };
}

export interface EventSetOwnID {
  type: WebsocketMessageType.SetOwnID;
  content: {
    ID: string;
  };
}

export interface EventDisconnectResponse {
  type: WebsocketMessageType.Disconnect;
  content: {
    ID: string;
  };
}

export interface EventConnectResponse {
  type: WebsocketMessageType.Connect;
  content: {
    ID: string;
  };
}

export interface EventLobbyMessage {
  type: WebsocketMessageType.LobbyMessage;
  content: ChatMessage;
}

export interface EventGlobalMessage {
  type: WebsocketMessageType.GlobalMessage;
  content: ChatMessage;
}

export interface EventSetUsername {
  type: WebsocketMessageType.SetUsername;
  content: {
    ID: string;
    Username: string;
  };
}
