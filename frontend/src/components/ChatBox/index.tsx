import Musa from "mini-framework";
import ChatMessage from "../../models/ChatMessage";
import Message from "../Message";
let message = "";


function ChatBox({ title, messages, sendMessage, ...props }) {
  return (
    <div
    { ...props }>
      <div className="SmallTitle">{title}</div>
      <div className="chatbox">
        <div className="messages">
          {messages.map((msg: ChatMessage) => (
            <Message msg={msg} />
          ))}
        </div>
        <form className="input" onSubmit={sendMessage}>
          <input
            id="chatBoxInput"
            value={message}
            type="text"
            placeholder="Type message.."
            name="message"
            autocomplete="off"
            required
            autofocus
          ></input>

          <button type="submit" class="btn">
            Send
          </button>
        </form>
      </div>
    </div>
  );
}

export default ChatBox;
