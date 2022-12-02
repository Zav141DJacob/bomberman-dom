import Musa from "mini-framework";
import { IndexToPlayerClass } from "../Player";

const flameIcon =
  "https://upload.wikimedia.org/wikipedia/commons/thumb/a/ae/Nuclear_symbol.svg/40px-Nuclear_symbol.svg.png";
const speedIcon =
  "https://d29fhpw069ctt2.cloudfront.net/icon/image/39214/preview.png";
const otherIcon = "https://cdn-icons-png.flaticon.com/512/523/523777.png";

function Dashboard({ players, dir }) {
  return (
    <div className="dashboard">
      {players.map((player, i) => (
        <PlayerInfo player={player} dir={dir} index={i} />
      ))}
    </div>
  );
}

function PlayerInfo({ player, index, dir }) {
  let dead = player.lives === 0;

  let extraStats = dead ? (
    <span>DEAD</span>
  ) : (
    <div>
      <div className="powerupIcons">
        {player.powerups.map((i) =>
          i === "flame" ? (
            <img
              style={{ position: "relative", width: "30px" }}
              src={flameIcon}
            />
          ) : i === "speed" ? (
            <img
              style={{ position: "relative", width: "30px" }}
              src={speedIcon}
            />
          ) : (
            <img
              style={{ position: "relative", width: "30px" }}
              src={otherIcon}
            />
          )
        )}
      </div>
      <div>
        Lives: <span style={{ color: "red" }}>{"❤".repeat(player.getLives())}</span>
      </div>
      <div>
        Bombs: {player.bombs_left}/{player.bombs_total}
      </div>
    </div>
  );

  let characterClass = dir[index].includes("stand") ? "Avatarstand" : "Avatar";

  let playerName =
    (player.name ? player.name : "Nameless") + (player.is_me ? " (YOU)" : "");

  let miscPlayerClass =
    dir[index] === "stand"
      ? "Characterstanddown"
      : dir[index].includes("stand")
      ? "Character" + dir[index]
      : "Characterstand" + dir[index];

  return (
    <div className="playerInfo">
      <div className="playerIcon">
        <div
          className={`Character ${miscPlayerClass}`}
          style={{ transform: "scale(2)" }}
        >
          <div className="Character_shadow pixelart" />
          <div
            className={`${characterClass} pixelart ${IndexToPlayerClass(
              index
            )}`}
          />
        </div>
      </div>
      {dead ? <img className="deathIcon" /> : <div></div>}
      <div className="playerStats">
        <span>Name: {playerName}</span>
        {extraStats}
      </div>
    </div>
  );
}

export default Dashboard;
