import Musa from "mini-framework";
import MapItem from ".";

function Explosion({ X, Y }) {
  return <MapItem X={X} Y={Y} type={"Explosion"} />;
}

export default Explosion;
