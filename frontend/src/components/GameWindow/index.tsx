import Musa from "mini-framework";
import Block from "../MapObject/Block";
import Bomb from "../MapObject/Bomb";
import Wall from "../MapObject/Wall";
import Player from "../Player";
import { toI, toX, toY } from "../../functions/coords";
import { mapObject, WINDOW_HEIGHT, WINDOW_WIDTH } from "../../models/Maps";
import { User } from "../../models/User";
import PowerUp from "../MapObject/PowerUp";
import Explosion from "../MapObject/Explosion";

function GameWindow({ gameMap, dir, players, winner, ...props }) {
  console.log(winner);
  var {className, ...props} = props;
  return (
    <div 
    className={ className ? "GameWindow " + className : "GameWindow" }
    { ...props }>
      {gameMap.map((elem: mapObject, i: number) => {
        switch (elem) {
          case mapObject.Wall: {
            return <Wall X={toX(i).toString()} Y={toY(i).toString()} />;
          }
          case mapObject.Block: {
            return <Block X={toX(i).toString()} Y={toY(i).toString()} />;
          }
          case mapObject.Bomb: {
            return (
              <div>
                <Bomb X={toX(i).toString()} Y={toY(i).toString()} />
                {players
                  .map((user: User, j: number) => {
                    if (user.isDead()) {
                      return;
                    }
                    if (toI(user.getX(), user.getY()) == i) {
                      return (
                        <Player
                          dir={dir[j]}
                          index={j}
                          style={{
                            left: toX(i) * (WINDOW_WIDTH / 13) + "px",
                            top: toY(i) * (WINDOW_HEIGHT / 11) + "px",
                          }}
                        />
                      );
                    }
                  })
                  .filter((elem) => elem !== undefined) ?? <div></div>}
              </div>
            );
          }
          case mapObject.BombPowerup: {
            return (
              <PowerUp
                type={"bomb"}
                X={toX(i).toString()}
                Y={toY(i).toString()}
              />
            );
          }
          case mapObject.FlamePowerup: {
            return (
              <PowerUp
                type={"flames"}
                X={toX(i).toString()}
                Y={toY(i).toString()}
              />
            );
          }
          case mapObject.SpeedPowerup: {
            return (
              <PowerUp
                type={"speed"}
                X={toX(i).toString()}
                Y={toY(i).toString()}
              />
            );
          }
          default:
            return (
              <div>
                {players
                  .map((user: User, j: number) => {
                    if (user.isDead()) {
                      return;
                    }
                    if (toI(user.getX(), user.getY()) == i) {
                      return (
                        <Player
                          dir={dir[j]}
                          index={j}
                          style={{
                            left: toX(i) * (WINDOW_WIDTH / 13) + "px",
                            top: toY(i) * (WINDOW_HEIGHT / 11) + "px",
                          }}
                        />
                      );
                    }
                  })
                  .filter((elem) => elem !== undefined) ?? <div></div>}
              </div>
            );
          case mapObject.Explosion:
            return <Explosion X={toX(i).toString()} Y={toY(i).toString()} />;
        }
      })}
      {winner === 42 ? (
        <div className="WinnerScreen">
          <p>There are no winners in war</p>
        </div>
      ) : winner.length > 0 ? (
        <div className="WinnerScreen">
          <p>The winner is {winner.length > 0 ? winner : "Nameless"}!</p>
        </div>
      ) : (
        <div></div>
      )}
    </div>
  );
}

function PotentialPlayer({ players, dir, index }) {
  console.log("made it inside potential");
  var i: number = parseInt(index);
  return (
    players
      .map((user: User, j: number) => {
        if (toI(user.getX(), user.getY()) == i) {
          return (
            <Player
              dir={dir[j]}
              style={{
                left: toX(i) * (WINDOW_WIDTH / 13) + "px",
                top: toY(i) * (WINDOW_HEIGHT / 11) + "px",
              }}
            />
          );
        }
      })
      .filter((elem) => elem !== undefined) ?? <div></div>
  );
}
export default GameWindow;
