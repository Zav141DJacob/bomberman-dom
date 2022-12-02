import Musa from "mini-framework";
import { User } from "../models/User";

export function UsersList({ users, ...props }) {
  if (users.length === 0) {
    return <div />;
  }
  var { className, ...props } = props;
  return (
    <div
    className={className ? "UsersList " + className : "UsersList"}>
      <ol>
        {users.map((user: User) => {
          let name =
            (user.name ? user.name : "Nameless") + (user.is_me ? " (YOU)" : "");

          return (
            <li>
              {name} ({user.id})
            </li>
          );
        })}
      </ol>
    </div>
  );
}

export default UsersList;
