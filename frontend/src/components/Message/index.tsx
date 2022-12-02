import Musa from "mini-framework";
import ChatMessage from "../../models/ChatMessage";

function Message({ msg }) {
  let message: ChatMessage = msg;

  let name = message.Username ? message.Username : "Nameless";
  let id = message.ID.slice(0, 5);

  let usernameFormat = `${name} (ID ${id})`;

  let date = new Date(Date.parse(message.Time.split(" +")[0]));
  let hours = date.getHours();
  let what =
    date.getMinutes() / 10 < 1 ? "0" + date.getMinutes() : date.getMinutes(); // ???

  let time = `${hours}:${what}`;

  return (
    <div className="Message">
      <div className="Infobar">
        <span className="Username">{usernameFormat}</span>
        <span className="Date">{time}</span>
      </div>
      <div className="ContentRow">{message.Content}</div>
    </div>
  );
}

export default Message;
