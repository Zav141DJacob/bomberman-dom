export const WIDTH = 13;
export const HEIGHT = 11;
export const WINDOW_WIDTH = WIDTH * 45; //585
export const WINDOW_HEIGHT = HEIGHT * 45; //495

export enum mapObject {
  Empty = 0, // Empty square
  Wall, // Indestructible wall
  Block, // Destructible block
  Bomb,

  BombPowerup,
  FlamePowerup,
  SpeedPowerup,

  Explosion,
}

const SIZE: number = 13 * 11;
export type MapInfo = FixedSizeArray<143, mapObject>;

// https://mstn.github.io/2018/06/08/fixed-size-arrays-in-typescript/
type FixedSizeArray<N extends number, T> = N extends 0
  ? never[]
  : {
      0: T;
      length: N;
    } & ReadonlyArray<T>;
