import { createContext, useEffect, useState } from "react";
import Cookies from "js-cookie";

const UserContext = createContext(null);

interface IUser {
  loggedIn: boolean;
}

function UserProvider({ children }) {
  const [userState, setUserState] = useState<IUser>({
    loggedIn: false,
  });

  useEffect(() => {
    async function loadUser() {
      try {
        if (Cookies.get("token")) {
          setUserState({ loggedIn: true });
        }
      } catch (err) {}
    }

    loadUser();
  }, []);

  return (
    <UserContext.Provider value={userState}>{children}</UserContext.Provider>
  );
}

const UserConsumer = UserContext.Consumer;

export default UserProvider;
export { UserConsumer, UserContext };
