import Musa from "mini-framework";
import MapObject from ".";

function Wall({ X, Y }) {
  return <MapObject X={X} Y={Y} type={"Wall"} />;
}

export default Wall;
