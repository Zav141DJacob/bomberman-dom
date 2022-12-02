import Musa from "mini-framework";
import MapObject from ".";

function Block({ X, Y }) {
  return <MapObject X={X} Y={Y} type={"Block"} />;
}

export default Block;
