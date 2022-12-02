import Musa from "mini-framework";

export function IndexToPlayerClass(index: number): string {
  switch (index) {
    case 0:
      return "Player1";
    case 1:
      return "Player2";
    case 2:
      return "Player3";
    default:
      return "Player4";
  }
}

function Player({ dir, ...props }) {
  var { index, ...props } = props;

  let playerClass = IndexToPlayerClass(index);

  let miscPlayerClass =
    dir === "stand" ? "Characterstand_spritesheet" : "Character_spritesheet";

  let miscPlayerClass2 =
    dir === "stand" ? "Characterstanddown" : "Character" + dir;

  return (
    <div className={`Character ${miscPlayerClass2}`} style={props.style}>
      <div className="Character_shadow pixelart" />
      <div className={`${miscPlayerClass} pixelart ${playerClass}`} />
      {props.children}
    </div>
  );
}

export default Player;
