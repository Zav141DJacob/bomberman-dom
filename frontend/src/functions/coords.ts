export function toX(i: number) {
  return i % 13;
}

export function toY(i: number) {
  return Math.floor(i / 13);
}

export function toI(x: number, y: number) {
  return y * 13 + x;
}
