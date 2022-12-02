import Musa from "mini-framework";
import { WINDOW_HEIGHT, WINDOW_WIDTH } from "../../models/Maps";

function MapObject({ X, Y, type }) {
  return (
    <div
      style={{
        left: parseInt(X) * (WINDOW_WIDTH / 13) + "px",
        top: parseInt(Y) * (WINDOW_HEIGHT / 11) + "px",
      }}
      className={type}
    ></div>
  );
}

export default MapObject;
