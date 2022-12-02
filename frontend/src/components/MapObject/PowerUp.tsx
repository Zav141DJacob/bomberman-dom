import Musa from "mini-framework";
import MapObject from ".";

function PowerUp({ X, Y, type }) {
  return <MapObject X={X} Y={Y} type={"powerup_" + type} />;
}

export default PowerUp;
