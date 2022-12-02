import Musa from "mini-framework";

const downKey =
  "https://icons.iconarchive.com/icons/chromatix/keyboard-keys/128/arrow-down-icon.png";
const leftKey =
  "https://icons.iconarchive.com/icons/chromatix/keyboard-keys/128/arrow-left-icon.png";
const rightKey =
  "https://icons.iconarchive.com/icons/chromatix/keyboard-keys/128/arrow-right-icon.png";
const upKey =
  "https://icons.iconarchive.com/icons/chromatix/keyboard-keys/128/arrow-up-icon.png";
const spaceKey =
  "https://raw.githubusercontent.com/q2apro/keyboard-keys-speedflips/master/single-keys-blank/200dpi/spacebar.png";

function Helper({ ...props }) {
    return (
        <div className="helper">
            <span style={{ fontSize: "24px", fontWeight: "bold" }}>Help</span>
            <span style={{ fontSize: "18px" }}>Keys</span>
            <div className="keys">
            <span className="movementlabel">Movement</span>
            <img className="downKey" src={downKey} />
            <img className="leftKey" src={leftKey} />
            <img className="rightKey" src={rightKey} />
            <img className="upKey" src={upKey} />
            <span className="spacelabel">Place bomb</span>
            <img className="spaceKey" style={{ width: "150px" }} src={spaceKey} />
            </div>
            <div className="powerups">
            <span style={{ fontSize: "18px" }}>
                Powerups
                <br />
            </span>
            <span>
                Flame - bigger explosion
                <br />
            </span>
            <span>
                Bombs - multibomb
                <br />
            </span>
            <span>
                Speed - faster movement
                <br />
            </span>
            <span>
                Bomb - larger arsenal
                <br />
            </span>
            </div>
        </div>
    )
}

export default Helper;