// export type User = {
//     name: string,
//     id: string,
//     in_current_lobby: boolean,
//     is_dead: boolean,
//     is_me: boolean,
//     x: number,
//     y: number,
// }
export class User {
  name: string;
  id: string;

  in_current_lobby: boolean;
  is_dead: boolean;
  is_me: boolean;
  lives: number;

  x: number;
  y: number;
  powerups: string[];

  bombs_left: number;
  bombs_total: number;

  constructor(
    id: string,
    is_me: boolean,
    in_current_lobby?: boolean,
    name?: string,
    lives?: number,
    is_dead?: boolean,
    x?: number,
    y?: number,
    bombs_left?: number,
    bombs_total?: number
  ) {
    this.id = id;
    this.is_me = is_me;
    this.in_current_lobby = in_current_lobby;
    this.name = name;
    this.lives = lives ?? 3;
    this.is_dead = is_dead;
    this.x = x;
    this.y = y;
    this.powerups = [];
    this.bombs_left = bombs_left ?? 1;
    this.bombs_total = bombs_total ?? 1;
  }

  getName(): string {
    return this.name;
  }
  getId(): string {
    return this.id;
  }
  isInCurrentLobby(): boolean {
    return this.in_current_lobby;
  }
  isDead(): boolean {
    return this.is_dead;
  }
  isMe(): boolean {
    return this.is_me;
  }
  getX(): number {
    return this.x;
  }
  getY(): number {
    return this.y;
  }
  getLives(): number {
    return this.lives <= 0 ? 0 : this.lives;
  }
  getBombsLeft(): number {
    return this.bombs_left;
  }
  getBombsTotal(): number {
    return this.bombs_total;
  }

  setName(v: string) {
    this.name = v;
  }
  setId(v: string) {
    this.id = v;
  }
  setInCurrentLobby(v: boolean) {
    this.in_current_lobby = v;
  }
  setIsDead(v: boolean) {
    this.is_dead = v;
  }
  setIsMe(v: boolean) {
    this.is_me = v;
  }
  setX(v: number) {
    this.x = v;
  }
  setY(v: number) {
    this.y = v;
  }
  setLives(v: number) {
    this.lives = v;
  }
  setBombsLeft(v: number) {
    this.bombs_left = v;
  }
  setBombsTotal(v: number) {
    this.bombs_total = v;
  }
}

// export function newUser(id = undefined, is_me = false) {
//     this.name = undefined;
//     this.id = id;
//     this.in_current_lobby = undefined;
//     this.is_dead = undefined;
//     this.is_me = is_me;
//     this.x = undefined;
//     this.y = undefined
// }
