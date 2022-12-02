import Musa from "mini-framework";
import MapObject from ".";

function Bomb({ X, Y }) {
  return <MapObject X={X} Y={Y} type={"Bomb"} />;
}

export default Bomb;
