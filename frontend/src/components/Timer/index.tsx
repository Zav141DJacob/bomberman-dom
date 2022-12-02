import Musa from "mini-framework";
import { User } from "../../models/User";

let timer2 = 0;

let interval: number;

export function resetTimer(n: number = 29) {
  timer2 = n
}

export function stopTimer() {
  if (interval) clearInterval(interval);
}

export function Timer({ users }) {
  if (!interval && users.filter((u: User) => u.isInCurrentLobby()).length > 1) {
    interval = setInterval(() => {
      let timerNode = document.getElementById("Timer");
      timerNode.innerHTML = timer2.toString();
      if (timer2 <= 10) {
        timerNode.style.color = "red"
      }
      timer2--
      if (timer2 == 0) {
        clearInterval(interval);
      }
    }, 1000);
  }
  return (
    <div id="Timer">
      {timer2}
    </div>
  );
}

